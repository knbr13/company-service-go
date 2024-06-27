FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
