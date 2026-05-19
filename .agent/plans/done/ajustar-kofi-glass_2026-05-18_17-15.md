# Plano: Ajustar Design Ko-fi com Efeito Liquid-Glass

## Pedido do usuário
O usuário quer um design específico inspirado no Ko-fi:
- Fundo: `Gold Lightest` (`#faf6ee`)
- Acento: `Green Accent` (`#00754A`)
- As divs não devem ser totalmente verdes (desfazer a mudança anterior onde os cards ficaram `#3d5229`).
- Cantos arredondados (mais circulares/curvos).
- Efeito `liquid-glass` (glassmorphism/vidro translúcido).

## Objetivo técnico
Reverter as superfícies (cards, divs) para um tom claro, mas com opacidade, aplicando filtros de desfoque (`backdrop-filter`) para criar o efeito "liquid-glass". Ajustar as variáveis globais para a paleta Ko-fi. Aumentar o `border-radius` de botões e cards.

## Estratégia de implementação

**1. Paleta de Cores (Tema Light):**
- `--bg`: `#faf6ee` (Gold Lightest)
- `--card`: `rgba(255, 255, 255, 0.65)` (Vidro translúcido)
- `--card-elevated`: `rgba(255, 255, 255, 0.85)`
- `--text`: `#1a2b22` (Texto escuro legível)
- `--text-secondary`: `#556b5e`
- `--accent`: `#00754A` (Green Accent)
- `--accent-soft`: `rgba(0, 117, 74, 0.1)`
- `--border`: `rgba(0, 117, 74, 0.15)` (Borda sutil esverdeada)
- `--nav-bg`: `rgba(250, 246, 238, 0.8)`

**2. Efeito Liquid-Glass (Classes CSS):**
- Atualizar `.card`, `.top-bar`, `.bottom-nav` para incluir:
  ```css
  backdrop-filter: blur(16px) saturate(180%);
  -webkit-backdrop-filter: blur(16px) saturate(180%);
  ```
- Para dar o volume do "liquid", usar uma sombra interna (inset shadow) sutil misturada com uma sombra externa:
  `box-shadow: 0 4px 24px rgba(0, 0, 0, 0.04), inset 0 1px 1px rgba(255, 255, 255, 0.6);`

**3. Cantos Arredondados:**
- `.card`: Aumentar `border-radius` para `24px` ou `28px`.
- `.btn`: Alterar `border-radius` para `100px` (estilo pill/cápsula).
- `.input`, `.perm-row`: Alterar para `16px` ou `20px`.
- `.stat-card`: `20px`.

## Passos detalhados
1.  **Editar Variáveis em `index.css`:** Substituir as cores do `[data-theme="light"]` pelas variáveis Ko-fi e RGBA para vidro. Retornar `--text` para escuro.
2.  **Atualizar Border Radius:** Procurar as classes `.card`, `.btn`, `.input`, `.perm-row` em `index.css` e aumentar seus valores de `border-radius`.
3.  **Aplicar Backdrop Filter:** Adicionar a regra de `backdrop-filter` na classe `.card`, que é o principal container. A `.bottom-nav` e `.top-bar` já devem ter algo parecido, mas garantirei que estejam intensificadas para o "liquid-glass".
4.  **Ajustar `useTheme.ts`:** Atualizar as cores de fundo do Telegram SDK para `#faf6ee`.

## Riscos e Impactos
- O efeito `backdrop-filter` pode consumir um pouco mais de processamento no mobile (CSS repaints), mas é amplamente suportado hoje em dia.
- O tema Dark precisará de um ajuste recíproco: fundo `#0b120f`, accent `#00754A` (ou um verde mais claro para contraste, como `#00a86b`), e cards translúcidos escuros `rgba(255,255,255,0.05)`.
