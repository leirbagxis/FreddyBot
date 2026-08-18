# Plano: auditar-codigo-completo

## Pedido do usuário
Realizar uma análise de todo o código para encontrar possíveis falhas.

## Objetivo
Produzir uma auditoria técnica baseada em evidências, cobrindo segurança, confiabilidade, corretude, concorrência, API, persistência e frontend, sem alterar o código de produção nesta etapa.

## Contexto atual
- Backend Go com Gin, telego, GORM/SQLite em desenvolvimento e Redis.
- Dashboard React/TypeScript/Vite integrado ao binário.
- Auditorias anteriores já corrigiram autorização de canal, agendamento, limites de payload e exposição de URL com token; esta revisão deve procurar recorrências e novos vetores.

## Arquivos analisados
- .agent/context.md
- README.md
- go.mod
- internal/api/routes/routes.go
- internal/api/auth/*.go
- internal/core/services/**/*.go
- internal/database/repositories/**/*.go
- internal/telegram/**/*.go
- dashboard/src/**/*

## Arquivos que poderão ser modificados
- Nenhum arquivo de produção nesta etapa.
- Este plano poderá ser movido para `.agent/plans/done/` após a entrega do relatório.
- `.agent/memory/memory.md` e `.agent/decisions.md` somente se a auditoria revelar uma decisão ou risco estrutural que precise ser preservado.

## Estratégia de implementação
Executar revisão estática orientada a riscos, complementar com verificações automatizadas existentes e rastrear cada achado até arquivo, fluxo afetado e condição de exploração. Separar defeitos confirmados de hipóteses que exigem reprodução.

## Passos detalhados

1. Ler documentação, configuração, ponto de entrada e rotas para mapear superfícies expostas.
2. Revisar autenticação, autorização, manipulação de credenciais, uploads/downloads e limites de requisição.
3. Examinar serviços, repositórios e fluxos Telegram para controle de posse, validação, transações, concorrência e tratamento de erros.
4. Examinar dashboard e cliente HTTP para autenticação, XSS, estados assíncronos, tipagem e compatibilidade com a API.
5. Executar testes, análise estática, type-check e build quando possível.
6. Pesquisar padrões recorrentes dos achados e entregar relatório priorizado, sem implementar correções.

## Riscos
- Parte dos fluxos depende de Telegram, Redis e credenciais indisponíveis localmente; esses trechos podem permanecer como risco não reproduzido.
- O worktree possui alterações locais anteriores; a auditoria será somente leitura para preservá-las.

## Impactos esperados
- Nenhuma alteração funcional nesta etapa.
- Lista priorizada de falhas confirmadas e recomendações verificáveis para uma futura correção.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npm run build
```

### Testes
```bash
GOCACHE=/tmp/freddybot-go-cache go test ./...
GOCACHE=/tmp/freddybot-go-cache go vet ./...
PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npx tsc --noEmit
```

### Execução
```bash
go run ./cmd/FreddyBot
```

## Rollback
Não há mudança de produção para reverter. O único artefato novo é este plano, que será preservado no histórico.

## Observações
- Nenhum deploy, commit, push, alteração de segredo ou acesso a serviços externos será realizado.
- A auditoria não alegará cobertura de comportamento que dependa de credenciais ou integrações não disponíveis.
