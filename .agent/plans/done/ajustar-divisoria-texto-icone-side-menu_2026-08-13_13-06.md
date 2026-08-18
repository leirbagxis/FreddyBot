# Plano: Ajustar Linha Divisória Alinhada do Texto até o Ícone da Direita em SideMenu.tsx

## Pedido do usuário
Ajustar a linha divisória entre os itens de navegação ("Meus Canais", "Templates", "Agendamentos", etc.) para que ela não corte a div inteira por baixo do ícone da esquerda. A linha deve começar exatamente do início do texto (`left-[46px]`) e ir até o final do ícone `>` da direita (`right-4`).

## Objetivo
Reestruturar as linhas divisórias dos itens dentro do card de navegação em `SideMenu.tsx`, posicionando um separador `left-[46px] right-4 h-px bg-white/10` de forma precisa sob o texto e o ícone da direita.

## Contexto atual
- Atualmente, as bordas `border-b border-white/10` estão aplicadas na tag `<button>`, atravessando a div de ponta a ponta inclusive por baixo do ícone da esquerda.

## Arquivos analisados
- `dashboard/src/components/SideMenu.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/SideMenu.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/SideMenu.tsx`**:
   - Remover a classe `border-b border-white/10` da tag `<button>`.
   - Adicionar uma div com posicionamento absoluto para a linha divisória em cada um dos itens (exceto o último):
     `<div className="absolute bottom-0 left-[46px] right-4 h-px bg-white/10 pointer-events-none" />`
   - Onde `46px` corresponde exatamente ao alinhamento do texto (padding esquerdo 16px + ícone 18px + gap 12px), e `right-4` (16px) corresponde à margem final do ícone `>`.

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/components/SideMenu.tsx`.
2. Executar `npm run build`.

## Riscos
- **Baixo**: Ajuste visual fino de posicionamento CSS em `SideMenu.tsx`.

## Impactos esperados
- Linhas divisórias elegantes e alinhadas que iniciam no primeiro caractere do texto do item e terminam na extremidade do ícone `>` da direita.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/components/SideMenu.tsx
```
