# Plano: corrigir-layout-e-sessao-postbuilder

## Pedido do usuário
1. Mover o botão "Enviar para Canais" para uma linha abaixo do botão "Compartilhar".
2. Corrigir o erro "Sessão expirada ou não encontrada" que ocorre ao clicar em "Enviar para Canais".

## Objetivo
Ajustar a interface do PostBuilder e permitir que as ações pós-salvamento funcionem mesmo após a exclusão do estado temporário de edição.

## Contexto atual
- Ao salvar a postagem (`pb-save`), o estado temporário (`PostBuilderState`) é deletado para liberar memória.
- O `CallbackHandler` possui uma verificação global que exige a existência do `state` para quase todos os comandos.
- Como o botão "Enviar para Canais" utiliza o ID da sessão salva (persistente), ele não deve depender do estado temporário de edição.

## Arquivos analisados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. **Corrigir Verificação de Sessão**: Atualizar a condição inicial do `CallbackHandler` para permitir que os novos prefixos (`pb-send-to-channels:` e `pb-send-apply:`) passem mesmo que o `state` seja `nil`.
2. **Ajustar Teclado**: Modificar a estrutura do `InlineKeyboardMarkup` no caso `pb-save` para colocar os botões em linhas separadas.

## Passos detalhados

1. Editar `internal/telegram/events/postBuilder/postBuilder.go`:
    - Na verificação inicial de `state == nil`, adicionar exceções para os prefixos de envio.
    - No caso `pb-save`, alterar a matriz `InlineKeyboard` para duas linhas.

## Riscos
- Nenhum risco técnico identificado.

## Impactos esperados
- Interface mais limpa com botões empilhados.
- Funcionamento correto do envio para canais após o salvamento da postagem.

## Como testar

### Execução
1. Criar e salvar uma postagem.
2. Verificar se o botão "Enviar para Canais" está abaixo de "Compartilhar".
3. Clicar em "Enviar para Canais" e verificar se a lista de canais aparece sem erros de sessão.
