# Plano: fix-dynamic-bot-username-placeholder

## Pedido do usuário
O usuário indicou que a variável/placeholder a ser utilizada no painel admin (e consequentemente no código) deve ser exata e especificamente `{usernameBot}`, e não o `{botUser}` que implementei anteriormente. A lógica de substituição ("fazer o parse") ao cadastrar um novo canal deve usar essa nomenclatura.

## Objetivo
Alterar a string de substituição dinâmica no código de vinculação de canal e atualizar a inicialização padrão no banco de dados para corresponder ao placeholder `{usernameBot}` preferido pelo usuário.

## Contexto atual
- `internal/database/database.go` usa o placeholder `{botUser}`.
- `internal/telegram/handlers/events/addChannel/addChannel.go` procura e substitui `{botUser}` pela string do nome de usuário do bot (`botInfo.Username`).

## Arquivos analisados
- `internal/database/database.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Estratégia de implementação
1. **Modificar Inicialização**: Atualizar a string em `internal/database/database.go` para conter `[t.me/legendasbot](https://t.me/{usernameBot})`.
2. **Modificar Parse**: Em `AddYesHandlerTelego` (`addChannel.go`), atualizar a regra de substituição `strings.ReplaceAll` para buscar por `{usernameBot}`.

## Passos detalhados
1. Editar `internal/database/database.go`:
   - Encontrar `GlobalDefaultCaption: "🐈‍⠀៹ [t.me/legendasbot](https://t.me/{botUser})  ‹",`
   - Substituir por `GlobalDefaultCaption: "🐈‍⠀៹ [t.me/legendasbot](https://t.me/{usernameBot})  ‹",`
2. Editar `internal/telegram/handlers/events/addChannel/addChannel.go`:
   - Encontrar `globalDefault = strings.ReplaceAll(globalDefault, "{botUser}", botInfo.Username)`
   - Substituir por `globalDefault = strings.ReplaceAll(globalDefault, "{usernameBot}", botInfo.Username)`
   - Encontrar `globalNewPack = strings.ReplaceAll(globalNewPack, "{botUser}", botInfo.Username)`
   - Substituir por `globalNewPack = strings.ReplaceAll(globalNewPack, "{usernameBot}", botInfo.Username)`

## Riscos
- Risco nulo.

## Impactos esperados
- O painel admin e os comandos do bot aceitarão a tag exata solicitada pelo usuário (`{usernameBot}`).
- Ao vincular um novo canal, o bot procurará essa string específica na configuração global para substituição.

## Compatibilidade
- Linux
- macOS
- Windows

## Como testar

### Build
`go build -o tmp/FreddyBot ./cmd/FreddyBot/`

### Execução
1. Atualizar manualmente no painel admin a legenda padrão para usar a tag `{usernameBot}` (já que o banco existente não será sobreescrito).
2. Adicionar o bot a um novo canal.
3. Verificar a tabela `default_captions` ou testar uma postagem para garantir que o `@` real do bot substituiu a tag.

## Rollback
Desfazer as strings de volta para `{botUser}` (não recomendado, pois foge da especificação do usuário).