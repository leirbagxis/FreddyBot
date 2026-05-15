# Plano: corrigir-botoes-reacao-postbuilder_2026-05-13_11-50.md

## Pedido do usuário
Os botões de voto/reação não funcionam corretamente quando configurados via Post Builder, embora funcionem em canais.

## Objetivo
Corrigir a funcionalidade de votos no Post Builder, garantindo que os contadores atualizem visualmente mesmo em mensagens compartilhadas via modo inline.

## Contexto atual
O bot não conseguia processar votos de mensagens inline porque o `inline_message_id` não era tratado adequadamente e o teclado original não estava disponível no callback.

## Arquivos analisados
- `internal/telegram/events/postBuilder/postBuilder.go`
- `internal/telegram/callbacks/vote/vote.go`
- `internal/telegram/events/loader.go`

## Arquivos modificados
- `internal/telegram/events/postBuilder/postBuilder.go`
- `internal/telegram/callbacks/vote/vote.go`
- `internal/telegram/events/loader.go`

## Estratégia de implementação
Implementar uma "ponte" de mapeamento usando o evento `ChosenInlineResult`. Ao enviar uma mensagem inline, o bot vincula o ID da mensagem à sessão no Redis. No callback do voto, o bot recupera essa sessão e reconstrói o teclado para atualização.

## Passos detalhados
1. Adicionar `ChosenInlineResultHandler` em `postBuilder.go`.
2. Registrar o handler em `events/loader.go`.
3. Atualizar `vote/vote.go` para suportar `inline_message_id` e reconstrução de teclado via cache.
4. Garantir que o `CallbackData` permaneça no formato `vote:emoji`.

## Riscos
- Dependência de ativação do Inline Feedback no BotFather.
- Expiração do cache no Redis pode tornar mensagens antigas ( > 24h) não atualizáveis visualmente.

## Impactos esperados
- Votos funcionando em mensagens do Post Builder.
- Feedback visual correto para o usuário.

## Compatibilidade
- Linux, Docker, Redis.

## Como testar
1. Criar post no Post Builder com emojis.
2. Compartilhar via modo inline.
3. Votar e verificar se o contador no botão atualiza.

## Rollback
Reverter commits/arquivos para o estado anterior e remover o handler de `ChosenInlineResult`.

## Observações
Implementado com sucesso durante a sessão de 13/05/2026.
