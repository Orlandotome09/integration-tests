# Build phase
FROM golang:1.19.3-alpine AS builder

LABEL maintainer="Bexs Digital"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go mod download

WORKDIR /app/src

RUN go build -o temis-compliance

# This phase will create a new docker image, only with the executable, leaving the source code on the previous image
FROM alpine

WORKDIR /app

COPY --from=builder /app/src/temis-compliance .

COPY --from=builder /app/src/docs .

CMD ["./temis-compliance"]