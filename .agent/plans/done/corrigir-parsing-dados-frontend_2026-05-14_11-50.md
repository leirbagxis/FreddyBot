# Plano: corrigir-parsing-dados-frontend_2026-05-14_11-50.md

## Pedido do usuário
A dashboard não lista os canais e redireciona ao acessar um canal específico, mesmo com o backend retornando os dados corretamente.

## Objetivo
Corrigir o parsing das respostas da API no frontend para acessar os dados através do campo `.data`, que é o padrão utilizado pelo backend (`APIResponse[T]`).

## Contexto atual
- O backend utiliza um wrapper `APIResponse` que coloca o conteúdo real dentro de uma chave `data`.
- O frontend em `App.tsx` está tentando acessar `response.channels` ou `dashRes.user` diretamente na raiz do objeto retornado.
- Como esses campos não existem na raiz (estão dentro de `data`), o frontend acha que não há canais ou dados, disparando redirecionamentos ou estados vazios.

## Arquivos analisados
- `dashboard/src/App.tsx`
- `internal/api/types/response.go`

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`

## Estratégia de implementação
1.  **Ajustar `isChannelsRoute` logic**: Garantir que `response.data.channels` seja usado.
2.  **Ajustar `channelId` logic**: Garantir que `dashRes.data` seja usado para popular o estado do dashboard.

## Passos detalhados

1.  **Modificar `dashboard/src/App.tsx`**
    - Na seção de `fetchUserChannels()`, mudar para: `const response = await fetchUserChannels();` e acessar `response.data`.
    - Na seção de `fetchDashboardData(channelId)`, mudar para: `const dashRes = await fetchDashboardData(channelId);` e usar `dashRes.data`.

## Riscos
- **Baixo**: Apenas ajuste de caminho de objeto JSON.

## Como testar
1. O dashboard deve carregar os canais na lista inicial.
2. Ao clicar em um canal, ele deve abrir as configurações sem redirecionar para a lista.

## Rollback
`git checkout dashboard/src/App.tsx`
