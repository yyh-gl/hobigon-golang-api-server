
# Hobigon Server

![Go](https://img.shields.io/badge/Go-1.14.4-blue.svg)
![GitHub Actions](https://img.shields.io/github/workflow/status/yyh-gl/hobigon-golang-api-server/Workflow%20for%20Golang)

**Hobby + [KANEGON](https://lh3.googleusercontent.com/proxy/7iv-F5GG-BP9nAEEc85VJ0Uh-lAF47GkRQqWIPKSem4r1QNNnrAeIHyUExGd-gWBxEehqi9k6SOBbe8F41VdKYJj5lIOULIQeSiCJsCKaDyUHhQ) = Hobigon**

Hobigon is API server and it is made by me for me.

I remade that was made in Ruby(Rails).
-> [Hobigon ver.Ruby repository](https://github.com/yyh-gl/hobigon-rails-api-server) (sorry, this is private repo) 


# Features

## API
Hobigon don't use WAF and use only [gorilla/mux](https://github.com/gorilla/mux).

(Use [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) 
until [This commit](https://github.com/yyh-gl/hobigon-golang-api-server/tree/b0c0fb3e52df7714593386840e64a9bf7f32f1a4))

### Features

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

### Features

- Slack Notification
  - Access ranking of [My Blog](https://yyh-gl.github.io/tech-blog/)'s posts.
  - Today's task list
  - Today's birthday people


# TODO
- Make Hobigon DDD-like
  - Improve anemic domain model
    - Slack
    - Task
