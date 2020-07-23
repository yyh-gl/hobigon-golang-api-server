package rest

type errorResponse struct {
	Error errResp `json:"error"`
}

type errResp struct {
	Detail string `json:"detail"`
}

var (
	errBadRequest = errorResponse{
		Error: errResp{
			Detail: "不正なリクエスト形式です",
		},
	}
	errInterServerError = errorResponse{
		Error: errResp{
			Detail: "サーバ内でエラーが発生しました",
		},
	}
)
