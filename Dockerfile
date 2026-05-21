# Estágio 1: Build da Dashboard (Frontend)
FROM node:24-alpine AS frontend-builder
WORKDIR /app/dashboard
COPY dashboard/package*.json ./
RUN npm install
COPY dashboard/ ./
RUN npm run build

# Estágio 2: Build do Servidor (Go)
FROM golang:1.25.7-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Copia o build do frontend gerado no estágio anterior
COPY --from=frontend-builder /app/dashboard/dist ./dashboard/dist

ARG GIT_HASH=unknown
RUN go build -ldflags "-X 'github.com/leirbagxis/FreddyBot/internal/utils.Version=${GIT_HASH}'" -o server ./cmd/FreddyBot

# Estágio 3: Imagem Final
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/dashboard/dist ./dashboard/dist
COPY --from=builder /app/config ./config

EXPOSE 7000

CMD ["./server"]
