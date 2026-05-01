# Variáveis
DASHBOARD_DIR=dashboard
GO_MAIN=cmd/FreddyBot/main.go
GIT_HASH=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X 'github.com/leirbagxis/FreddyBot/internal/utils.Version=$(GIT_HASH)'"

.PHONY: all build build-ui build-server clean help

# Comando padrão: Build e Run
all: build

# Build completo (Frontend + Backend) e executa o binário
build: build-ui build-server
	@echo "🤖 Iniciando o FreddyBot ($(GIT_HASH))..."
	./server

# Instala dependências do frontend (apenas se node_modules não existir)
$(DASHBOARD_DIR)/node_modules:
	@echo "📦 Instalando dependências do Dashboard..."
	cd $(DASHBOARD_DIR) && npm install

# Build do frontend (React/Vite)
build-ui: $(DASHBOARD_DIR)/node_modules
	@echo "🚀 Iniciando build do Dashboard..."
	cd $(DASHBOARD_DIR) && npm run build
	@echo "✅ Build do Dashboard concluído!"

# Build do binário Go
build-server:
	@echo "📦 Construindo o binário do FreddyBot..."
	go build $(LDFLAGS) -o server $(GO_MAIN)
	@echo "✅ Build do servidor concluído!"

# Executa o bot usando go run (após build do dashboard)
dev: build-ui
	@echo "🛠️ Iniciando o FreddyBot em modo desenvolvimento..."
	go run $(LDFLAGS) $(GO_MAIN)

# Limpa a pasta de build do dashboard e o binário
clean:
	@echo "🧹 Limpando arquivos gerados..."
	rm -rf $(DASHBOARD_DIR)/dist
	rm -f server

# Exibe ajuda
help:
	@echo "Uso: make [comando]"
	@echo ""
	@echo "Comandos:"
	@echo "  all          - Build total e inicia o bot (padrão)"
	@echo "  build        - Build total (UI + Server) e executa o binário"
	@echo "  dev          - Build da UI e executa via 'go run'"
	@echo "  build-ui     - Apenas faz o build do dashboard"
	@echo "  build-server - Apenas constrói o binário do bot"
	@echo "  clean        - Remove arquivos gerados (dist e binário)"
