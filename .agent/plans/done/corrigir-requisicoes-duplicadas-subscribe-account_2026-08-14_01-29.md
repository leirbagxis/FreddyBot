# Plano: Eliminar Requisições Duplicadas de /subscribe e /account

## Pedido do usuário
"Ele esta chamando a rota /subscribe e a rota /account duas vezes. A rota account da success true porem o data desconnected :
{
    "success": true,
    "data": {
        "status": "disconnected"
    }
}"

## Objetivo
1. **Eliminar as chamadas duplicadas às rotas `/api/subscribe` e `/api/account`**:
   - Implementar deduplicação de requisições em andamento (*in-flight promise deduplication*) nas funções `fetchSubscriptionStatus()` e `fetchAccountStatus()` em `api.ts`.
   - Passar a prop `hasPremium={hasPremiumAccess}` do `App.tsx` para o `<DashboardInicioTab>`, removendo o `useEffect` duplicado que re-buscava as rotas no carregamento da aba.
2. **Esclarecimento sobre a resposta de `/api/account`**:
   - O retorno `{ "success": true, "data": { "status": "disconnected" } }` é a resposta esperada da API quando o usuário autenticado não possui uma conta MTProto do Telegram conectada. O status passa para `"connected"` assim que a conta é vinculada.

## Contexto atual
- `App.tsx` possui um `useEffect` ao autenticar que dispara `Promise.all([fetchSubscriptionStatus(), fetchAccountStatus()])`.
- `DashboardInicioTab.tsx` possuía outro `useEffect` ao montar a aba que disparava simultaneamente o mesmo `Promise.all(...)`, resultando em 2 requisições HTTP paralelas para `/api/subscribe` e 2 requisições para `/api/account`.

## Arquivos analisados
- `dashboard/src/api.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/DashboardInicioTab.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/api.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/DashboardInicioTab.tsx`

## Estratégia de implementação

1. **`dashboard/src/api.ts`**:
   - Adicionar variáveis de promise em andamento (`subscriptionStatusPromise` e `accountStatusPromise`) em `fetchSubscriptionStatus` e `fetchAccountStatus`.
   - Caso uma requisição para a mesma rota já esteja em processamento no momento da chamada, reutilizar a mesma promise.

2. **`dashboard/src/components/DashboardInicioTab.tsx`**:
   - Atualizar a interface `DashboardInicioTabProps` para aceitar `hasPremium?: boolean`.
   - Utilizar a prop `hasPremium` direta e remover o `useEffect` interno que fazia requisições duplicadas a `/api/subscribe` e `/api/account`.

3. **`dashboard/src/App.tsx`**:
   - Passar `hasPremium={hasPremiumAccess}` na renderização de `<DashboardInicioTab />`.

4. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Riscos
- **Nenhum**: Otimização de rede sem alteração de comportamento funcional.

## Impactos esperados
- Redução de requisições de rede no carregamento da dashboard (apenas 1 chamada para `/api/subscribe` e 1 chamada para `/api/account`).
- Desempenho de carregamento mais rápido e menor carga no servidor backend.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/api.ts dashboard/src/App.tsx dashboard/src/components/DashboardInicioTab.tsx
```
