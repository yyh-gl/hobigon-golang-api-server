package usecase

import (
	"context"
	"fmt"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/analysis"
)

// Notification : Notification用ユースケースのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(ctx context.Context) (int, error)
	NotifyAccessRanking(ctx context.Context) (int, error)
}

type notification struct {
	tg gateway.Task
	sg gateway.Slack
}

// NewNotification : Notification用ユースケースを取得
func NewNotification(
	tg gateway.Task,
	sg gateway.Slack,
) Notification {
	return &notification{
		tg: tg,
		sg: sg,
	}
}

// TODO: 通知内容のコンテンツ数を返すようにする（ex. タスク一覧通知の場合はタスクの数）

// NotifyTodayTasksToSlack : 今日のタスク一覧をSlackに通知
// FIXME: Trello -> Notion への移行を突貫工事で作ったのでリファクタ推奨
func (n notification) NotifyTodayTasksToSlack(ctx context.Context) (int, error) {
	cautionAndToDoTasks, err := n.tg.FetchCautionAndToDoTasks(ctx)
	if err != nil {
		return 0, fmt.Errorf("taskGateway.FetchCautionTasks()内でのエラー: %w", err)
	}

	deadTasks, err := n.tg.FetchDeadTasks(ctx)
	if err != nil {
		return 0, fmt.Errorf("taskGateway.FetchDeadTasks()内でのエラー: %w", err)
	}

	if err := n.sg.SendTask(ctx, cautionAndToDoTasks, deadTasks); err != nil {
		return 0, fmt.Errorf("slackGateway.SendTask()内でのエラー: %w", err)
	}

	notifiedNum := len(cautionAndToDoTasks) + len(deadTasks)
	return notifiedNum, nil
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
