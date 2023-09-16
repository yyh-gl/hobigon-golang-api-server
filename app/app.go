package app

import "os"

// version : アプリケーションのバージョン情報（GitHubのReleasesに対応）
// build時に埋め込む
var version string

// contextKey : コンテキストにセットするさいのキー用の型
type contextKey int

const (
	// ContextKeyCLI : cli.Contextをcontext.Contextにセットするさいのキー
	ContextKeyCLI contextKey = iota

	// ContextKeyTraceId is context key for trace id
	ContextKeyTraceId
)

// IsDev : 実行環境がDevelopmentかどうかを確認
func IsDev() bool {
	return os.Getenv("APP_ENV") == "develop"
}

// IsTest : 実行環境がTestかどうかを確認
func IsTest() bool {
	return os.Getenv("APP_ENV") == "test"
}

// IsPrd : 実行環境がProductionかどうかを確認
func IsPrd() bool {
	return os.Getenv("APP_ENV") == "production"
}

func GetVersion() string {
	return version
}
