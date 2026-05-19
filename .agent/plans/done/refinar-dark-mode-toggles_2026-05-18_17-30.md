# Plano: Refinar Dark Mode e Animação de Toggles

## Pedido do usuário
O usuário solicitou um reajuste no modo escuro (Dark Mode) para acompanhar a nova estética e a adição de um efeito visual interativo quando os toggles (interruptores) são clicados.

## Objetivo técnico
1. Harmonizar o tema Dark com a estética Liquid-Glass inspirada no Ko-fi (fundo verde ultra-escuro, cards translúcidos, acento verde brilhante).
2. Adicionar uma animação tátil (squish/scale) aos elementos `.toggle` no CSS para melhorar o feedback de interação.

## Estratégia de implementação

**1. Ajustes no Tema Dark (`[data-theme="dark"]`):**
- O fundo atual é `#0b120f`. Vamos manter um tom ultra-escuro (quase preto, com toque verde), mas ajustar os cards para serem "glass" sobre ele.
- `--bg`: `#0d1310` (Dark Forest)
- `--card`: `rgba(255, 255, 255, 0.05)` (Vidro muito sutil)
- `--card-elevated`: `rgba(255, 255, 255, 0.08)`
- `--text`: `#f0f4f2`
- `--accent`: `#00a86b` (Ko-fi green ajustado para brilho no escuro)

**2. Efeito nos Toggles (`.toggle`):**
- Adicionar uma transição de escala quando clicado (estado `:active`).
- Esticar levemente o "botãozinho" interno (knob) durante o clique para dar um efeito "squishy" (tátil/líquido).
```css
.toggle:active {
  transform: scale(0.92);
}
.toggle:active::after {
  width: 22px; /* Efeito squish horizontal */
}
.toggle.on:active::after {
  transform: translateX(16px); /* Compensar o aumento de largura */
}
```

## Passos detalhados
1.  Editar `dashboard/src/index.css` e atualizar as variáveis do bloco `[data-theme="dark"]`.
2.  Procurar a definição da classe `.toggle` no mesmo arquivo e injetar as pseudo-classes de `:active`.
3.  Validar o build.

## Riscos
- O efeito de alargar o knob (`width: 22px`) precisa ter o translateX compensado quando o toggle estiver no estado `.on`, para que não vaze para fora da borda. Testarei os valores `translateX` cuidadosamente.