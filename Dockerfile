FROM golang:1.24.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_ENABLED=1

RUN go build -o ./bin/app ./cmd/app

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/bin/app ./app

WORKDIR /root/data/

CMD ["/root/app"]