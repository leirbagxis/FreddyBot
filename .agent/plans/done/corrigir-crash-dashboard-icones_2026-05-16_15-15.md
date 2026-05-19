# Plano: corrigir-crash-dashboard-icones_2026-05-16_15-15.md

## Pedido do usuário
Correção de tela preta ao acessar a aba de permissões.

## Objetivo técnico
Corrigir erro de execução (Runtime Error) causado pela falta da importação do ícone `Link2` no componente `App.tsx`.

## Contexto atual
A funcionalidade de Links Dinâmicos foi implementada usando o ícone `Link2` da biblioteca `lucide-react`, mas o símbolo não foi adicionado ao bloco de imports.

## Arquivos analisados
- `dashboard/src/App.tsx`

## Arquivos modificados
- `dashboard/src/App.tsx`

## Estratégia de implementação
Adicionar `Link2` ao import do `lucide-react`.

## Impactos esperados
- Restauração do acesso à aba de Permissões.
- Correção do crash visual (tela preta).
