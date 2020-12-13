package di

import (
	"log"

	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/http"
)

// ContainerAPI : API用DIコンテナ
type ContainerAPI struct {
	HandlerBlog         http.Blog
	HandlerBirthday     http.Birthday
	HandlerNotification http.Notification

	DB     *db.DB
	Logger *log.Logger
}

// ContainerCLI : CLI用DIコンテナ
type ContainerCLI struct {
	HandlerNotification cli.Notification

	DB     *db.DB
	Logger *log.Logger
}
