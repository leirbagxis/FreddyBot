# Plano: corrigir-panic-sticker-e-sync-admin-bot_2026-05-14_15-15.md

## Pedido do usuário
O bot deu panic ao enviar um sticker e o modo de admin/blacklist não está funcionando no bot.

## Objetivo
1. Corrigir o panic (segmentation fault) no `SetStickerSeparatorHandler`.
2. Garantir que as alterações de Admin/Blacklist feitas na Dashboard tenham efeito imediato no Bot, assim como fizemos na API.

## Contexto atual
- **Panic**: O erro ocorre em `stickerSeparator.go:112`. Analisando o código, se a sessão expira, o bot tenta chamar `b.AnswerCallbackQuery(ctx, ...)` passando `update.CallbackQuery.ID`. Porém, como o gatilho foi o envio de um **Sticker** (que é um `update.Message`), o campo `update.CallbackQuery` é `nil`, causando o panic. Além disso, `userId := update.Message.From.ID` pode falhar se o update for de outro tipo.
- **Admin/Blacklist**: O bot possui middlewares que já consultam o banco (`CheckBlacklistMiddleware` e `CheckMaintenceMiddleware`), mas o restante do bot (handlers comuns) não sabe se o usuário é Admin a menos que o handler peça explicitamente ao banco. Precisamos padronizar o uso de um "UserContext" no bot.

## Arquivos analisados
- `internal/telegram/callbacks/my_channel/stickerSeparator.go`
- `internal/middleware/blacklistMiddleware.go`
- `internal/middleware/maintenanceMiddleware.go`

## Arquivos que poderão ser modificados
- `internal/telegram/callbacks/my_channel/stickerSeparator.go`
- `internal/middleware/blacklistMiddleware.go`
- `internal/middleware/maintenanceMiddleware.go`

## Estratégia de implementação
1. **Fix Panic**: 
   - Adicionar verificações de `nil` em `SetStickerSeparatorHandler` para `update.Message`, `update.Message.From` e `update.Message.Sticker`.
   - Substituir `AnswerCallbackQuery` por `SendMessage` no caso de erro de cache, já que o input é uma mensagem (sticker).
2. **Fix Status Sync**:
   - Refatorar `getUpdateUserID` para ser mais robusto.
   - Garantir que `CheckBlacklistMiddleware` interrompa a execução corretamente.
   - Validar se outros handlers de Admin estão consultando o banco corretamente.

## Passos detalhados

1.  **Modificar `internal/telegram/callbacks/my_channel/stickerSeparator.go`**
    - Corrigir a função `SetStickerSeparatorHandler` para validar os campos do update antes de acessar.
    - Remover a chamada a `AnswerCallbackQuery` que causa o crash.

2.  **Modificar `internal/middleware/maintenanceMiddleware.go`**
    - Melhorar o `getUpdateUserID` para evitar nil dereference em casos raros.

3.  **Verificar `internal/middleware/blacklistMiddleware.go`**
    - Garantir que o bloqueio seja total para usuários na blacklist.

## Riscos
- **Baixo**: Correções de segurança e estabilidade.

## Como testar
1. Enviar um sticker para o bot em um momento aleatório (sem estar no fluxo de setar separador) -> Não deve crashar.
2. Banir um usuário pela Dashboard -> O bot deve parar de responder aos comandos dele imediatamente.
3. Promover um usuário a Admin -> Ele deve conseguir usar comandos de manutenção do bot imediatamente.

## Rollback
`git checkout ...`
