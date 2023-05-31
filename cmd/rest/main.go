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

	// Blog handlers
	blogCreateFunc := middleware.InstrumentPrometheus(diContainer.HandlerBlog.Create, "blog_create")
	r.HandleFunc("/api/v1/blogs", middleware.Attach(blogCreateFunc)).Methods(http.MethodPost)
	blogShowFunc := middleware.InstrumentPrometheus(diContainer.HandlerBlog.Show, "blog_show")
	r.HandleFunc("/api/v1/blogs/{title}", middleware.Attach(blogShowFunc)).Methods(http.MethodGet)
	blogLikeFunc := middleware.InstrumentPrometheus(diContainer.HandlerBlog.Like, "blog_like")
	r.HandleFunc("/api/v1/blogs/{title}/like", middleware.Attach(blogLikeFunc)).Methods(http.MethodPost)

	// Birthday handler
	birthdayCreateFunc := middleware.InstrumentPrometheus(diContainer.HandlerBirthday.Create, "birthday_create")
	r.HandleFunc("/api/v1/birthday", middleware.Attach(birthdayCreateFunc)).Methods(http.MethodPost)

	// Notification handlers
	notificationTaskFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerNotification.NotifyTodayTasksToSlack,
		"notification_notify_today_tasks_to_slack",
	)
	r.HandleFunc(
		"/api/v1/notifications/slack/tasks/today",
		middleware.Attach(notificationTaskFunc),
	).Methods(http.MethodPost)
	// TODO: 誕生日の人が複数いたときに対応
	notificationBirthdayFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerNotification.NotifyTodayBirthdayToSlack,
		"notification_notify_today_birthday_to_slack",
	)
	r.HandleFunc(
		"/api/v1/notifications/slack/birthdays/today",
		middleware.Attach(notificationBirthdayFunc),
	).Methods(http.MethodPost)
	notificationRankingFunc := middleware.InstrumentPrometheus(
		diContainer.HandlerNotification.NotifyAccessRankingToSlack,
		"notification_notify_access_ranking_to_slack",
	)
	r.HandleFunc(
		"/api/v1/notifications/slack/rankings/access",
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
