# Variáveis
DASHBOARD_DIR=dashboard
GO_MAIN=cmd/FreddyBot/main.go

.PHONY: all build-ui run clean help

# Comando padrão: Build do UI e executa o Bot
all: build-ui run

# Build do frontend (React/Vite)
build-ui:
	@echo "🚀 Iniciando build do Dashboard..."
	cd $(DASHBOARD_DIR) && npm run build
	@echo "✅ Build do Dashboard concluído!"

# Executa o backend Go
run:
	@echo "🤖 Iniciando o FreddyBot..."
	go run $(GO_MAIN)

# Limpa a pasta de build do dashboard
clean:
	@echo "🧹 Limpando arquivos gerados..."
	rm -rf $(DASHBOARD_DIR)/dist

# Exibe ajuda
help:
	@echo "Uso: make [comando]"
	@echo ""
	@echo "Comandos:"
	@echo "  all       - Build do dashboard e inicia o bot (padrão)"
	@echo "  build-ui  - Apenas faz o build do dashboard"
	@echo "  run       - Apenas inicia o bot Go"
	@echo "  clean     - Remove a pasta dist do dashboard"
