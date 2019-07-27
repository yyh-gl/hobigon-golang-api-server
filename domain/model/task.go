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
	Due         time.Time  `json:"due"`
}

func (t Task) GetJSTDue(utcDue time.Time) (jstDue time.Time) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstDue = utcDue.In(jst)
	return jstDue
}
