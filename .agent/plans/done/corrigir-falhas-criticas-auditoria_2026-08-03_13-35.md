# Plano: corrigir-falhas-criticas-auditoria

## Pedido do usuário
Corrigir completamente as falhas confirmadas na auditoria, preservando o `docker-compose.yml` atual por ser usado somente em desenvolvimento.

## Objetivo
Eliminar os vetores confirmados de falsificação de webhook, duplicação de tarefas agendadas, acesso premium expirado, conexões MTProto inválidas, cobrança inconsistente, corpos HTTP sem limite e exposição de token/XSS no dashboard; acrescentar testes de regressão e manter build, testes e verificação estática verdes.

## Contexto atual
- O backend Go usa Gin, telego, GORM e Redis; o dashboard é React/TypeScript/Vite.
- O webhook do Telegram é aceito sem secret token e lê o corpo inteiro em memória.
- Bot e API criam containers independentes, iniciando mais de um scheduler; os agendamentos ainda não têm uma operação de claim atômica.
- Assinaturas vencidas continuam ativas até ação manual, pois a expiração e os lembretes não são executados em segundo plano.
- O caminho alternativo MTProto grava uma conta habilitada mesmo sem uma sessão autenticada.
- A criação de cobrança não guarda a intenção (quantidade de canais e valor), impedindo validar pré-checkout e ativar extras com segurança.
- O `docker-compose.yml` contém ajustes exclusivamente de desenvolvimento e ficará inalterado por decisão do usuário.

## Arquivos analisados
- `.agent/context.md`
- `README.md`
- `cmd/FreddyBot/main.go`
- `pkg/config/config.go`
- `internal/api/api.go`
- `internal/api/controllers/authController.go`
- `internal/api/controllers/channelController.go`
- `internal/api/controllers/paymentController.go`
- `internal/container/appContainer.go`
- `internal/core/services/subscriptionService.go`
- `internal/core/services/connectedAccountService.go`
- `internal/database/models/*.go`
- `internal/database/repositories/scheduledPostRepository.go`
- `internal/telegram/client.go`
- `internal/telegram/mtproto/auth/auth.go`
- `internal/telegram/mtproto/admin/auth.go`
- `dashboard/src/components/AdminNoticeTab.tsx`

## Arquivos que poderão ser modificados
- `cmd/FreddyBot/main.go`
- `pkg/config/config.go`
- `internal/api/api.go`
- `internal/api/controllers/authController.go`
- `internal/api/controllers/channelController.go`
- `internal/api/controllers/paymentController.go`
- `internal/api/middleware/*.go`
- `internal/container/appContainer.go`
- `internal/core/services/subscriptionService.go`
- `internal/core/services/connectedAccountService.go`
- `internal/core/services/paymentService.go` e testes correspondentes
- `internal/core/services/scheduledPostService.go` e testes correspondentes
- `internal/database/models/*.go`
- `internal/database/repositories/*scheduled*` e novos repositórios de intenção de pagamento
- `internal/telegram/client.go`
- `internal/telegram/mtproto/auth/auth.go`
- `internal/telegram/mtproto/admin/auth.go`
- `dashboard/src/components/AdminNoticeTab.tsx` e testes relevantes
- Arquivos de teste Go/TypeScript necessários
- `.agent/memory/memory.md`, `.agent/decisions.md` e o histórico deste plano

## Estratégia de implementação
Centralizar a criação do container e a inicialização dos trabalhadores em um único ponto de composição. Proteger todas as fronteiras externas com autenticação, validação, limites de tamanho e operações persistentes verificáveis. Para cada fluxo que altera estado financeiro ou de postagem, usar uma transição condicional no banco para que somente um processo possa concluí-la. Nenhuma alteração será feita no compose de desenvolvimento.

## Passos detalhados

1. Confirmar as assinaturas das bibliotecas instaladas (telego, Gin e GORM) e os contratos de rotas/modelos antes das alterações; registrar o comportamento pretendido nos testes de regressão.
2. Adicionar configuração de secret token do webhook e falhar de modo seguro quando webhook estiver habilitado sem secret; enviar o secret ao Telegram, limitar o corpo, comparar o cabeçalho em tempo constante, tratar JSON inválido e responder com indisponibilidade temporária quando a fila estiver saturada.
3. Tornar a criação do `AppContainer` única no `main`, injetando-a no bot e na API; separar a partida dos processos em segundo plano e protegê-la contra chamadas duplicadas.
4. Implementar claim atômico para postagens agendadas com estado de processamento e recuperação de claims abandonados; adaptar o scheduler para enviar somente itens reivindicados e testar concorrência.
5. Corrigir a validade de assinatura para considerar `CurrentPeriodEnd`, iniciar manutenção periódica de expiração e lembretes a partir do container único e testar que acessos premium vencidos deixam de ser concedidos.
6. Remover os fluxos MTProto simulados: quando a autenticação real não estiver disponível, retornar erro explícito e não persistir conta habilitada; exigir sessão válida para considerar uma conexão ativa e para conceder recursos premium.
7. Criar uma intenção de pagamento persistente e expirada para cada invoice, contendo usuário, tipo, canais extras e valor esperado; validar payload, proprietário, valor e estado no pré-checkout e no pagamento, aplicar extras uma única vez e cobrir repetição/replay por testes.
8. Aplicar limite global aos corpos JSON da API e limites explícitos aos downloads/streams de mídia; devolver erro de tamanho em vez de truncar silenciosamente.
9. Remover o JWT da resposta JSON de login, mantendo apenas o cookie HttpOnly, e escapar o conteúdo de preview de avisos antes de qualquer renderização HTML, preservando a formatação permitida.
10. Pesquisar os padrões defeituosos em toda a base (webhooks sem secret, `io.ReadAll` externo, contas sem sessão, resposta de token e renderização HTML perigosa), corrigir recorrências dentro do escopo e adicionar/ajustar testes focados.
11. Executar testes Go, `go vet`, type-check e build do dashboard; revisar as alterações, atualizar memória e decisão arquitetural, e só então mover o plano para concluído.

## Riscos
- Exigir secret token em modo webhook demanda configurar a nova variável no ambiente de produção antes de atualizar o binário; a falha intencional é preferível a aceitar atualizações forjadas.
- A introdução de estados de processamento para agendamentos precisa recuperar tarefas interrompidas para não deixá-las presas após desligamento.
- A intenção de pagamento requer migração aditiva do GORM; instalações antigas receberão novas tabelas/colunas no startup.
- A retirada de stubs MTProto impedirá conexões que antes pareciam funcionar, mas que não tinham sessão válida e não podiam executar ações reais.
- Limites de upload/download precisam comportar mídias e payloads normais sem reintroduzir leitura ilimitada.

## Impactos esperados
- Atualizações de webhook não autenticadas deixam de alcançar comandos do bot.
- Uma postagem programada passa a ser enviada uma única vez mesmo com concorrência entre workers.
- Usuários vencidos deixam de ter recursos premium e pagamentos registram exatamente o que foi comprado.
- Integrações MTProto incompletas não geram acesso premium falso.
- API e dashboard reduzem a superfície de exaustão de memória, roubo de token e XSS administrativo.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
GOCACHE=/tmp/freddybot-go-cache go build ./cmd/FreddyBot
cd dashboard && PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npm run build
```

### Testes
```bash
GOCACHE=/tmp/freddybot-go-cache go test -count=1 ./...
GOCACHE=/tmp/freddybot-go-cache go vet ./...
cd dashboard && PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npx tsc --noEmit
```

### Execução
```bash
TELEGRAM_WEBHOOK_SECRET=uma-chave-segura go run ./cmd/FreddyBot
```

## Rollback
Reverter somente o conjunto de commits/alterações deste plano. As migrações serão aditivas: tabelas ou colunas novas podem permanecer sem afetar a versão anterior; não será executada migração destrutiva. O `docker-compose.yml` não será tocado.

## Observações
- Não haverá deploy, push, alteração de segredo existente nem mudança em infraestrutura externa.
- O teste de integração real com Telegram/MTProto permanecerá dependente de credenciais e ambiente externos; os contratos locais e caminhos de erro serão exercitados automaticamente.
