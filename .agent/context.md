# Contexto do Projeto — FreddyBot

## Stack Tecnológica
- **Linguagem:** Go (Golang) v1.24+
- **Framework Bot:** `github.com/go-telegram/bot`
- **Banco de Dados:** SQLite (GORM)
- **Cache/Sessão:** Redis
- **Dashboard:** React (TypeScript) + Vite
- **Proxy:** Nginx

## Arquitetura
- **CMD:** Ponto de entrada em `cmd/FreddyBot/main.go`.
- **Core Services:** Localizados em `internal/core/services`. É a camada obrigatória para TODA a regra de negócio e acesso a dados, servindo tanto à API quanto ao Bot.
- **Internal/API:** Handlers (Controllers) em `internal/api/controllers` que utilizam Core Services e respondem via `APIResponse[T]`.
- **Internal/Telegram:** Lógica do bot dividida em `commands`, `callbacks` e `events`, também utilizando exclusivamente Core Services.
- **Pipeline de Postagem:** Sistema modular em `internal/telegram/events/channelPost` para processar e enviar mensagens para canais.
- **Post Builder:** Criador de postagens interativo em `internal/telegram/events/postBuilder`.

## Fluxo do Sistema
1. Usuário interage com o bot via comandos ou enviando mídia.
2. Mídia detectada ativa o `PostBuilder`.
3. O `PostBuilder` utiliza Redis para manter o estado da criação.
4. Ao compartilhar via modo Inline, o bot utiliza mapeamento de `inline_message_id` para permitir atualizações de botões de voto.
