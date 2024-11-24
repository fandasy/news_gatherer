FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml ./config/config.yaml

ENV CONFIG_PATH=./config/config.yaml
ENV TG_TOKEN=
ENV VK_TOKEN=
ENV YA_GPT_TOKEN=

RUN chmod +x ./main

CMD ["./main"]
