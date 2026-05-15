# Plano: habilitar-sticker-postbuilder

## Pedido do usuário
O PostBuilder não é ativado ao enviar um sticker. O objetivo é permitir que stickers também iniciem o fluxo de criação, permitindo adicionar botões e reações (mesmo que não suportem título/corpo).

## Objetivo
Adicionar suporte a stickers como uma mídia válida no PostBuilder, permitindo que usuários criem postagens interativas a partir de um sticker.

## Contexto atual
- `matchPostBuilder` em `internal/telegram/events/loader.go` não inclui `update.Message.Sticker` na verificação de mídias.
- `Handler` em `internal/telegram/events/postBuilder/postBuilder.go` não possui a lógica para extrair o `FileID` de um sticker.

## Arquivos analisados
- `internal/telegram/events/loader.go`
- `internal/telegram/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/loader.go`
- `internal/telegram/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. **Atualizar o Matcher:** Incluir `update.Message.Sticker != nil` na lógica de detecção de mídias no arquivo `loader.go`.
2. **Atualizar o Handler:** Adicionar a detecção de stickers no `postBuilder.go`, capturando o `FileID` e definindo o `mediaType` como "sticker".

## Passos detalhados

1. Editar `internal/telegram/events/loader.go`:
    - Modificar `matchPostBuilder` para retornar `true` se `update.Message.Sticker != nil`.

2. Editar `internal/telegram/events/postBuilder/postBuilder.go`:
    - No `Handler`, adicionar um `else if update.Message.Sticker != nil` para definir `mediaID` e `mediaType = "sticker"`.

## Riscos
- **Conflito com Sticker Separador:** Como corrigimos anteriormente o `matchAwaitingSticker` para checar o Redis primeiro, o PostBuilder não deve interferir no fluxo de configuração de canal. No entanto, se o usuário enviar um sticker sem estar no fluxo de "Set Separador", o PostBuilder agora vai oferecer a criação de post. Isso é o comportamento esperado.

## Impactos esperados
- Ao enviar um sticker para o bot, ele oferecerá o botão "🛠️ Post Builder".
- Usuários poderão adicionar botões de link e reações a stickers.

## Como testar

### Execução
1. Abrir o bot.
2. Enviar um sticker aleatório.
3. Clicar em "🛠️ Post Builder".
4. Adicionar um botão de link.
5. Ver o Preview ou salvar/enviar inline.
