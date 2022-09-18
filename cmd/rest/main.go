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

func main() {
	// 依存関係を定義
	diContainer := initApp()
	defer func() { _ = diContainer.DB.Close() }()

	// ロガー設定
	// TODO: いちいちdi.Containerにバインドする意味があるのかもう一度検討
	app.Logger = diContainer.Logger

	r := mux.NewRouter()

	// Preflight handler
	r.PathPrefix("/").Handler(middleware.Attach(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, "pleflight")).Methods(http.MethodOptions)

	// Health Check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Blog handlers
	r.HandleFunc("/api/v1/blogs", middleware.Attach(
		diContainer.HandlerBlog.Create, "blog_create"),
	).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/blogs/{title}", middleware.Attach(
		diContainer.HandlerBlog.Show, "blog_show"),
	).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/blogs/{title}/like", middleware.Attach(
		diContainer.HandlerBlog.Like,
		"blog_like"),
	).Methods(http.MethodPost)

	// Birthday handler
	r.HandleFunc("/api/v1/birthday", middleware.Attach(
		diContainer.HandlerBirthday.Create, "birthday_create"),
	).Methods(http.MethodPost)

	// Notification handlers
	r.HandleFunc("/api/v1/notifications/slack/tasks/today", middleware.Attach(
		diContainer.HandlerNotification.NotifyTodayTasksToSlack,
		"notification_notify_today_tasks_to_slack"),
	).Methods(http.MethodPost)
	// TODO: 誕生日の人が複数いたときに対応
	r.HandleFunc("/api/v1/notifications/slack/birthdays/today", middleware.Attach(
		diContainer.HandlerNotification.NotifyTodayBirthdayToSlack,
		"notification_notify_today_birthday_to_slack")).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/notifications/slack/rankings/access", middleware.Attach(diContainer.HandlerNotification.NotifyAccessRankingToSlack, "notification_notify_access_ranking_to_slack")).Methods(http.MethodPost)

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
		fmt.Println(" ↳  Log File -> " + os.Getenv("LOG_PATH") + "/" + app.APILogFilename)
		fmt.Println("========================")
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
	fmt.Println("Server shutdown")
}
