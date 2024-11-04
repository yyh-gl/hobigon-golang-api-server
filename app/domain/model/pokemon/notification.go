package pokemon

import (
	"time"
)

type Notification struct {
	category string
	title    string
	date     string
}

func NewNotification(category, title, date string) Notification {
	return Notification{
		category: category,
		title:    title,
		date:     date,
	}
}

func (n Notification) Title() string {
	return n.title
}

func (n Notification) IsEventCategory() bool {
	return n.category == "イベント"
}

func (n Notification) IsReceivedInToday() bool {
	today := time.Now().Format("2006.1.2")
	return n.date == today
}

func (n Notification) IsReceivedInYesterday() bool {
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006.1.2")
	return n.date == yesterday
}
