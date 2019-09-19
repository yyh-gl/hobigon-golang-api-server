FROM golang:1.12.5-alpine AS build

WORKDIR /go/src/github.com/yyh-gl/hobigon-golang-api-server

RUN apk add --no-cache ca-certificates git make gcc libc-dev \
    && go get github.com/oxequa/realize \
    && go get github.com/golangci/golangci-lint/cmd/golangci-lint

COPY . .

RUN make build

FROM alpine AS app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/api/bin/api-server /app/api-server
COPY --from=build /tmp /tmp

EXPOSE 3000

CMD ["/app/api-server"]
