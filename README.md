
# Hobigon サーバ

**自宅APIサーバ**

-> [Railsで作ったもの](https://github.com/yyh-gl/hobigon-rails-api-server)をGolangで書き直したもの


# 実装機能

## API

- 今日のタスク通知
- [ブログ](https://yyh-gl.github.io/tech-blog/)いいね
- 誕生日通知くん
- アクセスランキング通知

## CLI

[urfave/cli](https://github.com/urfave/cli)を使って実装中


# TODO
- ドメインモデル貧血症の改善
  - Birthday, Blog については完了
- GORM タグをドメインモデルから除去
  - Birthday, Blog については完了
- 構造体初期化方法を下記方法で統一
  - `h := hoge{}`
- `github.com/pkg/errors` パッケージを削除して標準パッケージの `errors` のみを使うようにする
