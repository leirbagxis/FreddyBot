# Dockerfile
FROM golang:1.24.4 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY .env .env
RUN go mod download

COPY . .

RUN go build -o server ./cmd/FreddyBot

# Final image
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 7000

CMD ["./server"]
