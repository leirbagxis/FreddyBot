# Plano: Adicionar Botão "Cancelar" com Estilo Danger (Vermelho) no Menu de Salvamento

## Pedido do usuário
Adicionar um botão "Cancelar" no fim do menu de salvamento do PostBuilder com estilo vermelho (`Style: "danger"` nas propriedades do `InlineKeyboardButton`).

## Objetivo
Incluir o botão `❌ Cancelar` com estilo `danger` (cor vermelha nativa da Telegram Bot API) no menu de salvamento do post (tanto em `pb-save` quanto em `showSavedPostMenu`), permitindo ao usuário encerrar/cancelar o post salvo com edição limpa de mensagem.

## Contexto atual
- O menu de salvamento (`showSavedPostMenu` e `pb-save`) possui 4 botões: Compartilhar, Enviar para Canais, Agendar Envio e Salvar como Template.
- O tipo `telego.InlineKeyboardButton` possui o campo nativo `Style: "danger"` que aplica a cor vermelha aos botões em clientes Telegram compatíveis.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Adicionar o botão `❌ Cancelar` com `Style: "danger"`**:
   - Em `case "pb-save"` e `showSavedPostMenu`, adicionar a 5ª linha no teclado inline:
     ```go
     {
         {Text: "❌ Cancelar", CallbackData: "pb-cancel", Style: "danger"},
     }
     ```

2. **Garantir tratamento limpo em `case "pb-cancel"`**:
   - Ao clicar no botão `❌ Cancelar`, a função `pb-cancel` apaga o estado do PostBuilder e edita a mensagem atual (`EditMessageText`) para `"❌ Post Builder cancelado."`.

---

## Passos detalhados

### Passo 1: Atualizar o teclado inline em `case "pb-save"` e `showSavedPostMenu`
- Incluir a nova linha de botão contendo `Text: "❌ Cancelar"`, `CallbackData: "pb-cancel"`, `Style: "danger"`.

### Passo 2: Atualizar o handler `case "pb-cancel"`
- Garantir que `pb-cancel` edite a mensagem (`EditMessageText`) para `"❌ Post Builder cancelado."`.

### Passo 3: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. `Style: "danger"` é um campo padrão suportado pela biblioteca `telego` e Telegram Bot API.

## Impactos esperados
- Menu de salvamento ganha opção visível de cancelamento em vermelho (`danger`).
- Clique no botão edita a mensagem atual para o aviso de cancelamento sem poluir o chat.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Salvar uma postagem no PostBuilder -> Verificar a presença do botão "❌ Cancelar" em vermelho no fim do menu de salvamento.
2. Clicar em "❌ Cancelar" -> Verificar que a mensagem é editada para "❌ Post Builder cancelado.".
