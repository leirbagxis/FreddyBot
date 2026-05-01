# Variáveis
DASHBOARD_DIR=dashboard
GO_MAIN=cmd/FreddyBot/main.go
GIT_HASH=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X 'github.com/leirbagxis/FreddyBot/internal/utils.Version=$(GIT_HASH)'"

.PHONY: all build build-ui build-server run clean help

# Comando padrão: Apenas build
all: build

# Build completo (Frontend + Backend)
build: build-ui build-server

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

# Executa o bot usando o binário compilado
run: build-server
	@echo "🤖 Iniciando o FreddyBot ($(GIT_HASH))..."
	./server

# Atalho para build e run
dev: build run

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
	@echo "  build        - Build total (Frontend + Backend)"
	@echo "  build-ui     - Apenas faz o build do dashboard"
	@echo "  build-server - Apenas constrói o binário do bot"
	@echo "  run          - Inicia o bot Go (via go run)"
	@echo "  clean        - Remove arquivos gerados (dist e binário)"
