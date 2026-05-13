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
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install github.com/air-verse/air@latest
RUN go install github.com/google/wire/cmd/wire@latest
RUN make build target=rest version=$VERSION

FROM gcr.io/distroless/base

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/bin/api-server /app/api-server
COPY --from=builder /tmp /tmp

EXPOSE 3000

CMD ["/app/api-server"]
