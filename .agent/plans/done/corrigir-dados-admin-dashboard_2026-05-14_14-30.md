# Plano: corrigir-dados-admin-dashboard_2026-05-14_14-30.md

## Pedido do usuário
A dashboard admin não está mostrando os dados, os canais e usuários.

## Objetivo
Corrigir o parsing dos dados no frontend da Dashboard Admin. O backend retorna os dados envolvidos no padrão `APIResponse[T]` (dentro de uma chave `data`), mas o frontend está tentando acessar as chaves `users` e `channels` diretamente na raiz da resposta.

## Contexto atual
- O backend (`GetAdminOverview`) retorna: `{ "success": true, "data": { "users": [...], "channels": [...] } }`.
- O frontend (`fetchAdminDashboard` em `api.ts`) faz: `const data = await apiFetch('/api/admin/overview'); return { users: data.users, channels: data.channels };`.
- Como `data.users` e `data.channels` são `undefined` (pois estão dentro de `data.data`), o dashboard fica vazio.

## Arquivos analisados
- `dashboard/src/api.ts`
- `internal/api/controllers/adminController/getAllUserAdminController.go`
- `internal/api/types/response.go`

## Arquivos que poderão ser modificados
- `dashboard/src/api.ts`

## Estratégia de implementação
Ajustar a função `fetchAdminDashboard` em `dashboard/src/api.ts` para acessar corretamente a propriedade `.data` do objeto retornado pelo `apiFetch`.

## Passos detalhados

1.  **Modificar `dashboard/src/api.ts`**
    - Na função `fetchAdminDashboard`, extrair `users` e `channels` de `response.data`.

## Riscos
- **Inexistente**: É apenas uma correção de mapeamento de JSON.

## Como testar
1. Acessar a Dashboard Admin (`/admin/dash`).
2. Verificar se a contagem de usuários e canais no topo aparece corretamente.
3. Verificar se as listas de usuários e canais são preenchidas.

## Rollback
`git checkout dashboard/src/api.ts`
