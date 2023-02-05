package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func main() {
	app.NewLogger()

	diContainer := initApp()
	defer func() { _ = diContainer.DB.Close() }()

	cliApp := cli.NewApp()

	cliApp.Name = "Hobigon CLI"
	cliApp.Usage = "This app can execute some commands in Hobigon."
	cliApp.Version = "0.0.1"

	cliApp.Flags = []cli.Flag{}
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

	if err := cliApp.Run(os.Args); err != nil {
		app.Error(fmt.Errorf("cliApp.Run()内でのエラー: %w", err))
		os.Exit(1)
	}
}
