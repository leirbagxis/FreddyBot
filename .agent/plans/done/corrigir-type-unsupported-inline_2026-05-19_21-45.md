# Plano: corrigir-type-unsupported-inline

## Pedido do usuário
O usuário reportou um erro no modo inline ao tentar compartilhar uma postagem do PostBuilder: `telego: answerInlineQuery: api: 400 "Bad Request: can't parse inline query result: type \"\" is unsupported for the inline query result"`.

## Objetivo
Corrigir a ausência do campo `Type` nos structs de resposta `InlineQueryResult`, garantindo que o Telegram consiga parsear o resultado com sucesso. 

## Contexto atual
- A biblioteca `telego` requer que o campo `Type` das structs de `InlineQueryResult` seja preenchido (ex: `Type: "article"`).
- Atualmente, em várias partes do sistema, esse campo está sendo omitido. Isso faz com que a API do Telegram receba `type: ""` (vazio), o que gera o erro HTTP 400.
- O problema afeta o PostBuilder, o fluxo de "Claim" (recuperar canal) e o "Maintenance" (manutenção).

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/claimChannel/claim.go`
- `internal/middleware/maintenanceMiddlewareTelego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/claimChannel/claim.go`
- `internal/middleware/maintenanceMiddlewareTelego.go`

## Estratégia de implementação
1. **Adicionar o campo Type:** Atualizar todos os structs de `telego.InlineQueryResult` (como `InlineQueryResultArticle`, `InlineQueryResultCachedPhoto`, etc.) para incluir explicitamente o campo `Type` com a respectiva string exigida pela API do Telegram (`"article"`, `"photo"`, `"video"`, `"mpeg4_gif"`, `"audio"`, `"document"`, `"sticker"`).
2. **Revisar todos os handlers Inline:** Garantir que tanto o caso de erro (Postagem não encontrada/Manutenção) quanto o caso de sucesso contenham o tipo correto.

## Passos detalhados

1. **Em `internal/telegram/handlers/events/postBuilder/postBuilder.go`:**
   - No caso de erro (não encontrada), adicionar `Type: "article"` em `InlineQueryResultArticle`.
   - No switch `state.MediaType`, adicionar `Type: "photo"`, `Type: "video"`, `Type: "mpeg4_gif"`, `Type: "audio"`, `Type: "document"`, `Type: "sticker"` e `Type: "article"` (default) em suas respectivas structs.
   
2. **Em `internal/telegram/handlers/callbacks/claimChannel/claim.go`:**
   - Adicionar `Type: "article"` em todas as instâncias de `telego.InlineQueryResultArticle` (incluindo a função auxiliar `buildErrorArticleTelego`).

3. **Em `internal/middleware/maintenanceMiddlewareTelego.go`:**
   - Adicionar `Type: "article"` na estrutura `InlineQueryResultArticle` na função `sendMaintenanceResponseTelego`.

## Riscos
- Risco nulo. Esta é uma exigência estrita da API do Telegram e resolverá diretamente o HTTP 400.

## Impactos esperados
- O modo inline de todos os módulos voltará a funcionar corretamente.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Testes
1. Criar um post no PostBuilder e testar o modo inline: `@FreddyCaptionBot pb <id>`.
2. O preview deve ser exibido na tela, independente de ser foto, vídeo, sticker ou texto puro.

### Execução
```bash
go run cmd/FreddyBot/main.go
```

## Rollback
Desfazer as alterações nos três arquivos.