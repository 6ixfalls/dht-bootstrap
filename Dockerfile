FROM golang:1.21-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /dht-bootstrap

EXPOSE 4001

CMD ["/dht-bootstrap"]