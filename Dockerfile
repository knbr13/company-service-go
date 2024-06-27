FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

EXPOSE 8080

CMD ["./main"]
