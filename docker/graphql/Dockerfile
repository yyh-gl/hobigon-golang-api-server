FROM golang:1.26-alpine AS builder

WORKDIR /go/src/github.com/yyh-gl/hobigon-golang-graphql-server

ARG VERSION

ENV TZ="Asia/Tokyo"

COPY . .

RUN apk add --no-cache \
    ca-certificates \
    git \
    make \
    gcc \
    libc-dev \
    curl
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN go install github.com/air-verse/air@latest
RUN go install github.com/google/wire/cmd/wire@latest
RUN go install github.com/99designs/gqlgen@latest
RUN make build target=graphql version=$VERSION

FROM gcr.io/distroless/base

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/yyh-gl/hobigon-golang-api-server/cmd/graphql/bin/graphql-server /app/graphql-server
COPY --from=builder /tmp /tmp

EXPOSE 8080

CMD ["/app/graphql-server"]
