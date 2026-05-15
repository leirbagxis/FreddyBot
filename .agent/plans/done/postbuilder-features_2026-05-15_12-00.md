# Plano: postbuilder-features_2026-05-15_12-00.md

## Pedido do usuário
1. **Importação de Canal:** Adicionar um recurso no PostBuilder que permita ao usuário importar os dados (legenda padrão, reações e botões) de um dos seus canais cadastrados diretamente para a postagem que está sendo construída. Se ele não tiver canais, exibir um aviso.
2. **Gerenciamento de Botões:** Modificar o fluxo de botões no PostBuilder. Ao clicar em "Botões", em vez de pedir diretamente para adicionar um novo, exibir uma lista com os botões atuais. O usuário poderá excluir botões existentes clicando neles, ou clicar em um botão "➕ Adicionar Novo Botão".

## Objetivo
- Criar o fluxo `pb-import-channel` para consultar e listar canais e `pb-import-apply:{ID}` para aplicar os dados ao `PostBuilderState`.
- Criar o fluxo `pb-manage-buttons` para listar os botões do `PostBuilderState` com opções de exclusão (`pb-del-button:{index}`) e um botão para iniciar a adição (`pb-add-button`).

## Contexto atual
- O PostBuilder gerencia estado via Redis (`PostBuilderState`).
- O botão atual de "Botão" (`pb-add-button`) leva direto ao prompt de entrada de texto.
- `ChannelService.GetUserChannels` será usado para a listagem de canais.

## Arquivos analisados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Estratégia de implementação

### 1. Importação de Canal
- **Menu Principal:** Adicionar o botão `📥 Importar Canal` (`pb-import-channel`) na função `showMenu`.
- **Callback `pb-import-channel`:** Consultar `c.ChannelService.GetUserChannels`. Se vazio, `ShowAlert`. Se houver canais, listar botões (`pb-import-apply:{ID}`).
- **Callback `pb-import-apply:{ID}`:** Obter canal, sobrescrever `state.Body`, `state.Reactions`, `state.Buttons` e voltar ao menu.

### 2. Gerenciamento de Botões
- **Menu Principal:** Alterar a ação do botão "🔘 Botão" para apontar para `pb-manage-buttons`.
- **Callback `pb-manage-buttons`:** Exibir um teclado inline. 
  - Para cada botão no estado, criar um botão na interface do tipo `❌ [Nome do Botão]` com callback `pb-del-button:{index}`.
  - Adicionar um botão `➕ Adicionar Botão` apontando para a lógica original (`pb-add-button`).
  - Adicionar um botão `🔙 Voltar` (`pb-start`).
- **Callback `pb-del-button:{index}`:** Interceptar prefixo, extrair index, remover elemento do slice `state.Buttons`, salvar no cache e recarregar `pb-manage-buttons`.
- **Prompt de Adição (`pb-add-button`):** Manter como é hoje, mas ao terminar de adicionar o texto, voltar para a tela de gerenciamento de botões (`pb-manage-buttons`) em vez da tela principal, para facilitar adições múltiplas.

## Passos detalhados
1.  Atualizar o slice de botões no `showMenu` (Adicionar Importar e trocar callback de Botões).
2.  Implementar `showButtonManager(ctx, b, chatID, userID, c, state)` para exibir a lista de botões.
3.  Atualizar `CallbackHandler` para lidar com prefixos dinâmicos (`pb-import-apply:` e `pb-del-button:`).
4.  Atualizar `handleTextInput` no caso `awaiting_button` para chamar `showButtonManager` após o sucesso.

## Impactos esperados
- PostBuilder se tornará muito mais flexível, permitindo correções rápidas (exclusão de botões) e reaproveitamento de configurações de canais existentes.

## Compatibilidade
- Linux, Docker

## Como testar
1. Abrir PostBuilder.
2. Clicar em "📥 Importar Canal" e verificar a sobrescrita.
3. Clicar em "🔘 Botão". Verificar se abre a lista.
4. Adicionar botões e verificar se eles aparecem na lista.
5. Clicar nos botões da lista e verificar se são excluídos do preview.

## Rollback
Reverter o arquivo `postBuilder.go`.