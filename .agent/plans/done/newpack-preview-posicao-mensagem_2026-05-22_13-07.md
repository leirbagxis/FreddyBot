# Plano: newpack preview posicao mensagem

## Pedido do usuário
Ajustar o New Pack para respeitar a permissão de link preview nas mensagens e adicionar um novo toggle na dashboard do usuário para escolher se a mensagem do New Pack fica acima ou abaixo do sticker. Se ficar acima, o bot edita a mensagem "esperando sticker". Se ficar abaixo, o bot envia a mensagem do New Pack abaixo do sticker enviado e exclui a mensagem de espera.

## Objetivo
Controlar corretamente o preview de links conforme as permissões de mensagem e adicionar configuração persistente de posição da mensagem do New Pack.

## Contexto atual
O New Pack edita a mensagem de espera usando `EditMessageTextParams` e já calcula `disableLP` a partir de `perms.CanUseLinkPreview`. Porém o fluxo só existe no modo "acima". O novo modo "abaixo" precisa usar `SendMessageParams`, aplicar a mesma regra de `LinkPreviewOptions`, opcionalmente anexar botão nessa nova mensagem e apagar a mensagem de espera.

A dashboard do usuário já salva `newPackMessageButtons` e `newPackStickerButtons`. O modelo `Channel`, DTO, request e repository precisarão receber mais uma configuração persistente para posição da mensagem.

## Arquivos analisados
- internal/telegram/events/channelPost/newpack.go
- internal/telegram/events/channelPost/permissions.go
- internal/database/models/models.go
- internal/database/repositories/channel.go
- internal/api/types/captions.go
- internal/api/dto/dto.go
- dashboard/src/components/NewPackCaptionCard.tsx
- dashboard/src/App.tsx

## Arquivos que poderão ser modificados
- internal/database/models/models.go
- internal/database/repositories/channel.go
- internal/core/services/channels.go
- internal/core/services/captions.go
- internal/api/types/captions.go
- internal/api/dto/dto.go
- internal/api/dto/mapper.go
- internal/telegram/events/channelPost/newpack.go
- dashboard/src/types.ts
- dashboard/src/api.ts
- dashboard/src/App.tsx
- dashboard/src/components/NewPackCaptionCard.tsx
- dashboard/src/mockData.ts

## Estratégia de implementação
Adicionar um campo `NewPackMessagePosition` com valores `above` e `below`, default `above`. No frontend, exibir um controle binário abaixo dos toggles de botão. No backend, salvar esse campo junto com a legenda e toggles.

No handler Telegram:
- `above`: comportamento atual, editar a mensagem de espera e aplicar `LinkPreviewOptions` conforme permissão.
- `below`: enviar uma nova mensagem após o sticker, aplicar `LinkPreviewOptions` conforme permissão, anexar botão se `newPackMessageButtons` estiver ativo, e deletar a mensagem de espera. O botão do sticker continua controlado por `newPackStickerButtons`.

## Passos detalhados

1. Adicionar `NewPackMessagePosition *string` ao modelo `Channel` com default `above`.
2. Inicializar canais novos com posição `above`.
3. Expor `newPackMessagePosition` no DTO e mapper, usando default `above` quando nulo/vazio.
4. Adicionar o campo ao request `NewPackCaptionUpdateRequest`.
5. Atualizar repository/service para persistir posição quando enviada, validando apenas `above` ou `below`.
6. Atualizar tipos e API da dashboard.
7. Atualizar `NewPackCaptionCard` para receber/salvar posição e renderizar o controle "Mensagem acima/abaixo".
8. Atualizar `App.tsx` e `mockData.ts` para o novo campo.
9. Atualizar `newpack.go` para bifurcar envio por posição.
10. Garantir `LinkPreviewOptions` tanto em `EditMessageTextParams` quanto em `SendMessageParams` quando a permissão desativar preview.
11. Adicionar logs do modo escolhido e do link preview habilitado/desabilitado.
12. Rodar `gofmt` nos arquivos Go.
13. Rodar `npm run build`, `git diff --check`, e tentar build/testes Go.

## Riscos
- Médio risco por tocar modelo, API, frontend e fluxo Telegram.
- Apagar a mensagem de espera no modo abaixo pode falhar se o bot não puder deletar aquela mensagem; nesse caso deve logar erro e continuar.
- Em bancos/caches antigos, campo ausente deve cair para `above` para preservar comportamento atual.

## Impactos esperados
- Link preview passa a obedecer a permissão nos dois modos.
- Usuário consegue escolher se a mensagem do New Pack aparece acima/editada ou abaixo do sticker.
- Com default `above`, canais existentes mantêm o comportamento atual até o usuário mudar.

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

Teste manual:
1. Na dashboard, desativar link preview nas permissões de mensagem.
2. Configurar New Pack com link embutido.
3. Testar posição "acima": `/newpack` deve editar a mensagem de espera sem preview se permissão estiver off.
4. Testar posição "abaixo": `/newpack` deve apagar a mensagem de espera e enviar a mensagem abaixo do sticker, também sem preview se permissão estiver off.

## Rollback
Reverter os arquivos alterados. A coluna nova pode permanecer sem uso ou ser removida manualmente em migration controlada.

## Observações
O botão da mensagem no modo "abaixo" será anexado à nova mensagem enviada. O botão do sticker continua independente.
