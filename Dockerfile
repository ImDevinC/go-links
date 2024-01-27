FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.6-alpine AS backend

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/main ./cmd/main.go

FROM node:21-alpine AS frontend

WORKDIR /app

COPY frontend/ .

RUN npm ci && \
    npm run build

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

COPY --from=backend /app/main /main
COPY --from=frontend /app/build/ /

EXPOSE 8080

ENTRYPOINT [ "/main" ]

