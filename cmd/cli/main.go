package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

func main() {
	diContainer := initApp()
	defer func() {
		sqlDB, err := diContainer.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}()

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
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Error(context.Background(), fmt.Errorf("failed to cliApp.Run(): %w", err))
		os.Exit(1)
	}
}
