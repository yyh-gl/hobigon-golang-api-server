package gateway

import (
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// SlackGateway : 通知用のゲートウェイインターフェース
type SlackGateway interface {
	SendTask(todayTasks []model.Task, dueOverTasks []model.Task) error
	SendBirthday(birthday model.Birthday) error
	SendLikeNotify(blog model.Blog) error
	SendRanking(rankin string) error
}
