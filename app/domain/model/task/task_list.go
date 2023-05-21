package task

// List : タスクリストを表すドメインモデル
type List []Task

// GetTodayTasks : タスクリストから今日のタスクを取得
func (l List) GetTodayTasks() (todayTasks []Task) {
	for _, task := range l {
		if task.IsTodayTask() {
			todayTasks = append(todayTasks, task)
		}
	}
	return todayTasks
}
