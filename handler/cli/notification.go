package cli

import (
	"context"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

//////////////////////////////////////////////////
// NewSlackNotificationHandler
//////////////////////////////////////////////////

// SlackNotificationHandler : Slack 通知用のハンドラーインターフェース
type SlackNotificationHandler interface {
	NotifyTodayTasks(c *cli.Context) error
	NotifyTodayBirthday(c *cli.Context) error
	NotifyAccessRanking(c *cli.Context) error
}

type slackNotificationHandler struct {
	nu usecase.NotificationUseCase
}

// NewSlackNotificationHandler : Slack 通知用のハンドラーを取得
func NewSlackNotificationHandler(nu usecase.NotificationUseCase) SlackNotificationHandler {
	return &slackNotificationHandler{
		nu: nu,
	}
}

// NotifyTodayTasks : 今日のタスク一覧を Slack に通知
func (snh slackNotificationHandler) NotifyTodayTasks(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := snh.nu.NotifyTodayTasksToSlack(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

// NotifyTodayBirthday : 今日誕生日の人を Slack に通知
func (snh slackNotificationHandler) NotifyTodayBirthday(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := snh.nu.NotifyTodayBirthdayToSlack(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

// NotifyAccessRanking : アクセスランキングを Slack に通知
func (snh slackNotificationHandler) NotifyAccessRanking(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := snh.nu.NotifyAccessRanking(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}
