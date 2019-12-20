FROM golang:1.13.5-alpine as build
WORKDIR /src
ARG APP_NAME=accounts

COPY go.mod go.sum      ./
COPY vendor             ./vendor
COPY internal/pkg       ./internal/pkg
COPY internal/$APP_NAME ./internal/$APP_NAME
COPY cmd/$APP_NAME      ./cmd/$APP_NAME

ENV CGO_ENABLED 0
RUN go test  -mod=vendor ./internal/...
RUN go test  -mod=vendor ./cmd/...
RUN go build -mod=vendor -o /app/$APP_NAME ./cmd/$APP_NAME


FROM alpine:3.10.2
WORKDIR /app
COPY --from=build /app/accounts ./
ENTRYPOINT ["./accounts"]
