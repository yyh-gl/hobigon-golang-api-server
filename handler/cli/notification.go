package cli

import (
	"context"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

//////////////////////////////////////////////////
// NewNotificationHandler
//////////////////////////////////////////////////

// NotificationHandler : Slack 通知用のハンドラーインターフェース
type NotificationHandler interface {
	NotifyTodayTasksToSlack(c *cli.Context) error
	NotifyTodayBirthdayToSlack(c *cli.Context) error
	NotifyAccessRankingToSlack(c *cli.Context) error
}

type notificationHandler struct {
	nu usecase.NotificationUseCase
}

// NewNotificationHandler : Slack 通知用のハンドラーを取得
func NewNotificationHandler(nu usecase.NotificationUseCase) NotificationHandler {
	return &notificationHandler{
		nu: nu,
	}
}

//////////////////////////////////////////////////
// NotifyTodayTasksToSlack
//////////////////////////////////////////////////

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (snh notificationHandler) NotifyTodayTasksToSlack(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := snh.nu.NotifyTodayTasksToSlack(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

//////////////////////////////////////////////////
// NotifyTodayBirthdayToSlack
//////////////////////////////////////////////////

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (snh notificationHandler) NotifyTodayBirthdayToSlack(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := snh.nu.NotifyTodayBirthdayToSlack(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

//////////////////////////////////////////////////
// NotifyAccessRankingToSlack
//////////////////////////////////////////////////

// NotifyAccessRankingToSlack : アクセスランキングを Slack に通知
func (snh notificationHandler) NotifyAccessRankingToSlack(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := snh.nu.NotifyAccessRanking(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}
