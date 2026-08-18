# Plano: Redesign da Dashboard de Configurações Mobile do FreddyBot (Estilo Telegram Mini App)

## Pedido do Usuário
Redesenhar a dashboard de configuração de canais do FreddyBot no mobile mantendo ~80% do conceito e arquitetura de informação original, mas refinando o espaçamento, a hierarquia visual, o contraste e a usabilidade de toque para que pareça um **Telegram Mini App nativo**.

---

## Requisitos de Design e Diretrizes Técnicas

### 1. Cores e Hierarquia Visual
- **Fundo da Aplicação**: Navy escuro inspirado no Telegram (`#0f172a` / `#0e1621` ou `var(--bg)`).
- **Cartões de Seção**: Navy sutilmente mais claro (`#17212b` ou `var(--card)`), cantos arredondados (`rounded-2xl`), sem aninhamentos excessivos (Estrutura: `Page -> Section Card -> Content`).
- **Azul Telegram (`#2481cc`)**: Reservado para ações primárias (ex: "Salvar Reações"), estados ativos e cápsula selecionada na navegação inferior.
- **Paleta de Ícones por Funcionalidade**:
  - **Azul**: Captions (`CaptionCard.tsx`)
  - **Amarelo / Laranja**: New Pack Caption (`NewPackCaptionCard.tsx`)
  - **Roxo**: Reações / Votos Grid (`ReactionsCard.tsx`)
  - **Verde**: Reações Nativas do Telegram (`NativeReactionsCard.tsx`)

### 2. Reformulação dos Cards Principais

#### A. Caption Padrão (`CaptionCard.tsx`)
- Ícone de documento azul (`FileText`).
- Título **Caption Padrão** e subtítulo *"Aplicada em todas as mensagens"*.
- Botão de edição alinhado à direita.
- Área de prévia embaixo exibindo a legenda atual com links destacados no azul Telegram (`text-[#2481cc]`).

#### B. New Pack Caption (`NewPackCaptionCard.tsx`)
- Ícone de pacote amarelo/laranja (`Package`).
- Título **New Pack Caption** e subtítulo *"Template para novo pack"*.
- Botão de edição à direita.
- Área de prévia destacando visualmente as variáveis (ex: `$name`, `$title`, `$link`, `$count`) com badge/cor de destaque de acento (`text-[#2481cc]` / `bg-accent/15`).
- Preservar símbolos decorativos Unicode e formatação.

#### C. Reações / Votos (Grid) (`ReactionsCard.tsx`)
- Ícone roxo (`SmilePlus`).
- Título **Reações / Votos (Grid)** e subtítulo *"Adicione até 5 emojis para votação rápida."*.
- **5 slots interativos de emojis**:
  - Slots vazios: quadrado com borda tracejada (`border-2 border-dashed border-border`), `+` centralizado, target de toque confortável (`size-12`), com estados hover/pressed.
  - Slots preenchidos: emoji em destaque em tamanho grande com botão para limpar/substituir.
- Botão proeminente de largura total: **Salvar Reações** em azul Telegram (`bg-[#2481cc] hover:bg-[#1f70b2] text-white font-bold h-12 rounded-xl`).

#### D. Divisor "OUTRAS OPÇÕES"
- Linha sutil com rótulo pequeno mutado: `OUTRAS OPÇÕES` (sem ondas decorativas entre todas as seções).

#### E. Reações Nativas do Telegram (`NativeReactionsCard.tsx`)
- Ícone verde (`Heart`).
- Configuração em linha compacta: ícone + título **Reações Nativas do Telegram** + descrição *"Use as reações nativas do Telegram em vez do grid."* + Switch alinhado à direita.

### 3. Navegação Inferior (`TabBar.tsx` & `index.css`)
- Item selecionado dentro de uma cápsula/pill azul arredondada (`[ Legendas ]` -> `bg-[#2481cc] text-white font-semibold rounded-full px-3.5 py-1.5 flex items-center gap-2`).
- Sensação de app mobile nativo do Telegram.
- Fixed/sticky no rodapé com padding inferior garantido na página (`pb-28`) para que o conteúdo NUNCA seja sobreposto pela barra.

---

## Arquivos que serão modificados

- `dashboard/src/components/CaptionCard.tsx` (redesenhar card no padrão azul)
- `dashboard/src/components/NewPackCaptionCard.tsx` (redesenhar card no padrão amarelo/laranja com destaque em variáveis)
- `dashboard/src/components/ReactionsCard.tsx` (5 slots interativos de emojis + botão Salvar Reações azul Telegram)
- `dashboard/src/components/NativeReactionsCard.tsx` (redesenhar card no padrão verde compacto)
- `dashboard/src/components/CaptionPreview.tsx` (highlighting de variáveis `$name`, `$title`, `$link`, `$count` e links Telegram)
- `dashboard/src/components/TabBar.tsx` (cápsula azul no item ativo)
- `dashboard/src/App.tsx` (divisores sutis "OUTRAS OPÇÕES" e espaçamento `space-y-5 pb-28`)
- `dashboard/src/index.css` (estilos do tema Telegram dark mini app, animações e feedback de toque)

---

## Passos Detalhados

### Passo 1 — Atualizar `CaptionPreview.tsx`
- Adicionar suporte a destaque visual elegante de variáveis (ex: `$name`, `$title`, `$link`, `$count`) no preview das legendas.

### Passo 2 — Redesenhar `CaptionCard.tsx`
- Aplicar ícone de documento azul, estrutura de header limpa, área de preview sutil e botão de editar.

### Passo 3 — Redesenhar `NewPackCaptionCard.tsx`
- Aplicar ícone de pacote amarelo/laranja, formulário moderno de opções e preview com variáveis destacadas.

### Passo 4 — Redesenhar `ReactionsCard.tsx`
- Implementar os 5 slots de emoji interativos (borda tracejada, toque largo, emoji em grande tamanho) e o botão primário "Salvar Reações" em azul Telegram (`#2481cc`).

### Passo 5 — Redesenhar `NativeReactionsCard.tsx`
- Aplicar ícone verde, linha compacta com Switch à direita, e opções expansíveis de modo (Fixo / Aleatório).

### Passo 6 — Redesenhar `TabBar.tsx` e Rodapé no `App.tsx` & `index.css`
- Atualizar a barra de navegação inferior com a cápsula azul (`bg-[#2481cc] text-white rounded-full`) para a aba ativa.
- Adicionar divisor "OUTRAS OPÇÕES" entre as seções principais e opções secundárias.
- Garantir `pb-28` em todas as páginas para evitar sobreposição da navegação.

---

## Como Testar
```bash
cd dashboard && npm run build
```
- Testar compilação limpa do Vite.
- Verificar visualmente no mobile que o layout possui feedback de toque responsivo, zero sobreposição no rodapé e estética Telegram Mini App impecável.

## Rollback
```bash
git checkout dashboard/src/components/CaptionCard.tsx dashboard/src/components/NewPackCaptionCard.tsx dashboard/src/components/ReactionsCard.tsx dashboard/src/components/NativeReactionsCard.tsx dashboard/src/components/CaptionPreview.tsx dashboard/src/components/TabBar.tsx dashboard/src/App.tsx dashboard/src/index.css
```
