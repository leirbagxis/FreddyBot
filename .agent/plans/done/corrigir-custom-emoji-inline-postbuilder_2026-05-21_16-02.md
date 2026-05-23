# Plano: corrigir-custom-emoji-inline-postbuilder

## Pedido do usuário
No PostBuilder, o preview mostra os botões com custom emoji corretamente, mas depois de salvar e usar o inline `@FreddyCaptionBot pb <id>`, o resultado aparece sem o custom emoji.

## Objetivo
Garantir que os botões do resultado inline salvo preservem `custom_emoji_id` e renderizem o mesmo teclado exibido no preview.

## Contexto atual
- O preview usa `sendFinalPostTelego`, que monta o teclado com `buildPostBuilderURLButton`.
- O resultado inline usa `InlineHandlerTelego`, que também monta o teclado com `buildPostBuilderURLButton`.
- A sessão salva usa `SavePostBuilderSession`, que serializa `PostBuilderState` no Redis.
- O usuário mostrou uma payload Redis com `buttons[].custom_emoji_id` preenchido, então o dado existe.
- A divergência pode vir de:
  - sessão salva antes de normalizar o texto do botão;
  - renderização inline usando `Text` com fallback/emoji textual incorreto;
  - cache do resultado inline do Telegram reaproveitando uma resposta anterior sem o campo do ícone.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/cache/cache.go`
- `internal/cache/types.go`
- Payload Redis enviada pelo usuário

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação
Centralizar a normalização dos botões do PostBuilder e aplicá-la em três pontos: captura do botão, salvamento da sessão e renderização do teclado. Para o inline, responder a query como pessoal e com cache zero para reduzir reaproveitamento de resultado antigo pelo Telegram.

## Passos detalhados

1. Ajustar `buildPostBuilderURLButton` para, quando `CustomEmojiID` existir:
   - remover o emoji textual inicial do label se ele ainda estiver presente;
   - preservar o texto restante;
   - usar `" "` quando não sobrar texto;
   - setar `IconCustomEmojiID`.
2. Adicionar helper para normalizar todos os botões de um `PostBuilderState`.
3. Antes de `SavePostBuilderSession`, salvar uma cópia do estado com botões normalizados.
4. Manter a renderização do preview e inline usando o mesmo helper.
5. No `AnswerInlineQuery`, definir `IsPersonal: true` junto com `CacheTime: 0`.
6. Rodar `gofmt`.
7. Tentar `go build` e registrar limitação se o Go local continuar sem `compile`.

## Riscos
- Se um botão com `CustomEmojiID` começar intencionalmente com outro emoji Unicode, esse primeiro emoji será removido. No fluxo do Telegram, esse emoji textual é o fallback do custom emoji, então o comportamento é esperado para este caso.
- O cache do cliente Telegram pode ainda mostrar resultado antigo por alguns instantes se a mesma query já tiver sido usada antes.

## Impactos esperados
- Preview e resultado inline passam a usar a mesma normalização.
- Sessões novas salvas não carregam o emoji textual duplicado no `text`.
- Sessões antigas do Redis ainda renderizam corretamente, porque a limpeza também acontece no helper de renderização.
- Sem alteração em schema, API ou Dashboard.

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
make dev
```

Teste manual:
1. Criar PostBuilder com botão cujo texto começa com custom emoji.
2. Clicar em preview e confirmar o botão com custom emoji.
3. Clicar em salvar.
4. Usar `@FreddyCaptionBot pb <id>`.
5. Confirmar que o resultado inline enviado mantém o custom emoji no botão.

## Rollback
Reverter as alterações em `internal/telegram/handlers/events/postBuilder/postBuilder.go`.

## Observações
- Não será alterado o contrato do Redis.
- A correção deve cobrir tanto sessões novas quanto payloads já existentes.
