# Plano: Migração de Biblioteca Telegram (go-telegram/bot -> telego)

## Pedido do usuário
Migrar a biblioteca do bot para `github.com/mymmrac/telego` para aproveitar as atualizações mais recentes da API do Telegram (v8.0+), mantendo a funcionalidade total e garantindo que o código continue compilando durante o processo.

## Objetivo técnico
Substituir `github.com/go-telegram/bot` por `github.com/mymmrac/telego` em todas as camadas (Client, Middlewares, Handlers e Services) sem alterar a estrutura do banco de dados e mantendo a estabilidade do sistema.

## Contexto atual
- O bot utiliza `go-telegram/bot` v1.19.0.
- A biblioteca está integrada no `AppContainer`, serviços, middlewares e dezenas de handlers.
- O sistema usa GORM para DB e Redis para cache/sessão.

## Arquivos analisados
- `cmd/FreddyBot/main.go`
- `internal/telegram/client.go`
- `internal/container/appContainer.go`
- `internal/telegram/events/loader.go`
- `internal/middleware/saveUserMiddleware.go` (e outros middlewares)
- `internal/core/services/channels.go`

## Arquivos que poderão ser modificados
- `go.mod` / `go.sum`
- `internal/container/appContainer.go`
- `internal/telegram/client.go`
- `internal/telegram/loader.go` (e sub-loaders)
- Praticamente todos os arquivos em `internal/telegram/handlers/...`
- Todos os arquivos em `internal/middleware/...`
- `internal/core/services/channels.go` (e outros que usam o bot)

## Estratégia de implementação

Para garantir que o bot continue funcionando e o código compile, usaremos uma **estratégia de coexistência temporária e refatoração incremental**:

1.  **Fase 1: Preparação:** Adicionar `telego` ao projeto e criar uma interface/wrapper no `AppContainer` que suporte ambos os clientes ou prepare o terreno para a troca.
2.  **Fase 2: Refatoração do Container e Client:** Alterar o `AppContainer` para aceitar o cliente `telego`. Criar o novo ponto de entrada em `internal/telegram/client.go`.
3.  **Fase 3: Migração de Middlewares:** Converter os middlewares globais para o formato do `telego`.
4.  **Fase 4: Migração Incremental de Handlers:** Converter comandos, callbacks e eventos um por um. Como o `telego` usa uma estrutura de roteamento diferente, faremos a transição de todos os registros no `loader.go`.
5.  **Fase 5: Limpeza:** Remover as dependências da biblioteca antiga.

## Passos detalhados

### Parte 1: Infraestrutura e Container (O "Coração")
1. Adicionar `github.com/mymmrac/telego` ao `go.mod`.
2. Modificar `internal/container/appContainer.go`:
    - Adicionar `TelegoBot *telego.Bot` à struct.
    - Manter `Bot *bot.Bot` temporariamente se necessário para compilação.
3. Criar a nova inicialização em `internal/telegram/client_telego.go` (ou substituir `client.go` se decidirmos fazer o "big bang" controlado na inicialização).

### Parte 2: Middlewares e Loader
1. Migrar `SaveUserMiddleware`, `BlacklistMiddleware`, etc.
2. O `telego` usa `telego.HandlerFunc`, que recebe `(bot *telego.Bot, update telego.Update)`.
3. Criar o novo `loader.go` para o `telego`.

### Parte 3: Handlers (Comandos e Callbacks)
1. Converter os handlers simples primeiro (ex: `/start`, `/help`).
2. Implementar a lógica de `MatchFunc` no `telego` (que usa `telego.Predicate`).

### Parte 4: PostBuilder e ChannelPost (As partes complexas)
1. Converter a lógica de captura de mídia.
2. Converter o mapeamento de `inline_message_id`.

## Riscos
- **Incompatibilidade de Tipos:** `models.Message` do antigo vs `telego.Message` do novo.
- **Diferença de Comportamento:** A forma como os middlewares e handlers são executados.
- **Bad Request:** Mudanças sutis em como os parâmetros são enviados (especialmente ParseMode e Teclados).

## Impactos esperados
- Melhor suporte a novas features do Telegram.
- Código mais moderno e tipagem mais forte (telego é conhecido por isso).
- **Sem alteração no banco de dados.**

## Compatibilidade
- Linux/Docker/Windows (Go nativo).

## Como testar

### Build
```bash
go build -o main ./cmd/FreddyBot/main.go
```

### Execução (Modo Polling para teste rápido)
```bash
./main
```

## Rollback
- Reverter os commits de migração e voltar para a branch `main` estável.

## Observações
- Vou focar primeiro em fazer o `AppContainer` e o `client.go` aceitarem o `telego`, garantindo que o bot ligue e responda um comando simples de "Pong".
