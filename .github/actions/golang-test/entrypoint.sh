#!/bin/bash

APP_DIR="/go/src/github.com/${GITHUB_REPOSITORY}/"

# shellcheck disable=SC2164
mkdir -p "${APP_DIR}" && cp -r ./ "${APP_DIR}" && cd "${APP_DIR}"

export GO111MODULE=on
go mod tidy
go mod verify

if [[ "$1" == "lint" ]]; then
    echo "############################"
    echo "# Running GolangCI-Lint... #"
    echo "############################"
    golangci-lint --version
    echo
    mv ./cmd/api/wire_gen.go ./cmd/api/wire_gen.go.tmp && echo "mv ./cmd/api/wire_gen.go ./cmd/api/wire_gen.go.tmp\n"
    golangci-lint run ./...
    mv ./cmd/api/wire_gen.go.tmp ./cmd/api/wire_gen.go && echo "mv ./cmd/api/wire_gen.go.tmp ./cmd/api/wire_gen.go\n"
fi
