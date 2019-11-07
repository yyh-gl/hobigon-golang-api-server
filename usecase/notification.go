package usecase

import (
	"context"
	"os"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/task"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/service"
)

//////////////////////////////////////////////////
// NewNotificationUseCase
//////////////////////////////////////////////////

// NotificationUseCase : 通知用のユースケースインターフェース
type NotificationUseCase interface {
	NotifyTodayTasksToSlack(ctx context.Context) error
	NotifyTodayBirthdayToSlack(ctx context.Context) error
	NotifyAccessRanking(ctx context.Context) error
}

type notificationUseCase struct {
	tg gateway.TaskGateway
	sg gateway.SlackGateway
	br repository.BirthdayRepository
	ns service.NotificationService
	rs service.RankingService
}

// NewNotificationUseCase : 通知用のユースケースを取得
func NewNotificationUseCase(
	tg gateway.TaskGateway,
	sg gateway.SlackGateway,
	br repository.BirthdayRepository,
	ns service.NotificationService,
	rs service.RankingService,
) NotificationUseCase {
	return &notificationUseCase{
		tg: tg,
		sg: sg,
		br: br,
		ns: ns,
		rs: rs,
	}
}

//////////////////////////////////////////////////
// NotifyTodayTasksToSlack
//////////////////////////////////////////////////

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (nu notificationUseCase) NotifyTodayTasksToSlack(ctx context.Context) error {
	var todayTasks []task.Task
	var dueOverTasks []task.Task

	// TODO: ビジネスロジックを結構持ってしまっているのでドメインモデルに落とし込んでいく
	boardIDList := [3]string{os.Getenv("MAIN_BOARD_ID"), os.Getenv("TECH_BOARD_ID"), os.Getenv("WORK_BOARD_ID")}
	for _, boardID := range boardIDList {
		lists, err := nu.tg.GetListsByBoardID(ctx, boardID)
		if err != nil {
			return errors.Wrap(err, "taskGateway.GetListsByBoardID()内でのエラー")
		}

		for _, list := range lists {
			// TODO: 今後必要があれば動的に変更できる仕組みを追加
			if list.Name == "TODO" || list.Name == "WIP" {
				taskList, dueOverTaskList, err := nu.tg.GetTasksFromList(ctx, *list)
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
	if err := nu.tg.MoveToWIP(ctx, todayTasks); err != nil {
		return errors.Wrap(err, "taskGateway.MoveToWIP(todayTasks)内でのエラー")
	}

	// 期限切れのタスクを WIP リストに移動
	if err := nu.tg.MoveToWIP(ctx, dueOverTasks); err != nil {
		return errors.Wrap(err, "taskGateway.MoveToWIP(dueOverTasks)内でのエラー")
	}

	// 今日および期限切れのタスクを Slack に通知
	if err := nu.sg.SendTask(ctx, todayTasks, dueOverTasks); err != nil {
		return errors.Wrap(err, "slackGateway.SendTask()内でのエラー")
	}

	return nil
}

//////////////////////////////////////////////////
// NotifyTodayBirthdayToSlack
//////////////////////////////////////////////////

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (nu notificationUseCase) NotifyTodayBirthdayToSlack(ctx context.Context) error {
	// 今日の誕生日情報を取得
	today := time.Now().Format("0102")
	birthday, err := nu.br.SelectByDate(ctx, today)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "birthdayRepository.SelectByDate()内でのエラー")
	}

	// 誕生日情報を Slack に通知
	if birthday != nil {
		err = nu.ns.SendTodayBirthdayToSlack(ctx, *birthday)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "notificationService.SendTodayBirthdayToSlack()内でのエラー")
		}
	}

	return nil
}

//////////////////////////////////////////////////
// NotifyAccessRanking
//////////////////////////////////////////////////

// NotifyAccessRanking : アクセスランキングを Slack に通知
func (nu notificationUseCase) NotifyAccessRanking(ctx context.Context) error {
	// アクセスランキングの結果を取得
	// TODO: エクセルに出力して解析とかしたい
	// TODO: アウトプット再検討
	rankingMsg, _, err := nu.rs.GetAccessRanking(ctx)
	if err != nil {
		return errors.Wrap(err, "infra.GetAccessRanking()内でのエラー")
	}

	// アクセスランキングの結果を Slack に通知
	err = nu.sg.SendRanking(ctx, rankingMsg)
	if err != nil {
		return errors.Wrap(err, "slackGateway.SendRanking()内でのエラー")
	}

	return nil
}
