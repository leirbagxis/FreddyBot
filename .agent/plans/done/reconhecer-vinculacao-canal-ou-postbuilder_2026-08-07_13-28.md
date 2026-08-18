# Plano: Reconhecimento Inteligente de Vinculação de Canal vs PostBuilder

## Pedido do usuário
Quando o usuário encaminha uma mensagem/post de um canal no privado:
1. **Se o bot ESTÁ no canal, É ADMIN, mas NÃO ESTÁ CONFIGURADO no banco**: Enviar a mensagem de vinculação de canal ("Deseja cadastrar o canal X no bot?").
2. **Se o bot ESTÁ no canal, É ADMIN e JÁ ESTÁ CONFIGURADO no banco**: Não enviar vinculação e ativar o **PostBuilder** ("✨ Conteúdo detectado! Deseja usar o Post Builder...").
3. **Se o bot NÃO é admin do canal**: Ativar o **PostBuilder** diretamente.

## Objetivo
Implementar a verificação de 3 condições ao receber mensagens encaminhadas de canais, direcionando para a vinculação do canal ou para o PostBuilder conforme a configuração do canal no banco de dados.

## Contexto atual
- `GetChannelByID` permite consultar se o canal já está configurado no banco de dados.
- `GetChatMember` permite verificar se o bot é admin no canal encaminhado.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`
- `internal/middleware/checkAddBotMiddlewareTelego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Estratégia de implementação

1. **Checagem no recebimento de mensagem encaminhada (`ProcessIncomingContentTelego`)**:
   - Ao receber mensagem encaminhada de canal (`ForwardOrigin` do tipo `MessageOriginChannel`):
     - Obter `channelID`.
     - Verificar se o canal já existe no banco (`c.ChannelService.GetChannelByID`).
     - Se **NÃO** existe no banco:
       - Checar se o bot é membro e possui permissões de admin (`b.GetChatMember`).
       - Se for admin no canal: enviar a mensagem de vinculação de canal (`toadd-require-message` com botões Sim/Não) e interromper.
     - Se o canal **JÁ existe no banco** ou se o bot **NÃO é admin**:
       - Prosseguir para o **PostBuilder**, extraindo mídias, legenda, formatação HTML, emojis e botões, exibindo a prompt do `🛠️ Post Builder`.

2. **Centralizar e organizar no `addChannel.go` e `postBuilder.go`**:
   - Exportar helper `SendAddChannelPromptTelego` em `addChannel.go` para ser chamado quando o bot for admin mas o canal não estiver configurado.

---

## Passos detalhados

### Passo 1: Criar helper `SendAddChannelPromptTelego` em `addChannel.go`
- Enviar a mensagem `toadd-require-message` com os botões `Sim` / `Não`.

### Passo 2: Atualizar `ProcessIncomingContentTelego` em `postBuilder.go`
- Se for mensagem encaminhada de canal, checar existência no DB e permissões de admin.
- Se for admin e não estiver no DB ➔ Chamar `SendAddChannelPromptTelego`.
- Se já estiver no DB ou não for admin ➔ Processar PostBuilder normalmente.

### Passo 3: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. Mantém a separação perfeita de regras.

## Impactos esperados
- Canais não configurados onde o bot é admin oferecem vinculação.
- Canais já configurados (e mídias gerais) abrem o PostBuilder diretamente.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Encaminhar mensagem de um canal onde o bot É admin mas NÃO está configurado ➔ Bot manda prompt de vinculação de canal.
2. Encaminhar mensagem de um canal onde o bot É admin e JÁ ESTÁ configurado ➔ Bot abre o PostBuilder.
3. Encaminhar mensagem de canal onde o bot NÃO é admin ➔ Bot abre o PostBuilder.
