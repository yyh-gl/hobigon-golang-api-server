package task

// Task : タスク用のドメインモデル
type Task struct {
	Title       string
	Description string
	Due         Date
	Board       string
	List        string
	ShortURL    string
}

func NewTask() {

}
