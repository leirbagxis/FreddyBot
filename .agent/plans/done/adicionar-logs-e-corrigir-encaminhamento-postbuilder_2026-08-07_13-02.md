# Plano: Adicionar Logs Detalhados e Corrigir Interceptação de Mensagens Encaminhadas de Canais no PostBuilder

## Pedido do usuário
1. O usuário relata que ao enviar/encaminhar uma mensagem citando um canal (no qual o bot não é admin e não está configurado), o bot continua não detectando como mídia/post para o PostBuilder.
2. Adicionar logs para rastrear o fluxo exato de cada mensagem recebida.

## Causa Raiz Identificada
- Em `loader_telego.go`, a regra `forwardedGroup := bh.Group(matchForwardedChannelTelego())` é registrada **antes** do handler do PostBuilder.
- Qualquer mensagem encaminhada de um canal acionava a regra `matchForwardedChannelTelego()`, que redirecionava para `addchannel.AskAddChannelHandlerTelego` (solicitação de adição de canal).
- Como o framework Telego interrompe o processamento após o primeiro handler correspondente, o handler `postbuilder.HandlerTelego` **nunca era executado** para mensagens encaminhadas de canais.

## Objetivo
1. Adicionar logs estruturados (`logger.Bot` e `logger.Info`) no recebimento de mensagens, identificando `ChatID`, `UserID`, `ForwardOrigin`, `MediaType` e decisão de roteamento.
2. Garantir que mensagens encaminhadas de canais (mesmo onde o bot não é admin e não está configurado) ativem o PostBuilder para extração de mídias, textos, emojis e botões.

## Arquivos analisados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Arquivos que poderão ser modificados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Estratégia de implementação

1. **Adicionar Logs de Rastreamento**:
   - Em `loader_telego.go` (nos predicados `matchForwardedChannelTelego` e `matchPostBuilderTelego`), registrar no log o recebimento da atualização e a decisão de cada filtro.
   - Em `postBuilder.go` (`HandlerTelego`), registrar no log a recepção da mensagem, o tipo de mídia/texto identificado e o estado resultante.
   - Em `addChannel.go` (`AskAddChannelHandlerTelego`), registrar no log quando uma mensagem encaminhada for processada para convite de canal.

2. **Permitir que Mensagens Encaminhadas de Canais Acionem o PostBuilder**:
   - Ajustar `AskAddChannelHandlerTelego` em `addChannel.go` para, após tratar o fluxo de convite do canal (ou se o canal não for adicionado), permitir que a mensagem encaminhada também prossiga/dispare o PostBuilder para a criação da postagem.

---

## Passos detalhados

### Passo 1: Adicionar logs em `loader_telego.go` e `postBuilder.go`
- Incluir mensagens de log descritivas informando ID do usuário, origem do encaminhamento e se a mensagem foi aceita para o PostBuilder.

### Passo 2: Ajustar integração entre `AskAddChannelHandlerTelego` e `PostBuilder`
- Em `AskAddChannelHandlerTelego`, se a mensagem encaminhada contiver conteúdo de post (mídia ou legenda/texto), invocar também a lógica do PostBuilder para extrair as propriedades do post e oferecer o botão `🛠️ Post Builder`.

### Passo 3: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. Preserva o fluxo de adição de canais enquanto integra o PostBuilder.

## Impactos esperados
- Logs claros no console/arquivo mostrando exatamente o caminho percorrido por cada mensagem.
- Suporte total a mensagens encaminhadas de qualquer canal no PostBuilder.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Encaminhar uma mensagem/mídia de um canal não configurado no bot no chat privado.
2. Observar os logs detalhados no terminal/arquivo de log.
3. Verificar se o bot oferece a criação da postagem no PostBuilder com o botão `🛠️ Post Builder`.
