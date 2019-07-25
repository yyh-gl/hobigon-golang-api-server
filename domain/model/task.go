package model

import (
	"github.com/adlio/trello"
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
	Due         time.Time  `json:"due"`
}

func (t Task) getJSTDue(utcDue time.Time) (jstDue time.Time) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstDue = utcDue.In(jst)
	return jstDue
}

func ConvertToTasksModel(trelloCards []*trello.Card) (tasks []*Task) {
	for _, card := range trelloCards {
		task := new(Task)
		task.Title       = card.Name
		task.Description = card.Desc
		if card.Due != nil {
			task.Due = task.getJSTDue(*card.Due)
		}
		tasks = append(tasks, task)
	}
	return tasks
}
