# Plano: corrigir seguranca e estabilidade

## Pedido do usuário
Corrigir todos os problemas identificados na auditoria técnica do repositório: vulnerabilidades de autenticação/autorização, falhas funcionais do dashboard, bloqueios de testes e riscos operacionais.

## Objetivo
Eliminar os três vetores críticos de acesso indevido, tornar as rotas sensíveis verificáveis por testes de regressão e restaurar a qualidade de compilação/testes do backend e do dashboard sem alterar o desenho administrativo já existente.

## Contexto atual
- Backend Go com Gin, GORM, Redis e Bot API Telegram; regras de negócio devem permanecer nos Core Services.
- Dashboard React/TypeScript compila pelo Vite, mas `tsc --noEmit` reporta 70 erros atualmente.
- `go test ./...` está bloqueado por um teste do parser que chama `GetMessage`, função removida em favor de `GetMessageTelego`.
- A árvore de trabalho possui alterações locais de UI anteriores; elas serão preservadas e não fazem parte desta correção, exceto quando um arquivo já alterado precisar de ajuste funcional explícito.

## Arquivos analisados
- internal/api/controllers/authController.go
- internal/api/auth/middleware.go
- internal/api/routes/routes.go
- internal/api/controllers/userController.go
- internal/api/controllers/schedulerController.go
- internal/api/handlers/ping.go
- internal/core/services/channels.go
- internal/core/services/scheduler.go
- internal/database/repositories/channel.go
- pkg/parser/parser.go
- pkg/parser/parser_test.go
- dashboard/src/App.tsx
- dashboard/src/api.ts
- dashboard/src/types.ts
- dashboard/src/components/Toast.tsx
- dashboard/src/components/ScheduleTab.tsx
- dashboard/src/components/UserTemplatesManager.tsx
- dashboard/src/components/ButtonGrid.tsx
- dashboard/src/components/ConnectedAccountCard.tsx

## Arquivos que poderão ser modificados
- internal/api/controllers/authController.go
- internal/api/controllers/userController.go
- internal/api/controllers/schedulerController.go
- internal/api/handlers/ping.go
- internal/core/services/channels.go
- internal/core/services/scheduler.go
- internal/api/types/user.go
- arquivos de teste novos ou existentes em `internal/api/controllers/`, `internal/api/auth/`, `internal/core/services/` e `pkg/parser/`
- pkg/parser/parser_test.go
- dashboard/src/App.tsx
- dashboard/src/api.ts
- dashboard/src/types.ts
- dashboard/src/components/ScheduleTab.tsx
- dashboard/src/components/UserTemplatesManager.tsx
- dashboard/src/components/ConnectedAccountCard.tsx
- dashboard/src/components/AdminLogsTab.tsx
- dashboard/src/components/AdminNoticeTab.tsx
- dashboard/src/components/ButtonGrid.tsx ou `dashboard/src/components/ui/button.tsx`, conforme a API existente determine o ajuste mínimo
- dashboard/tsconfig*.json e dashboard/vite.config.ts, somente se necessários para corrigir compatibilidade de compilação já observada
- .agent/memory/memory.md
- .agent/decisions.md

## Estratégia de implementação
1. Corrigir confiança em dados controlados pelo cliente: o login validará o `user.id` extraído do JSON assinado; transferência e agendamento usarão o usuário autenticado e a propriedade real do canal, com exceção explícita apenas para admin/owner se esse for o comportamento atual desejado.
2. Levar a autorização para a camada de serviço quando ela for regra de negócio, mantendo os controllers como validadores de entrada/saída.
3. Adicionar limites ao endpoint anônimo de logs e preservar a captura de falhas pré-login.
4. Reparar os contratos TypeScript e usos em runtime; não apenas silenciar erros de compilação.
5. Criar testes de regressão para os vetores críticos e rodar as verificações completas.

## Passos detalhados

1. Ajustar `Login` para decodificar o campo `user` do initData validado e comparar `tgUser.ID == req.UserID`; adicionar casos de igualdade, ID por prefixo e payload inválido.
2. Remover `oldOwnerId` como fonte de autoridade da transferência. Obter `userID` e `role` do contexto, validar posse do canal para usuário comum e limitar a transferência de admin/owner à política existente; atualizar request, controller e serviço/repositório conforme necessário. Criar testes de posse, bypass administrativo autorizado e tentativa por terceiro.
3. Validar propriedade/role do canal antes de criar um agendamento e reforçar a regra no serviço para evitar futuros chamadores inseguros. Criar teste que prove que o scheduler não pode ser criado para canal alheio.
4. Limitar o corpo de `/api/log/client-error`, registrar de modo truncado/estruturado e aplicar proteção simples contra abuso compatível com o middleware atual.
5. Corrigir o teste do parser para chamar a API pública atual e evitar dependência de diretório de trabalho/global `sync.Once` que torne a suíte instável.
6. Corrigir os contratos de resposta API no dashboard e os tipos ausentes/inconsistentes (`ButtonType`, campos de conta/usuário e valores anuláveis).
7. Corrigir falhas de runtime do dashboard: invocação de toast, reordenação de templates, uso inválido de `asChild` e demais erros de TypeScript reportados. Remover imports e variáveis inativas apenas nos arquivos tocados pela correção.
8. Avaliar e corrigir as opções de TypeScript apenas se a biblioteca/target justificarem; não aumentar o alvo para ocultar incompatibilidades sem compatibilidade confirmada.
9. Executar formatação Go, testes unitários e verificação de segurança/qualidade: `go test ./...`, `go vet ./...`, `npx tsc --noEmit` e `npm run build`.
10. Atualizar memória e decisão técnica com a regra de que identidade e propriedade nunca vêm do payload do cliente; mover o plano para `done` após verificação.

## Riscos
- Mudanças em login podem invalidar clientes que enviem IDs ou JSON fora do formato Telegram esperado; os testes cobrirão o formato assinado válido.
- Definir bypass de admin/owner para transferência/agendamento exige preservar a política de administração existente e não conceder privilégios ao usuário comum.
- A correção de tipos do dashboard pode revelar inconsistências adicionais de payload hoje mascaradas por `any`.
- Limite agressivo de logs pode descartar diagnósticos úteis; será escolhido um limite suficiente para erros de navegador comuns.

## Impactos esperados
- Usuários não poderão autenticar, transferir canais ou agendar posts em nome de terceiros.
- O endpoint público de diagnóstico deixará de aceitar corpos ilimitados.
- O dashboard passará no typecheck e deixará de falhar ao exibir toasts e mover botões de templates.
- A suíte Go voltará a executar integralmente.

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
PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npm run dev
go run ./cmd/FreddyBot
```

## Rollback
Reverter somente os arquivos listados neste plano para o estado anterior às alterações, preservando as mudanças locais de UI que já existiam. As decisões e a memória podem receber uma nota de reversão sem apagar o histórico.

## Observações
- Não haverá deploy, push, mudança de infraestrutura, migração destrutiva nem alteração de credenciais.
- Achados de desempenho do bundle serão tratados somente no necessário para passar as verificações; code splitting amplo fica fora deste escopo por não ser correção de falha funcional ou segurança.
