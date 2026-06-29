package task

import "time"

// List : タスクリストを表すドメインモデル
type List []Task

// FilterByStatus : ステータスが一致するタスクを抽出
func (l List) FilterByStatus(s Status) List {
	var result List
	for _, t := range l {
		if t.Status == s {
			result = append(result, t)
		}
	}
	return result
}

// GetDeadlineApproachingTasks : 期限が近づいているタスクを抽出（今日 ≤ Due ≤ 今日+7日）
func (l List) GetDeadlineApproachingTasks(now time.Time) List {
	var result List
	for _, t := range l {
		if t.IsDeadlineApproaching(now) {
			result = append(result, t)
		}
	}
	return result
}

// GetDueOverTasks : 期限切れのタスクを抽出
func (l List) GetDueOverTasks(now time.Time) List {
	var result List
	for _, t := range l {
		if t.IsDueOver(now) {
			result = append(result, t)
		}
	}
	return result
}

// GetTodayTasks : タスクリストから今日のタスクを取得
func (l List) GetTodayTasks(now time.Time) (todayTasks []Task) {
	for _, task := range l {
		if task.IsTodayTask(now) {
			todayTasks = append(todayTasks, task)
		}
	}
	return todayTasks
}
