FROM golang:1.21.1 AS builder

WORKDIR /app

RUN go env -w GOMODCACHE=/root/.cache/go-build

COPY src/go.mod src/go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY src/ ./

RUN go build -o app

FROM alpine:3.18.4

RUN apk update && apk upgrade

RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

RUN adduser -D appuser

USER appuser

WORKDIR /app

COPY --from=builder /app/app .

ENV LISTEN_PORT=9999

EXPOSE 9999

CMD ["./app"]