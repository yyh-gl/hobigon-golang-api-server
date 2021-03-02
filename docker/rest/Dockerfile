FROM golang:1.16.0-alpine AS builder

WORKDIR /go/src/github.com/yyh-gl/hobigon-golang-api-server

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
RUN make build target=api

FROM gcr.io/distroless/base

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/bin/api-server /app/api-server
COPY --from=builder /tmp /tmp

EXPOSE 3000

CMD ["/app/api-server"]
