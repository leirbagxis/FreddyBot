# Plano: ajustar-montagem-legenda-midia_2026-05-16_15-30.md

## Pedido do usuário
Evitar que a legenda original da mídia seja substituída pela legenda padrão.

## Objetivo técnico
Alterar a estratégia de montagem de legendas de mídia de "replace" para "append".

## Contexto atual
O bot estava configurado para ignorar a legenda original em fotos/vídeos se houvesse uma legenda padrão definida.

## Arquivos analisados
- `internal/telegram/events/channelPost/stage_transform.go`

## Arquivos modificados
- `internal/telegram/events/channelPost/stage_transform.go`

## Estratégia de implementação
Utilizar a função `composeMessage` com a estratégia `append` também para o tipo de mensagem mídia, unindo o conteúdo original com a legenda padrão do banco de dados.

## Impactos esperados
- Preservação do texto enviado pelo usuário.
- Legenda do bot adicionada como rodapé.
