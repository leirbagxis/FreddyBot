# Plano: Animação Fluida de Troca de Tema

## Pedido do usuário
Adicionar uma animação ao trocar entre modo claro e modo escuro (toggle).

## Objetivo técnico
Implementar uma transição suave e visualmente agradável que afete tanto o ícone do botão quanto o fundo da aplicação, sem comprometer a performance.

## Estratégia de implementação

**1. Animação Global (Cor e Fundo):**
- Adicionar uma transição de `0.4s` nas propriedades de `background-color` e `color` no seletor global (`html, body, #root`). Isso fará com que a mudança de "clarear/escurecer" seja gradual em vez de instantânea.

**2. Animação do Ícone (Toggle):**
- Adicionar um efeito de rotação e escala ao ícone dentro do botão `.theme-switch`.
- Quando o tema mudar, o ícone fará um giro de 360 graus e uma leve pulsação.

**3. Otimização:**
- Usar apenas transições CSS nativas (GPU accelerated) para garantir que a animação rode a 60fps mesmo em celulares.

## Passos detalhados

1. **Editar `index.css`**:
   - Atualizar `html, body, #root` com a transição global.
   - Adicionar estilos de transição para o SVG do `.theme-switch`.
   - Criar uma animação de "morph" visual usando escala e rotação.

2. **Verificar Build**:
   - Garantir que não há erros de sintaxe.

## Riscos
- Se a transição for muito lenta, pode parecer "pesado". Usarei um tempo equilibrado (400ms).

## Como testar
- Clicar no botão de Sol/Lua.
- O fundo deve mudar suavemente.
- O ícone deve girar elegantemente durante a troca.
