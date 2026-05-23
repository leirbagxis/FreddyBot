# Plano: corrigir-botao-emoji-custom-postbuilder

## Pedido do usuario
Corrigir o PostBuilder para que, ao adicionar um botao com emoji customizado, nao sejam exibidos dois botoes/duas representacoes: o emoji customizado e o botao padrao.

## Objetivo
Garantir que botoes de URL do PostBuilder com `IconCustomEmojiID` renderizem apenas o icone customizado esperado, sem manter o emoji textual original duplicado no `Text` do botao.

## Contexto atual
- `PostBuilderButton` possui `Text`, `URL` e `CustomEmojiID`.
- No fluxo `awaiting_button`, o nome do botao e salvo com a primeira linha inteira digitada pelo usuario.
- Se a primeira linha contem um emoji customizado, o codigo extrai `CustomEmojiID`, mas mantem o texto original em `Text`.
- Na renderizacao, o botao e criado com `Text: btn.Text` e `IconCustomEmojiID: btn.CustomEmojiID`.
- Resultado provavel: Telegram mostra o emoji/texto padrao junto do icone customizado.

## Arquivos analisados
- `AGENTS.md`
- `.agent/context.md`
- `.agent/plans/done/add-custom-emoji-buttons-postbuilder_2026-05-20_11-30.md`
- `.agent/plans/done/add-custom-emoji-postbuilder_2026-05-20_11-30.md`
- `internal/cache/types.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderao ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estrategia de implementacao
Centralizar a criacao dos botoes inline do PostBuilder em uma pequena funcao helper. Essa funcao montara o `telego.InlineKeyboardButton` de URL e, quando houver `CustomEmojiID`, removera do `Text` a parte textual correspondente ao emoji customizado. Se o texto restante ficar vazio, usara um espaco minimo, igual ao padrao ja usado nas reacoes customizadas.

## Passos detalhados

1. Adicionar helper em `postBuilder.go` para montar `telego.InlineKeyboardButton` a partir de `cache.PostBuilderButton`.
2. O helper deve:
   - preencher `URL`;
   - quando `CustomEmojiID` estiver vazio, manter `Text` original;
   - quando `CustomEmojiID` existir, remover o placeholder/texto do emoji customizado do label quando possivel;
   - se o label ficar vazio, usar `" "` como texto minimo;
   - preencher `IconCustomEmojiID`.
3. Substituir as duas renderizacoes atuais de botoes do PostBuilder pelo helper:
   - envio final direto;
   - envio inline.
4. Manter intacto o comportamento das reacoes customizadas, que ja usa `Text = " "` quando ha `eid:`.
5. Rodar `gofmt` no arquivo alterado.
6. Tentar rodar build/testes Go. Se o ambiente continuar com `go: no such tool "vet"`, registrar a limitacao.

## Riscos
- Se um usuario digitou texto alem do emoji customizado no nome do botao, a remocao precisa preservar esse texto.
- Como entidades do Telegram usam offsets UTF-16, uma remocao por offset pode ser arriscada; a implementacao deve ser conservadora e preferir remover apenas o primeiro rune/placeholder quando houver `CustomEmojiID`.

## Impactos esperados
- Botao com somente emoji customizado passa a mostrar apenas o icone customizado.
- Botao com emoji customizado mais texto continua mostrando o icone customizado e o texto restante.
- Sem impacto no Dashboard, banco, API ou botoes padrao do canal.

## Compatibilidade
- Linux: compativel
- macOS: compativel
- Windows: compativel
- Docker: sem impacto esperado
- CI/CD: depende apenas de build Go

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
1. Abrir o PostBuilder.
2. Adicionar botao enviando uma primeira linha com um emoji customizado e uma segunda linha com URL.
3. Enviar/previewar o post.
4. Confirmar que aparece somente o botao com icone customizado, sem duplicar o botao/texto padrao.

## Rollback
Reverter as alteracoes em `internal/telegram/handlers/events/postBuilder/postBuilder.go`.

## Observacoes
- Nao sera feita alteracao de schema ou migracao.
- Nao sera alterado o formato salvo em Redis.
- Mudanca deve ser pequena e localizada na renderizacao dos botoes do PostBuilder.
