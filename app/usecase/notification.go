package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/pokemon"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

// Notification : Notification用ユースケースのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(ctx context.Context) (int, error)
	NotifyPokemonEvent(ctx context.Context) (int, error)
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
func (n notification) NotifyTodayTasksToSlack(ctx context.Context) (int, error) {
	now := time.Now()

	activeTasks, err := n.tg.FetchActiveTasks(ctx)
	if err != nil {
		return 0, fmt.Errorf("taskGateway.FetchActiveTasks()内でのエラー: %w", err)
	}

	toDoTasks := activeTasks.FilterByStatus(task.StatusToDo)

	var updatedTasks task.List
	for _, t := range toDoTasks.GetDeadlineApproachingTasks(now) {
		if err := n.tg.UpdateTaskStatus(ctx, t, task.StatusDoing); err != nil {
			log.Error(ctx, err)
			continue
		}
		updatedTasks = append(updatedTasks, t)
	}

	keyTasks := append(activeTasks.FilterByStatus(task.StatusDoing), updatedTasks...)
	deadTasks := toDoTasks.GetDueOverTasks(now)

	if err := n.sg.SendTasks(ctx, keyTasks, deadTasks); err != nil {
		return 0, fmt.Errorf("slackGateway.SendTasks()内でのエラー: %w", err)
	}

	return len(keyTasks) + len(deadTasks), nil
}

// NotifyPokemonEvent : Notify event notifications about Pokemon card to Slack.
func (n notification) NotifyPokemonEvent(ctx context.Context) (int, error) {
	notifications, err := crawlNotifications()
	if err != nil {
		return 0, err
	}

	events := extractNewEventNotifications(notifications)
	if err := n.sg.SendPokemonEvents(ctx, events); err != nil {
		return 0, err
	}

	return len(events), nil
}

func crawlNotifications() ([]pokemon.Notification, error) {
	c := colly.NewCollector()

	var events []pokemon.Notification
	c.OnHTML("li.List_item a div.List_body", func(e *colly.HTMLElement) {
		strs := strings.Split(e.Text, "\n")
		events = append(events, pokemon.NewNotification(
			strings.TrimSpace(strs[1]),
			strings.TrimSpace(strs[2]),
			strings.TrimSpace(strs[3]),
		))
	})

	if err := c.Visit("https://www.pokemon-card.com/info"); err != nil {
		return nil, err
	}

	return events, nil
}

func extractNewEventNotifications(notifications []pokemon.Notification) []pokemon.Notification {
	now := time.Now()
	existenceMap := make(map[string]struct{})
	events := make([]pokemon.Notification, 0, len(notifications))
	for _, n := range notifications {
		if !n.IsEventCategory() {
			continue
		}

		if _, ok := existenceMap[n.Title()]; ok {
			continue
		}

		if n.IsReceivedInToday(now) || n.IsReceivedInYesterday(now) {
			events = append(events, n)
			existenceMap[n.Title()] = struct{}{}
		}
	}
	return events
}
