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
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest"
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
	r.PathPrefix("/").Handler(wrapHandler(preflightHandler)).Methods(http.MethodOptions)

	// Health Check
	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Blog handlers
	r.HandleFunc("/api/v1/blogs", wrapHandler(diContainer.HandlerBlog.Create)).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/blogs/{title}", wrapHandler(diContainer.HandlerBlog.Show)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/blogs/{title}/like", wrapHandler(diContainer.HandlerBlog.Like)).Methods(http.MethodPost)

	// Birthday handler
	r.HandleFunc("/api/v1/birthday", wrapHandler(diContainer.HandlerBirthday.Create)).Methods(http.MethodPost)

	// Notification handlers
	r.HandleFunc("/api/v1/notifications/slack/tasks/today", wrapHandler(diContainer.HandlerNotification.NotifyTodayTasksToSlack)).Methods(http.MethodPost)
	// TODO: 誕生日の人が複数いたときに対応
	r.HandleFunc("/api/v1/notifications/slack/birthdays/today", wrapHandler(diContainer.HandlerNotification.NotifyTodayBirthdayToSlack)).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/notifications/slack/rankings/access", wrapHandler(diContainer.HandlerNotification.NotifyAccessRankingToSlack)).Methods(http.MethodPost)

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

// wrapHandler : 全ハンドラー共通処理
func wrapHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエスト内容をログ出力
		// TODO: Bodyの内容を出力
		app.Logger.Println("[AccessLog] " + r.Method + " " + r.URL.String())

		// CORS用ヘッダーを付与
		switch {
		case app.IsPrd():
			w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
		case app.IsDev() || app.IsTest():
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3001")
		}
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		// リクエスト数をカウント by Prometheus
		rest.CountRequest(r.Method, r.URL.Path)

		h.ServeHTTP(w, r)
	}
}

// preflightHandler : preflight用のハンドラー
func preflightHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
