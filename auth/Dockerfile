FROM golang:1.23-alpine

WORKDIR /app

COPY . ./
RUN go mod download
RUN go mod verify

RUN go build -o /auth ./cmd/auth/main.go
RUN go build -o /migrate ./cmd/migrator/main.go
