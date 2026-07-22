package line

import (
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	modelLine "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
)

const lineMessagesConfigPath = "config/line_messages.yaml"

type lineMessageConfig struct{}

// NewLINEMessageConfig : LINEMessageConfig用ゲートウェイを取得
func NewLINEMessageConfig() gateway.LINEMessageConfig {
	return &lineMessageConfig{}
}

// GetMessage : 設定ファイルからメッセージキーに対応する文言を取得
func (c lineMessageConfig) GetMessage(_ context.Context, messageKey string) (string, error) {
	data, err := os.ReadFile(lineMessagesConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to os.ReadFile(%s): %w", lineMessagesConfigPath, err)
	}

	var messages modelLine.LINEMessages
	if err := yaml.Unmarshal(data, &messages); err != nil {
		return "", fmt.Errorf("failed to yaml.Unmarshal(): %w", err)
	}

	return messages.FindMessageFor(messageKey)
}
