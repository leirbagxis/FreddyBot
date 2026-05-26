# Plano: alterar-mensagem-adeus-canal

## Pedido do usuário
Substituir a mensagem enviada antes do bot sair de um canal por: `Ate breve, bye 👋`.

## Objetivo
Alterar somente o texto de despedida enviado pelo serviço de desconexão de canal, sem mudar o fluxo de saída, remoção do banco ou invalidação de cache.

## Contexto atual
A mensagem de despedida fica em `internal/core/services/channels.go`, dentro de `ChannelService.DisconnectChannel`. O fluxo atual:
1. Define `farewellMsg` com uma mensagem longa.
2. Tenta enviar a mensagem no canal.
3. Executa `LeaveChat`.
4. Remove o canal do banco/cache via `DeleteChannel`.

Esse método é chamado por rotas/API, comandos admin e fluxo de remoção detectado por `my_chat_member`.

## Arquivos analisados
- `internal/core/services/channels.go`
- `internal/middleware/checkAddBotMiddlewareTelego.go`
- `internal/api/controllers/channelController.go`
- `internal/api/controllers/adminController/auditController.go`
- `internal/telegram/handlers/commands/admin/admin_channels.go`
- `internal/telegram/handlers/callbacks/my_channel/delete_channel.go`

## Arquivos que poderão ser modificados
- `internal/core/services/channels.go`

## Estratégia de implementação
Trocar apenas o valor de `farewellMsg` para `Ate breve, bye 👋`, mantendo o envio best-effort e o restante do fluxo intacto.

## Passos detalhados
1. Alterar a string `farewellMsg` em `ChannelService.DisconnectChannel`.
2. Rodar `gofmt` no arquivo.
3. Tentar `go build ./cmd/FreddyBot/main.go` e registrar resultado.

## Riscos
- Baixo risco, pois é alteração de texto.
- O ambiente Go local já apresentou erro de toolchain (`no such tool "compile"`), então a validação Go pode continuar bloqueada localmente.

## Impactos esperados
- Antes de sair do canal, o bot enviará `Ate breve, bye 👋`.
- Nenhuma mudança em banco, API, dashboard ou lógica de remoção.

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
```

### Execução
```bash
# Desconectar um canal pela dashboard ou comando admin e verificar a mensagem enviada antes do LeaveChat.
```

## Rollback
Restaurar a string anterior de `farewellMsg` em `internal/core/services/channels.go`.

## Observações
O texto pedido está sem acento em `Ate`; será mantido exatamente como solicitado.
