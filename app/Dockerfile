FROM golang:1.23-alpine AS builder

RUN apk upgrade -U
RUN apk add --no-cache \
        python3 \
        pipx \
        ffmpeg \
        curl \
        && pipx install wheel \
        && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
        && chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /app
COPY . .
RUN go mod tidy

RUN go build -o /main ./cmd/

ENTRYPOINT ["/main"]
