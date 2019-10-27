package cli

import (
	"context"

	"github.com/urfave/cli"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

// NotifyTodayTaskToSlackHandler は今日のタスク一覧を Slack に通知
func NotifyTodayTasksToSlackHandler(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := usecase.NotifyTodayTasksToSlackUseCase(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

// NotifyTodayBirthdayToSlackHandler は今日誕生日の人を Slack に通知
func NotifyTodayBirthdayToSlackHandler(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := usecase.NotifyTodayBirthdayToSlackUseCase(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

// NotifyAccessRankingHandler はアクセスランキングを Slack に通知
func NotifyAccessRankingToSlackHandler(c *cli.Context) error {
	logger := app.Logger

	ctx := context.TODO()
	ctx = context.WithValue(ctx, app.CliContextKey, c)

	if err := usecase.NotifyAccessRankingUseCase(ctx); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}
