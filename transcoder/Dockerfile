FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o transcoder .

FROM alpine:latest

WORKDIR /app

RUN apk upgrade && apk add --no-cache ffmpeg

COPY --from=builder /app/transcoder .

ENTRYPOINT ["/app/transcoder"]
