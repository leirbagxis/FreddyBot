# Plano: changelog commit correcoes media root

## Pedido do usuario
Atualizar o changelog e fazer um git commit das correcoes recentes.

## Objetivo
Documentar e commitar as correcoes feitas apos o commit anterior:
- preservacao de formatacao em media groups;
- correcao do alvo da legenda em albums;
- fallback convidado para a rota raiz da dashboard.

## Contexto atual
- O ultimo commit foi `3aec333 feat(admin): expand dashboard controls and post automation`.
- Ha mudancas pendentes em:
  - `dashboard/src/App.tsx`
  - `internal/telegram/events/channelPost/dispatch_telego.go`
  - `internal/telegram/events/channelPost/stage_transform_telego.go`
- Ha novos planos em `.agent/plans/done/`.
- Existe um arquivo `Release` nao rastreado, gerado por build local, que nao deve entrar no commit.
- O `CHANGELOG.md` ja possui `1.5.0 - 2026-05-23`.

## Arquivos analisados
- `CHANGELOG.md`
- `dashboard/src/App.tsx`
- `internal/telegram/events/channelPost/dispatch_telego.go`
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `git status --short`
- `git diff --stat`

## Arquivos que poderao ser modificados
- `CHANGELOG.md`

## Estrategia de implementacao
Adicionar uma secao patch `1.5.1 - 2026-05-23` no topo do changelog, descrevendo:
- media groups preservam formatacao de caption;
- a legenda final do album e aplicada na mesma midia que tinha legenda original;
- a rota `/` renderiza dashboard convidada sem canais e sem erro de init data invalido.

Depois, stagear somente os arquivos relevantes e os planos `.agent/plans/done/` correspondentes, ignorando o binario local `Release`.

## Passos detalhados

1. Inserir `1.5.1 - 2026-05-23` no `CHANGELOG.md`.
2. Rodar `npm run build` no dashboard.
3. Rodar `git diff --check`.
4. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando falha se o toolchain local continuar sem `vet`/`compile`.
5. Stagear:
   - `CHANGELOG.md`
   - `dashboard/src/App.tsx`
   - `internal/telegram/events/channelPost/dispatch_telego.go`
   - `internal/telegram/events/channelPost/stage_transform_telego.go`
   - os tres planos `.agent/plans/done/` recentes
6. Confirmar que `Release` nao esta staged.
7. Criar commit profissional, por exemplo:
   `fix(bot): preserve album captions and add guest dashboard`

## Riscos
- O binario `Release` nao rastreado deve ficar fora do commit.
- Go build/test pode continuar indisponivel pelo toolchain local.

## Impactos esperados
- Changelog reflete as correcoes recentes.
- Commit separado e focado em bugfixes.
- Worktree fica com apenas artefatos nao rastreados, se houver.

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
go build ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./...
git diff --check
```

### Execucao
```bash
make run
```

## Rollback
Usar `git revert <commit>` se for necessario desfazer apos o commit.

## Observacoes
- Nao sera usado `git clean`; o arquivo `Release` sera apenas ignorado no stage.
