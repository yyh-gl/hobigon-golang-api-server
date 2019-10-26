package usecase

import (
	"context"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/service"
)

//////////////////////////////////////////////////
// NotifyTodayTasksToSlackUseCase
//////////////////////////////////////////////////

// NotifyTodayTasksToSlackUseCase は今日のタスク一覧を Slack に通知
func NotifyTodayTasksToSlackUseCase(ctx context.Context) error {
	taskGateway := gateway.NewTaskGateway()
	slackGateway := gateway.NewSlackGateway()

	var todayTasks []model.Task
	var dueOverTasks []model.Task

	// TODO: ビジネスロジックを結構持ってしまっているのでドメインモデルに落とし込んでいく
	boardIDList := [3]string{os.Getenv("MAIN_BOARD_ID"), os.Getenv("TECH_BOARD_ID"), os.Getenv("WORK_BOARD_ID")}
	for _, boardID := range boardIDList {
		lists, err := taskGateway.GetListsByBoardID(boardID)
		if err != nil {
			return errors.Wrap(err, "taskGateway.GetListsByBoardID()内でのエラー")
		}

		for _, list := range lists {
			// TODO: 今後必要があれば動的に変更できる仕組みを追加
			if list.Name == "TODO" || list.Name == "WIP" {
				taskList, dueOverTaskList, err := taskGateway.GetTasksFromList(*list)
				if err != nil {
					return errors.Wrap(err, "taskGateway.GetTasksFromList()内でのエラー")
				}

				switch list.Name {
				case "TODO":
					// TODOリストからは今日のタスクのみ出力
					tasks := taskList.GetTodayTasks()
					todayTasks = append(todayTasks, tasks...)
				case "WIP":
					// WIPリストにあるタスクは全て出力
					todayTasks = append(todayTasks, taskList.Tasks...)
				}

				// 期限切れタスクは問答無用で通知
				dueOverTasks = append(dueOverTasks, dueOverTaskList.Tasks...)
			}
		}
	}

	// 今日のタスクを WIP リストに移動
	if err := taskGateway.MoveToWIP(todayTasks); err != nil {
		return errors.Wrap(err, "taskGateway.MoveToWIP(todayTasks)内でのエラー")
	}

	// 期限切れのタスクを WIP リストに移動
	if err := taskGateway.MoveToWIP(dueOverTasks); err != nil {
		return errors.Wrap(err, "taskGateway.MoveToWIP(dueOverTasks)内でのエラー")
	}

	// 今日および期限切れのタスクを Slack に通知
	if err := slackGateway.SendTask(todayTasks, dueOverTasks); err != nil {
		return errors.Wrap(err, "slackGateway.SendTask()内でのエラー")
	}

	return nil
}

//////////////////////////////////////////////////
// NotifyTodayBirthdayToSlackUseCase
//////////////////////////////////////////////////

// NotifyTodayBirthdayToSlackUseCase は今日誕生日の人を Slack に通知
func NotifyTodayBirthdayToSlackUseCase(ctx context.Context) error {
	birthdayRepository := repository.NewBirthdayRepository()
	slackGateway := gateway.NewSlackGateway()
	notificationService := service.NewNotificationService(slackGateway, nil)

	// 今日の誕生日情報を取得
	today := time.Now().Format("0102")
	birthday, err := birthdayRepository.SelectByDate(today)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "birthdayRepository.SelectByDate()内でのエラー")
	}

	// 誕生日情報を Slack に通知
	err = notificationService.SendTodayBirthdayToSlack(birthday)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "notificationService.SendTodayBirthdayToSlack()内でのエラー")
	}

	return nil
}

//////////////////////////////////////////////////
// NotifyAccessRankingUseCase
//////////////////////////////////////////////////

// NotifyAccessRankingUseCase はアクセスランキングを Slack に通知
func NotifyAccessRankingUseCase(ctx context.Context) error {
	slackGateway := gateway.NewSlackGateway()

	// アクセスランキングの結果を取得
	// TODO: エクセルに出力して解析とかしたい
	// TODO: アウトプット再検討
	rankingMsg, _, err := service.GetAccessRanking()
	if err != nil {
		return errors.Wrap(err, "infra.GetAccessRanking()内でのエラー")
	}

	// アクセスランキングの結果を Slack に通知
	err = slackGateway.SendRanking(rankingMsg)
	if err != nil {
		return errors.Wrap(err, "slackGateway.SendRanking()内でのエラー")
	}

	return nil
}
