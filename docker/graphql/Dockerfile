FROM golang:1.16.0-alpine AS builder

WORKDIR /go/src/github.com/yyh-gl/hobigon-golang-graphql-server

ENV TZ="Asia/Tokyo"

COPY . .

RUN apk add --no-cache \
    ca-certificates \
    git \
    make \
    gcc \
    libc-dev \
    curl
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.30.0
RUN go get -u github.com/cosmtrek/air
RUN go get -u github.com/google/wire/cmd/wire
RUN make build target=graphql

FROM gcr.io/distroless/base

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/graphql/bin/graphql-server /app/graphql-server
COPY --from=builder /tmp /tmp

EXPOSE 8080

CMD ["/app/graphql-server"]
