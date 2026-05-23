# Plano: atualizar changelog e commit

## Pedido do usuario
Atualizar o CHANGELOG.md e fazer um git commit profissional com todas as atualizacoes acumuladas.

## Objetivo
Documentar as mudancas recentes em uma nova entrada do changelog e criar um commit organizado, com mensagem clara e incluindo as atualizacoes reais do projeto.

## Contexto atual
- Existem varias alteracoes acumuladas desde a ultima versao do CHANGELOG.md.
- O CHANGELOG.md esta na versao 1.4.0 de 2026-05-20.
- Ha mudancas em dashboard admin, PostBuilder, NewPack, configuracoes globais, middleware, Dockerfile, Makefile, API, banco e cache.
- Ha varios planos novos em .agent/plans/done/ que representam as mudancas feitas.
- O arquivo binario rastreado FreddyBot aparece como removido.
- Existe um package-lock.json novo na raiz com conteudo vazio/minimo, provavelmente criado por acidente ao rodar npm fora de dashboard. O lockfile correto do frontend e dashboard/package-lock.json.
- O Go local esta com toolchain quebrado para build/test (`go: no such tool "compile"` e `go: no such tool "vet"`), entao essa limitacao deve ser registrada no resultado.

## Arquivos analisados
- CHANGELOG.md
- git status --short
- git diff --stat
- Makefile
- Dockerfile
- dashboard/package.json
- package-lock.json da raiz
- principais diffs de backend/frontend

## Arquivos que poderao ser modificados
- CHANGELOG.md
- package-lock.json da raiz sera removido se ainda existir

## Estrategia de implementacao
1. Criar uma nova secao no topo do CHANGELOG.md, provavelmente `1.5.0 - 2026-05-23`.
2. Resumir as mudancas por categorias profissionais:
   - Admin Dashboard
   - PostBuilder
   - NewPack
   - Configuracoes globais/admin
   - Seguranca/manutencao
   - Build/Docker
   - Fixes
3. Remover o package-lock.json da raiz por ser artefato acidental e redundante.
4. Rodar validacoes possiveis:
   - npm run build em dashboard
   - git diff --check
   - go test/go build se possivel, registrando falha do toolchain se repetir
5. Revisar `git status --short`.
6. Fazer `git add` das atualizacoes reais.
7. Criar commit com mensagem profissional, por exemplo:
   `feat(admin): expand dashboard controls and post automation`

## Passos detalhados

1. Editar CHANGELOG.md com a nova versao 1.5.0.
2. Remover package-lock.json da raiz, mantendo dashboard/package-lock.json.
3. Rodar `npm run build` em dashboard.
4. Rodar `git diff --check`.
5. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`; se falharem pelo toolchain local, nao mascarar a falha.
6. Conferir status final.
7. Stagear as mudancas reais, incluindo:
   - arquivos modificados do projeto;
   - remocao do binario FreddyBot rastreado;
   - novos planos `.agent/plans/done/`;
   - CHANGELOG.md atualizado.
8. Garantir que package-lock.json da raiz nao entre no commit.
9. Fazer o commit.

## Riscos
- O commit sera grande porque agrega varias entregas feitas nesta sessao.
- Incluir os planos `.agent` aumenta o volume do commit, mas preserva rastreabilidade conforme AGENTS.md.
- A remocao do binario FreddyBot e uma mudanca relevante; ela parece coerente com o Makefile/Dockerfile atual, que usam Release, mas sera incluida explicitamente.
- O build/test Go pode continuar indisponivel pelo toolchain local quebrado.

## Impactos esperados
- CHANGELOG.md passa a refletir as entregas recentes.
- O repositorio fica com um commit unico e descritivo das atualizacoes acumuladas.
- O package-lock.json acidental da raiz nao polui o historico.

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
Antes do commit, reverter as edicoes em CHANGELOG.md e restaurar qualquer arquivo removido por engano. Depois do commit, usar `git revert <commit>` se for necessario desfazer de forma segura.

## Observacoes
- Nao sera feito `git reset --hard` nem `git clean`.
- A remocao do package-lock.json da raiz sera feita pontualmente, porque o arquivo nao e rastreado e foi identificado como artefato acidental.
