# Plano: ignorar edicoes manutencao blacklist

## Pedido do usuário
Garantir que, quando o sistema estiver em manutenção ou quando o dono do canal estiver na blacklist, o bot ignore solicitações de edição/processamento de postagens, sem tentar tratar essas postagens.

## Objetivo
Reforçar o fluxo de segurança para que postagens de canal e edições de postagens de canal não sejam processadas indevidamente durante manutenção ou quando o proprietário do canal estiver bloqueado.

## Contexto atual
- O handler principal de postagens de canal processa `channel_post`.
- O polling já solicita `edited_channel_post` em `AllowedUpdates`.
- O `StagePreflightTelego` já interrompe silenciosamente `channel_post` quando:
  - o servidor está em manutenção;
  - o canal pertence a um usuário em blacklist.
- O middleware de blacklist já deixa `channel_post` e `edited_channel_post` passarem para o fluxo específico de canal.
- O middleware de manutenção só ignora `channel_post`, mas não ignora explicitamente `edited_channel_post`.
- Não há handler registrado para `edited_channel_post`, então hoje ele tende a não ser tratado; mesmo assim, o comportamento deve ficar explícito para evitar regressões.

## Arquivos analisados
- `internal/telegram/loader_telego.go`
- `internal/telegram/events/channelPost/stage_preflight_telego.go`
- `internal/telegram/events/channelPost/channelPost.go`
- `internal/middleware/blacklistMiddlewareTelego.go`
- `internal/middleware/maintenanceMiddlewareTelego.go`
- `internal/middleware/utilsTelego.go`
- `internal/telegram/client.go`

## Arquivos que poderão ser modificados
- `internal/middleware/maintenanceMiddlewareTelego.go`
- `internal/telegram/events/channelPost/channelPost.go`

## Estratégia de implementação
Aplicar uma proteção em duas camadas:

1. No middleware de manutenção, tratar `edited_channel_post` igual `channel_post`: deixar passar sem resposta de manutenção, porque postagens de canal são decididas pelo pipeline específico.
2. No handler de postagens de canal, adicionar um retorno explícito para `edited_channel_post`, documentando que edições de posts de canal não devem ser processadas pelo pipeline de edição automática.

O fluxo já existente do `StagePreflightTelego` continuará responsável por bloquear `channel_post` em manutenção e blacklist de dono do canal.

## Passos detalhados

1. Alterar `CheckMaintenanceMiddlewareTelego` para retornar `ctx.Next(upt)` quando `upt.EditedChannelPost != nil`.
2. Adicionar guarda explícita em `channelpost.HandlerTelego` para ignorar `update.EditedChannelPost`.
3. Manter `StagePreflightTelego` sem mudança, pois ele já bloqueia `channel_post` em manutenção e blacklist.
4. Rodar `gofmt` nos arquivos alterados.
5. Rodar `go test ./...` quando possível.
6. Rodar `go build ./cmd/FreddyBot/main.go` quando possível.
7. Rodar `git diff --check`.

## Riscos
- Impacto baixo e restrito ao roteamento de updates de canal.
- Se no futuro for desejado processar `edited_channel_post`, será necessário criar um fluxo próprio com preflight equivalente.

## Impactos esperados
- Postagens novas de canal continuam sendo ignoradas em manutenção e blacklist pelo preflight atual.
- Edições de postagens de canal ficam explicitamente ignoradas e não recebem resposta de manutenção.
- Reduz risco de regressão caso alguém registre `edited_channel_post` no handler no futuro.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./...
git diff --check
```

### Execução
```bash
go run ./cmd/FreddyBot/main.go
```

## Rollback
Reverter as alterações em:
- `internal/middleware/maintenanceMiddlewareTelego.go`
- `internal/telegram/events/channelPost/channelPost.go`

## Observações
- O comportamento de blacklist para dono do canal já existe em `StagePreflightTelego`.
- O comportamento de manutenção para `channel_post` também já existe em `StagePreflightTelego`.
- Esta mudança é principalmente de garantia explícita para `edited_channel_post` e proteção contra regressões.
