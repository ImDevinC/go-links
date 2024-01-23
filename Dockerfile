FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/main.go

FROM scratch

COPY --from=builder /app/main /main

EXPOSE 8080

ENTRYPOINT [ "/main" ]

