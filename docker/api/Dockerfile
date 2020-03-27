FROM golang:1.14.1-alpine AS build

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
RUN robo api-build-prod

FROM alpine AS app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/api/bin/api-server /app/api-server
COPY --from=build /tmp /tmp

EXPOSE 3000

CMD ["/app/api-server"]
