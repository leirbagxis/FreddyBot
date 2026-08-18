# Plano: Remover Bordas das Divs e do Avatar no SideMenu

## Pedido do usuário
1. **Remover borda do avatar/foto de perfil**: Tirar a classe `border-2 border-accent/20` da imagem/avatar circular do usuário.
2. **Remover bordas de todas as divs internas**:
   - Tirar `border border-white/15` da div "Minha Conta".
   - Tirar `border border-white/15` da div do card de navegação.
   - Tirar `border border-white/15` do botão de alternar tema.

## Objetivo
Deixar o design do `SideMenu.tsx` totalmente clean e sem bordas (borderless), mantendo os cartões arredondados com o fundo translúcido `bg-white/5`.

## Contexto atual
- `SideMenu.tsx` possui bordas `border border-white/15` nas divs internas e `border-2 border-accent/20` ao redor do avatar do usuário.

## Arquivos analisados
- `dashboard/src/components/SideMenu.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/SideMenu.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/SideMenu.tsx`**:
   - No avatar do usuário: alterar `border-2 border-accent/20` para remover a borda.
   - Na div "Minha Conta": remover `border border-white/15`.
   - Na div do card de navegação: remover `border border-white/15`.
   - No botão "Alternar Tema": remover `border border-white/15`.

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/components/SideMenu.tsx`.
2. Executar `npm run build`.

## Riscos
- **Baixo**: Ajuste visual simples em `SideMenu.tsx`.

## Impactos esperados
- Menu lateral totalmente sem bordas, com visual limpo e contemporâneo baseado em cartões com fundo suave `bg-white/5`.

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
