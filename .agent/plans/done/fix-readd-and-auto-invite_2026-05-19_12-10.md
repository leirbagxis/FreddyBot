# Plano: Correção de Re-vínculo e Convite Automático no PV

## Pedido do usuário
1. Erro `UNIQUE constraint failed` ao tentar adicionar um canal que já foi excluído anteriormente.
2. O bot não está mandando a mensagem de confirmação no PV automaticamente quando é adicionado como administrador em um canal.

## Objetivo
1. Resolver o erro de integridade referencial garantindo que a exclusão de um canal remova todos os seus dados dependentes (Cascade Delete).
2. Restaurar o gatilho automático que envia a mensagem de "Deseja vincular este canal?" no chat privado do usuário que adicionou o bot ao canal.

## Contexto atual
- O SQLite não estava com `foreign_keys = ON`, então o GORM deletava o canal mas deixava registros órfãos em `default_captions`.
- No `loader_telego.go`, os eventos de `MyChatMember` (quando o bot é adicionado) estão sendo filtrados pelo middleware, mas o fluxo de envio de mensagem proativa no PV parece não estar sendo disparado corretamente ou falta verificação se o canal já existe.

## Arquivos analisados
- `internal/database/database.go`
- `internal/middleware/checkAddBotMiddlewareTelego.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`
- `internal/telegram/loader_telego.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go` (concluir ativação de FK)
- `internal/telegram/handlers/events/addChannel/addChannel.go`
- `internal/middleware/checkAddBotMiddlewareTelego.go`

## Estratégia de implementação

1. **Integridade do Banco:**
   - Confirmar a ativação de `PRAGMA foreign_keys = ON;` em `database.go`. (Ação preventiva para novos vínculos).
   
2. **Convite Automático no PV:**
   - Modificar `AskAddChannelHandlerTelego` em `addChannel.go` para verificar se o canal já está cadastrado no banco de dados. Se já existir, não manda a mensagem (evita spam ao atualizar permissões).
   - Garantir que o middleware `CheckAddBotMiddlewareTelego` apenas prossiga para o handler de convite quando for uma *nova* promoção a administrador (evitar disparos em cada mudança de permissão simples).

## Passos detalhados

1. **Ajustar `addChannel.go`:**
   - Na função `AskAddChannelHandlerTelego`, adicionar uma busca rápida `c.ChannelService.GetChannelByID`. Se o canal já existir, retornar `nil` sem enviar mensagem.
   
2. **Ajustar Middleware (`checkAddBotMiddlewareTelego.go`):**
   - Na função `handleMyChatMemberTelego`, refinar a lógica de detecção de "Novo Admin" para garantir que o handler subsequente seja chamado.
   
3. **Validar Fluxo:**
   - O `loader_telego.go` já agrupa `AnyMyChatMember` com o middleware e o handler `AskAddChannelHandlerTelego`. A lógica deve fluir:
     `Update -> Middleware (valida admin/perms) -> Handler (valida existência no DB e envia PV)`.

## Riscos
- O evento `MyChatMember` pode ser disparado múltiplas vezes se o Telegram reenviar o update. A verificação de existência no DB mitiga o risco de mensagens duplicadas.

## Impactos esperados
- Canais excluídos poderão ser re-adicionados sem erros de banco de dados.
- Experiência de usuário fluida: adicionou o bot no canal -> recebeu o link de confirmação no PV na hora.

## Como testar
1. Adicionar o bot a um canal novo como admin.
2. Verificar se o bot enviou a mensagem no PV.
3. Confirmar a adição.
4. Excluir o canal pelo painel.
5. Adicionar o bot novamente ao mesmo canal.
6. Verificar se não ocorreu erro de `UNIQUE constraint` e se o convite no PV chegou novamente.

## Rollback
- Reverter as alterações no middleware e no handler de adição.
