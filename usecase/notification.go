package usecase

import (
	"context"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/infra/service"

	"github.com/yyh-gl/hobigon-golang-api-server/infra"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

// NotifyTodayBirthdayToSlackUseCase は今日誕生日の人を Slack に通知
func NotifyTodayBirthdayToSlackUseCase(ctx context.Context) error {
	birthdayRepository := repository.NewBirthdayRepository()
	slackGateway := gateway.NewSlackGateway()
	notificationService := service.NewNotificationService(slackGateway)

	// 今日の誕生日情報を取得
	today := time.Now().Format("0102")
	birthday, err := birthdayRepository.SelectByDate(today)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "birthdayRepository.SelectByDate()内でのエラー")
	}

	// 誕生日情報を Slack に通知
	err = notificationService.SendBirthdayToSlack(birthday)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "notificationService.SendBirthdayToSlack()内でのエラー")
	}

	return nil
}

// NotifyAccessRankingUseCase はアクセスランキングを Slack に通知
func NotifyAccessRankingUseCase(ctx context.Context) error {
	slackGateway := gateway.NewSlackGateway()

	// アクセスランキングの結果を取得
	// TODO: エクセルに出力して解析とかしたい
	// TODO: アウトプット再検討
	rankingMsg, _, err := infra.GetAccessRanking()
	if err != nil {
		return errors.Wrap(err, "infra.GetAccessRanking()内でのエラー")
	}

	// アクセスランキングの結果を Slack に通知
	err = slackGateway.SendRanking(rankingMsg)
	if err != nil {
		return errors.Wrap(err, "slackGateway.SendRanking()内でのエラー")
	}

	return nil
}
