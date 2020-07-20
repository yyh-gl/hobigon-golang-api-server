package rest

type errorResponse2 struct {
	Error errResp `json:"error"`
}

type errResp struct {
	Detail string `json:"detail"`
}

var (
	errBadRequest = errorResponse2{
		Error: errResp{
			Detail: "不正なリクエスト形式です",
		},
	}
	errInterServerError = errorResponse2{
		Error: errResp{
			Detail: "サーバ内でエラーが発生しました",
		},
	}
)
