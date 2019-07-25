VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

.PHONY: build mod
help: ## この文章を表示します。
	# http://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## apiサーバをbuildします
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/api ./cmd/api

mod: ## packageをdownloadします
	go mod download

gen: ## go generateを実行します
	go generate ./...

migrate: ## 全てのテーブルをDROPして新しくMigrateし、初期データを挿入します
	go run cmd/migrate/*.go

clean_all: ## dockerコンテナを綺麗にします
	docker ps -q | xargs docker stop
	docker ps -q -a | xargs docker rm
	docker images -q | xargs docker rmi -f
