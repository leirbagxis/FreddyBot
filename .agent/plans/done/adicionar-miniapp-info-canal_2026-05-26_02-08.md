# Plano: adicionar-miniapp-info-canal

## Pedido do usuário
Adicionar no comando `/info [channelid]` um botão de miniapp que abre diretamente a dashboard daquele canal, mantendo o acesso restrito a admin do bot ou owner.

## Objetivo
Exibir as informações atuais do canal e anexar um botão `WebApp` apontando para `/dashboard/<channelID>`, garantindo que o comando só rode para owner ou usuários com `is_admin`.

## Contexto atual
O comando `/info` está em `internal/telegram/handlers/commands/admin/admin_channels.go` e já gera a URL da miniapp com `auth.GenerateMiniAppUrl`. Porém, essa URL aparece no texto da mensagem, não como botão.

Em `internal/telegram/loader_telego.go`, os comandos administrativos estão agrupados por `matchOwnerTelego()`, então atualmente o `/info` funciona apenas para o owner. A API e o dashboard já distinguem `owner` e `admin`, e admins têm acesso liberado às rotas de canal pelo middleware HTTP.

## Arquivos analisados
- `internal/telegram/handlers/commands/admin/admin_channels.go`
- `internal/telegram/loader_telego.go`
- `internal/api/auth/signature.go`
- `internal/middleware/checkAdminMiddlewareTelego.go`
- `internal/api/controllers/authController.go`
- `internal/api/auth/middleware.go`
- `internal/api/api.go`
- `dashboard/src/App.tsx`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/commands/admin/admin_channels.go`
- `internal/telegram/loader_telego.go`

## Estratégia de implementação
1. Criar um predicado Telegram para admin ou owner, consultando `UserService` quando não for o owner.
2. Usar esse predicado apenas no comando `/info`, mantendo os demais comandos administrativos no fluxo atual.
3. Remover a URL crua do corpo da mensagem do `/info` ou manter o campo textual de forma limpa, e adicionar `InlineKeyboardMarkup` com um botão `WebApp`.
4. Gerar a URL com `auth.GenerateMiniAppUrl`, usando o ID de quem executou o comando para que a autenticação do miniapp fique coerente com o usuário que clicou.
5. Manter `ParseMode` HTML e desativar preview de link na mensagem.

## Passos detalhados
1. Ajustar `loader_telego.go` para registrar `/info` em um grupo ou handler separado com permissão owner/admin.
2. Implementar `matchAdminOrOwnerTelego(c)` reaproveitando `config.OwnerID` e `c.UserService.GetUserByID`.
3. Alterar `GetInfoChannelHandlerTelego` para montar um botão `telego.InlineKeyboardButton` com `WebApp: &telego.WebAppInfo{URL: miniAppURL}`.
4. Confirmar que o comando continua buscando canal e dono do canal como antes.
5. Rodar `gofmt`, `git diff --check`, `go test ./...` e `go build ./cmd/FreddyBot/main.go` quando possível.

## Riscos
- Se o Telegram WebApp exigir domínio configurado no BotFather, o botão pode aparecer mas não abrir corretamente em ambientes sem domínio HTTPS válido.
- Como o comando deixará de ser apenas owner para aceitar admins, admins conseguirão abrir qualquer canal pela dashboard, o que já é compatível com a política atual das rotas admin/owner.
- A toolchain Go local pode impedir `go test` ou `go build`, como ocorreu em validações recentes.

## Impactos esperados
- `/info <channelID>` continua retornando os dados do canal.
- A mensagem passa a ter um botão de miniapp para abrir a dashboard daquele canal.
- Admins do bot, além do owner, podem usar `/info`.
- Demais comandos administrativos permanecem com o comportamento atual.

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
/info <channelID>
```

Validar com:
- owner do bot
- usuário `is_admin = true`
- usuário comum, que não deve receber o comando

## Rollback
Reverter as alterações em `internal/telegram/handlers/commands/admin/admin_channels.go` e `internal/telegram/loader_telego.go`, voltando o `/info` para o grupo `matchOwnerTelego()` e sem botão WebApp.

## Observações
Há alterações pendentes de tarefa anterior no worktree, incluindo normalização de URLs de botões e o arquivo não rastreado `Release`. Essas mudanças não devem ser revertidas nem incluídas indevidamente neste plano.
