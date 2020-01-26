package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func main() {
	// 依存関係を定義
	diContainer := initApp()
	defer func() { _ = diContainer.DB.Close() }()

	// ロガー設定
	// TODO: いちいちdi.Containerにバインドする意味があるのかもう一度検討
	app.Logger = diContainer.Logger

	cliApp := cli.NewApp()

	cliApp.Name = "Hobigon CLI"
	cliApp.Usage = "This app can execute some commands in Hobigon."
	cliApp.Version = "0.0.1"

	// コマンドオプションを設定
	cliApp.Flags = []cli.Flag{}

	// コマンドを設定
	cliApp.Commands = []cli.Command{
		{
			Name:    "notify-today-tasks",
			Aliases: []string{"ntt"},
			Usage:   "Notify the today's tasks to Slack",
			Action:  diContainer.HandlerNotification.NotifyTodayTasksToSlack,
		},
		{
			Name:    "notify-today-birthday",
			Aliases: []string{"ntb"},
			Usage:   "Notify the today's birthday to Slack",
			Action:  diContainer.HandlerNotification.NotifyTodayBirthdayToSlack,
		},
		{
			Name:    "notify-access-ranking",
			Aliases: []string{"nar"},
			Usage:   "Notify the access ranking to Slack",
			Action:  diContainer.HandlerNotification.NotifyAccessRankingToSlack,
		},
	}

	app.Logger.Print("[CLI-ExecuteLog] $ hobi " + os.Args[1])

	if err := cliApp.Run(os.Args); err != nil {
		app.Logger.Panic(errors.Wrap(err, "cliApp.Run()内でのエラー"))
	}
}
