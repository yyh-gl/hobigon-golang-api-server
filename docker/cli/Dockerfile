FROM golang:1.26-alpine AS builder

WORKDIR /go/src/github.com/yyh-gl/hobigon-golang-api-server

ARG VERSION

ENV TZ="Asia/Tokyo"

COPY . .

RUN apk add --no-cache \
    ca-certificates \
    git \
    make \
    gcc \
    libc-dev
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
RUN go install github.com/air-verse/air@latest
RUN go install github.com/google/wire/cmd/wire@latest
RUN make build target=cli version=$VERSION

FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/cli/bin/hobi /app/hobi
COPY --from=builder /go/src/github.com/yyh-gl/hobigon-golang-api-server/config /app/config
COPY --from=builder /tmp /tmp

ENTRYPOINT ["/app/hobi"]
