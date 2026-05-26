# Plano: atualizar-changelog-e-commit

## Pedido do usuário
Atualizar o changelog com todas as últimas entregas e criar um commit git profissional com as mudanças.

## Objetivo
Registrar no `CHANGELOG.md` um resumo fiel das alterações recentes e consolidar as mudanças atuais em um commit descritivo.

## Contexto atual
O repositório já possui várias mudanças acumuladas relacionadas a:
- logs persistentes por canal na dashboard admin
- botão de abrir logs no comando `/info`
- ajuste da mensagem de despedida ao remover o bot do canal
- correções e melhorias de PostBuilder, NewPack, links, dashboard e infraestrutura

O `CHANGELOG.md` atual termina em `1.5.1 - 2026-05-23` e ainda não documenta as entregas mais recentes.

## Arquivos analisados
- `CHANGELOG.md`
- `git status --short`
- `.agent/context.md`

## Arquivos que poderão ser modificados
- `CHANGELOG.md`
- possivelmente nenhum outro arquivo de código, apenas staging/commit

## Estratégia de implementação
1. Criar uma nova seção no topo do `CHANGELOG.md` com a próxima versão patch, documentando as entregas recentes em categorias claras.
2. Incluir somente mudanças já implementadas no código, sem inventar funcionalidades novas.
3. Revisar o diff para evitar registrar artefatos irrelevantes.
4. Fazer commit git com mensagem profissional e descritiva.

## Passos detalhados
1. Atualizar `CHANGELOG.md` com a nova entrada de versão e subseções relevantes.
2. Validar o conteúdo do changelog com uma leitura final.
3. Criar o commit git com as mudanças atuais do trabalho.
4. Informar ao usuário o hash do commit e o resumo do que foi incluído.

## Riscos
- O worktree está bem sujo e contém mudanças anteriores; é preciso evitar incluir artefatos não desejados como binários ou arquivos temporários.
- Se houver dúvidas sobre o escopo do commit, pode ser necessário separar mudanças por grupos antes de finalizar.

## Impactos esperados
- O changelog passa a refletir as entregas recentes do projeto.
- O histórico git fica consolidado em um commit único e legível.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
npm run build
```

### Testes
```bash
go test ./...
```

### Execução
```bash
# Sem execução adicional obrigatória; foco em documentação e commit.
```

## Rollback
Reverter a alteração no `CHANGELOG.md` e desfazer o commit com `git revert` se necessário.

## Observações
Ainda existe limitação local no toolchain Go (`no such tool "compile"` / `vet`) observada em validações recentes, então o commit deve ser feito mesmo que a validação Go permaneça bloqueada no ambiente.
