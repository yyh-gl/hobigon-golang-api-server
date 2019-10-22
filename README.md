
# Hobigon サーバ

**自宅APIサーバ**

-> [Railsで作ったもの](https://github.com/yyh-gl/hobigon-rails-api-server)をGolangで書き直したもの

現状、コード記述量を減らすために Application 層を省略して、Handler 層に統合している。  
今後処理が複雑そうになりそうだったら分離する。
-> Applicaton 層と Handler 層の分離を開始（通知系 API からスタート）

# 実装機能

- 今日のタスク通知
- [ブログ](https://yyh-gl.github.io/tech-blog/)いいね
- 誕生日通知くん
- アクセスランキング通知

# TODO
- ドメインモデル貧血症の改善
- JSON タグの場所
- Application 層の追加
  - 今まで省略してきたがだんだんきつくなってきたのでそろそろ…
