# Plano: corrigir-lista-canais-parsing_2026-05-14_11-55.md

## Pedido do usuário
A dashboard principal ainda não mostra a lista de canais do usuário.

## Objetivo
Corrigir a extração da lista de canais no frontend. O backend retorna o array de canais diretamente no campo `.data`, e não dentro de `.data.channels`.

## Contexto atual
- No plano anterior, tentei acessar `response.data.channels`, mas como o backend retorna um array diretamente em `data`, esse campo resulta em `undefined`.
- Isso faz com que a lista de canais no frontend fique vazia.

## Arquivos analisados
- `dashboard/src/App.tsx`
- `internal/api/controllers/userController.go`

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`

## Estratégia de implementação
Ajustar a lógica de captura de canais para verificar se `response.data` é um array e usá-lo diretamente se for o caso.

## Passos detalhados

1.  **Modificar `dashboard/src/App.tsx`**
    - Alterar a linha de `channelsData` para:
      ```typescript
      const channelsData = Array.isArray(response?.data) ? response.data : (response?.data?.channels || response?.channels || []);
      ```

## Riscos
- **Nulo**: Ajuste de lógica de parsing.

## Como testar
1. Acessar a dashboard principal (`/me/channels`).
2. A lista de canais deve aparecer.

## Rollback
`git checkout dashboard/src/App.tsx`
