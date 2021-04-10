package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/analysis"
)

// Notification : Notification用ユースケースのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(ctx context.Context) (int, error)
	NotifyTodayBirthdayToSlack(ctx context.Context) (int, error)
	NotifyAccessRanking(ctx context.Context) (int, error)
}

type notification struct {
	tg gateway.Task
	sg gateway.Slack
	r  repository.Birthday
}

// NewNotification : Notification用ユースケースを取得
func NewNotification(
	tg gateway.Task,
	sg gateway.Slack,
	r repository.Birthday,
) Notification {
	return &notification{
		tg: tg,
		sg: sg,
		r:  r,
	}
}

// TODO: 通知内容のコンテンツ数を返すようにする（ex. タスク一覧通知の場合はタスクの数）

// NotifyTodayTasksToSlack : 今日のタスク一覧をSlackに通知
func (n notification) NotifyTodayTasksToSlack(ctx context.Context) (int, error) {
	var todayTasks []model.Task
	var dueOverTasks []model.Task

	// TODO: ビジネスロジックを結構持ってしまっているのでドメインサービスに移す
	boardIDList := [3]string{os.Getenv("MAIN_BOARD_ID"), os.Getenv("TECH_BOARD_ID"), os.Getenv("WORK_BOARD_ID")}
	for _, boardID := range boardIDList {
		lists, err := n.tg.GetListsByBoardID(ctx, boardID)
		if err != nil {
			return 0, fmt.Errorf("taskGateway.GetListsByBoardID()内でのエラー: %w", err)
		}

		for _, list := range lists {
			// TODO: 今後必要があれば動的に変更できる仕組みを追加
			if list.Name == "TODO" || list.Name == "WIP" {
				taskList, dueOverTaskList, err := n.tg.GetTasksFromList(ctx, *list)
				if err != nil {
					return 0, fmt.Errorf("taskGateway.GetTasksFromList()内でのエラー: %w", err)
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

	// 今日のタスクをWIPリストに移動
	if err := n.tg.MoveToWIP(ctx, todayTasks); err != nil {
		return 0, fmt.Errorf("taskGateway.MoveToWIP(todayTasks)内でのエラー: %w", err)
	}

	// 期限切れのタスクをWIPリストに移動
	if err := n.tg.MoveToWIP(ctx, dueOverTasks); err != nil {
		return 0, fmt.Errorf("taskGateway.MoveToWIP(dueOverTasks)内でのエラー: %w", err)
	}

	// 今日および期限切れのタスクをSlackに通知
	if err := n.sg.SendTask(ctx, todayTasks, dueOverTasks); err != nil {
		return 0, fmt.Errorf("slackGateway.SendTask()内でのエラー: %w", err)
	}

	notifiedNum := len(todayTasks) + len(dueOverTasks)
	return notifiedNum, nil
}

// NotifyTodayBirthdayToSlack : 今日誕生日の人をSlackに通知
func (n notification) NotifyTodayBirthdayToSlack(ctx context.Context) (int, error) {
	// 今日の誕生日情報を取得
	today := time.Now().Format("0102")
	birthdayList, err := n.r.FindAllByDate(ctx, today)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return 0, ErrBirthdayNotFound
		}
		return 0, fmt.Errorf("birthdayRepository.SelectByDate()内でのエラー: %w", err)
	}

	// 誕生日情報を Slack に通知
	err = n.sg.SendBirthday(ctx, birthdayList)
	if err != nil {
		return 0, fmt.Errorf("notificationService.SendBirthday()内でのエラー: %w", err)
	}

	return len(birthdayList), nil
}

// NotifyAccessRanking : アクセスランキングをSlackに通知
func (n notification) NotifyAccessRanking(ctx context.Context) (int, error) {
	// アクセスランキングの結果を取得
	// NOTE: シンプルさのためにinfraを直参照
	rankingMsg, notifiedNum, err := analysis.GetAccessRanking(ctx)
	if err != nil {
		return 0, fmt.Errorf("infra.GetAccessRanking()内でのエラー: %w", err)
	}

	// アクセスランキングの結果を Slack に通知
	err = n.sg.SendRanking(ctx, rankingMsg)
	if err != nil {
		return 0, fmt.Errorf("slackGateway.SendRanking()内でのエラー: %w", err)
	}

	return notifiedNum, nil
}
