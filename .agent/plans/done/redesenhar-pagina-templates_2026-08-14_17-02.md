# Plano: Redesenhar Página "Meus Templates" (Templates Page UI Refinement)

## Pedido do usuário
Redesenhar e refinar a página "Meus Templates" (`UserTemplatesManager.tsx`) para o formato nativo de Telegram Mini App, melhorando a hierarquia visual, usabilidade do editor de botões inline, acionamento de criação compacto, visualização colapsada/expandida de cards individuais e confirmação segura de exclusão.

## Objetivo
Transformar a estrutura de cards aninhados atual em um layout limpo, escovado e escalável onde cada template é seu próprio card independente com suporte a:
- Ação compacta de criação `+ Novo template`
- Cabeçalho de card clicável com nome curto, preview de legenda, badge de contagem de botões, lixeira e chevron de expansão (`›` / `⌃`)
- Editor expandido com seções claras: **Informações do template** (nome curto + textarea multilinha de legenda) e **Botões Inline** (controles explícitos de Linhas `[-] X [+]` e Colunas `[-] Y [+]` + Grade visual interativa com slots `+` pontilhados)
- Botão de salvamento visível **Salvar alterações**
- Ação destrutiva **Excluir template** em vermelho com confirmação modal (`ConfirmModal`)

## Contexto atual
- O componente `UserTemplatesManager.tsx` hoje utiliza um container `<Card>` gigante com formulário de criação fixo e grande que consome muito espaço vertical.
- As opções de grid e botões utilizam seletores amontoados sem hierarquia clara.
- Os cards de template colapsados/expandidos precisam de alinhamento tátil mobile.

## Arquivos analisados
- `dashboard/src/components/UserTemplatesManager.tsx` — Gerenciador de templates e editor de botões
- `dashboard/src/components/ConfirmModal.tsx` — Modal de confirmação reutilizável
- `dashboard/src/App.tsx` — Container de rota da página Meus Templates

## Arquivos que poderão ser modificados
- `dashboard/src/components/UserTemplatesManager.tsx`
- `dashboard/src/App.tsx`

## Estratégia de implementação

### 1. Ajuste do Cabeçalho da Página (`App.tsx`)
- Substituir o cabeçalho simples por:
  - Título: **`Meus Templates`** com ícone de camadas/templates (`Layers`).
  - Subtítulo mutado: `"Crie e gerencie seus templates de legendas e botões."`.

### 2. Seção de Templates e Criação Compacta (`UserTemplatesManager.tsx`)
- **Seção de Cabeçalho**: `Meus Templates` acompanhado de badge sutil de contagem total.
- **Ação Compacta de Criação**:
  - Botão/Card retrátil `+ Novo template` com subtítulo *"Crie um novo template de legenda e botões."*.
  - Ao ser clicado, expande suavemente o campo de input do nome curto (`Ex: promo_blackfriday`) e o botão `+ Adicionar`.

### 3. Card de Template (Estado Colapsado)
- Cada template vira um card independente com borda suave (`rounded-2xl border border-border/80 bg-card p-3.5`).
- **Linha do Cabeçalho (Clicável inteira)**:
  - Ícone `#` azul em container quadrado arredondado.
  - Título principal em destaque: **`promo`**.
  - Preview secundário da legenda truncado abaixo.
  - Badge sutil com quantidade de botões configurados (ex: `1`).
  - Ícone de lixeira para exclusão.
  - Chevron indicativo de estado: `›` (colapsado) e `⌃` (expandido).

### 4. Conteúdo Expandido do Template
- **Seção 1: Informações do Template**:
  - Título da seção: `Informações do template`.
  - Campo **Nome curto** (`[ promo ]`).
  - Campo **Legenda (opcional)** com dica de variáveis (`Use variáveis como {nome}`).
  - Textarea multilinha para edição da legenda.

- **Seção 2: Botões Inline**:
  - Título da seção com ícone azul de grade: **`▦ Botões Inline`**.
  - Descrição: *"Configure a grade de botões da mensagem."*.
  - Controles explícitos de layout:
    - `Linhas [ - ] 4 [ + ]`
    - `Colunas [ - ] 3 [ + ]`
  - Grade visual interativa:
    - Botões existentes exibem rótulo e URL.
    - Slots vazios com borda pontilhada suave, cantos arredondados e botão `+` centralizado.
  - Dica sutil: *"💡 Dica: Clique em um botão para editar. Arraste para reordenar."*.

- **Seção 3: Salvamento e Exclusão**:
  - Botão primário em Azul Telegram: `✓ Salvar alterações`.
  - Link destrutivo em vermelho no rodapé: `Excluir template`.
  - Confirmação via `ConfirmModal`: *"Excluir template 'promo'? Essa ação não poderá ser desfeita."*.

## Passos detalhados
1. Atualizar a página `App.tsx` para formatar o cabeçalho de "Meus Templates" com descrição e ícone.
2. Refatorar o componente `UserTemplatesManager.tsx` implementando a nova estrutura de formulário retrátil de criação, cards independentes de template, editores de legenda multilinha e grade de botões inline com seletores explícitos `[-] X [+]`.
3. Integrar o `ConfirmModal` para proteção contra exclusão acidental.
4. Executar `cd dashboard && npm run build` para garantir compilação limpa.

## Riscos
- Nenhum risco ao backend. Todas as rotas da API (`listUserCaptionTemplates`, `createUserCaptionTemplate`, `updateUserCaptionTemplate`, `deleteUserCaptionTemplate`, etc.) continuam idênticas.

## Impactos esperados
- Experiência fluida e 100% otimizada para Telegram Mini App no celular.
- Clareza máxima na criação e edição de botões e legendas.
- Navegação tátil confortável sem amontoamento de cards.

## Compatibilidade
- Linux, macOS, Windows, Mobile (iOS/Android)

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
Restaurar `UserTemplatesManager.tsx` e `App.tsx`.
