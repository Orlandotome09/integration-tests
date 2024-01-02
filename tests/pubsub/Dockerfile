# Build phase
FROM golang:1.19.3-alpine

COPY . /app

WORKDIR /app

RUN go mod download

RUN go build -o temis-query-pubsub-config

RUN chmod +x temis-query-pubsub-config

ENTRYPOINT ["./temis-query-pubsub-config"]