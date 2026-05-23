# Plano: dashboard raiz convidado

## Pedido do usuario
Ao acessar a rota padrao da dashboard, por exemplo `localhost:7000/`, nao deve aparecer erro de init data invalido. Em vez disso, deve renderizar a dashboard/lista de canais como convidado, sem canais, com nome `Convidado`.

## Objetivo
Criar um fallback controlado para a rota raiz `/`, fora do Telegram, renderizando a mesma experiencia de lista de canais vazia sem exigir autenticação Telegram.

## Contexto atual
- A dashboard chama `login(initData, userID)` logo ao carregar.
- Fora do Telegram, `initData` fica vazio e `userID` fica `0`.
- Em producao, isso cai em erro de autenticacao (`init data invalido`).
- A rota `/me/channels` ja renderiza a lista de canais quando autenticada.
- A rota `/` hoje nao e considerada `isChannelsRoute()`, entao mesmo que houvesse dados, algumas partes da tela de canais ficam ocultas.
- O card vazio ja existe quando nao ha `channel`, mas o header/greeting de canais so aparece em `/me/channels`.

## Arquivos analisados
- `dashboard/src/App.tsx`
- `dashboard/src/mockData.ts`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`

## Arquivos que poderao ser modificados
- `dashboard/src/App.tsx`

## Estrategia de implementacao
Adicionar uma funcao `isRootRoute()` e tratar `/` como modo convidado:
- se a rota for `/`, autenticar localmente como convidado sem chamar `login`;
- preencher `data` com `channel: null` e usuario `Convidado`;
- usuario convidado tera `channels: []`;
- considerar a rota raiz como tela de canais para renderizar o greeting e o titulo `Canais Encontrados`;
- manter `/admin/dash`, `/me/channels` e `/dashboard/:id` com o fluxo atual de autenticacao.

## Passos detalhados

1. Criar helper `isRootRoute()` em `App.tsx`.
2. Criar um objeto local de usuario convidado no efeito de autenticacao quando `isRootRoute()` for verdadeiro.
3. Antes de chamar `login`, se for rota raiz:
   - setar `authState` como `authenticated`;
   - setar `data` com `channel: null` e user convidado sem canais;
   - setar `loading` como `false`;
   - retornar.
4. Ajustar a booleana de renderizacao para tratar `isRootRoute()` como tela de canais.
5. Usar essa booleana para exibir greeting, titulo e lista vazia na rota raiz.
6. Garantir que o nome exibido seja `Convidado`.
7. Rodar `npm run build` no dashboard.
8. Rodar `git diff --check`.

## Riscos
- Impacto restrito ao frontend da dashboard.
- A rota raiz nao tera permissoes para editar canais, porque nao havera canais carregados.
- Se no futuro a rota `/` tiver outro uso, esse fallback devera ser revisto.

## Impactos esperados
- Acessar `localhost:7000/` renderiza uma tela amigavel de canais vazia.
- Nao aparece mais erro de init data invalido na rota raiz.
- O restante das rotas continua exigindo autenticacao Telegram.

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

### Execucao
```bash
make run
```

## Rollback
Reverter as alteracoes em `dashboard/src/App.tsx`.

## Observacoes
- Nao sera alterada a API de autenticacao.
- O fallback sera exclusivo para a rota `/`.
