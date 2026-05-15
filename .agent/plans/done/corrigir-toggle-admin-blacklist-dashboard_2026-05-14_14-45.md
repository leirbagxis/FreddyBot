# Plano: corrigir-toggle-admin-blacklist-dashboard_2026-05-14_14-45.md

## Pedido do usuário
O tornar admin não muda quando o é setado de admin.

## Objetivo
Corrigir a atualização do estado local no frontend após promover/remover admin ou adicionar/remover da blacklist. O frontend está tentando acessar as propriedades `isAdmin` e `isBlacklisted` diretamente na raiz do objeto de resposta, mas elas estão dentro da propriedade `data`.

## Contexto atual
- O backend retorna um objeto `APIResponse` onde o dado real está em `data`.
- Exemplo Admin: `{ "success": true, "data": { "isAdmin": true } }`.
- No frontend (`AdminDashboard.tsx`), o código faz: `setLocalUsers(prev => prev.map(u => u.id === uid ? { ...u, is_admin: res.isAdmin } : u));`.
- Como `res.isAdmin` é `undefined`, o estado não é atualizado corretamente (ou é setado como undefined).

## Arquivos analisados
- `dashboard/src/components/AdminDashboard.tsx`
- `internal/api/controllers/adminController/getAllUserAdminController.go`

## Arquivos que poderão ser modificados
- `dashboard/src/components/AdminDashboard.tsx`

## Estratégia de implementação
Ajustar as funções `handleToggleAdmin` e `handleToggleBlacklist` para acessar as propriedades corretas dentro de `res.data`.

## Passos detalhados

1.  **Modificar `dashboard/src/components/AdminDashboard.tsx`**
    - Em `handleToggleAdmin`, usar `res.data.isAdmin`.
    - Em `handleToggleBlacklist`, usar `res.data.isBlacklisted`.

## Riscos
- **Inexistente**: Correção trivial de mapeamento de objeto.

## Como testar
1. Acessar a Dashboard Admin -> Usuários.
2. Clicar em um usuário.
3. Clicar em "Tornar Admin" ou "Add Blacklist".
4. O botão deve mudar de cor e texto imediatamente após o sucesso da requisição.

## Rollback
`git checkout dashboard/src/components/AdminDashboard.tsx`
