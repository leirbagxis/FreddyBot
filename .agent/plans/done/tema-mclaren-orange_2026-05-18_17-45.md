# Plano: Tema Papaya Orange (McLaren)

## Pedido do usuário
O usuário deseja alterar a cor de acento do verde (Ko-fi) para o "Laranja McLaren" (Papaya Orange), mantendo a estética Liquid-Glass e os cantos arredondados estabelecidos anteriormente.

## Objetivo técnico
Substituir todas as instâncias da cor de acento verde pelas variações do Laranja McLaren (`#FF8000`), ajustando as variáveis RGB correspondentes para os efeitos de hover e soft-background.

## Estratégia de implementação

**1. Paleta de Cores Papaya Orange:**
A cor icônica da McLaren é aproximadamente `#FF8000` (RGB: 255, 128, 0).
Vou usar essa cor como `--accent` para ambos os temas (Light e Dark), pois o laranja vibrante contrasta muito bem tanto no fundo creme (Light) quanto no fundo escuro (Dark).

**2. Atualização no Tema Light (`[data-theme="light"]`):**
- `--accent`: `#FF8000`
- `--accent-rgb`: `255, 128, 0`
- `--accent-soft`: `rgba(255, 128, 0, 0.08)`
- `--accent-hover`: `rgba(255, 128, 0, 0.12)`
- `--link`: `#FF8000`
- `--border-active`: `#FF8000`
- (Opcional) Reverter `--success` para um verde padrão (`#2f9e44`), pois havíamos alterado para o verde Ko-fi, mas "sucesso" em laranja não faz sentido semântico.

**3. Atualização no Tema Dark (`[data-theme="dark"]`):**
- `--accent`: `#FF8A00` (Um tom ligeiramente mais quente para o escuro).
- `--accent-rgb`: `255, 138, 0`
- `--accent-soft`: `rgba(255, 138, 0, 0.12)`
- `--accent-hover`: `rgba(255, 138, 0, 0.18)`
- `--link`: `#FF8A00`
- `--border-active`: `#FF8A00`
- Fundo: O fundo atual Dark é `Deep Forest` (`#0d1310`). Com o tema McLaren, um fundo "Carbon Fiber" (cinza carvão ultra-escuro) combina melhor do que verde. Vou ajustar `--bg` para `#0f1115` e `--nav-bg` correspondentemente.

## Passos detalhados
1.  Editar `dashboard/src/index.css`:
    - Atualizar as variáveis `--accent` e relacionadas no tema Light para Laranja.
    - Reverter `--success` no tema Light para `#2f9e44`.
    - No tema Dark, atualizar os acentos para Laranja e os fundos base (`--bg`, `--nav-bg`, `--input-bg`, `--toggle-bg`) para tons de cinza carvão/carbono.
2.  Editar `dashboard/src/hooks/useTheme.ts`:
    - Atualizar as cores de fundo do SDK do Telegram no tema Dark para o novo cinza carbono (`#0f1115`).
3.  Apresentar o plano e, após aprovação, implementar.

## Riscos
Nenhum risco funcional. O contraste do botão primário (Laranja) com texto branco (`--accent-text: #ffffff`) é visualmente agradável e remete diretamente aos carros da equipe de F1.