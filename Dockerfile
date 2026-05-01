# Dockerfile
FROM golang:1.24.4 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ARG GIT_HASH=unknown
RUN go build -ldflags "-X 'github.com/leirbagxis/FreddyBot/internal/utils.Version=${GIT_HASH}'" -o server ./cmd/FreddyBot

# Final image
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 7000

CMD ["./server"]