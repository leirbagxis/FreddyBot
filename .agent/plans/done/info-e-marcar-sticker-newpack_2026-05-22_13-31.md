# Plano: info e marcar sticker newpack

## Pedido do usuário
Além de esconder a mensagem de ajuda das variáveis atrás do ícone de informação, adicionar mais uma opção: quando o usuário selecionar `Mensagem abaixo`, aparecer um botão `Marcar Sticker`. Se ativo, ao enviar a legenda New Pack abaixo do sticker, o bot deve marcar/responder o sticker enviado.

## Objetivo
Melhorar a UX do card de New Pack e permitir que a mensagem enviada abaixo do sticker seja enviada como resposta ao sticker quando configurado.

## Contexto atual
Existe um plano anterior pendente apenas para o ícone de ajuda: `.agent/plans/pending/info-variaveis-newpack_2026-05-22_13-26.md`. Este novo plano substitui aquele escopo, adicionando também a funcionalidade `Marcar Sticker`.

`NewPackCaptionCard` mostra a ajuda de variáveis sempre visível e já possui controles de botão na mensagem, botão no sticker e posição `Mensagem acima/abaixo`.

No backend, `newpack.go` já possui o modo `below`, que usa `SendMessageParams` para enviar a mensagem abaixo do sticker e depois tenta apagar a mensagem de espera. O Telego suporta `ReplyParameters` em `SendMessageParams`, com `MessageID` e `AllowSendingWithoutReply`, então marcar o sticker pode ser feito configurando `ReplyParameters.MessageID = post.MessageID` quando o novo toggle estiver ativo.

## Arquivos analisados
- dashboard/src/components/NewPackCaptionCard.tsx
- internal/telegram/events/channelPost/newpack.go
- internal/database/models/models.go
- internal/api/types/captions.go
- internal/database/repositories/channel.go
- internal/core/services/captions.go
- /home/malbs/go/pkg/mod/github.com/mymmrac/telego@v1.9.0/types.go
- /home/malbs/go/pkg/mod/github.com/mymmrac/telego@v1.9.0/methods.go

## Arquivos que poderão ser modificados
- dashboard/src/components/NewPackCaptionCard.tsx
- dashboard/src/types.ts
- dashboard/src/api.ts
- dashboard/src/App.tsx
- dashboard/src/mockData.ts
- internal/database/models/models.go
- internal/api/dto/dto.go
- internal/api/dto/mapper.go
- internal/api/types/captions.go
- internal/database/repositories/channel.go
- internal/core/services/channels.go
- internal/core/services/captions.go
- internal/telegram/events/channelPost/newpack.go

## Estratégia de implementação
Adicionar um novo campo persistido `newPackReplyToSticker`, com default `false`, porque é uma funcionalidade nova e não deve alterar comportamento de canais existentes.

Na dashboard:
- Adicionar estado/prop `replyToSticker`.
- Mostrar o controle `Marcar Sticker` apenas quando `messagePosition === 'below'`.
- Esconder a mensagem de variáveis por padrão e mostrar ao clicar no ícone `Info`.

No backend:
- Expor/salvar `newPackReplyToSticker` na API.
- No modo `below`, se ativo, preencher `sendParams.ReplyParameters` apontando para `post.MessageID`.
- Manter o comportamento `above` sem usar reply, porque esse modo edita a mensagem de espera.

## Passos detalhados

1. Adicionar `NewPackReplyToSticker *bool` ao modelo `Channel` com default `false`.
2. Inicializar canais novos com `false`.
3. Expor `newPackReplyToSticker` no DTO/mapper com default `false` quando ausente.
4. Adicionar `NewPackReplyToSticker *bool` ao request de atualização de New Pack.
5. Atualizar repository/service para salvar esse campo quando enviado.
6. Atualizar tipos e API frontend para carregar/enviar `newPackReplyToSticker`.
7. Atualizar `App.tsx` para passar o valor ao card e atualizar estado local após salvar.
8. Atualizar `NewPackCaptionCard`:
   - adicionar `showHelp` para o ícone de informação;
   - renderizar a mensagem de variáveis só quando `showHelp` estiver ativo;
   - adicionar botão/toggle `Marcar Sticker` apenas quando posição for `below`.
9. Atualizar `newpack.go` para aplicar `ReplyParameters` no `SendMessageParams` do modo `below` quando `newPackReplyToSticker` estiver ativo.
10. Adicionar logs indicando se a mensagem abaixo marcou o sticker.
11. Rodar `gofmt` nos arquivos Go.
12. Rodar `npm run build`.
13. Rodar `git diff --check`, `go build` e `go test` quando possível.

## Riscos
- Médio risco por tocar modelo, API, dashboard e handler Telegram.
- Se o Telegram não encontrar o sticker para resposta, `AllowSendingWithoutReply` deve permitir enviar mesmo assim.
- Como `Marcar Sticker` só aparece no modo abaixo, não deve afetar o fluxo acima.

## Impactos esperados
- A ajuda de variáveis aparece somente ao clicar no ícone de info.
- Ao selecionar `Mensagem abaixo`, aparece o toggle `Marcar Sticker`.
- Se `Marcar Sticker` estiver ativo, a mensagem New Pack abaixo será enviada como resposta ao sticker.
- Canais antigos mantêm comportamento anterior com `newPackReplyToSticker=false`.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
npm run build
go build ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./...
```

### Execução
```bash
npm run dev
go run ./cmd/FreddyBot/main.go
```

Teste manual:
1. Abrir dashboard na aba de legendas e editar New Pack.
2. Clicar no ícone de informação e confirmar que a ajuda aparece/some.
3. Selecionar `Mensagem abaixo` e confirmar que `Marcar Sticker` aparece.
4. Ativar `Marcar Sticker`, salvar e testar `/newpack` com sticker de pack público.
5. Confirmar que a mensagem enviada abaixo aparece respondendo/marcando o sticker.

## Rollback
Reverter os arquivos alterados. A coluna nova pode permanecer sem uso ou ser removida manualmente em migration controlada.

## Observações
O plano anterior `info-variaveis-newpack_2026-05-22_13-26.md` deve ser considerado substituído por este plano e não precisa ser implementado separadamente.
