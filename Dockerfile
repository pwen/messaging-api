FROM golang:1.15-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY *.go .
RUN go build -o /go-msg-svr

EXPOSE 8080

CMD ["/go-msg-svr"]