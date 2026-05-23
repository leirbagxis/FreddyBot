# Plano: corrigir alvo legenda media group

## Pedido do usuario
Corrigir o bug introduzido no ajuste anterior: em media group, o bot passou a deixar uma legenda sem formatacao em uma foto e colocar a legenda final em outra. O esperado e ter apenas uma legenda no album, preservando a formatacao original e adicionando a legenda do bot.

## Objetivo
Garantir que, em albums de foto/video, a edicao final da legenda aconteca na mesma mensagem do grupo que ja possui a legenda original. Se nenhuma mensagem tiver legenda, manter fallback para a primeira mensagem.

## Contexto atual
- `StageMediaGroupingTelego` salva cada item do album com `HasCaption`, `Caption` e `CaptionEntities`.
- `StageTransformTelego` foi corrigido para usar `m.CaptionEntities` da mensagem que tem legenda.
- `ProcessMediaGroupDispatchTelego`, porem, ainda edita sempre `pCtx.GroupMessages[0]`.
- Quando a legenda original esta em uma mensagem diferente da primeira, a legenda original permanece nessa mensagem e a legenda final e aplicada em outra, gerando duas legendas no album.

## Arquivos analisados
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `internal/telegram/events/channelPost/stage_media_telego.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`

## Arquivos que poderao ser modificados
- `internal/telegram/events/channelPost/dispatch_telego.go`

## Estrategia de implementacao
No dispatch de media group para foto/video:
- escolher como alvo a primeira mensagem com `HasCaption == true`;
- se nenhuma tiver legenda, usar `GroupMessages[0]`;
- editar apenas essa mensagem com `pCtx.FormattedText`.

Isso preserva o comportamento anterior para albums sem legenda e corrige albums cuja legenda esta em qualquer item que nao seja o primeiro.

## Passos detalhados

1. Criar helper local ou bloco simples para selecionar a mensagem alvo do album.
2. Em `ProcessMediaGroupDispatchTelego`, substituir `targetMessage := pCtx.GroupMessages[0]` pela selecao da mensagem com legenda.
3. Manter `EditMessageCaption` e `ParseMode HTML` como ja estao.
4. Rodar `gofmt` no arquivo alterado.
5. Rodar `git diff --check`.
6. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando a falha local se o toolchain continuar sem `vet`/`compile`.

## Riscos
- Impacto restrito a albums de foto/video.
- Se o album nao tiver legenda original, o comportamento continua editando a primeira midia.
- Se a legenda original estiver em qualquer outra midia, a correcao evita duplicidade de legenda.

## Impactos esperados
- Album com legenda em negrito continua com uma unica legenda.
- A legenda final fica na mesma foto/video que ja tinha a legenda original.
- A formatacao original e preservada e a legenda do bot e adicionada no mesmo caption.

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
Reverter a alteracao em `internal/telegram/events/channelPost/dispatch_telego.go`.

## Observacoes
- A alteracao anterior em `StageTransformTelego` ainda e necessaria para preservar `CaptionEntities`.
- Esta correcao completa o fluxo, alinhando o texto formatado com a mensagem editada.
