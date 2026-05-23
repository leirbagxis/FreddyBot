# Plano: fallback-emoji-normal-inline-postbuilder

## Pedido do usuário
Como o inline nao tem suporte confiavel a custom emoji nos botoes, o resultado inline do PostBuilder deve enviar os emojis normais/fallback no texto do botao.

## Objetivo
Manter custom emoji real no preview/envio direto e usar o emoji normal no resultado via `@FreddyCaptionBot pb <id>`.

## Contexto atual
- O preview usa `sendFinalPostTelego` e deve continuar usando `IconCustomEmojiID`.
- O resultado inline usa `InlineHandlerTelego`.
- A alteracao anterior passou a remover o emoji textual do botao quando existe `CustomEmojiID`, o que deixa o inline sem emoji normal para mostrar.
- A payload original tinha exatamente o dado necessario para fallback inline: `text: "🤖 Legendas BOT"` e `custom_emoji_id`.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/cache/types.go`

## Arquivos que poderao ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estrategia de implementacao
Separar dois modos de montagem de botoes:
- direto: remove o emoji textual inicial e usa `IconCustomEmojiID`;
- inline: preserva `Text` como veio do Redis, incluindo o emoji normal, e nao depende do `IconCustomEmojiID`.

Tambem parar de normalizar/remover o emoji textual antes de salvar a sessao, para que o inline tenha fallback Unicode.

## Passos detalhados

1. Remover ou deixar de usar `normalizePostBuilderStateButtons` no `pb-save`.
2. Criar/ajustar helper `buildPostBuilderURLButton` para receber um modo de renderizacao.
3. No preview/envio direto (`sendFinalPostTelego`), usar modo direto:
   - `Text` sem emoji textual inicial;
   - `IconCustomEmojiID` preenchido.
4. No inline (`InlineHandlerTelego`), usar modo inline:
   - `Text` preservado com emoji normal;
   - sem depender de `IconCustomEmojiID`.
5. Na captura de novos botoes, manter o texto original no estado quando houver custom emoji, para preservar o fallback inline.
6. Rodar `gofmt`.
7. Tentar `go build` e registrar limitacao se o ambiente Go continuar sem `compile`.

## Riscos
- Preview e inline vao ter aparencia diferente por limitacao do Telegram: preview com custom emoji, inline com emoji normal.
- Sessoes criadas depois da correcao anterior podem ja estar salvas sem o emoji normal; essas precisarao ser recriadas para ter fallback no inline.

## Impactos esperados
- Novos posts salvos pelo PostBuilder terao resultado inline com emojis normais nos botoes.
- Preview/envio direto continuara com custom emoji real.
- Sem alteracao em banco, API ou Dashboard.

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

### Execucao
```bash
make dev
```

Teste manual:
1. Criar um novo PostBuilder com botao iniciado por custom emoji.
2. Ver preview: deve mostrar custom emoji.
3. Salvar.
4. Usar `@FreddyCaptionBot pb <id>`.
5. Confirmar que o resultado inline mostra emoji normal no texto do botao.

## Rollback
Reverter as alteracoes em `internal/telegram/handlers/events/postBuilder/postBuilder.go`.

## Observacoes
- Para sessoes ja salvas sem emoji textual fallback, o usuario deve salvar novamente o post.
