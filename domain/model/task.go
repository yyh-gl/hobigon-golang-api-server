package model

import (
	"time"
)

type Board struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Desc           string `json:"desc"`
	Closed         bool   `json:"closed"`
	IDOrganization string `json:"idOrganization"`
	Pinned         bool   `json:"pinned"`
	URL            string `json:"url"`
	ShortURL       string `json:"shortUrl"`
}

type List struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Task struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Due         *time.Time `json:"due"`
}

type TaskList struct {
	Tasks []Task
}

func (t Task) GetJSTDue(utcDue *time.Time) *time.Time {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstDue := utcDue.In(jst)
	return &jstDue
}

func (tl TaskList) GetTodayTasks() (todayTasks []Task) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, jst)
	tomorrow := time.Date(today.Year(), today.Month(), today.AddDate(0, 0, 1).Day(), 23, 59, 59, 0, jst)

	for _, task := range tl.Tasks {
		// task.Due が NULL でなく、今日 または 今日以前 のタスクを抽出
		if task.Due != nil && (task.Due.Before(today) || (task.Due.After(today) && task.Due.Before(tomorrow))) {
			todayTasks = append(todayTasks, task)
		}
	}
	return todayTasks
}
