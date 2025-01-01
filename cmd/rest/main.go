package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest/middleware"
	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"
)

func main() {
	ctx := context.Background()

	log.NewLogger()

	diContainer := initApp()
	defer func() { _ = diContainer.DB.Close() }()

	router := newRouter(diContainer)

	s := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Info(ctx, "server start >> http://localhost"+s.Addr)
		middleware.CountUpRunningVersion(app.GetVersion())
		errCh <- s.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		log.Error(ctx, fmt.Errorf("received error signal: %w", err))
	case sig := <-sigCh:
		log.Info(ctx, "received signal: "+sig.String())
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Error(ctx, fmt.Errorf("graceful shutdown is failed: %w", err))
	}
	middleware.CountDownRunningVersion(app.GetVersion())
	log.Info(ctx, "server shutdown is succeeded")
}

func newRouter(diContainer *di.ContainerAPI) *mux.Router {
	r := mux.NewRouter()

	// Preflight handler
	r.PathPrefix("/").Handler(middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).Methods(http.MethodOptions)

	// Health Check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Debug Handlers
	debugGetFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(debugGetFunc, "/api/debug", "debug_get"),
	).Methods(http.MethodGet)

	debugPostFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(debugPostFunc, "/api/debug", "debug_post"),
	).Methods(http.MethodPost)

	// Blog handlers
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerBlog.Create,
			"/api/v1/blogs",
			"blog_create",
		),
	).Methods(http.MethodPost)

	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerBlog.Show,
			"/api/v1/blogs/{title}",
			"blog_show",
		),
	).Methods(http.MethodGet)

	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerBlog.Like,
			"/api/v1/blogs/{title}/like",
			"blog_like",
		),
	).Methods(http.MethodPost)

	// Notification handlers
	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerNotification.NotifyTodayTasksToSlack,
			"/api/v1/notifications/slack/tasks/today",
			"today_tasks_notification_send",
		),
	).Methods(http.MethodPost)

	r.HandleFunc(
		middleware.CreateHandlerFuncWithMiddleware(
			diContainer.HandlerNotification.NotifyPokemonEventToSlack,
			"/api/v1/notifications/slack/pokemon/events",
			"pokemon_events_notification_send",
		),
	).Methods(http.MethodPost)

	r.Handle("/metrics", promhttp.Handler())

	return r
}
