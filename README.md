
# Hobigon Server

![Go](https://img.shields.io/badge/Go-1.13.8-blue.svg)
![GitHub Actions](https://img.shields.io/github/workflow/status/yyh-gl/hobigon-golang-api-server/Workflow%20for%20Golang)

**Hoby + Kanegon = Hobigon**

Hobigon is API server and it is made by me for me.

I remade that was made in Ruby(Rails).
-> [Hobigon ver.Ruby repository](https://github.com/yyh-gl/hobigon-rails-api-server) (sorry, this is private repo) 


# Features

## API
Hobigon don't use WAF but use only [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter).

- Slack Notification
  - Access ranking of [My Blog](https://yyh-gl.github.io/tech-blog/)'s posts.
  - Today's task list
  - Today's birthday people
- [My Blog](https://yyh-gl.github.io/tech-blog/)
  - Create post data
  - Get post data
  - Like👍 post
- Birthday
  - Create birthday data

## CLI
Hobigon use [urfave/cli](https://github.com/urfave/cli).

- Slack Notification
  - Access ranking of [My Blog](https://yyh-gl.github.io/tech-blog/)'s posts.
  - Today's task list
  - Today's birthday people


# TODO
- Make Hobigon DDD-like
  - Improve anemic domain model
    - Ranking
    - Slack
    - Task
