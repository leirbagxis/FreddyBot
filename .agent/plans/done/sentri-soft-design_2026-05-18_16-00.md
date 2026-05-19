# Plano: sentri-soft-design_2026-05-18_16-00.md

## Pedido do usuário
Suavizar as cores da dashboard para não parecer "modo daltônico" (alto contraste excessivo), mantendo a essência do `@DESIGN.md`.

## Objetivo
Criar a versão "Sentri Soft": trocar o Verde Lima (`Electric Lime`) pelo Roxo Vibrante (`Accent Violet`) como cor principal de interação, e usar o fundo de forma menos agressiva.

## Estratégia de implementação
1.  **Troca de Protagonismo:** 
    - A cor `--accent` principal passará a ser `#6a5fc1` (o roxo do design system).
    - O `#c2ef4e` (verde lima) será movido para uma variável secundária `--lime` e usado apenas em pequenos detalhes (ícones específicos ou badges de sucesso).
2.  **Ajuste de Superfícies:**
    - No modo Dark, clarear levemente o fundo dos cards para dar mais profundidade.
    - No modo Light, usar o roxo como acento principal ao invés do preto absoluto.

## Passos detalhados (CSS)

1.  **Modificar index.css:**
    - Atualizar `[data-theme="light"]`: `--accent` muda de `#150f23` para `#6a5fc1`.
    - Atualizar `[data-theme="dark"]`: `--accent` muda de `#c2ef4e` para `#7c7fff` (um violeta mais visível no escuro).
    - Adicionar `--accent-vibrant: #c2ef4e` para usos pontuais.
    - Ajustar `--bg` e `--card` no dark para tons de cinza-roxo menos "vazios".

## Riscos
- Perder a "assinatura" da marca Sentri, mas o ganho em conforto visual compensa para uma dashboard de uso prolongado.

## Como testar
- Navegar pelas abas e ver se os botões e estados "On" parecem mais integrados e menos destacados "a força".
