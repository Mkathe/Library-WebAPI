# syntax=docker/dockerfile:1
FROM golang:tip-bullseye AS builder
LABEL maintainer="Magzhan Yedilbayev"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

FROM gcr.io/distroless/static

COPY --from=builder /app/main /app/main

EXPOSE 3000
ENTRYPOINT ["/app/main"]
