FROM golang:1.21.1-alpine AS builder

WORKDIR /app

COPY src/go.mod ./

RUN go mod download

COPY src/ ./

RUN go build -o app

RUN adduser -D appuser

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd

USER appuser

WORKDIR /app

COPY --from=builder /app/app .

ENV LISTEN_PORT=9999

EXPOSE 9999

CMD ["./app"]