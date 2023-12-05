FROM golang:alpine AS builder


# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git bash && mkdir -p /build/epos_bot

WORKDIR /build/epos_bot

RUN go install github.com/cortesi/modd/cmd/modd@latest

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download -json

COPY . .

RUN mkdir -p /app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w -extldflags="-static"' -o /app/epos_bot

FROM golang:alpine AS prod
COPY --from=builder /app/epos_bot /app/epos_bot
WORKDIR /build/epos_bot
COPY ./configs/config.json /build/epos_bot/configs/config.json

LABEL org.opencontainers.image.description A docker image for the epos_bot telegram bot.

RUN chmod -R 755 /app

ENTRYPOINT ["/app/epos_bot"]
