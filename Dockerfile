FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o event-processor ./cmd/cli

RUN ls -la /app/event-processor

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/event-processor /event-processor

RUN chmod +x /event-processor

COPY --from=builder /app/kit/config /app/kit/config

ENV CONF_DIR=/app/kit/config
ENV SCOPE=stage

RUN chmod -R 755 /app/kit/config

CMD ["/event-processor"]