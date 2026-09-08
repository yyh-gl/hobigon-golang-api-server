// =============================================================================
// テストリスト（Canon TDD Step 3）
// 対象: notification.NotifyTodayTasksToSlack — Dead Tasks抽出元をactiveTasksに変更
//
// 正常系・振り分けロジック（nowFuncが2024-06-14 22:00 UTC＝06-15 07:00JSTを返す）:
//   - Doing・Deadline昨日 → Dead Tasksに含まれ、Key Tasksに含まれず、UpdateTaskStatusは呼ばれない
//   - To Do・Deadline昨日 → Dead Tasksに含まれ、Key Tasksに含まれず、UpdateTaskStatus(Doing)が1回呼ばれる
//   - To Do・Deadline今日 → Key Tasksに含まれ、UpdateTaskStatus(Doing)が1回呼ばれる
//   - To Do・Deadline8日後 → どちらにも含まれず、UpdateTaskStatusは呼ばれない
//   - To Do・Deadline nil → どちらにも含まれず、UpdateTaskStatusは呼ばれない
//   - Doing・Deadline nil → Key Tasksに含まれる
//   - Doing・Deadline今日 → Key Tasksに含まれ、UpdateTaskStatusは呼ばれない
//   - 戻り値の件数がlen(Key Tasks)+len(Dead Tasks)と一致する
// 混在ケース（決定1の核心）:
//   - 混在リストで同一タスクIDがKey TasksとDead Tasksの両方に現れない
//   - UpdateTaskStatus(Doing)要求が発生するIDは「To Doかつ(期限切れ or 期限間近)」と一致する
// 境界値:
//   - activeTasksが空のとき、Key/Deadともに空でSendTasksが呼ばれ、戻り値は0
// 異常系:
//   - FetchActiveTasksがエラーを返すとき(0, err)を返しSendTasksは呼ばれない
//   - SendTasksがエラーを返すとき(0, err)を返す
//   - To Do期限切れのUpdateTaskStatus失敗時もDead Tasksに掲載され続ける
//   - To Do期限間近のUpdateTaskStatus失敗時もKey Tasksに掲載され続ける（Dead Tasks側と対称）
//
// =============================================================================

package usecase

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	blogmodel "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/pokemon"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

func TestMain(m *testing.M) {
	log.NewLogger()
	os.Exit(m.Run())
}

type taskUpdate struct {
	id     string
	status task.Status
}

type fakeTaskGateway struct {
	tasks         task.List
	fetchErr      error
	updateErrByID map[string]error
	updates       []taskUpdate
}

func (f *fakeTaskGateway) FetchActiveTasks(context.Context) (task.List, error) {
	if f.fetchErr != nil {
		return nil, f.fetchErr
	}
	return f.tasks, nil
}

func (f *fakeTaskGateway) UpdateTaskStatus(_ context.Context, t task.Task, status task.Status) error {
	f.updates = append(f.updates, taskUpdate{id: t.ID, status: status})
	if err, ok := f.updateErrByID[t.ID]; ok {
		return err
	}
	return nil
}

type fakeSlackGateway struct {
	sendErr   error
	sent      bool
	keyTasks  []task.Task
	deadTasks []task.Task
}

func (f *fakeSlackGateway) SendTasks(_ context.Context, todayTasks []task.Task, dueOverTasks []task.Task) error {
	f.sent = true
	f.keyTasks = todayTasks
	f.deadTasks = dueOverTasks
	return f.sendErr
}

func (f *fakeSlackGateway) SendLikeNotification(context.Context, blogmodel.Blog) error {
	return nil
}

func (f *fakeSlackGateway) SendPokemonEvents(context.Context, []pokemon.Notification) error {
	return nil
}

func setNowFunc(t *testing.T, now time.Time) {
	t.Helper()
	nowFunc = func() time.Time { return now }
	t.Cleanup(func() { nowFunc = time.Now })
}

func containsTaskID(tasks []task.Task, id string) bool {
	for _, t := range tasks {
		if t.ID == id {
			return true
		}
	}
	return false
}

func countUpdatesForID(updates []taskUpdate, id string) int {
	count := 0
	for _, u := range updates {
		if u.id == id {
			count++
		}
	}
	return count
}

func TestNotifyTodayTasksToSlack_DoingOverdueTask(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	yesterday := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "doing-overdue", Status: task.StatusDoing, Deadline: &yesterday},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if !containsTaskID(sg.deadTasks, "doing-overdue") {
		t.Error("doing-overdue should be in Dead Tasks")
	}
	if containsTaskID(sg.keyTasks, "doing-overdue") {
		t.Error("doing-overdue should not be in Key Tasks")
	}
	if countUpdatesForID(tg.updates, "doing-overdue") != 0 {
		t.Errorf("UpdateTaskStatus should not be called for doing-overdue, got %d calls", countUpdatesForID(tg.updates, "doing-overdue"))
	}
}

func TestNotifyTodayTasksToSlack_ToDoOverdueTask(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	yesterday := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "todo-overdue", Status: task.StatusToDo, Deadline: &yesterday},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if !containsTaskID(sg.deadTasks, "todo-overdue") {
		t.Error("todo-overdue should be in Dead Tasks")
	}
	if containsTaskID(sg.keyTasks, "todo-overdue") {
		t.Error("todo-overdue should not be in Key Tasks")
	}
	if got := countUpdatesForID(tg.updates, "todo-overdue"); got != 1 {
		t.Errorf("UpdateTaskStatus calls for todo-overdue = %d, want 1", got)
	}
}

func TestNotifyTodayTasksToSlack_ToDoTodayTask(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	today := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "todo-today", Status: task.StatusToDo, Deadline: &today},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if !containsTaskID(sg.keyTasks, "todo-today") {
		t.Error("todo-today should be in Key Tasks")
	}
	if got := countUpdatesForID(tg.updates, "todo-today"); got != 1 {
		t.Errorf("UpdateTaskStatus calls for todo-today = %d, want 1", got)
	}
}

func TestNotifyTodayTasksToSlack_ToDoFarFutureTask(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	in8Days := time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "todo-far", Status: task.StatusToDo, Deadline: &in8Days},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if containsTaskID(sg.keyTasks, "todo-far") || containsTaskID(sg.deadTasks, "todo-far") {
		t.Error("todo-far should not be in Key Tasks nor Dead Tasks")
	}
	if got := countUpdatesForID(tg.updates, "todo-far"); got != 0 {
		t.Errorf("UpdateTaskStatus calls for todo-far = %d, want 0", got)
	}
}

func TestNotifyTodayTasksToSlack_ToDoNoDeadlineTask(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "todo-nil", Status: task.StatusToDo, Deadline: nil},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if containsTaskID(sg.keyTasks, "todo-nil") || containsTaskID(sg.deadTasks, "todo-nil") {
		t.Error("todo-nil should not be in Key Tasks nor Dead Tasks")
	}
	if got := countUpdatesForID(tg.updates, "todo-nil"); got != 0 {
		t.Errorf("UpdateTaskStatus calls for todo-nil = %d, want 0", got)
	}
}

func TestNotifyTodayTasksToSlack_DoingNoDeadlineTask(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "doing-nil", Status: task.StatusDoing, Deadline: nil},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if !containsTaskID(sg.keyTasks, "doing-nil") {
		t.Error("doing-nil should be in Key Tasks")
	}
}

func TestNotifyTodayTasksToSlack_DoingTodayTask(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	today := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "doing-today", Status: task.StatusDoing, Deadline: &today},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if !containsTaskID(sg.keyTasks, "doing-today") {
		t.Error("doing-today should be in Key Tasks")
	}
	if got := countUpdatesForID(tg.updates, "doing-today"); got != 0 {
		t.Errorf("UpdateTaskStatus calls for doing-today = %d, want 0", got)
	}
}

func TestNotifyTodayTasksToSlack_ReturnValueMatchesTotalCount(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	yesterday := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)
	today := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "doing-overdue", Status: task.StatusDoing, Deadline: &yesterday},
		{ID: "todo-today", Status: task.StatusToDo, Deadline: &today},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	got, err := n.NotifyTodayTasksToSlack(context.Background())
	if err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}
	want := len(sg.keyTasks) + len(sg.deadTasks)
	if got != want {
		t.Errorf("NotifyTodayTasksToSlack() = %d, want %d", got, want)
	}
}

func TestNotifyTodayTasksToSlack_MixedList(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	yesterday := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)
	today := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{tasks: task.List{
		{ID: "doing-overdue", Status: task.StatusDoing, Deadline: &yesterday},
		{ID: "todo-overdue", Status: task.StatusToDo, Deadline: &yesterday},
		{ID: "todo-today", Status: task.StatusToDo, Deadline: &today},
		{ID: "doing-today", Status: task.StatusDoing, Deadline: &today},
		{ID: "todo-nil", Status: task.StatusToDo, Deadline: nil},
	}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	for _, id := range []string{"doing-overdue", "todo-overdue", "todo-today", "doing-today", "todo-nil"} {
		if containsTaskID(sg.keyTasks, id) && containsTaskID(sg.deadTasks, id) {
			t.Errorf("%q appears in both Key Tasks and Dead Tasks", id)
		}
	}

	wantUpdated := map[string]bool{"todo-overdue": true, "todo-today": true}
	for _, u := range tg.updates {
		if !wantUpdated[u.id] {
			t.Errorf("unexpected UpdateTaskStatus call for %q", u.id)
		}
	}
	for id := range wantUpdated {
		if countUpdatesForID(tg.updates, id) != 1 {
			t.Errorf("UpdateTaskStatus calls for %q = %d, want 1", id, countUpdatesForID(tg.updates, id))
		}
	}
}

func TestNotifyTodayTasksToSlack_EmptyActiveTasks(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))

	tg := &fakeTaskGateway{tasks: task.List{}}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	got, err := n.NotifyTodayTasksToSlack(context.Background())
	if err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}
	if !sg.sent {
		t.Error("SendTasks should be called even with empty active tasks")
	}
	if len(sg.keyTasks) != 0 || len(sg.deadTasks) != 0 {
		t.Errorf("keyTasks/deadTasks should be empty, got %v / %v", sg.keyTasks, sg.deadTasks)
	}
	if got != 0 {
		t.Errorf("NotifyTodayTasksToSlack() = %d, want 0", got)
	}
}

func TestNotifyTodayTasksToSlack_FetchActiveTasksError(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	fetchErr := errors.New("fetch failed")

	tg := &fakeTaskGateway{fetchErr: fetchErr}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	got, err := n.NotifyTodayTasksToSlack(context.Background())
	if err == nil {
		t.Fatal("NotifyTodayTasksToSlack() error = nil, want error")
	}
	if got != 0 {
		t.Errorf("NotifyTodayTasksToSlack() = %d, want 0", got)
	}
	if sg.sent {
		t.Error("SendTasks should not be called when FetchActiveTasks fails")
	}
}

func TestNotifyTodayTasksToSlack_SendTasksError(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))

	tg := &fakeTaskGateway{tasks: task.List{}}
	sg := &fakeSlackGateway{sendErr: errors.New("send failed")}
	n := NewNotification(tg, sg, nil, nil)

	got, err := n.NotifyTodayTasksToSlack(context.Background())
	if err == nil {
		t.Fatal("NotifyTodayTasksToSlack() error = nil, want error")
	}
	if got != 0 {
		t.Errorf("NotifyTodayTasksToSlack() = %d, want 0", got)
	}
}

func TestNotifyTodayTasksToSlack_ToDoOverdueUpdateFails(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	yesterday := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{
		tasks: task.List{
			{ID: "todo-overdue", Status: task.StatusToDo, Deadline: &yesterday},
		},
		updateErrByID: map[string]error{"todo-overdue": errors.New("update failed")},
	}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if !containsTaskID(sg.deadTasks, "todo-overdue") {
		t.Error("todo-overdue should still be listed in Dead Tasks even when UpdateTaskStatus fails")
	}
}

func TestNotifyTodayTasksToSlack_ToDoApproachingUpdateFails(t *testing.T) {
	setNowFunc(t, time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC))
	today := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	tg := &fakeTaskGateway{
		tasks: task.List{
			{ID: "todo-today", Status: task.StatusToDo, Deadline: &today},
		},
		updateErrByID: map[string]error{"todo-today": errors.New("update failed")},
	}
	sg := &fakeSlackGateway{}
	n := NewNotification(tg, sg, nil, nil)

	if _, err := n.NotifyTodayTasksToSlack(context.Background()); err != nil {
		t.Fatalf("NotifyTodayTasksToSlack() error = %v", err)
	}

	if !containsTaskID(sg.keyTasks, "todo-today") {
		t.Error("todo-today should still be listed in Key Tasks even when UpdateTaskStatus fails (symmetric with Dead Tasks)")
	}
}
