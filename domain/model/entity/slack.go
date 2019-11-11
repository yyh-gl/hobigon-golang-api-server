package entity

import (
	"os"
	"strconv"
)

// Slack : Slack 用のドメインモデル
// TODO: ドメイン貧血症を治す
type Slack struct {
	Username string
	Channel  string
	Text     string
}

// GetWebHookURL : WebHook URL を取得
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

// CreateTaskMessage : タスク通知用のメッセージを作成
func (s Slack) CreateTaskMessage(todayTasks []Task, dueOverTasks []Task) string {
	var mainTodayTasks []Task
	var techTodayTasks []Task
	var workTodayTasks []Task
	for _, task := range todayTasks {
		switch task.Board {
		case "Main":
			mainTodayTasks = append(mainTodayTasks, task)
		case "Tech":
			techTodayTasks = append(techTodayTasks, task)
		case "Work":
			workTodayTasks = append(workTodayTasks, task)
		}
	}

	var mainDueOverTasks []Task
	var techDueOverTasks []Task
	var workDueOverTasks []Task
	for _, task := range dueOverTasks {
		switch task.Board {
		case "Main":
			mainDueOverTasks = append(mainDueOverTasks, task)
		case "Tech":
			techDueOverTasks = append(techDueOverTasks, task)
		case "Work":
			workDueOverTasks = append(workDueOverTasks, task)
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
	for _, task := range mainTodayTasks {
		if task.IsTodayTask() {
			body += "\n:cubimal_chick:【M" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		} else {
			body += "\n:small_orange_diamond:【M" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		}
		body += "    :curly_loop: _" + task.ShortURL + " _\n"

		if task.Due == nil {
			body += "    :alarm_clock: `なるはや`\n\n"
		} else {
			body += "    :alarm_clock: `" + task.Due.Format("2006年1月2日 15時04分まで") + "`\n\n"
		}

		index++
	}

	for _, task := range mainDueOverTasks {
		body += "\n:space_invader:【M" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		body += "    :curly_loop: _" + task.ShortURL + " _\n"
		body += "    :alarm_clock: `" + task.Due.Format("2006年1月2日 15時04分まで") + "` :face_palm:\n\n"

		index++
	}

	// Techボードのタスクをメッセージ化
	if len(techTodayTasks) > 0 || len(techDueOverTasks) > 0 {
		body += "\n          :mario2::dash: :kanji-waza::kanji-jutsu::kanji-mukau::kanji-ue: :mario2::dash:\n\n"
	}
	index = 1
	for _, task := range techTodayTasks {
		if task.IsTodayTask() {
			body += "\n:cubimal_chick:【T" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		} else {
			body += "\n:small_orange_diamond:【T" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		}
		body += "    :curly_loop: _" + task.ShortURL + " _\n"

		if task.Due == nil {
			body += "    :alarm_clock: `なるはや`\n\n"
		} else {
			body += "    :alarm_clock: `" + task.Due.Format("2006年1月2日 15時04分まで") + "`\n\n"
		}

		index++
	}

	for _, task := range techDueOverTasks {
		body += "\n:space_invader:【T" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		body += "    :curly_loop: _" + task.ShortURL + " _\n"
		body += "    :alarm_clock: `" + task.Due.Format("2006年1月2日 15時04分まで") + "` :face_palm:\n\n"

		index++
	}

	// Workボードのタスクをメッセージ化
	if len(workTodayTasks) > 0 || len(workDueOverTasks) > 0 {
		body += "\n          :mario2::dash: :en-w::en-o::en-r::en-k: :mario2::dash:\n\n"
	}
	index = 1
	for _, task := range workTodayTasks {
		if task.IsTodayTask() {
			body += "\n:cubimal_chick:【W" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		} else {
			body += "\n:small_orange_diamond:【W" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		}
		body += "    :curly_loop: _" + task.ShortURL + " _\n"

		if task.Due == nil {
			body += "    :alarm_clock: `なるはや`\n\n"
		} else {
			body += "    :alarm_clock: `" + task.Due.Format("2006年1月2日 15時04分まで") + "`\n\n"
		}

		index++
	}

	for _, task := range workDueOverTasks {
		body += "\n:space_invader:【W" + strconv.Itoa(index) + "】 *_" + task.Title + "_*\n"
		body += "    :curly_loop: _" + task.ShortURL + " _\n"
		body += "    :alarm_clock: `" + task.Due.Format("2006年1月2日 15時04分まで") + "` :face_palm:\n\n"

		index++
	}

	return header + body
}
