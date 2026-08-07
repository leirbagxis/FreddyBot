# Plano: Remover Bloqueio de Middleware em Encaminhamentos e Ativar PostBuilder Direto

## Causa Raiz Diagnosticada pelos Logs
Nos logs fornecidos:
```
2026/08/07 13:06:30 [BOT] Predicate: matchForwardedChannelTelego = true para UserID=7595607953
```
- A mensagem foi capturada pelo predicado `matchForwardedChannelTelego()`.
- O grupo `forwardedGroup` aplicava o middleware `middleware.CheckAddBotMiddlewareTelego(c)`.
- Como o bot **não era admin** no canal de origem da mensagem encaminhada, a chamada `b.GetChatMember` dentro do middleware retornava erro e interrompia a execução (`return nil`), impedindo que `AskAddChannel` e o `PostBuilder` fossem executados. O middleware simplesmente engolia o update!

## Objetivo
- Remover a interceptação e o bloqueio do `forwardedGroup` em `loader_telego.go`.
- Permitir que mensagens encaminhadas de **qualquer** canal (seja o bot admin ou não, cadastrado ou não) fluam diretamente para `postbuilder.HandlerTelego`, extraindo mídias, legenda, emojis premium e botões.

## Arquivos analisados
- `internal/telegram/loader_telego.go`
- `internal/middleware/checkAddBotMiddlewareTelego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/loader_telego.go`

## Estratégia de implementação

1. **Remover a regra `forwardedGroup` em `loader_telego.go`**:
   - Remover as linhas:
     ```go
     forwardedGroup := bh.Group(matchForwardedChannelTelego())
     forwardedGroup.Use(middleware.CheckAddBotMiddlewareTelego(c))
     forwardedGroup.Handle(addchannel.AskAddChannelHandlerTelego(c))
     ```
   - O cadastro de canais continua funcionando 100% via `AnyMyChatMember()` (quando o bot é adicionado como admin em um canal) e via comando `/add`.

2. **Garantir que mensagens encaminhadas fluam diretamente para o PostBuilder**:
   - Com a remoção do middleware bloqueador, toda mensagem encaminhada no chat privado será avaliada por `matchPostBuilderTelego` e processada por `postbuilder.HandlerTelego`.

---

## Passos detalhados

### Passo 1: Atualizar `loader_telego.go`
- Remover o grupo `forwardedGroup` e o predicado `matchForwardedChannelTelego()`.

### Passo 2: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. O fluxo de adição de canais via promoção de admin (`AnyMyChatMember`) permanece intocado.

## Impactos esperados
- Correção definitiva: qualquer mensagem/mídia encaminhada no privado (de qualquer canal, onde o bot é ou não admin) ativa instantaneamente o PostBuilder.
- Fim de mensagens descartadas silenciosamente no middleware.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Encaminhar mensagem/mídia de qualquer canal no chat privado com o bot.
2. Observar o log: `Predicate: matchPostBuilderTelego = true (conteudo/midia)` e `PostBuilder: Processando mensagem recebida...`.
3. Verificar a resposta imediata do bot com `"✨ Conteúdo detectado! Deseja usar o Post Builder..."` e o botão `"🛠️ Post Builder"`.
