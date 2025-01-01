package di

import (
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest"
)

// ContainerAPI : API用DIコンテナ
type ContainerAPI struct {
	HandlerBlog         rest.Blog
	HandlerNotification rest.Notification

	DB *db.DB
}

// ContainerCLI : CLI用DIコンテナ
type ContainerCLI struct {
	HandlerNotification cli.Notification

	DB *db.DB
}
