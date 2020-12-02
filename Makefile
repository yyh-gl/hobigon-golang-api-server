.DEFAULT_GOAL := help
.PHONY: help
help: ## helpを表示
	@echo '  see: https://git.dmm.com/dmm-app/pointclub-api'
	@echo ''
	@grep -E '^[%/0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build
build: ## build target=[api, cli, graphql] env=[local, prd]
	@if [ -z ${target} ]; then \
		echo 'targetを指定してください。'; \
		exit 1; \
	fi
	@if [ -z ${env} ]; then \
		echo 'envを指定してください。'; \
		exit 1; \
	fi
	
	@if [ ${target} = api -a ${env} = local ]; then \
		echo 'build api for local'; \
 		GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./docker/rest/bin/api-server ./cmd/rest/main.go ./cmd/rest/wire_gen.go; \
	fi
	@if [ ${target} = api -a ${env} = prd ]; then \
		echo 'build api for prd'; \
   		GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/rest/bin/api-server ./cmd/rest/main.go ./cmd/rest/wire_gen.go; \
   	fi
	
	@if [ ${target} = cli -a ${env} = local ]; then \
		echo 'build cli for local'; \
		GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./docker/cli/bin/hobi ./cmd/cli/main.go ./cmd/cli/wire_gen.go; \
	fi
	@if [ ${target} = cli -a ${env} = prd ]; then \
		echo 'build cli for prd'; \
		GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/cli/bin/hobi ./cmd/cli/main.go ./cmd/cli/wire_gen.go; \
	fi
	
	@if [ ${target} = graphql -a ${env} = local ]; then \
		echo 'build graphql for local'; \
		GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./docker/graphql/bin/graphql-server ./cmd/graphql/main.go; \
	fi
	@if [ ${target} = graphql -a ${env} = prd ]; then \
		echo 'build graphql for prd'; \
		GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/graphql/bin/graphql-server ./cmd/graphql/main.go; \
	fi
	
	@exit 0;

.PHONY: build-all
build-all: ## build-all env=[local, prd]
	@if [ -z ${env} ]; then \
		echo 'envを指定してください。'; \
		exit 1; \
	fi
	
	@if [ ${env} = local ]; then \
		make build target=api env=local; \
		make build target=cli env=local; \
		make build target=graphql env=local; \
	fi
	
	@if [ ${env} = prd ]; then \
		make build target=api env=prd; \
		make build target=cli env=prd; \
		make build target=graphql env=prd; \
	fi

.PHONY: wire-all
wire-all: ## all wire gen
	cd ./cmd/rest && wire
	cd ./cmd/cli && wire
	cd ./test && wire

.PHONY: deploy
deploy: ## deploy to prd
	make build-all env=prd
	git add ./cmd/rest/bin/api-server ./cmd/cli/bin/hobi ./cmd/graphql/bin/graphql-server
	git commit -m "[`date +'%Y/%m/%d %H:%M:%S'`] 最新版デプロイ"
	git push origin master

.PHONY: test
test: ## go test
	APP_ENV=test go test ./...

.PHONY: lint
lint: ## lint
	docker-compose exec -T rest golangci-lint --timeout 5m0s run ./...

