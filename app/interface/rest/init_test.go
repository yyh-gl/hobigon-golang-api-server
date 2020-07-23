package rest_test

import (
	"os"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/cmd/api/di"

	"github.com/gorilla/mux"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/test"
)

var (
	DIContainer *di.ContainerAPI
	Router      *mux.Router
)

func TestMain(m *testing.M) {
	DIContainer = test.InitTestApp()
	defer func() { _ = DIContainer.DB.Close() }()

	// TODO: いちいちdi.Containerにバインドする意味があるのかもう一度検討
	app.Logger = DIContainer.Logger

	Router = mux.NewRouter()
	Router.HandleFunc("/api/v1/blogs/{title}", DIContainer.HandlerBlog.Show)

	code := m.Run()

	os.Exit(code)
}
