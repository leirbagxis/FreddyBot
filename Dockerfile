FROM golang:1.24.4-alpine

WORKDIR /app

RUN apk add --no-cache gcc g++ musl-dev

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN go build -o app ./cmd/FreddyBot

CMD ["./app"]
