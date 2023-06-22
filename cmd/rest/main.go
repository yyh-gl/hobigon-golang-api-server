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
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest/middleware"
)

// version : アプリケーションのバージョン情報（GitHubのReleasesに対応）
// build時に埋め込む
var version string

func main() {
	app.NewLogger()

	diContainer := initApp()
	defer func() { _ = diContainer.DB.Close() }()

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
	r.HandleFunc("/api/debug", middleware.Attach(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).Methods(http.MethodGet)
	r.HandleFunc("/api/debug", middleware.Attach(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})).Methods(http.MethodPost)

	// Blog handlers
	blogCreatePath := "/api/v1/blogs"
	blogCreateFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerBlog.Create,
		"blog_create",
		blogCreatePath,
	)
	r.HandleFunc(blogCreatePath, middleware.Attach(blogCreateFunc)).Methods(http.MethodPost)
	blogShowPath := "/api/v1/blogs/{title}"
	blogShowFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerBlog.Show,
		"blog_show",
		blogShowPath,
	)
	r.HandleFunc(blogShowPath, middleware.Attach(blogShowFunc)).Methods(http.MethodGet)
	blogLikePath := "/api/v1/blogs/{title}/like"
	blogLikeFunc := middleware.InstrumentPrometheus(diContainer.HandlerBlog.Like, "blog_like", blogLikePath)
	r.HandleFunc(blogLikePath, middleware.Attach(blogLikeFunc)).Methods(http.MethodPost)

	// Birthday handler
	birthdayCreatePath := "/api/v1/birthday"
	birthdayCreateFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerBirthday.Create,
		"birthday_create",
		birthdayCreatePath,
	)
	r.HandleFunc(birthdayCreatePath, middleware.Attach(birthdayCreateFunc)).Methods(http.MethodPost)

	// Notification handlers
	notificationTaskPath := "/api/v1/notifications/slack/tasks/today"
	notificationTaskFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerNotification.NotifyTodayTasksToSlack,
		"notification_notify_today_tasks_to_slack",
		notificationTaskPath,
	)
	r.HandleFunc(
		notificationTaskPath,
		middleware.Attach(notificationTaskFunc),
	).Methods(http.MethodPost)
	// TODO: 誕生日の人が複数いたときに対応
	notificationBirthdayPath := "/api/v1/notifications/slack/birthdays/today"
	notificationBirthdayFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerNotification.NotifyTodayBirthdayToSlack,
		"notification_notify_today_birthday_to_slack",
		notificationBirthdayPath,
	)
	r.HandleFunc(
		notificationBirthdayPath,
		middleware.Attach(notificationBirthdayFunc),
	).Methods(http.MethodPost)
	notificationRankingPath := "/api/v1/notifications/slack/rankings/access"
	notificationRankingFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerNotification.NotifyAccessRankingToSlack,
		"notification_notify_access_ranking_to_slack",
		notificationRankingPath,
	)
	r.HandleFunc(
		notificationRankingPath,
		middleware.Attach(notificationRankingFunc),
	).Methods(http.MethodPost)

	r.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Println("========================")
		fmt.Println("Server Start >> http://localhost" + s.Addr)
		fmt.Println("========================")
		middleware.CountUpRunningVersion(version)
		errCh <- s.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		fmt.Println("Error happened:", err.Error())
	case sig := <-sigCh:
		fmt.Println("Signal received:", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Graceful shutdown failed:", err.Error())
	}
	middleware.CountDownRunningVersion(version)
	fmt.Println("Server shutdown")
}
