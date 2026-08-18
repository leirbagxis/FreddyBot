# Plano: corrigir vulnerabilidades da auditoria

## Pedido do usuário
Corrigir todas as falhas confirmadas na auditoria mais recente: vazamento do token Telegram por redirecionamento de foto, erro de agendamento, leituras de corpo sem limite, crescimento do rate limiter, enumeração de perfis e validação fraca de horários.

## Objetivo
Eliminar a exposição de credenciais e os riscos de acesso/consumo de recursos, assegurar que alterações de horário não antecipem publicações e tornar as regras novas cobertas por testes de regressão.

## Contexto atual
- A rota `/api/channel/:channelId/photo` herda autenticação JWT, mas não a autorização do canal e redireciona para `telego.Bot.FileDownloadURL`, cuja URL contém o token do bot.
- O scheduler persiste `next_run_at` como o horário atual quando recebe apenas `scheduleTime`.
- Dois handlers de permissões fazem `io.ReadAll` sem limite e sem necessidade funcional.
- O rate limiter de logs mantém chaves de IP expiradas para sempre.
- A consulta de usuário numérica usa `GetChat` do Telegram e retorna o objeto bruto para qualquer usuário autenticado.
- `parseHHMM` aceita valores inválidos e os normaliza silenciosamente via `time.Date`.

## Arquivos analisados
- internal/api/routes/routes.go
- internal/api/controllers/captionController.go
- internal/api/controllers/permissionController.go
- internal/api/controllers/schedulerController.go
- internal/api/controllers/userController.go
- internal/api/handlers/ping.go
- internal/core/services/scheduler.go
- internal/database/repositories/scheduled_post.go
- dashboard/src/components/DashboardInicioTab.tsx
- /home/gabriel/go/pkg/mod/github.com/mymmrac/telego@v1.9.0/bot.go

## Arquivos que poderão ser modificados
- internal/api/routes/routes.go
- internal/api/controllers/captionController.go
- internal/api/controllers/permissionController.go
- internal/api/controllers/schedulerController.go
- internal/api/controllers/userController.go
- internal/api/handlers/ping.go
- internal/core/services/scheduler.go
- internal/database/repositories/scheduled_post.go
- dashboard/src/components/DashboardInicioTab.tsx, somente se a foto protegida exigir ajuste de carregamento
- testes novos ou existentes em `internal/api/controllers/`, `internal/api/handlers/` e `internal/core/services/`
- .agent/memory/memory.md
- .agent/decisions.md

## Estratégia de implementação
1. Substituir o redirecionamento externo da foto por um proxy binário no servidor, sem expor a URL que contém o token; colocar a rota sob `AuthorizeChannel` para usuário comum, preservando o cookie enviado por `<img>` em mesma origem.
2. Centralizar cálculo/validação de horário no `SchedulerService`. Edições de apenas `scheduleTime` recalcularão o próximo disparo conforme o tipo do agendamento, sem usar `time.Now()` como valor persistido.
3. Remover as leituras redundantes de corpo dos handlers de permissões e adicionar limite explícito onde for necessário.
4. Limpar janelas expiradas do rate limiter e limitar a cardinalidade da estrutura em memória.
5. Restringir busca numérica de usuários à base local e devolver DTO, preservando a busca de destinatário que já iniciou o bot.

## Passos detalhados

1. Mover a rota de foto para o grupo protegido por canal e alterar o controller para buscar o arquivo Telegram no servidor, com timeout, limite de tamanho e `Content-Type` seguro; nunca enviar a URL de download ao navegador.
2. Criar testes que verifiquem que a resposta de foto não possui `Location` com token e que o middleware de canal protege a rota.
3. Remover `io.ReadAll` e a reconstrução de `Request.Body` dos dois handlers de permissões; criar teste de payload grande para confirmar rejeição sem cópia integral em memória.
4. Adicionar coleta de entradas expiradas e teto de clientes ao rate limiter; testar expiração e limite sem crescer o mapa indefinidamente.
5. Alterar `GetUserInfo` para consultar usuários persistidos e responder somente o DTO público; cobrir ID não cadastrado e usuário cadastrado.
6. Fazer o parser de `HH:MM` retornar erro para formato/range inválido, validar horário nos fluxos de criação e edição e cobrir casos como `99:99`.
7. Refatorar `UpdateScheduleTime` para atualizar apenas os campos solicitados e recalcular `next_run_at` de agendamentos recorrentes quando necessário; testar que uma edição de horário não agenda para o instante atual.
8. Executar formatação, `go test ./...`, `go vet ./...`, `npx tsc --noEmit` e `npm run build`; atualizar documentação e mover o plano para concluído.

## Riscos
- A foto de canal depende de cookie de mesma origem para `<img>`; será validado no fluxo atual da dashboard antes de remover a rota global.
- O proxy de mídia precisa impor tamanho e timeout para não trocar vazamento de token por consumo excessivo de memória.
- Restringir busca numérica a usuários locais pode rejeitar destinatários que nunca interagiram com o bot, comportamento coerente com a regra já comunicada na transferência.
- Recalcular recorrência requer preservar os tipos `once`, `daily`, `weekly` e `queue` existentes.

## Impactos esperados
- O token Telegram não será exposto ao navegador e a rota de foto respeitará a propriedade do canal.
- Horários inválidos serão rejeitados e edições não dispararão posts prematuramente.
- Endpoints de permissões e logs terão uso de memória limitado.
- Dados de perfis de terceiros deixarão de ser enumeráveis pela API.

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
go test ./...
go vet ./...
PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npx tsc --noEmit
```

### Execução
```bash
go run ./cmd/FreddyBot
```

## Rollback
Reverter exclusivamente os arquivos listados para o estado anterior às correções, mantendo as alterações locais de UI não relacionadas. Se a credencial Telegram já tiver sido exposta, o rollback de código não substitui a rotação do token.

## Observações
- Não haverá deploy, commit, push ou alteração automática de credenciais.
- A rotação de `TELEGRAM_BOT_TOKEN`, caso a rota tenha alcançado produção, exige operação de segredo e redeploy autorizados separadamente.
