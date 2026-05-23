# Plano: preservar formatacao media group

## Pedido do usuario
Corrigir o processamento de media groups para preservar a formatacao da legenda original do usuario. Exemplo: um album com 3 fotos e legenda em negrito deve continuar com o negrito depois que o bot adicionar a legenda configurada.

## Objetivo
Garantir que as entidades de legenda (`CaptionEntities`) do item do album que possui legenda sejam usadas na transformacao, independentemente de qual mensagem do grupo fechou o timer de processamento.

## Contexto atual
- O media group e agregado em `StageMediaGroupingTelego`.
- Cada item salvo em `GroupMessages` ja carrega `Caption` e `CaptionEntities`.
- Em `StageTransformTelego`, quando `pCtx.IsMediaGroup` e verdadeiro, o codigo procura o primeiro item com legenda.
- Porem, ele so copia as entidades se `post.MessageID == m.MessageID`, usando o `ChannelPost` do update que fechou o timer.
- Em albums, o update final pode ser outra foto/video sem legenda; nesse caso `baseText` e preservado, mas `entities` fica vazio.
- Sem `entities`, `ProcessTextWithFormattingTelego` nao recria o HTML de negrito/italico/link da legenda original.

## Arquivos analisados
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `internal/telegram/events/channelPost/stage_media_telego.go`
- `internal/telegram/events/channelPost/pipeline_telego.go`
- `internal/telegram/events/channelPost/formatting_telego.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`

## Arquivos que poderao ser modificados
- `internal/telegram/events/channelPost/stage_transform_telego.go`

## Estrategia de implementacao
No fluxo de media group, usar sempre `m.CaptionEntities` do `MediaMessageTelego` que possui a legenda, em vez de depender de `pCtx.Update.ChannelPost.CaptionEntities`.

Isso aproveita a informacao que ja e coletada em `StageMediaGroupingTelego` e corrige albums onde a legenda e de uma mensagem diferente daquela usada como update final.

## Passos detalhados

1. Alterar `StageTransformTelego` no bloco `pCtx.IsMediaGroup`.
2. Ao encontrar `m.HasCaption`, definir:
   - `baseText = m.Caption`
   - `entities = m.CaptionEntities`
3. Remover a condicao que compara `post.MessageID` com `m.MessageID`.
4. Rodar `gofmt` no arquivo alterado.
5. Rodar `go test ./...` quando possivel.
6. Rodar `go build ./cmd/FreddyBot/main.go` quando possivel.
7. Rodar `git diff --check`.

## Riscos
- Impacto restrito ao processamento de legendas em albums.
- Se algum item do album carregar entidades invalidas do Telegram, a funcao existente `ProcessEntitiesOnlyTelego` ja faz limites por UTF-16 e ignora ranges invalidos.
- O Go local pode continuar falhando por toolchain ausente (`compile`/`vet`).

## Impactos esperados
- Legendas de media group preservam negrito, italico, links embutidos, underline, spoiler e custom emoji.
- A legenda configurada do bot continua sendo anexada como antes.
- O comportamento de envio/edicao do album nao muda.

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

### Execucao
```bash
make run
```

## Rollback
Reverter a alteracao em `internal/telegram/events/channelPost/stage_transform_telego.go`.

## Observacoes
- Nao sera necessario alterar o formato salvo em `MediaMessageTelego`, porque ele ja possui `CaptionEntities`.
