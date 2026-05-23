# Plano: filtro minimo canais admin

## Pedido do usuário
Alterar o filtro de quantidade de canais na dashboard admin para mostrar usuários com a quantidade informada ou mais. Exemplo: ao digitar `3`, mostrar usuários com 3, 4, 5 canais em diante.

## Objetivo
Trocar a comparação exata do filtro de canais por uma comparação de mínimo, mantendo o comportamento de busca por nome/ID.

## Contexto atual
- O filtro fica em `dashboard/src/components/AdminDashboard.tsx`.
- O estado usado é `adminChannelCountFilter`.
- Atualmente o filtro usa `(u.channels?.length || 0) === parseInt(adminChannelCountFilter, 10)`.
- Isso retorna somente usuários com exatamente a quantidade digitada.

## Arquivos analisados
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/App.tsx`
- `dashboard/src/types.ts`

## Arquivos que poderão ser modificados
- `dashboard/src/components/AdminDashboard.tsx`

## Estratégia de implementação
Alterar o cálculo `matchesCount` para comparar se a quantidade de canais do usuário é maior ou igual ao número informado. Também ajustar o placeholder do input para indicar que o filtro é "a partir de" uma quantidade.

## Passos detalhados

1. Em `AdminDashboard.tsx`, calcular o valor numérico do filtro com `parseInt`.
2. Quando o filtro estiver vazio ou inválido, manter todos os usuários.
3. Quando houver número válido, retornar usuários com `channels.length >= filtro`.
4. Ajustar o placeholder para algo como `Mostrar a partir de qtd. de canais`.
5. Rodar build da dashboard.
6. Rodar `git diff --check`.

## Riscos
- Impacto restrito ao frontend da dashboard admin.
- Usuários podem perceber mudança de semântica do filtro; o placeholder deve reduzir ambiguidade.

## Impactos esperados
- Digitar `3` mostra usuários com 3 ou mais canais.
- Busca por nome/ID continua funcionando junto com o filtro.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
cd dashboard
npm run build
```

### Testes
```bash
git diff --check
```

### Execução
```bash
cd dashboard
npm run dev
```

## Rollback
Reverter a alteração em `dashboard/src/components/AdminDashboard.tsx`, voltando a comparação para igualdade exata.

## Observações
- Não há necessidade de alterar backend ou API.
