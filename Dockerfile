FROM golang:1.17.6-alpine AS build
WORKDIR /tmp/app

COPY . .
RUN apk add --no-cache git && \
  go mod download && \
  go build -o main ./

FROM alpine:latest
WORKDIR /app

COPY --from=build /tmp/app/main /app/main

ENTRYPOINT ["/app/main"]
