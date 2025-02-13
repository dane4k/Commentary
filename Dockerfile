FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main ./app/cmd/main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD ["/app/main"]
LABEL authors="danya"
