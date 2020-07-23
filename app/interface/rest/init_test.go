package rest_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/di"

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
	// Blog handlers
	Router.HandleFunc("/api/v1/blogs", DIContainer.HandlerBlog.Create).Methods(http.MethodPost)
	Router.HandleFunc("/api/v1/blogs/{title}", DIContainer.HandlerBlog.Show).Methods(http.MethodGet)
	Router.HandleFunc("/api/v1/blogs/{title}/like", DIContainer.HandlerBlog.Like).Methods(http.MethodPost)

	// Birthday handler
	Router.HandleFunc("/api/v1/birthday", DIContainer.HandlerBirthday.Create).Methods(http.MethodPost)

	os.Exit(m.Run())
}
