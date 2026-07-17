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
	NotifyPokemonEventToSlack(c *cli.Context) error
	NotifyCoopPaymentReminderToLine(c *cli.Context) error
	NotifySeisenkanPaymentReminderToLine(c *cli.Context) error
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

// NotifyPokemonEventToSlack : ポケモンカードのイベント情報をSlackに通知
func (n notification) NotifyPokemonEventToSlack(c *cli.Context) error {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.ContextKeyCLI, c)

	if _, err := n.u.NotifyPokemonEvent(ctx); err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// NotifyCoopPaymentReminderToLine : 生協の入金リマインドをLINEに通知
func (n notification) NotifyCoopPaymentReminderToLine(c *cli.Context) error {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.ContextKeyCLI, c)

	if err := n.u.NotifyCoopPaymentReminderToLine(ctx); err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// NotifySeisenkanPaymentReminderToLine : 生鮮館の入金リマインドをLINEに通知
func (n notification) NotifySeisenkanPaymentReminderToLine(c *cli.Context) error {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.ContextKeyCLI, c)

	if err := n.u.NotifySeisenkanPaymentReminderToLine(ctx); err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}
