FROM golang:1.23-alpine

WORKDIR /app

COPY . .
RUN go mod download
RUN go mod verify

RUN go build -o /main ./cmd/

CMD ["/main"]
