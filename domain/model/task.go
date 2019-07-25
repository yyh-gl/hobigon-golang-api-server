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
	Due         *time.Time `json:"due"`
}

func ConvertToTasksModel(trelloCards []*trello.Card) (tasks []*Task) {
	for _, card := range trelloCards {
		task := Task{
			Title: card.Name,
			Description: card.Desc,
			Due: card.Due,
		}
		tasks = append(tasks, &task)
	}
	return tasks
}
