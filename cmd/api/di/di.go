package di

import (
	"log"

	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/rest"
)

// ContainerAPI : API用DIコンテナ
type ContainerAPI struct {
	HandlerBlog         rest.Blog
	HandlerBirthday     rest.Birthday
	HandlerNotification rest.Notification

	DB     *db.DB
	Logger *log.Logger
}

// ContainerCLI : CLI用DIコンテナ
type ContainerCLI struct {
	HandlerNotification cli.NotificationHandler

	DB     *db.DB
	Logger *log.Logger
}
