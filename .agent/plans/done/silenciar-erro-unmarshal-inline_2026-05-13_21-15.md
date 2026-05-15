# Plano: silenciar-erro-unmarshal-inline_2026-05-13_21-15.md

## Pedido do usuário
O bot está funcionando e editando o teclado, mas registra um erro de log: `cannot unmarshal bool into Go value of type models.Message`.

## Objetivo
Silenciar o erro de unmarshal que ocorre ao editar o teclado de mensagens inline, pois a API do Telegram retorna um booleano `true` (sucesso) em vez de um objeto `Message` para esse tipo de mensagem, o que causa confusão no log embora a operação seja bem-sucedida.

## Contexto atual
- A biblioteca `go-telegram/bot` tenta sempre decodificar a resposta de `EditMessageReplyMarkup` como `models.Message`.
- Para mensagens inline, o Telegram retorna `true`, resultando no erro de unmarshal.
- O usuário confirma que a edição do teclado está funcionando corretamente.

## Arquivos analisados
- `internal/telegram/callbacks/vote/vote.go`

## Arquivos que poderão ser modificados
- `internal/telegram/callbacks/vote/vote.go`

## Estratégia de implementação
1.  **Ajustar Tratamento de Erro**: No `vote.Handler`, modificar a verificação de erro após `b.EditMessageReplyMarkup`.
2.  **Filtrar Erro Específico**: Se a mensagem for inline e o erro contiver `cannot unmarshal bool`, ignorar o log de erro, pois indica sucesso na operação.

## Passos detalhados

1.  **Modificar `internal/telegram/callbacks/vote/vote.go`**
    - Localizar a chamada `b.EditMessageReplyMarkup`.
    - Atualizar a condicional de erro para:
      ```go
      if err != nil {
          errStr := err.Error()
          isNotModified := strings.Contains(errStr, "message is not modified")
          isUnmarshalBool := strings.Contains(errStr, "cannot unmarshal bool")

          if !isNotModified && !isUnmarshalBool {
              logger.Error("VOTE", "Erro ao editar teclado: %v", err)
          }
      }
      ```

## Riscos
- **Baixo**: Apenas evita logs de "falso erro" que ocorrem em operações bem-sucedidas de mensagens inline.

## Como testar
1. Reiniciar o bot.
2. Votar em uma mensagem inline.
3. Verificar que o teclado atualiza e o erro de unmarshal não aparece mais no log.

## Rollback
`git checkout internal/telegram/callbacks/vote/vote.go`
