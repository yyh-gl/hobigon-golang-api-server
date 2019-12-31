package rest

// Response : REST API 用の共通エラーレスポンス
type errorResponse struct {
	Detail string `json:"detail"`
}
