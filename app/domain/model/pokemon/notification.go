package pokemon

import (
	"strings"
	"time"
)

type Notification struct {
	category string
	title    string
	date     string
}

// importantEventKeywords are the keywords to judge whether the event is important or not.
var importantEventKeywords = []string{
	"シティリーグ",
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

func (n Notification) IsImportantEvent() bool {
	for _, keyword := range importantEventKeywords {
		if strings.Contains(n.title, keyword) {
			return true
		}
	}
	return false
}
