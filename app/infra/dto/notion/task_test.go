// =============================================================================
// テストリスト（Canon TDD Step 4）
// 対象: NotionTaskDTO.ToTaskListDomainModel(ctx) — 時刻付きDeadline対応、パース失敗時のフォールバック
//
// 正常系（パース）:
//   - Deadline.Date.Start = "2024-06-15"（日付のみ）のとき、Deadlineは2024-06-15 00:00 UTCと一致する
//   - Deadline.Date.Start = "2024-06-15T10:00:00.000+09:00"（小数秒付き）のとき、
//     Deadlineは2024-06-15 10:00 +09:00と同一時刻を指す
//   - Deadline.Date.Start = "2024-06-15T10:00:00+09:00"（小数秒なし）のときも同様
//   - Deadline.Date.Start = "2024-06-15T01:00:00.000Z"（Z表記・UTC）のときも同様
//   - Deadline.Date.Start = ""（未設定）のとき、Deadlineはnil
//   - タイトルが空のページはスキップされる（既存挙動の回帰確認）
//   - Resultsが空のとき、返るtask.Listの長さは0
//   - 正常系1件について、ID・Title・Status・ShortURLが入力の値を正しく転写している
//
// 異常系・フォールバック:
//   - Deadline.Date.Start = "invalid"のとき、タスクは残りDeadlineはnil
//   - Deadline.Date.Start = "2024/06/15"（不正な区切り文字）のとき、タスクは残りDeadlineはnil
//   - Deadline.Date.Start = "2024-06"（不完全な日付）のとき、タスクは残りDeadlineはnil
//   - パース失敗タスク1件と正常タスク1件が混在するとき、continueで正常タスクごと消えず2件とも残る
//
// =============================================================================

package notion_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto/notion"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

func TestMain(m *testing.M) {
	log.NewLogger()
	os.Exit(m.Run())
}

func resultJSON(id, title, deadlineStart, status, url string) string {
	titleBlocks := "[]"
	if title != "" {
		titleBlocks = fmt.Sprintf(`[{"plain_text": %q}]`, title)
	}
	return fmt.Sprintf(`{
		"id": %q,
		"properties": {
			"Deadline": {"date": {"start": %q}},
			"Status": {"select": {"name": %q}},
			"Name": {"title": %s}
		},
		"url": %q
	}`, id, deadlineStart, status, titleBlocks, url)
}

func buildDTO(t *testing.T, results ...string) notion.NotionTaskDTO {
	t.Helper()
	body := fmt.Sprintf(`{"results": [%s]}`, joinComma(results))
	var dto notion.NotionTaskDTO
	if err := json.Unmarshal([]byte(body), &dto); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return dto
}

func joinComma(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}

func TestToTaskListDomainModel_DateOnlyDeadline(t *testing.T) {
	dto := buildDTO(t, resultJSON("task-1", "Task 1", "2024-06-15", "To Do", "https://notion.so/task-1"))

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	want := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if got[0].Deadline == nil || !got[0].Deadline.Equal(want) {
		t.Errorf("Deadline = %v, want %v", got[0].Deadline, want)
	}
}

func TestToTaskListDomainModel_DateTimeDeadlineWithMillis(t *testing.T) {
	dto := buildDTO(t, resultJSON("task-1", "Task 1", "2024-06-15T10:00:00.000+09:00", "To Do", "https://notion.so/task-1"))

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	want := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)
	if got[0].Deadline == nil || !got[0].Deadline.Equal(want) {
		t.Errorf("Deadline = %v, want %v", got[0].Deadline, want)
	}
}

func TestToTaskListDomainModel_DateTimeDeadlineWithoutMillis(t *testing.T) {
	dto := buildDTO(t, resultJSON("task-1", "Task 1", "2024-06-15T10:00:00+09:00", "To Do", "https://notion.so/task-1"))

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	want := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)
	if got[0].Deadline == nil || !got[0].Deadline.Equal(want) {
		t.Errorf("Deadline = %v, want %v", got[0].Deadline, want)
	}
}

func TestToTaskListDomainModel_DateTimeDeadlineUTCZ(t *testing.T) {
	dto := buildDTO(t, resultJSON("task-1", "Task 1", "2024-06-15T01:00:00.000Z", "To Do", "https://notion.so/task-1"))

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	want := time.Date(2024, 6, 15, 1, 0, 0, 0, time.UTC)
	if got[0].Deadline == nil || !got[0].Deadline.Equal(want) {
		t.Errorf("Deadline = %v, want %v", got[0].Deadline, want)
	}
}

func TestToTaskListDomainModel_NoDeadline(t *testing.T) {
	dto := buildDTO(t, resultJSON("task-1", "Task 1", "", "To Do", "https://notion.so/task-1"))

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if got[0].Deadline != nil {
		t.Errorf("Deadline = %v, want nil", got[0].Deadline)
	}
}

func TestToTaskListDomainModel_SkipsEmptyTitle(t *testing.T) {
	dto := buildDTO(t, resultJSON("task-1", "", "2024-06-15", "To Do", "https://notion.so/task-1"))

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

func TestToTaskListDomainModel_EmptyResults(t *testing.T) {
	dto := buildDTO(t)

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

func TestToTaskListDomainModel_FieldTranscription(t *testing.T) {
	dto := buildDTO(t, resultJSON("task-1", "Task 1", "2024-06-15", "Doing", "https://notion.so/task-1"))

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	tsk := got[0]
	if tsk.ID != "task-1" {
		t.Errorf("ID = %q, want %q", tsk.ID, "task-1")
	}
	if tsk.Title != "Task 1" {
		t.Errorf("Title = %q, want %q", tsk.Title, "Task 1")
	}
	if tsk.Status.String() != "Doing" {
		t.Errorf("Status = %q, want %q", tsk.Status.String(), "Doing")
	}
	if tsk.ShortURL != "https://notion.so/task-1" {
		t.Errorf("ShortURL = %q, want %q", tsk.ShortURL, "https://notion.so/task-1")
	}
}

func TestToTaskListDomainModel_InvalidDeadlineIsKeptWithNilDeadline(t *testing.T) {
	tests := []struct {
		name  string
		start string
	}{
		{"完全に不正な文字列", "invalid"},
		{"不正な区切り文字", "2024/06/15"},
		{"不完全な日付", "2024-06"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := buildDTO(t, resultJSON("task-1", "Task 1", tt.start, "To Do", "https://notion.so/task-1"))

			got := dto.ToTaskListDomainModel(context.Background())

			if len(got) != 1 {
				t.Fatalf("len(got) = %d, want 1", len(got))
			}
			if got[0].Deadline != nil {
				t.Errorf("Deadline = %v, want nil", got[0].Deadline)
			}
		})
	}
}

func TestToTaskListDomainModel_InvalidTaskDoesNotDropValidTask(t *testing.T) {
	dto := buildDTO(t,
		resultJSON("task-invalid", "Invalid Task", "invalid", "To Do", "https://notion.so/task-invalid"),
		resultJSON("task-valid", "Valid Task", "2024-06-15", "To Do", "https://notion.so/task-valid"),
	)

	got := dto.ToTaskListDomainModel(context.Background())

	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
}
