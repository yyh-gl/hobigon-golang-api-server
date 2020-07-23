package rest

// errorResponse : API用の共通エラーレスポンス構造体
type errorResponse struct {
	Error errResp `json:"error"`
}

type errResp struct {
	Detail string `json:"detail"`
}

var (
	// errNotFound : 404 Not Found
	errNotFound = errorResponse{
		Error: errResp{
			Detail: "該当するリソースが存在しません",
		},
	}

	// errBadRequest : 404 Bad Request
	errBadRequest = errorResponse{
		Error: errResp{
			Detail: "不正なリクエスト形式です",
		},
	}

	// errInterServerError : 500 Internal Server Error
	errInterServerError = errorResponse{
		Error: errResp{
			Detail: "サーバ内でエラーが発生しました",
		},
	}
)
