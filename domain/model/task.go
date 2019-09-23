package model

import (
	"time"
)

// TODO: ドメイン貧血症を治す
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
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Due           *time.Time  `json:"due"`
	Board         string      `json:"board"`
	List          string      `json:"list"`
	ShortURL      string      `json:"short_url"`
	OriginalModel interface{} `json:"-"`
}

type TaskList struct {
	Tasks []Task
}

func (t Task) GetJSTDue(utcDue *time.Time) *time.Time {
	jst := getJSTNow()
	jstDue := utcDue.In(jst)
	return &jstDue
}

func (t Task) IsDueOver() (isDueOver bool) {
	jst := getJSTNow()
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, jst)
	return !t.Due.Equal(todayStart) && t.Due.Before(todayStart)
}

func (t Task) IsTodayTask() (isTodayTask bool) {
	jst := getJSTNow()
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, jst)
	todayEnd := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, jst)
	if t.Due != nil && t.Due.After(todayStart) && t.Due.Before(todayEnd) {
		return true
	}
	return false
}

func (tl TaskList) GetTodayTasks() (todayTasks []Task) {
	for _, task := range tl.Tasks {
		if task.IsTodayTask() {
			todayTasks = append(todayTasks, task)
		}
	}
	return todayTasks
}

func getJSTNow() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
