package task

// List : タスクリストを表すドメインモデル
type List struct {
	Tasks []Task
}

// GetTodayTasks : タスクリストから今日のタスクを取得
func (tl List) GetTodayTasks() (todayTasks []Task) {
	for _, task := range tl.Tasks {
		if task.IsTodayTask() {
			todayTasks = append(todayTasks, task)
		}
	}
	return todayTasks
}
