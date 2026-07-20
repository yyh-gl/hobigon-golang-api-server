package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

func main() {
	log.NewLogger()

	diContainer := initApp()

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
			Name:    "notify-pokemon-event",
			Aliases: []string{"npe"},
			Usage:   "Notify the Pokémon card event to Slack",
			Action:  diContainer.HandlerNotification.NotifyPokemonEventToSlack,
		},
		{
			Name:    "notify-to-line",
			Aliases: []string{"ntl"},
			Usage:   "Notify a message to a LINE bot, selected by --bot-key and --message-key",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "bot-key",
					Usage:    "target LINE bot key (e.g. son)",
					Required: true,
				},
				cli.StringFlag{
					Name:     "message-key",
					Usage:    "message key defined in config/line_messages.yaml (e.g. coop, seisenkan)",
					Required: true,
				},
			},
			Action: diContainer.HandlerNotification.NotifyToLINE,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Error(context.Background(), fmt.Errorf("failed to cliApp.Run(): %w", err))
		os.Exit(1)
	}
}
