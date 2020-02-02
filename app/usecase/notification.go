package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/service"
)

// Notification : Notification用ユースケースのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(ctx context.Context) error
	NotifyTodayBirthdayToSlack(ctx context.Context) error
	NotifyAccessRanking(ctx context.Context) error
}

type notification struct {
	tg gateway.Task
	sg gateway.Slack
	r  repository.Birthday
	rs service.Ranking
}

// NewNotification : Notification用ユースケースを取得
func NewNotification(
	tg gateway.Task,
	sg gateway.Slack,
	r repository.Birthday,
	rs service.Ranking,
) Notification {
	return &notification{
		tg: tg,
		sg: sg,
		r:  r,
		rs: rs,
	}
}

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (n notification) NotifyTodayTasksToSlack(ctx context.Context) error {
	var todayTasks []model.Task
	var dueOverTasks []model.Task

	// TODO: ビジネスロジックを結構持ってしまっているのでドメインモデルに落とし込んでいく
	boardIDList := [3]string{os.Getenv("MAIN_BOARD_ID"), os.Getenv("TECH_BOARD_ID"), os.Getenv("WORK_BOARD_ID")}
	for _, boardID := range boardIDList {
		lists, err := n.tg.GetListsByBoardID(ctx, boardID)
		if err != nil {
			return fmt.Errorf("taskGateway.GetListsByBoardID()内でのエラー: %w", err)
		}

		for _, list := range lists {
			// TODO: 今後必要があれば動的に変更できる仕組みを追加
			if list.Name == "TODO" || list.Name == "WIP" {
				taskList, dueOverTaskList, err := n.tg.GetTasksFromList(ctx, *list)
				if err != nil {
					return fmt.Errorf("taskGateway.GetTasksFromList()内でのエラー: %w", err)
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
	if err := n.tg.MoveToWIP(ctx, todayTasks); err != nil {
		return fmt.Errorf("taskGateway.MoveToWIP(todayTasks)内でのエラー: %w", err)
	}

	// 期限切れのタスクを WIP リストに移動
	if err := n.tg.MoveToWIP(ctx, dueOverTasks); err != nil {
		return fmt.Errorf("taskGateway.MoveToWIP(dueOverTasks)内でのエラー: %w", err)
	}

	// 今日および期限切れのタスクを Slack に通知
	if err := n.sg.SendTask(ctx, todayTasks, dueOverTasks); err != nil {
		return fmt.Errorf("slackGateway.SendTask()内でのエラー: %w", err)
	}

	return nil
}

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (n notification) NotifyTodayBirthdayToSlack(ctx context.Context) error {
	// 今日の誕生日情報を取得
	today := time.Now().Format("0102")
	birthday, err := n.r.SelectByDate(ctx, today)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("birthdayRepository.SelectByDate()内でのエラー: %w", err)
	}
	if birthday == nil {
		return nil
	}

	// 誕生日情報を Slack に通知
	err = n.sg.SendBirthday(ctx, *birthday)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("notificationService.SendBirthday()内でのエラー: %w", err)
	}

	return nil
}

// NotifyAccessRanking : アクセスランキングを Slack に通知
func (n notification) NotifyAccessRanking(ctx context.Context) error {
	// アクセスランキングの結果を取得
	// TODO: アウトプット再検討->エクセルに出力して解析とかしたい
	rankingMsg, _, err := n.rs.GetAccessRanking(ctx)
	if err != nil {
		return fmt.Errorf("infra.GetAccessRanking()内でのエラー: %w", err)
	}

	// アクセスランキングの結果を Slack に通知
	err = n.sg.SendRanking(ctx, rankingMsg)
	if err != nil {
		return fmt.Errorf("slackGateway.SendRanking()内でのエラー: %w", err)
	}

	return nil
}
