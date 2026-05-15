# Plano: corrigir-saida-sticker-postbuilder

## Pedido do usuário
O bot tenta inserir texto (legenda) ao enviar um sticker no PostBuilder, mas stickers não aceitam legendas no Telegram. O suporte para stickers deve se limitar apenas a botões e reações.

## Objetivo
Ajustar a saída do PostBuilder (Preview e Envio Inline) para tratar stickers corretamente, removendo a tentativa de enviar legendas e garantindo que apenas o sticker e o teclado (botões/reações) sejam enviados.

## Contexto atual
- `sendFinalPost` em `internal/telegram/events/postBuilder/postBuilder.go` cai no caso `default` para stickers, tentando enviar `SendMessage` com o texto da legenda.
- `InlineHandler` também cai no `default` (`InlineQueryResultArticle`), o que transforma o sticker em uma mensagem de texto simples ao ser enviado via modo inline.

## Arquivos analisados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. **Atualizar `sendFinalPost`:** Adicionar o caso `case "sticker"` que utiliza `b.SendSticker` sem o campo `Caption`.
2. **Atualizar `InlineHandler`:** Adicionar o caso `case "sticker"` que utiliza `models.InlineQueryResultCachedSticker`.

## Passos detalhados

1. Editar `internal/telegram/events/postBuilder/postBuilder.go`:
    - Na função `sendFinalPost`, adicionar:
      ```go
      case "sticker":
          _, err := b.SendSticker(ctx, &bot.SendStickerParams{
              ChatID:      chatID,
              Sticker:     &models.InputFileString{Data: state.MediaFileID},
              ReplyMarkup: kb,
          })
          // ... tratamento de erro
      ```
    - Na função `InlineHandler`, adicionar:
      ```go
      case "sticker":
          result = &models.InlineQueryResultCachedSticker{
              ID:           id,
              StickerFileID: state.MediaFileID,
              ReplyMarkup:  kb,
          }
      ```

## Riscos
- Nenhum risco técnico identificado, apenas a limitação funcional (esperada) de que stickers não terão texto acompanhando-os.

## Impactos esperados
- Stickers serão enviados corretamente como stickers (animados ou estáticos).
- Botões e reações aparecerão abaixo do sticker enviado.
- Mensagens de aviso (WARN) no log sobre mídias não reconhecidas desaparecerão para stickers.

## Como testar

### Execução
1. Enviar um sticker para o bot.
2. Iniciar o PostBuilder.
3. Adicionar botões e reações.
4. Clicar em "👁️ Preview". -> O bot deve enviar o sticker com os botões.
5. Salvar e enviar via modo inline. -> O resultado deve ser o sticker com os botões.
