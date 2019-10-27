package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	myCLI "github.com/yyh-gl/hobigon-golang-api-server/handler/cli"
)

func main() {
	// システム共通で使用するものを用意
	//  -> logger, DB
	app.Init(app.CLiLogFilename)

	logger := app.Logger

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
			Action:  myCLI.NotifyTodayTasksToSlackHandler,
		},
		{
			Name:    "notify-today-birthday",
			Aliases: []string{"ntb"},
			Usage:   "Notify the today's birthday to Slack",
			Action:  myCLI.NotifyTodayBirthdayToSlackHandler,
		},
		{
			Name:    "notify-access-ranking",
			Aliases: []string{"nar"},
			Usage:   "Notify the access ranking to Slack",
			Action:  myCLI.NotifyAccessRankingToSlackHandler,
		},
	}

	logger.Print("[CLI-ExecuteLog] $ hobi " + os.Args[1])

	if err := cliApp.Run(os.Args); err != nil {
		panic("cliApp.Run内でのエラー")
	}
}
