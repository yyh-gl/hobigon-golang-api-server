package cli

import (
	"context"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

// Notification : Notification用CLIサービスのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(c *cli.Context) error
	NotifyTodayBirthdayToSlack(c *cli.Context) error
	NotifyAccessRankingToSlack(c *cli.Context) error
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
	ctx = context.WithValue(ctx, app.CLIContextKey, c)

	if _, err := n.u.NotifyTodayTasksToSlack(ctx); err != nil {
		app.Error(err)
		return err
	}
	return nil
}

// NotifyTodayBirthdayToSlack : 今日誕生日の人をSlackに通知
func (n notification) NotifyTodayBirthdayToSlack(c *cli.Context) error {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CLIContextKey, c)

	if _, err := n.u.NotifyTodayBirthdayToSlack(ctx); err != nil {
		app.Error(err)
		return err
	}
	return nil
}

// NotifyAccessRankingToSlack : アクセスランキングをSlackに通知
func (n notification) NotifyAccessRankingToSlack(c *cli.Context) error {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CLIContextKey, c)

	if _, err := n.u.NotifyAccessRanking(ctx); err != nil {
		app.Error(err)
		return err
	}
	return nil
}
