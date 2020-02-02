
# Hobigon サーバ

**自宅APIサーバ**

-> [Railsで作ったもの](https://github.com/yyh-gl/hobigon-rails-api-server)をGolangで書き直したもの


# 機能

## API

WAFは使わずに [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) だけを使って実装

- Slack通知
  - 本日のタスク一覧
  - アクセスランキング
  - 誕生日
- [ブログ](https://yyh-gl.github.io/tech-blog/)いいね

## CLI

[urfave/cli](https://github.com/urfave/cli) を使って実装

- Slack通知
  - 本日のタスク一覧
  - アクセスランキング
  - 誕生日


# TODO
- ドメインモデル貧血症の改善
  - Ranking
  - Slack
  - Task
- 通知系APIのレスポンスで通知したコンテンツの数を返すようにする
