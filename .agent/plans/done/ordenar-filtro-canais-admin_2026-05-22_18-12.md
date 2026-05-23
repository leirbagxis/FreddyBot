# Plano: ordenar filtro canais admin

## Pedido do usuário
Quando filtrar usuários da dashboard admin por quantidade mínima de canais, a lista deve aparecer em sequência pela quantidade de canais. Exemplo: `3`, depois `4`, depois `5`, depois `6`, e assim por diante.

## Objetivo
Ordenar os usuários filtrados por quantidade de canais crescente quando o filtro mínimo estiver ativo.

## Contexto atual
- O filtro de quantidade mínima já retorna usuários com `channels.length >= valor`.
- A lista ainda preserva a ordem original recebida do backend.
- Por isso, ao filtrar por `3`, pode aparecer um usuário com 12 canais antes de outro com 4 canais.

## Arquivos analisados
- `dashboard/src/components/AdminDashboard.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/AdminDashboard.tsx`

## Estratégia de implementação
No `useMemo` de `filteredUsers`, separar o filtro da ordenação:
- calcular `minChannelCount`;
- filtrar por nome/ID e quantidade mínima;
- se `minChannelCount` for válido, ordenar por quantidade de canais crescente;
- em empate, ordenar por nome e depois por ID para manter resultado previsível.

## Passos detalhados

1. Alterar `filteredUsers` para armazenar o resultado filtrado em uma variável.
2. Detectar se o filtro numérico está ativo.
3. Ordenar uma cópia do resultado filtrado por `channels.length`.
4. Usar nome e ID como desempate.
5. Manter a ordem atual quando o filtro de quantidade estiver vazio/inválido.
6. Rodar build da dashboard.
7. Rodar `git diff --check`.

## Riscos
- Impacto restrito ao frontend da dashboard admin.
- A ordenação muda apenas quando o filtro de quantidade está ativo.

## Impactos esperados
- Com filtro `3`, a lista aparece primeiro com usuários de 3 canais, depois 4, 5, 6 etc.
- A paginação/carregar mais continua usando a lista já ordenada.

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
Reverter a alteração no `useMemo` de `filteredUsers` em `dashboard/src/components/AdminDashboard.tsx`.

## Observações
- Não há alteração de backend ou API.
