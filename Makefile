VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

.PHONY: build mod
help: ## ヘルプ表示
	# http://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## API サーバを Linux 用にビルド
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o cmd/api/bin/api-server cmd/api/main.go

deploy: ## ビルド後にデプロイ
	make build && git add ./cmd/api/bin/api-server && git commit -m "[`date +'%Y/%m/%d %H:%M:%S'`] 最新版ビルド" && git push origin master

mod: ## package をダウンロード
	go mod download
