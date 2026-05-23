# Plano: atualizar makefile dockerfile

## Pedido do usuario
Atualizar o Makefile e o Dockerfile, verificando tambem se as versoes de Go, Node e npm usadas no Dockerfile fazem sentido para o projeto.

## Objetivo
Deixar os comandos locais e o build Docker consistentes com o estado atual do projeto, evitando divergencia de nomes de binario e melhorando reprodutibilidade do build.

## Contexto atual
- O projeto usa backend Go em cmd/FreddyBot/main.go.
- O dashboard usa Vite em dashboard/.
- O go.mod declara Go 1.25.7.
- O Dockerfile atual usa golang:1.25.7-alpine, alinhado ao go.mod.
- A checagem publica de tags do Docker Hub em 2026-05-22 mostra que ja existe linha golang 1.26.x alpine, mas migrar para ela implicaria alterar/validar o projeto para Go 1.26. Por padrao, este plano mantem a versao do go.mod.
- O Dockerfile atual usa node:24-alpine; a checagem publica de tags do Docker Hub em 2026-05-22 mostra a linha Node 24 Alpine disponivel, incluindo Alpine 3.23.
- O ambiente local esta com Node v20.19.2 e npm 9.2.0.
- O package.json nao declara engines.
- Existe dashboard/package-lock.json, entao npm ci pode ser usado no Dockerfile.
- O Makefile foi parcialmente ajustado para gerar/executar Release.
- O alvo clean ainda remove server, deixando o binario Release para tras.
- O Dockerfile ainda gera e executa um binario chamado server.
- O Dockerfile ja injeta internal/utils.Version via GIT_HASH.

## Arquivos analisados
- Makefile
- Dockerfile
- docker-compose.yml
- go.mod
- dashboard/package.json
- dashboard/package-lock.json
- Docker Hub tags oficiais de node e golang consultadas em 2026-05-22

## Arquivos que poderao ser modificados
- Makefile
- Dockerfile

## Estrategia de implementacao
Atualizar o Makefile para:
- declarar phony targets usados;
- separar build de run para evitar que todo build execute o bot automaticamente;
- manter um alvo run explicito;
- corrigir clean para remover Release;
- adicionar alvos uteis de Docker: docker-build e docker-run;
- usar variaveis para imagem, binario e flags.

Atualizar o Dockerfile para:
- usar npm ci no build do dashboard;
- manter Go alinhado ao go.mod com golang:1.25.7-alpine;
- manter Node na linha 24 Alpine e preferir tag explicita com Alpine 3.23 se disponivel;
- compilar o mesmo binario Release;
- copiar/executar Release na imagem final;
- melhorar cache de dependencias mantendo copias separadas de manifests.

## Passos detalhados

1. Ajustar variaveis do Makefile: BINARY, IMAGE, PORT, GIT_HASH e LDFLAGS.
2. Trocar build para apenas compilar UI e servidor.
3. Criar alvo run para executar Release.
4. Manter all apontando para build.
5. Corrigir clean para remover dashboard/dist e Release.
6. Adicionar docker-build com build-arg GIT_HASH.
7. Adicionar docker-run expondo porta 7000.
8. Atualizar help com os novos comandos.
9. No Dockerfile, confirmar as tags base de Go e Node contra go.mod, package-lock.json e tags oficiais consultadas.
10. No Dockerfile, trocar npm install por npm ci.
11. No Dockerfile, compilar para Release e ajustar CMD.
12. Rodar npm run build no dashboard.
13. Tentar rodar go build ./cmd/FreddyBot/main.go.
14. Tentar rodar docker build se o ambiente permitir.
15. Rodar git diff --check.

## Riscos
- Mudar make build para nao executar automaticamente pode alterar o habito atual, mas e o comportamento mais esperado para um alvo chamado build.
- npm ci exige package-lock.json; o lockfile existe hoje, entao a troca e segura.
- Fixar uma tag mais explicita de Node pode exigir atualizacao manual futura; usar node:24-alpine3.23 preserva linha principal e base Alpine explicita.
- Migrar Docker para Go 1.26 sem alterar go.mod pode mascarar incompatibilidades; este plano mantem Go 1.25.7.
- docker build pode falhar por falta de rede ou Docker indisponivel no ambiente.
- O Go local ja apresentou falhas de toolchain (compile/vet ausentes), entao o build Go pode nao validar localmente.

## Impactos esperados
- Comandos locais ficam mais previsiveis.
- Docker e Makefile passam a usar o mesmo nome de binario.
- Build do dashboard no Docker fica mais reprodutivel.
- make clean passa a limpar o binario correto.
- Dockerfile fica com decisao explicita de versoes: Go segue go.mod; Node segue linha 24 Alpine.

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
go build -ldflags "-X github.com/leirbagxis/FreddyBot/internal/utils.Version=$(git rev-parse --short HEAD)" -o Release ./cmd/FreddyBot/main.go
docker build --build-arg GIT_HASH=$(git rev-parse --short HEAD) -t freddybot:local .
```

### Testes
```bash
git diff --check
```

### Execucao
```bash
make run
docker run --rm -p 7000:7000 --env-file .env freddybot:local
```

## Rollback
Reverter as alteracoes em:
- Makefile
- Dockerfile

## Observacoes
- O plano nao altera docker-compose.yml, porque o pedido citou apenas Makefile e Dockerfile.
- Ha dashboard/package-lock.json, entao npm ci sera usado.
- Fontes consultadas: Docker Hub tags oficiais de node e golang.
