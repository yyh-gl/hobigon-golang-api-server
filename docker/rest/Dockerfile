FROM golang:1.15.0-alpine AS build

WORKDIR /go/src/github.com/yyh-gl/hobigon-golang-api-server

ENV TZ="Asia/Tokyo"
ENV GO111MODULE=on

COPY . .

RUN apk add --no-cache \
    ca-certificates \
    git \
    make \
    gcc \
    libc-dev
RUN go get -u github.com/cosmtrek/air
RUN go get -u github.com/yyh-gl/realize
# TODO: なぜかGolangCI-Lintのgo getに失敗するようになったので調査する
#RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
RUN go get -u github.com/google/wire/cmd/wire
RUN go get -u github.com/tj/robo
RUN robo api-build-prod

FROM gcr.io/distroless/base

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/bin/api-server /app/api-server
COPY --from=build /tmp /tmp

EXPOSE 3000

CMD ["/app/api-server"]
