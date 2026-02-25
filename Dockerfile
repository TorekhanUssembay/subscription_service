FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o subscription_service ./cmd

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/subscription_service .

RUN apk add --no-cache ca-certificates

EXPOSE 8080

CMD ["./subscription_service"]