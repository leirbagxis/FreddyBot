# FreddyBot

FreddyBot é um bot de automação para Telegram focado em gerenciamento de legendas e criação de postagens personalizadas.

## 🚀 Funcionalidades

- **Legendas Automáticas**: Configure legendas padrão para canais.
- **Post Builder Interativo**: Crie postagens personalizadas com mídias, títulos formatados e botões dinâmicos.
- **Suporte Multi-Banco**: Funciona com SQLite3 (Dev) e PostgreSQL (Prod).
- **Gerenciamento de Canais**: Adicione, configure e transfira a posse de canais facilmente.

## 🛠️ Instalação e Execução

### Pré-requisitos
- Go 1.24+
- Redis (instalado e rodando)
- Token do Telegram (obtido via [@BotFather](https://t.me/BotFather))

### Configuração
1. Clone o repositório.
2. Copie o arquivo `.env-example` para `.env`:
   ```bash
   cp .env-example .env
   ```
3. Preencha as variáveis no `.env`:
   - `TELEGRAM_BOT_TOKEN`: Seu token do bot.
   - `REDIS_HOST`: Endereço do Redis (ex: `localhost:6379`).
   - `APP_ENV`: Use `dev` para SQLite ou `prod` para Postgres.
   - `DATABASE_FILE`: Caminho para o arquivo `.db` (se usar SQLite) ou string de conexão Postgres.

### Rodando o Bot
```bash
go run cmd/FreddyBot/main.go
```

## 🏗️ Estrutura do Projeto
- `cmd/`: Ponto de entrada da aplicação.
- `internal/`: Lógica central do bot, banco de dados e cache.
- `pkg/`: Configurações e utilitários globais.
- `config/`: Arquivos de mensagens YAML.

## 🛠️ Post Builder
Para usar o Post Builder, basta enviar qualquer mídia (Foto, Vídeo, Áudio, GIF ou Documento) no privado do bot. Um menu aparecerá para você construir sua postagem com Título, Corpo, Rodapé e Botões.
