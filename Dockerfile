FROM golang:1.21.6-alpine AS backend

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN go build -o /app/main ./cmd/main.go

FROM node:21-alpine AS frontend

WORKDIR /app

COPY frontend/ .

RUN npm ci && \
    npm run build

FROM gcr.io/distroless/base

COPY --from=backend /app/main /main
COPY --from=frontend /app/build/ /

EXPOSE 8080

ENTRYPOINT [ "/main" ]

