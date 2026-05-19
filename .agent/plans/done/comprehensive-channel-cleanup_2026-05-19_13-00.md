# Plano: comprehensive-channel-cleanup

## Pedido do usuário
O usuário deseja garantir que, ao remover um canal, todos os registros relacionados sejam removidos do banco de dados, do Redis e da memória RAM local.

## Objetivo
Centralizar e fortalecer a lógica de limpeza de cache (Redis e RAM) no serviço principal de canais, assegurando que exclusões feitas tanto pelo usuário quanto por administradores limpem adequadamente o estado do sistema.

## Contexto atual
- O banco de dados já deleta os registros usando exclusão em cascata (`DeleteChannelWithRelations`).
- O Redis é limpo para os metadados do canal via `InvalidateChannel`.
- A exclusão de sessões ativas do usuário no Redis (`DeleteAllUserSessionsBySuffix`) só é chamada na exclusão iniciada pelo dono, não pelo admin.
- O cache local na memória RAM (`go-cache`) é limpo individualmente para metadados, mas não há um limpador de sufixo para sessões locais.

## Arquivos analisados
- `internal/cache/cache.go`
- `internal/core/services/channels.go`
- `internal/telegram/handlers/callbacks/my_channel/delete_channel.go`

## Arquivos que poderão ser modificados
- `internal/cache/cache.go`
- `internal/core/services/channels.go`
- `internal/telegram/handlers/callbacks/my_channel/delete_channel.go`

## Estratégia de implementação
1. **Reforçar CacheService**:
   - Atualizar `DeleteAllUserSessionsBySuffix` para também iterar e remover as chaves correspondentes da memória RAM (`localCache`).
2. **Centralizar Limpeza**:
   - Mover a chamada de `DeleteAllUserSessionsBySuffix` para dentro de `ChannelService.DeleteChannel` ou garantir que a desconexão do canal limpe o cache do dono de forma consistente.
3. **Limpeza do Handler Antigo**:
   - Remover a chamada redundante de limpeza do Redis no manipulador de callback (`ConfirmDeleteChannelHandlerTelego`), deixando a responsabilidade para o serviço `ChannelService`.

## Passos detalhados
1. Em `internal/cache/cache.go`, na função `DeleteAllUserSessionsBySuffix`:
   - Adicionar lógica para iterar por `localCache.Items()` e deletar se a chave terminar com o sufixo (ex: `:123`).
2. Em `internal/core/services/channels.go`, na função `DeleteChannel`:
   - Adicionar a chamada `s.cache.DeleteAllUserSessionsBySuffix(ctx, userID)` logo após `s.cache.InvalidateChannel`.
3. Em `internal/telegram/handlers/callbacks/my_channel/delete_channel.go`:
   - Remover `_, _ = c.CacheService.DeleteAllUserSessionsBySuffix(context.Background(), userId)` do handler de confirmação, pois isso passará a ser feito de forma automática pelo serviço core.

## Riscos
- Risco de limpar sessões de forma muito agressiva. Mas como a regra atual já era "limpar sessões do usuário ao deletar canal", a mudança apenas garante que isso ocorra em todas as vias (admin e owner) e incluindo a RAM.

## Impactos esperados
- Não haverá mais sessões "fantasmas" no Redis ou RAM ("Selecione o canal" apontando para o vazio) ao excluir um canal, não importa como seja feita a exclusão.

## Compatibilidade
- Linux
- macOS
- Windows

## Como testar

### Build
`go build -o tmp/FreddyBot ./cmd/FreddyBot/`

### Execução
1. Criar e vincular um canal.
2. Selecionar o canal para gerenciar.
3. Usar o comando `/remove [id]` via Admin para deletá-lo.
4. Tentar interagir com o bot (ou verificar o Redis) e assegurar que as chaves antigas (`selected_channel`, etc.) foram limpas.

## Rollback
Desfazer as alterações de refatoração no serviço de canais e cache.