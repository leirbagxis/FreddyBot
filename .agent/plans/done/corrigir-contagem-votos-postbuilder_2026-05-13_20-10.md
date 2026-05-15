# Plano: corrigir-contagem-votos-postbuilder_2026-05-13_20-10.md

## Pedido do usuário
A contagem de votos do PostBuilder não está aparecendo. O usuário quer que a contagem apareça apenas a partir de 1 voto (ex: "👍 1"), sem mostrar o "0".

## Objetivo
Garantir que os botões de reação do PostBuilder sejam atualizados corretamente com a contagem de votos assim que o primeiro voto for computado.

## Contexto atual
- O código de `vote.go` já tenta mostrar a contagem apenas se `count > 0`.
- O problema relatado sugere que a atualização (`EditMessageReplyMarkup`) não está ocorrendo ou a contagem está vindo zerada/incorreta.
- Mensagens inline dependem do `pb_inline_map` para serem atualizadas.

## Arquivos analisados
- `internal/telegram/callbacks/vote/vote.go`
- `internal/telegram/events/postBuilder/postBuilder.go`
- `cmd/FreddyBot/main.go`

## Arquivos que poderão ser modificados
- `internal/telegram/callbacks/vote/vote.go`
- `cmd/FreddyBot/main.go` (se faltar configuração de updates)

## Estratégia de implementação
1.  **Melhorar Logs e Verificação em `vote.go`**: Adicionar logs detalhados para entender por que o teclado não está sendo editado em mensagens inline.
2.  **Verificar Configuração do Bot**: Garantir que o bot está configurado para receber `chosen_inline_result`, que é essencial para o mapeamento de mensagens inline.
3.  **Refinar a lógica de reconstrução**: Garantir que o teclado reconstruído para mensagens inline use as contagens atuais do banco de dados já na primeira renderização após o mapeamento.

## Passos detalhados

1.  **Verificar `cmd/FreddyBot/main.go`**
    - Garantir que `WithAllowedUpdates` inclua `chosen_inline_result`.

2.  **Refinar `internal/telegram/callbacks/vote/vote.go`**
    - Ajustar a reconstrução do teclado inline para que, ao montar os botões de reação, ele já consulte o mapa de `counts` obtido do serviço de votos.
    - Atualmente, o código reconstrói o teclado com os emojis limpos e *depois* tenta iterar para atualizar. Vou simplificar para que a reconstrução já use os dados atualizados.

## Riscos
- **Baixo:** A lógica de exibição (apenas > 0) será preservada conforme solicitado.

## Impactos esperados
- Votos em mensagens do PostBuilder (inline ou preview) passarão a exibir a contagem imediatamente após o primeiro voto.

## Compatibilidade
- Go-telegram/bot v1.19.0

## Como testar
1. Iniciar o bot e criar um post no PostBuilder.
2. Usar o Preview para votar -> verificar se aparece "emoji 1".
3. Compartilhar via modo inline -> votar -> verificar se aparece "emoji 1".

## Rollback
`git checkout internal/telegram/callbacks/vote/vote.go cmd/FreddyBot/main.go`
