# Plano: melhorar newpack dashboard

## Pedido do usuário
Adicionar suporte ao comando `/newpack` além de `!newpack`; na dashboard do usuário, abaixo do campo de legenda do new pack, adicionar dois botões/toggles para controlar se a mensagem do bot terá botão e se o sticker enviado terá botão. Também adicionar variável de quantidade de stickers no pack e corrigir o uso de Markdown com `$link`, por exemplo `[click here to add]($link)`.

## Objetivo
Melhorar o fluxo de New Pack para aceitar dois comandos, permitir configurar onde os botões do pack aparecem, expor a quantidade de stickers como variável de template e garantir que links Markdown com variáveis sejam renderizados corretamente.

## Contexto atual
O handler `newpack.go` aceita apenas `!newpack` via regex `^!newpack`. Quando recebe o sticker, ele busca o sticker set, renderiza `channel.NewPackCaption` com `$titulo`, `$name` e `$link`, edita a mensagem original e adiciona o mesmo botão tanto na mensagem editada quanto no sticker enviado.

A dashboard do usuário renderiza `NewPackCaptionCard` com apenas a legenda e salva via `updateNewPackCaption`. O cliente frontend envia `{ newPackCaption }`, mas o controller usa `CaptionDefaultUpdateRequest` com JSON `caption`, então vale corrigir o endpoint para aceitar os dois formatos enquanto adiciona os novos campos.

O modelo `Channel` ainda não possui flags para controlar botões do New Pack. Como o comportamento atual sempre adiciona botões nos dois lugares, as novas flags devem ter default `true` para manter compatibilidade em bancos existentes e novos canais.

## Arquivos analisados
- internal/telegram/events/channelPost/newpack.go
- internal/telegram/events/channelPost/utils_v2.go
- internal/database/models/models.go
- internal/database/repositories/channel.go
- internal/core/services/captions.go
- internal/api/types/captions.go
- internal/api/controllers/captionController.go
- internal/api/dto/dto.go
- internal/api/dto/mapper.go
- dashboard/src/types.ts
- dashboard/src/api.ts
- dashboard/src/App.tsx
- dashboard/src/components/NewPackCaptionCard.tsx

## Arquivos que poderão ser modificados
- internal/database/models/models.go
- internal/database/repositories/channel.go
- internal/core/services/channels.go
- internal/core/services/captions.go
- internal/api/types/captions.go
- internal/api/controllers/captionController.go
- internal/api/dto/dto.go
- internal/api/dto/mapper.go
- internal/telegram/events/channelPost/newpack.go
- dashboard/src/types.ts
- dashboard/src/api.ts
- dashboard/src/App.tsx
- dashboard/src/components/NewPackCaptionCard.tsx

## Estratégia de implementação
Adicionar dois campos booleanos no canal: um para botão na mensagem de new pack e outro para botão no sticker recebido. Os campos terão default `true` para preservar o comportamento atual. Atualizar DTO/API/frontend para carregar e salvar esses valores junto da legenda.

No handler do bot, aceitar `!newpack` e `/newpack`, calcular a quantidade de stickers do pack e renderizar novas variáveis. O template passará a aceitar `$count`, `$total`, `$stickers`, `$title`, `$titulo`, `$name` e `$link`. A renderização deve substituir variáveis antes de converter Markdown para HTML, mantendo `[texto]($link)` funcional.

## Passos detalhados

1. Adicionar `NewPackMessageButtons` e `NewPackStickerButtons` ao modelo `Channel` com default `true`.
2. Atualizar `CreateChannelWithDefaults` para inicializar os dois campos como `true`.
3. Atualizar DTO e mapper para expor os dois campos na API.
4. Criar/ajustar request de atualização de New Pack para aceitar `caption`, `newPackCaption`, `newPackMessageButtons` e `newPackStickerButtons`.
5. Atualizar service e repository para salvar legenda e flags em uma única operação, invalidando cache do canal.
6. Atualizar o frontend `Channel` type com os novos campos.
7. Atualizar `updateNewPackCaption` para enviar legenda e flags.
8. Atualizar `App.tsx` para passar flags ao card e refletir o estado salvo localmente.
9. Atualizar `NewPackCaptionCard` para mostrar dois toggles abaixo do editor da legenda, com cópias curtas e claras.
10. Atualizar dica de variáveis do card para incluir nome, link e quantidade de stickers.
11. Atualizar `newpack.go` para aceitar `!newpack` e `/newpack`.
12. Atualizar `renderNewPackTemplate` para aceitar quantidade de stickers e substituir as novas variáveis.
13. Aplicar flags no envio dos botões: mensagem editada só recebe botão se `NewPackMessageButtons` estiver ativo; sticker só recebe botão se `NewPackStickerButtons` estiver ativo.
14. Rodar `gofmt` nos arquivos Go alterados.
15. Rodar build da dashboard.
16. Rodar `git diff --check`, build Go e testes Go quando possível.

## Riscos
- Médio risco por tocar modelo, API, dashboard e handler do bot.
- Bancos existentes terão novas colunas adicionadas pelo AutoMigrate; é preciso garantir default `true` para preservar comportamento.
- Se o frontend antigo enviar apenas `newPackCaption`, o backend deve continuar aceitando.
- O Markdown ainda depende do conversor local `DetectParseMode`, então a correção deve focar em substituir variáveis antes da conversão e evitar quebrar o HTML gerado.

## Impactos esperados
- `/newpack` passa a funcionar junto de `!newpack`.
- Usuário consegue escolher pela dashboard se a mensagem do bot e/ou o sticker terão o botão do pack.
- Templates podem usar a quantidade de stickers no pack.
- `[click here to add]($link)` deve virar link clicável para o pack.
- Com defaults `true`, canais antigos continuam com comportamento equivalente ao atual, salvo quando o usuário desativar algum toggle.

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
npm run build
```

### Testes
```bash
go test ./...
```

### Execução
```bash
go run ./cmd/FreddyBot/main.go
```

Teste manual esperado:
1. Abrir dashboard do usuário na aba de legendas.
2. Editar New Pack Caption para `[click here to add]($link) - $count stickers`.
3. Alternar os dois toggles e salvar.
4. No canal, testar `!newpack` e `/newpack` com um sticker de pack público.
5. Conferir se o botão aparece apenas nos lugares selecionados e se o link Markdown aponta para o pack.

## Rollback
Reverter os arquivos alterados. Se as colunas novas já tiverem sido criadas no banco, elas podem permanecer sem uso ou ser removidas manualmente em uma migration controlada se necessário.

## Observações
A quantidade de stickers virá do retorno de `GetStickerSet`. Canais existentes devem preservar o comportamento atual com botões ativados por padrão.
