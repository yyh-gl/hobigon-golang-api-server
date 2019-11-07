package slack

import (
	"os"
	"strconv"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/task"
)

// TODO: ドメイン貧血症を治す
type Slack struct {
	Username string
	Channel  string
	Text     string
}

func (s Slack) GetWebHookURL() (webHookURL string) {
	switch s.Channel {
	case "00_today_tasks":
		webHookURL = os.Getenv("WEBHOOK_URL_TO_00")
	case "51_tech_blog":
		webHookURL = os.Getenv("WEBHOOK_URL_TO_51")
	case "2019新卒技術_雑談":
		webHookURL = os.Getenv("WEBHOOK_URL_TO_SHINSOTSU")
	}
	return webHookURL
}

func (s Slack) CreateTaskMessage(todayTasks []task.Task, dueOverTasks []task.Task) string {
	var mainTodayTasks []task.Task
	var techTodayTasks []task.Task
	var workTodayTasks []task.Task
	for _, t := range todayTasks {
		switch t.Board {
		case "Main":
			mainTodayTasks = append(mainTodayTasks, t)
		case "Tech":
			techTodayTasks = append(techTodayTasks, t)
		case "Work":
			workTodayTasks = append(workTodayTasks, t)
		}
	}

	var mainDueOverTasks []task.Task
	var techDueOverTasks []task.Task
	var workDueOverTasks []task.Task
	for _, t := range dueOverTasks {
		switch t.Board {
		case "Main":
			mainDueOverTasks = append(mainDueOverTasks, t)
		case "Tech":
			techDueOverTasks = append(techDueOverTasks, t)
		case "Work":
			workDueOverTasks = append(workDueOverTasks, t)
		}
	}

	header := ":pencil2::pencil2::pencil2: 今日のタスク :pencil2::pencil2::pencil2:\n\n"
	header += "           :mega: 総タスク数 " + strconv.Itoa(len(todayTasks)+len(dueOverTasks)) + " 個 :mega:\n"
	header += "   :space_invader: → 期限切れタスク数" + strconv.Itoa(len(dueOverTasks)) + "個 :space_invader:\n\n"

	body := ""

	// Mainボードのタスクをメッセージ化
	if len(mainTodayTasks) > 0 || len(mainDueOverTasks) > 0 {
		body += "\n           :mario2::dash: :kata-me::kata-i::kata-nn: :mario2::dash:\n"
	}
	index := 1
	for _, t := range mainTodayTasks {
		if t.IsTodayTask() {
			body += "\n:cubimal_chick:【M" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		} else {
			body += "\n:small_orange_diamond:【M" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		}
		body += "    :curly_loop: _" + t.ShortURL + " _\n"

		if t.Due == nil {
			body += "    :alarm_clock: `なるはや`\n\n"
		} else {
			body += "    :alarm_clock: `" + t.Due.Format("2006年1月2日 15時04分まで") + "`\n\n"
		}

		index++
	}

	for _, t := range mainDueOverTasks {
		body += "\n:space_invader:【M" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		body += "    :curly_loop: _" + t.ShortURL + " _\n"
		body += "    :alarm_clock: `" + t.Due.Format("2006年1月2日 15時04分まで") + "` :face_palm:\n\n"

		index++
	}

	// Techボードのタスクをメッセージ化
	if len(techTodayTasks) > 0 || len(techDueOverTasks) > 0 {
		body += "\n          :mario2::dash: :kanji-waza::kanji-jutsu::kanji-mukau::kanji-ue: :mario2::dash:\n\n"
	}
	index = 1
	for _, t := range techTodayTasks {
		if t.IsTodayTask() {
			body += "\n:cubimal_chick:【T" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		} else {
			body += "\n:small_orange_diamond:【T" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		}
		body += "    :curly_loop: _" + t.ShortURL + " _\n"

		if t.Due == nil {
			body += "    :alarm_clock: `なるはや`\n\n"
		} else {
			body += "    :alarm_clock: `" + t.Due.Format("2006年1月2日 15時04分まで") + "`\n\n"
		}

		index++
	}

	for _, t := range techDueOverTasks {
		body += "\n:space_invader:【T" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		body += "    :curly_loop: _" + t.ShortURL + " _\n"
		body += "    :alarm_clock: `" + t.Due.Format("2006年1月2日 15時04分まで") + "` :face_palm:\n\n"

		index++
	}

	// Workボードのタスクをメッセージ化
	if len(workTodayTasks) > 0 || len(workDueOverTasks) > 0 {
		body += "\n          :mario2::dash: :en-w::en-o::en-r::en-k: :mario2::dash:\n\n"
	}
	index = 1
	for _, t := range workTodayTasks {
		if t.IsTodayTask() {
			body += "\n:cubimal_chick:【W" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		} else {
			body += "\n:small_orange_diamond:【W" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		}
		body += "    :curly_loop: _" + t.ShortURL + " _\n"

		if t.Due == nil {
			body += "    :alarm_clock: `なるはや`\n\n"
		} else {
			body += "    :alarm_clock: `" + t.Due.Format("2006年1月2日 15時04分まで") + "`\n\n"
		}

		index++
	}

	for _, t := range workDueOverTasks {
		body += "\n:space_invader:【W" + strconv.Itoa(index) + "】 *_" + t.Title + "_*\n"
		body += "    :curly_loop: _" + t.ShortURL + " _\n"
		body += "    :alarm_clock: `" + t.Due.Format("2006年1月2日 15時04分まで") + "` :face_palm:\n\n"

		index++
	}

	return header + body
}
