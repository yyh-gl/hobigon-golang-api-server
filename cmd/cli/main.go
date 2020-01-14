package main

func main() {
	//// ロガー初期化
	//app.Logger = app.NewLogger(app.CLILogFilename)
	//logger := app.Logger
	//
	//cliApp := cli.NewApp()
	//
	//cliApp.Name = "Hobigon CLI"
	//cliApp.Usage = "This app can execute some commands in Hobigon."
	//cliApp.Version = "0.0.1"
	//
	//// コマンドオプションを設定
	//cliApp.Flags = []cli.Flag{}
	//
	//// 依存関係を定義
	//taskGateway := igateway.NewTaskGateway()
	//slackGateway := igateway.NewSlackGateway()
	//
	//birthdayRepository := irepository.NewBirthdayRepository()
	//
	//notificationService := iservice.NewNotificationService(slackGateway)
	//rankingService := iservice.NewRankingService()
	//
	//notificationUseCase := usecase.NewNotificationUseCase(taskGateway, slackGateway, birthdayRepository, notificationService, rankingService)
	//notificationHandler := myCLI.NewNotificationHandler(notificationUseCase)
	//
	//// コマンドを設定
	//cliApp.Commands = []cli.Command{
	//	{
	//		Name:    "notify-today-tasks",
	//		Aliases: []string{"ntt"},
	//		Usage:   "Notify the today's tasks to Slack",
	//		Action:  notificationHandler.NotifyTodayTasksToSlack,
	//	},
	//	{
	//		Name:    "notify-today-birthday",
	//		Aliases: []string{"ntb"},
	//		Usage:   "Notify the today's birthday to Slack",
	//		Action:  notificationHandler.NotifyTodayBirthdayToSlack,
	//	},
	//	{
	//		Name:    "notify-access-ranking",
	//		Aliases: []string{"nar"},
	//		Usage:   "Notify the access ranking to Slack",
	//		Action:  notificationHandler.NotifyAccessRankingToSlack,
	//	},
	//}
	//
	//logger.Print("[CLI-ExecuteLog] $ hobi " + os.Args[1])
	//
	//if err := cliApp.Run(os.Args); err != nil {
	//	logger.Panic(errors.Wrap(err, "cliApp.Run()内でのエラー"))
	//}
}
