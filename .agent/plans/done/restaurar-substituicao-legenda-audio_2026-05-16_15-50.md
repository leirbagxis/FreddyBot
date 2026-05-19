# Plano: restaurar-substituicao-legenda-audio_2026-05-16_15-50.md

## Pedido do usuário
Para mensagens de áudio (músicas), o bot deve substituir a legenda original pela legenda do bot, ao invés de adicionar embaixo (append).

## Objetivo técnico
Diferenciar o tratamento de legendas no `StageTransform` especificamente para o tipo de mensagem `audio`, aplicando a estratégia de "replace" ao invés de "append".

## Contexto atual
A alteração anterior unificou o comportamento de todas as mídias para "append". Precisamos agora abrir uma exceção para o tipo `audio`.

## Arquivos analisados
- `internal/telegram/events/channelPost/stage_transform.go`
- `internal/telegram/events/channelPost/types.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/stage_transform.go`

## Estratégia de implementação
No `StageTransform`, dentro do bloco de montagem final (Final Assembly):
1. Verificar se `pCtx.MessageType == MessageTypeAudio`.
2. Se for áudio, se houver uma legenda no banco (`dbCaption`), atribuir diretamente a `pCtx.FormattedText = dbCaption`.
3. Para os demais tipos de mídia, manter a lógica de `append`.

## Passos detalhados

1. Abrir `internal/telegram/events/channelPost/stage_transform.go`.
2. Modificar o bloco de decisão da etapa 5 (Final Assembly).
3. Adicionar a condição específica para `MessageTypeAudio`.

## Riscos
- Nenhum risco técnico identificado. É uma mudança de regra de negócio simples e isolada.

## Impactos esperados
- Mensagens de áudio/música ignorarão qualquer legenda original do arquivo e exibirão apenas a legenda configurada no bot.
- Fotos e vídeos continuarão com o comportamento de preservação da legenda original (append).

## Como testar

### Build
```bash
go build ./cmd/FreddyBot
```

### Execução/Teste
1. Enviar um arquivo de áudio (MP3) para o canal com uma legenda qualquer.
2. O bot deve editar a mensagem e deixar **apenas** a legenda padrão do bot.
3. Enviar uma foto com legenda. O bot deve manter a legenda original e adicionar a do bot embaixo.

## Rollback
Reverter a alteração no `stage_transform.go` para a lógica anterior.
