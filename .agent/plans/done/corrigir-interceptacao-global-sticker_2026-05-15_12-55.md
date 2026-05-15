# Plano: corrigir-interceptacao-global-sticker

## Pedido do usuário
O bot está interceptando todos os stickers enviados, tentando configurá-los como "Sticker Separador", o que gera erros de "Sessão Expirada" quando o usuário não está nesse fluxo.

## Objetivo
Garantir que o handler de configuração de sticker separador seja acionado APENAS quando o usuário estiver explicitamente no fluxo de configuração (aguardando o envio do sticker).

## Contexto atual
- No arquivo `internal/telegram/callbacks/loader.go`, a função `matchAwaitingSticker` retorna `true` para qualquer mensagem que contenha um sticker:
  ```go
  func matchAwaitingSticker(update *models.Update) bool {
      return update.Message != nil && update.Message.From != nil && !update.Message.From.IsBot && update.Message.Sticker != nil
  }
  ```
- Isso faz com que o `SetStickerSeparatorHandler` seja executado sempre. Dentro dele, o bot tenta buscar o estado no cache (`GetAwaitingStickerSeparator`) e, ao não encontrar (ou estar expirado), envia a mensagem de erro "⌛ Seção Expirada".

## Arquivos analisados
- `internal/telegram/callbacks/loader.go`
- `internal/telegram/callbacks/my_channel/stickerSeparator.go`

## Arquivos que poderão ser modificados
- `internal/telegram/callbacks/loader.go`
- `internal/telegram/callbacks/my_channel/stickerSeparator.go`

## Estratégia de implementação
1. **Refatorar o Matcher:** Mover a verificação de estado do cache para dentro do `matchAwaitingSticker`. Dessa forma, o matcher só retornará `true` se o usuário tiver um sticker E estiver no estado de espera no Redis.
2. **Limpeza do Handler:** Remover a mensagem de erro de "Sessão Expirada" do `SetStickerSeparatorHandler` no início da função, pois o matcher agora garantirá que a sessão existe.

## Passos detalhados

1. Editar `internal/telegram/callbacks/loader.go`:
    - Injetar o `container.AppContainer` ou acessar o cache para verificar se o usuário está aguardando um sticker no `matchAwaitingSticker`. 
    - *Nota:* Como `matchAwaitingSticker` é usado em `RegisterHandlerMatchFunc`, precisaremos criar uma closure ou garantir acesso ao container.

2. Alternativa mais simples e limpa:
    - Modificar `matchAwaitingSticker` para ser uma factory que recebe o container:
      ```go
      func matchAwaitingSticker(c *container.AppContainer) bot.MatchFunc {
          return func(update *models.Update) bool {
              if update.Message == nil || update.Message.Sticker == nil || update.Message.From == nil {
                  return false
              }
              // Verificar no cache se o usuário está esperando sticker
              channelId, _ := c.CacheService.GetAwaitingStickerSeparator(context.Background(), update.Message.From.ID)
              return channelId != 0
          }
      }
      ```

3. Editar `internal/telegram/callbacks/my_channel/stickerSeparator.go`:
    - Remover o log de erro e a mensagem de "Seção Expirada" no início de `SetStickerSeparatorHandler`, pois se o código chegar ali, o matcher já validou a existência da sessão.

## Riscos
- **Performance:** Uma consulta ao Redis em cada sticker enviado. Como o Redis é extremamente rápido e o volume de stickers privados tende a ser baixo, o risco é mínimo.
- **Contexto:** Usar `context.Background()` no matcher pode ser necessário se o `update` não prover um contexto, mas o ideal é usar o contexto da requisição se disponível.

## Impactos esperados
- Usuários poderão enviar stickers normalmente no chat privado com o bot sem receber mensagens de erro.
- O fluxo de configuração de separador continuará funcionando perfeitamente.

## Como testar

### Execução
1. Abrir o bot.
2. Enviar um sticker aleatório (sem estar no menu de configuração). -> O bot deve ignorar.
3. Ir em: Canais -> Selecionar Canal -> Sticker Separador -> Configurar.
4. Enviar um sticker. -> O bot deve salvar como separador.
