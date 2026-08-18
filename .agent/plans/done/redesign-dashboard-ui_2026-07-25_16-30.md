# Plano: redesign-dashboard-ui

## Pedido do usuário
Reajustar o design do dashboard do FreddyBot, que está "feio, divs dentro de divs".

## Objetivo
Reduzir aninhamento desnecessário de divs, melhorar hierarquia visual, consistência de espaçamento e adicionar personalidade visual.

## Contexto atual
- Dashboard React (Tailwind v4 + shadcn/ui) para Telegram WebApp
- Layout: top bar + main content + bottom tab bar
- Card pattern repetitivo: `Card > CardContent > flex > section-icon + text + subtext`
- Muitos `style={{}}` inline no App.tsx (mais de 100 props inline)
- Hierarquia visual plana, mesmo padding em tudo
- Sem identidade visual clara — parece Bootstrap genérico
- App.tsx tem 1407 linhas com toda a lógica e renderização misturadas

## Arquivos analisados
- dashboard/src/App.tsx (1407 linhas)
- dashboard/src/index.css (2504 linhas)
- dashboard/src/components/CaptionCard.tsx
- dashboard/src/components/NewPackCaptionCard.tsx
- dashboard/src/components/ReactionsCard.tsx
- dashboard/src/components/NativeReactionsCard.tsx
- dashboard/src/components/WaveDivider.tsx (PerfLine)
- dashboard/src/components/TabBar.tsx
- dashboard/src/components/ButtonGrid.tsx
- dashboard/src/components/UserTemplatesManager.tsx
- dashboard/src/components/CaptionPreview.tsx
- tailwind.config (via index.css @theme)

## Arquivos que poderão ser modificados
- dashboard/src/App.tsx
- dashboard/src/index.css
- dashboard/src/components/CaptionCard.tsx
- dashboard/src/components/NewPackCaptionCard.tsx
- dashboard/src/components/ReactionsCard.tsx
- dashboard/src/components/NativeReactionsCard.tsx
- dashboard/src/components/TabBar.tsx
- dashboard/src/components/WaveDivider.tsx

## Estratégia de implementação

Direção: **"Telegram Native Pro"** — premium, limpo, com identidade visual forte.

### Pilares do redesign:

1. **Remover aninhamento**: Criar componentes `ActionCard` e `SectionRow` que encapsulam o pattern repetitivo `Card > CardContent > flex > icon + text + chevron` em uma única props interface.

2. **Reduzir style={} inline**: Mover para classes CSS e Tailwind utilitárias.

3. **Melhorar espaçamento**: Padding generoso (p-5/p-6), gap consistente, hierarquia via tamanho/weight.

4. **Cards com identidade**: Borda sutil, glass-morphism, hover states com elevantamento.

5. **Animações**: Staggered reveals, micro-interações em hover/active.

6. **Tipografia**: Headings maiores e mais impactantes, corpo mais limpo.

## Passos detalhados

### Passo 1: index.css — novos componentes utilitários
- `.action-card` — card limpo com padding, hover, transição
- `.action-card-icon` — container de ícone padronizado
- `.section-title` / `.section-desc` — tipografia padronizada
- `.page-section` — espaçamento vertical entre seções
- Reduzir regras CSS mortas ou duplicadas

### Passo 2: App.tsx — refatorar renderização dos cards
- Substituir todos os `Card > CardContent > flex > section-icon + text + chevron` por `action-card`
- Substituir inline styles por classes Tailwind
- Agrupar seções com `page-section`
- Adicionar animate-stagger-in nos cards

### Passo 3: Componentes individuais
- `CaptionCard` — reduzir aninhamento, melhorar preview toggle
- `NewPackCaptionCard` — alinhar com novo padrão visual
- `ReactionsCard` — alinhar com novo padrão
- `NativeReactionsCard` — alinhar com novo padrão
- `TabBar` — melhorar visual dos tabs (ícone + label mais integrados)

### Passo 4: Limpeza
- Remover estilos inline redundantes
- Consolidar classes repetidas
- Verificar build

## Riscos
- Quebrar layout em dispositivos móveis (Telegram WebView)
- Perder estados de foco/acessibilidade
- Regressão em temas (light/dark/telegram)

## Impactos esperados
- Código mais limpo e manutenível
- UI mais consistente e agradável
- Menos divs aninhadas
- Melhor performance (menos style={} recalcs)

## Compatibilidade
- Linux (dev)
- Telegram WebView (iOS/Android)
- Desktop browsers

## Como testar

### Build
```bash
cd dashboard && npx tsc --noEmit && npx vite build
```

### Dev
```bash
cd dashboard && npx vite --host
```

## Rollback
```bash
git checkout -- dashboard/src/
```

## Observações
- Foco em mobile-first (Telegram WebApp é mobile)
- Preservar todos os temas (light, dark, telegram)
- Não alterar lógica de negócio — só apresentação
