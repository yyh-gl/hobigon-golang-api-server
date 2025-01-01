package cli

import (
	"context"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

// Notification : Notification用CLIサービスのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(c *cli.Context) error
}

type notification struct {
	u usecase.Notification
}

// NewNotification : Notification用CLIサービスを取得
func NewNotification(u usecase.Notification) Notification {
	return &notification{
		u: u,
	}
}

// NotifyTodayTasksToSlack : 今日のタスク一覧をSlackに通知
func (n notification) NotifyTodayTasksToSlack(c *cli.Context) error {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.ContextKeyCLI, c)

	if _, err := n.u.NotifyTodayTasksToSlack(ctx); err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}
