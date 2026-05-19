# Plano: Efeito Liquid Glass (Alta Performance)

## Pedido do usuário
Adicionar o efeito "Liquid Glass" à dashboard.

## Objetivo técnico
Implementar o visual de vidro líquido (transparência, brilho de borda e desfoque sutil) mantendo a fluidez de 60fps. O segredo é usar um desfoque muito baixo e focar nos reflexos de borda.

## Estratégia de implementação

**1. Definição do Estilo Liquid Glass:**
- **Fundo:** Usar `rgba` com transparência moderada.
- **Borda:** Uma borda fina (1px) com cor de destaque suave.
- **Reflexo (Specular):** Adicionar um `box-shadow: inset 0 1px 1px rgba(255,255,255,0.3)` para simular o brilho no topo do vidro.
- **Blur:** Manter o `backdrop-filter: blur(6px)`. É o limite para não causar lag no mobile e ainda parecer vidro.

**2. Aplicação nos Elementos:**
- Aplicar às classes `.card`, `.top-bar`, `.bottom-nav` e `.welcome-card`.

**3. Otimização de GPU:**
- Garantir que as transições não afetem o `backdrop-filter` durante animações para evitar quedas de frames.

## Passos detalhados

1. **Atualizar `.card` no `index.css`**: Inserir os reflexos internos e o blur otimizado.
2. **Atualizar `.top-bar` e `.bottom-nav`**: Sincronizar o efeito de vidro líquido em toda a interface.
3. **Ajustar Variáveis**: Refinar as cores de `--card` para que a transparência realce o efeito.

## Riscos
- Se o usuário achar que o blur de 6px ainda pesa, precisaremos reduzir para 4px. Começarei com 6px para um visual "premium".

## Como testar
- Verificar se os cards parecem peças de vidro flutuando sobre o fundo.
- Notar o brilho sutil na borda superior dos elementos.
- Confirmar que a rolagem continua lisa.
