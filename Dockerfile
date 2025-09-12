FROM golang:1.25-bookworm as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app ./cmd/trtl3/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 7713

CMD ["./app"]
