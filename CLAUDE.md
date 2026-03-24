# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# テスト実行
make test                          # APP_ENV=test go test ./...
go test ./app/presentation/rest/... # 特定パッケージのテスト

# ビルド
make build target=rest             # REST API サーバ
make build target=cli              # CLI
make build target=graphql          # GraphQL サーバ
make build-all                     # 全ターゲット一括ビルド

# Lint（Docker Compose 経由）
make lint

# コード生成（Docker Compose 起動後）
make wire-all                      # Wire DI コード再生成（全ターゲット + test/ ）
make gqlgen                        # GraphQL コード再生成

# ローカル開発（Air ホットリロード）
docker compose up
```

## Architecture

Clean Architecture + DDD 戦術パターンを採用。依存の方向は外側から内側へ（presentation → usecase → domain ← infra）。

```
app/
├── domain/        # コアビジネスロジック（外部依存なし）
│   ├── model/     # ドメインモデル（Blog, Task, Slack, Pokemon）
│   ├── repository/# リポジトリインターフェース
│   └── gateway/   # 外部サービスインターフェース
├── infra/         # 外部依存の実装
│   ├── dao/       # repository/gateway インターフェースの実装（GORM/MySQL）
│   ├── db/        # DB 接続
│   └── dto/       # 外部サービス（Notion など）のデータ変換
├── usecase/       # アプリケーション層（ユースケースの実装）
└── presentation/
    ├── rest/      # REST ハンドラ（Gorilla Mux, ポート 3000）
    ├── graphql/   # GraphQL リゾルバ（ポート 8081）
    └── cli/       # CLI コマンド（urfave/cli）

cmd/
├── rest/          # REST エントリポイント + Wire DI 定義
├── cli/           # CLI エントリポイント + Wire DI 定義
└── graphql/       # GraphQL エントリポイント

test/              # テスト用 Wire DI 定義とユーティリティ
```

## Key Conventions

**DI（依存性注入）:** [Wire](https://github.com/google/wire) で自動生成。`wire.go` を編集後、`make wire-all` で `wire_gen.go` を再生成する。

**ORM:** GORM v1（`github.com/jinzhu/gorm`）。テスト環境は SQLite、開発・本番は MySQL を使用（`APP_ENV` で切り替え）。

**GraphQL:** スキーマは `app/presentation/graphql/*.graphqls`。コード変更後 `make gqlgen` で `generated/` を再生成する。

**環境変数:** `APP_ENV` で環境を制御（`test` / `develop` / 本番）。秘密情報は `.secret_env` に記述（Git 管理外）。

**TODO:** 貧血ドメインモデルの改善が進行中（Slack, Task）。新規実装は DDD 戦術パターン（Value Object, Entity など）を意識する。
