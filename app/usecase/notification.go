package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/analysis"
)

// ErrBirthdayNotFound : 該当Birthdayが存在しないエラー
var ErrBirthdayNotFound = errors.New("birthday is not found")

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
// FIXME: Trello -> Notion への移行を突貫工事で作ったのでリファクタ推奨
func (n notification) NotifyTodayTasksToSlack(ctx context.Context) (int, error) {
	cautionTasks, err := n.tg.FetchCautionTasks(ctx)
	if err != nil {
		return 0, fmt.Errorf("taskGateway.FetchCautionTasks()内でのエラー: %w", err)
	}

	deadTasks, err := n.tg.FetchDeadTasks(ctx)
	if err != nil {
		return 0, fmt.Errorf("taskGateway.FetchDeadTasks()内でのエラー: %w", err)
	}

	if err := n.sg.SendTask(ctx, cautionTasks, deadTasks); err != nil {
		return 0, fmt.Errorf("slackGateway.SendTask()内でのエラー: %w", err)
	}

	notifiedNum := len(cautionTasks) + len(deadTasks)
	return notifiedNum, nil
}

// NotifyTodayBirthdayToSlack : 今日誕生日の人をSlackに通知
func (n notification) NotifyTodayBirthdayToSlack(ctx context.Context) (int, error) {
	// 今日の誕生日情報を取得
	today := time.Now().Format("0102")
	birthdayList, err := n.r.FindAllByDate(ctx, today)
	if err != nil {
		if errors.Is(err, repository.ErrBirthdayRecordNotFound) {
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
