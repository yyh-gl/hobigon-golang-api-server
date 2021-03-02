.DEFAULT_GOAL := help
.PHONY: help
help: ## helpを表示
	@echo '  see: https://git.dmm.com/dmm-app/pointclub-api'
	@echo ''
	@grep -E '^[%/0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build
build: ## build target=[api, cli, graphql]
	@if [ -z ${target} ]; then \
		echo 'targetを指定してください。'; \
		exit 1; \
	fi
	
	@if [ ${target} = api ]; then \
		echo 'build api'; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/rest/bin/api-server ./cmd/rest; \
	fi
	
	@if [ ${target} = cli ]; then \
		echo 'build cli'; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/cli/bin/hobi ./cmd/cli; \
	fi
	
	@if [ ${target} = graphql ]; then \
		echo 'build graphql'; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/graphql/bin/graphql-server ./cmd/graphql; \
	fi
	
	@exit 0;

.PHONY: build-all
build-all: ## build-all
	make build target=api
	make build target=cli
	make build target=graphql

.PHONY: wire-all
wire-all: ## all wire gen
	docker-compose exec -T -w /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/rest rest wire
	docker-compose exec -T -w /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/cli rest wire
	docker-compose exec -T -w /go/src/github.com/yyh-gl/hobigon-golang-api-server/test rest wire

.PHONY: test
test: ## go test
	APP_ENV=test go test ./...

.PHONY: lint
lint: ## lint
	docker-compose exec -T rest golangci-lint --timeout 5m0s run ./...

