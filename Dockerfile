FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/main.go

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]

