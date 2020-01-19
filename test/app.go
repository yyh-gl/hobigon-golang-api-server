package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// Client : REST Handlerのテストに関する機能を提供するインターフェース
type Client struct {
	app         *httprouter.Router
	DIContainer *di.ContainerAPI
}

func NewClient() *Client {
	// 依存関係を定義
	diContainer := initTestApp()

	// ロガー設定
	// TODO: いちいちdi.Containerにバインドする意味があるのかもう一度検討
	app.Logger = diContainer.Logger

	// ルーティング設定
	r := httprouter.New()

	return &Client{
		app:         r,
		DIContainer: diContainer,
	}
}

func (c Client) Get(handler func(http.ResponseWriter, *http.Request), path string) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, path, handler)

	req, _ := http.NewRequest(http.MethodGet, path, nil)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func (c Client) Post(handler func(http.ResponseWriter, *http.Request), path string, json string) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, path, handler)

	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewBuffer([]byte(json)))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
