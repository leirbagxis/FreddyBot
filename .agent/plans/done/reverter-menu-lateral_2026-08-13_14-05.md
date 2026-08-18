# Plano: Reverter Alterações do Menu Lateral (SideMenu.tsx)

## Pedido do usuário
"voce modicou o menu lateral, pqp. volta tudo como tava"

## Objetivo
Restaurar integralmente o componente `SideMenu.tsx` e suas integrações no `App.tsx` para o estado anterior às alterações recentes (restaurar o visual translúcido `bg-white/5`, a estrutura de cards e a propriedade do botão de alternância de tema no SideMenu).

## Contexto atual
- Nas últimas alterações, o `SideMenu.tsx` teve o fundo dos seus containers alterados de `bg-white/5` para `bg-white shadow-sm` e o botão de tema removido.
- O usuário solicita voltar o menu lateral exatamente como estava antes.

## Arquivos analisados
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/App.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/App.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/SideMenu.tsx`**:
   - Reverter as divs do bloco "Minha Conta" e de navegação de volta para `bg-white/5` (com hover `bg-white/10`).
   - Restaurar as props `theme: string` e `toggleTheme: () => void` na interface `SideMenuProps`.
   - Restaurar o botão de alternância de tema na parte inferior do SideMenu com os ícones `Send`, `Sun` e `Moon`.

2. **`dashboard/src/App.tsx`**:
   - Passar novamente as props `theme={theme}` e `toggleTheme={toggleTheme}` na chamada do `<SideMenu />`.

3. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Riscos
- **Nenhum**: Restauração exata do estado anterior do menu lateral a pedido expresso do usuário.

## Impactos esperados
- Menu lateral 100% restaurado ao estado e visual translúcido original (`bg-white/5`), mantendo o design aprovado anteriormente.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/components/SideMenu.tsx dashboard/src/App.tsx
```
