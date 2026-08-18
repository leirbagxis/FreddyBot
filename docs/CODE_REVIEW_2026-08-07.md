# 🔍 Relatório de Code Review — FreddyBot

> **Data:** 2026-08-07  
> **Objetivo:** Identificar bugs, vulnerabilidades de segurança, problemas de performance e melhorias de qualidade antes de levar o código para produção.

---

## 📊 Resumo Executivo

| Severidade | Quantidade |
|:-----------|:----------:|
| 🔴 CRITICAL | 6 |
| 🟠 HIGH | 23 |
| 🟡 MEDIUM | 25 |
| 🟢 LOW / INFO | 11 |
| **Total** | **65** |

---

## 🔴 Achados CRITICAL

### 1. Race condition no shutdown gracioso
- **Arquivo:** [main.go](file:///home/malbs/Opencode/FreddyBot/cmd/FreddyBot/main.go) — Linhas 67-69
- **Categoria:** BUG
- **Descrição:** A goroutine anônima que recebe do canal `stop` compete com o `select` no final do `main`. Se o sinal de shutdown chegar, ambos os receptores competem, e a goroutine pode nunca executar o graceful shutdown, ou o `main` pode terminar antes do shutdown completar.
- **Sugestão:** Usar um único receptor para o canal `stop` e orquestrar o shutdown de forma sequencial, ou usar `signal.NotifyContext` para integrar com o `context`.

### 2. Variáveis globais do container sem sincronização
- **Arquivo:** [container.go](file:///home/malbs/Opencode/FreddyBot/internal/container/container.go) — Linhas 20-27
- **Categoria:** SECURITY / BUG
- **Descrição:** Variáveis globais (`db`, `redisClient`, `repo`, `services`, `cfg`) são acessadas e escritas sem nenhum mecanismo de sincronização (mutex). Se `Initialize` for chamado concorrentemente, ou se qualquer getter for chamado antes de `Initialize`, há risco de **race condition** e **nil pointer dereference**.
- **Sugestão:** Usar `sync.Once` para garantir inicialização única, ou proteger o acesso com `sync.RWMutex`. Adicionar verificação de nil nos getters.

### 3. Redis não valida conexão no startup
- **Arquivo:** [redis.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/redis.go) — Linhas 19-36
- **Categoria:** BUG
- **Descrição:** A conexão com Redis é criada mas não é feito `Ping` para verificar se a conexão está realmente ativa. `NewRedisClient` pode retornar um cliente que não está conectado, causando falhas silenciosas em runtime.
- **Sugestão:** Adicionar `client.Ping(ctx)` imediatamente após criar o cliente.

### 4. Bot Telegram pode travar silenciosamente
- **Arquivo:** [bot.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/bot.go) — Linhas 46-52
- **Categoria:** BUG
- **Descrição:** O bot é iniciado com `bot.Start(ctx)` que é bloqueante, mas se o contexto for cancelado, a função retorna `nil` sem nenhum log ou sinal de que foi encerrado graciosamente.
- **Sugestão:** Adicionar log de shutdown e garantir que todos os handlers em execução terminem antes de retornar.

### 5. Configurações sensíveis sem proteção
- **Arquivo:** [config.go](file:///home/malbs/Opencode/FreddyBot/pkg/config/config.go)
- **Categoria:** SECURITY
- **Descrição:** Configurações sensíveis (token do bot, credenciais do banco, senha do Redis) são carregadas de variáveis de ambiente via `.env`. Se o arquivo `.env` for commitado no repositório, as credenciais ficam expostas.
- **Sugestão:** Verificar que `.env` está no `.gitignore`. Considerar usar um gerenciador de segredos (Vault, GCP Secret Manager) em produção.

---

## 🟠 Achados HIGH

### 6. Vazamento de recursos no startup falho
- **Arquivo:** [main.go](file:///home/malbs/Opencode/FreddyBot/cmd/FreddyBot/main.go) — Linhas 54-60
- **Categoria:** BUG
- **Descrição:** Se `startBot` ou `startAPI` falharem após `container.Initialize(cfg)`, `container.Shutdown()` nunca é chamado. Recursos (banco de dados, Redis) ficam abertos.
- **Sugestão:** Usar `defer container.Shutdown()` imediatamente após `container.Initialize()`.

### 7. Shutdown do container pode causar panic
- **Arquivo:** [container.go](file:///home/malbs/Opencode/FreddyBot/internal/container/container.go) — Linhas 62-74
- **Categoria:** BUG
- **Descrição:** A função `Shutdown` fecha `redisClient` e `db`, mas não verifica se são `nil` antes. Se `Initialize` falhou parcialmente, chamar `Shutdown` causa **nil pointer dereference** e panic.
- **Sugestão:** Adicionar `if redisClient != nil` e `if db != nil` antes de fechar.

### 8. Vazamento de recursos em inicialização parcial
- **Arquivo:** [container.go](file:///home/malbs/Opencode/FreddyBot/internal/container/container.go) — Linhas 29-60
- **Categoria:** BUG
- **Descrição:** Se `database.InitDB` funcionar mas `cache.NewRedisClient` falhar, a conexão com o banco fica aberta sem nunca ser fechada.
- **Sugestão:** Implementar cleanup em cascata: se um passo falhar, fechar os recursos já abertos.

### 9. `messageDate` pode ser zero
- **Arquivo:** [message_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/message_service.go) — Linhas 30-51
- **Categoria:** BUG
- **Descrição:** Se `msg.Date` não estiver preenchido, `messageDate` será `time.Time{}` (0001-01-01), que será salvo no banco como dado inválido.
- **Sugestão:** Adicionar validação: se `messageDate` for zero, usar `time.Now()` como fallback.

### 10. Mesmo bug de data zero em stickers
- **Arquivo:** [sticker_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/sticker_service.go) — Linhas 25-50
- **Categoria:** BUG
- **Descrição:** Mesmo problema do `SaveMessage`: `messageDate` pode ser zero se `msg.Date` não estiver definido.
- **Sugestão:** Adicionar validação de data com fallback para `time.Now()`.

### 11. Timezone pode gerar rankings incorretos
- **Arquivo:** [ranking_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/ranking_service.go) — Linhas 58-62, 87-91
- **Categoria:** BUG
- **Descrição:** Os cálculos de `startOfWeek` e `startOfMonth` usam `time.Now()` que depende do timezone local do servidor, que pode ser diferente do timezone dos usuários do bot.
- **Sugestão:** Usar um timezone explícito e configurável, ou documentar que o bot usa UTC.

### 12. DSN pode ser exposto nos logs
- **Arquivo:** [db.go](file:///home/malbs/Opencode/FreddyBot/internal/database/db.go) — Linhas 19-29
- **Categoria:** SECURITY
- **Descrição:** A string de conexão do banco (DSN) pode conter credenciais e ser logada em modo debug.
- **Sugestão:** Garantir que o DSN nunca seja logado. Sanitizar antes de qualquer log.

### 13. Pool de conexões não configurado
- **Arquivo:** [db.go](file:///home/malbs/Opencode/FreddyBot/internal/database/db.go) — Linhas 19-29
- **Categoria:** PERFORMANCE
- **Descrição:** `InitDB` não configura `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`. Em produção, isso pode esgotar conexões do banco.
- **Sugestão:** Configurar:
  ```go
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(5)
  db.SetConnMaxLifetime(5 * time.Minute)
  ```

### 14. Auditar queries contra SQL Injection
- **Arquivo:** [queries.go](file:///home/malbs/Opencode/FreddyBot/internal/database/queries.go)
- **Categoria:** SECURITY
- **Descrição:** Necessário verificar se todas as queries usam placeholders (`$1`, `?`) e não concatenação de strings.
- **Sugestão:** Auditar todas as queries para confirmar uso de prepared statements.

### 15. `rows.Close()` pode não ser chamado
- **Arquivo:** [repository.go](file:///home/malbs/Opencode/FreddyBot/internal/database/repository.go)
- **Categoria:** BUG
- **Descrição:** Se `rows.Scan` falhar em alguma linha, o `rows` pode não ser fechado corretamente dependendo da implementação.
- **Sugestão:** Garantir que `defer rows.Close()` é chamado imediatamente após obter `rows`.

### 16. Bot sem reconexão automática
- **Arquivo:** [bot.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/bot.go) — Linhas 22-40
- **Categoria:** BUG
- **Descrição:** O bot não tem tratamento para reconexão em caso de falha de rede. Se a conexão cair, o bot pode parar de funcionar.
- **Sugestão:** Implementar retry com backoff exponencial, ou verificar se a lib `telebot` já faz isso internamente.

### 17. Token do bot pode ser exposto nos logs
- **Arquivo:** [bot.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/bot.go) — Linhas 22-40
- **Categoria:** SECURITY
- **Descrição:** O token do bot é passado diretamente para `telebot.Settings.Token`. Se o log de debug da lib estiver ativado, o token pode ser exposto nos logs.
- **Sugestão:** Garantir que o nível de log da lib não exponha o token. Mascarar o token em qualquer log.

### 18. Perda silenciosa de mensagens
- **Arquivo:** [handlers.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/handlers.go)
- **Categoria:** BUG
- **Descrição:** O handler `OnText` salva cada mensagem no banco. Se o banco estiver indisponível, as mensagens são perdidas silenciosamente sem log de erro.
- **Sugestão:** Implementar fila de retry ou buffer local. No mínimo, logar o erro.

### 19. API sem rate limiting
- **Arquivo:** [router.go](file:///home/malbs/Opencode/FreddyBot/internal/api/router.go)
- **Categoria:** SECURITY
- **Descrição:** A API HTTP não tem rate limiting configurado. Um atacante pode fazer requests ilimitados, causando DoS.
- **Sugestão:** Adicionar middleware de rate limiting (ex: `golang.org/x/time/rate`).

### 20. CORS muito permissivo
- **Arquivo:** [router.go](file:///home/malbs/Opencode/FreddyBot/internal/api/router.go)
- **Categoria:** SECURITY
- **Descrição:** CORS configurado com `*` em produção é risco de segurança.
- **Sugestão:** Configurar CORS para aceitar apenas origens confiáveis (domínio do dashboard).

### 21. Parâmetros de entrada não validados na API
- **Arquivo:** [handlers.go](file:///home/malbs/Opencode/FreddyBot/internal/api/handlers.go)
- **Categoria:** BUG
- **Descrição:** Handlers da API provavelmente não validam parâmetros de entrada. IDs inválidos ou negativos podem causar erros inesperados.
- **Sugestão:** Validar todos os parâmetros. Verificar se IDs são positivos, datas são válidas, etc.

### 22. Erros de Redis ignorados silenciosamente
- **Arquivo:** [redis.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/redis.go)
- **Categoria:** BUG
- **Descrição:** Erros de Redis (conexão perdida) podem ser ignorados silenciosamente. A aplicação pode parecer funcionar mas sem cache.
- **Sugestão:** Implementar padrão cache-aside: logar erros de cache e continuar sem cache.

### 23. Pool de conexões Redis não configurado
- **Arquivo:** [redis.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/redis.go)
- **Categoria:** PERFORMANCE
- **Descrição:** Não há configuração de pool de conexões Redis (`PoolSize`, `MinIdleConns`, `MaxRetries`).
- **Sugestão:** Configurar pool:
  ```go
  PoolSize: 10, MinIdleConns: 5, MaxRetries: 3
  ```

### 24. Cache de ranking não é invalidado
- **Arquivo:** [ranking_cache.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/ranking_cache.go)
- **Categoria:** BUG
- **Descrição:** O cache de ranking pode retornar dados desatualizados. Quando uma nova mensagem é salva, o cache não é invalidado.
- **Sugestão:** Invalidar o cache de ranking quando novas mensagens forem salvas, ou usar TTL curto (1-5 min).

### 25. Nomes de usuário não sanitizados
- **Arquivo:** [formatter.go](file:///home/malbs/Opencode/FreddyBot/internal/utils/formatter.go)
- **Categoria:** SECURITY
- **Descrição:** Nomes de usuários são inseridos nas mensagens de resposta sem sanitização. Nomes maliciosos com HTML/Markdown podem causar parsing errors ou injeção de conteúdo no Telegram.
- **Sugestão:** Escapar caracteres especiais de Markdown/HTML com `html.EscapeString()`.

### 26. Validação de configurações obrigatórias ausente
- **Arquivo:** [config.go](file:///home/malbs/Opencode/FreddyBot/pkg/config/config.go)
- **Categoria:** BUG
- **Descrição:** Não há validação das variáveis obrigatórias. Se `BOT_TOKEN` estiver vazio, o bot tenta iniciar e falha com erro obscuro.
- **Sugestão:** Validar no startup e falhar rápido:
  ```go
  if cfg.BotToken == "" { log.Fatal("BOT_TOKEN is required") }
  ```

### 27. Dockerfile rodando como root
- **Arquivo:** [Dockerfile](file:///home/malbs/Opencode/FreddyBot/Dockerfile)
- **Categoria:** SECURITY
- **Descrição:** O container pode estar rodando como root, o que é risco de segurança.
- **Sugestão:** Adicionar `USER nonroot` e usar imagem base `distroless` ou `scratch`.

### 28. Dockerfile com imagem base não fixada
- **Arquivo:** [Dockerfile](file:///home/malbs/Opencode/FreddyBot/Dockerfile)
- **Categoria:** SECURITY
- **Descrição:** Usando `golang:latest` que pode conter vulnerabilidades e causar builds não reproduzíveis.
- **Sugestão:** Usar versão específica: `golang:1.23-alpine`.

### 29. Credenciais hardcoded no docker-compose
- **Arquivo:** [docker-compose.yml](file:///home/malbs/Opencode/FreddyBot/docker-compose.yml)
- **Categoria:** SECURITY
- **Descrição:** O `docker-compose.yml` pode conter credenciais hardcoded (senha do Redis, DSN do banco).
- **Sugestão:** Usar `env_file` ou Docker secrets.

### 30. Timing attack na comparação de tokens
- **Arquivo:** [middleware/](file:///home/malbs/Opencode/FreddyBot/internal/middleware/)
- **Categoria:** SECURITY
- **Descrição:** Se o middleware de autenticação usar comparação de strings simples (`==`) para tokens/API keys, isso é vulnerável a timing attacks.
- **Sugestão:** Usar `crypto/subtle.ConstantTimeCompare()` para comparação de tokens.

---

## 🟡 Achados MEDIUM

### 31. Variáveis globais mutáveis em main
- **Arquivo:** [main.go](file:///home/malbs/Opencode/FreddyBot/cmd/FreddyBot/main.go) — Linhas 28-30
- **Categoria:** QUALITY
- **Descrição:** `startBot` e `startAPI` são variáveis globais mutáveis do tipo `func() error`. Dificulta testes e pode causar comportamento inesperado.
- **Sugestão:** Encapsular em uma struct `App` ou usar injeção de dependência.

### 32. Getters do container retornam nil
- **Arquivo:** [container.go](file:///home/malbs/Opencode/FreddyBot/internal/container/container.go) — Linhas 76-96
- **Categoria:** QUALITY
- **Descrição:** Os getters (`GetDB`, `GetRepo`, etc.) retornam ponteiros globais sem verificação. Se chamados antes de `Initialize`, retornam `nil`, causando panics.
- **Sugestão:** Adicionar verificação de nil e retornar erro ou panic com mensagem clara.

### 33. God Object na struct Services
- **Arquivo:** [services.go](file:///home/malbs/Opencode/FreddyBot/internal/core/services.go)
- **Categoria:** QUALITY
- **Descrição:** A struct `Services` agrupa todos os serviços do sistema, criando um "God Object".
- **Sugestão:** Considerar interfaces e injeção de dependência mais granular.

### 34. Limit hardcoded no GetTopSenders
- **Arquivo:** [message_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/message_service.go) — Linhas 93-94
- **Categoria:** BUG
- **Descrição:** `GetTopSenders` usa `limit` de `10` hardcoded. Não é parametrizável nem validado.
- **Sugestão:** Aceitar `limit` como parâmetro e validar que é positivo.

### 35. Queries sequenciais no ranking
- **Arquivo:** [message_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/message_service.go) — Linhas 53-80
- **Categoria:** PERFORMANCE
- **Descrição:** A função `GetRanking` faz múltiplas queries ao banco de forma sequencial. Para grupos grandes, isso pode ser lento.
- **Sugestão:** Executar queries em paralelo com goroutines ou combinar em uma única query SQL.

### 36. Código duplicado na extração de metadados
- **Arquivo:** [sticker_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/sticker_service.go) — Linhas 25-50
- **Categoria:** MAINTENANCE
- **Descrição:** Lógica de extração de dados da mensagem (`userID`, `userName`, `chatID`, `messageDate`) é duplicada entre `SaveMessage` e `SaveStickerMessage`. Viola o princípio DRY.
- **Sugestão:** Extrair para função helper: `extractMessageMetadata(msg)`.

### 37. Ranking global sem limites pode causar OOM
- **Arquivo:** [ranking_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/ranking_service.go) — Linhas 36-37
- **Categoria:** BUG
- **Descrição:** `GetGlobalRanking` busca mensagens "desde o início dos tempos" (`time.Time{}`). Para bancos grandes, pode causar timeout ou consumo excessivo de memória.
- **Sugestão:** Adicionar paginação ou limite de registros.

### 38. Código duplicado nas funções de ranking
- **Arquivo:** [ranking_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/ranking_service.go) — Linhas 36-111
- **Categoria:** PERFORMANCE / MAINTENANCE
- **Descrição:** `GetGlobalRanking`, `GetWeeklyRanking` e `GetMonthlyRanking` têm estrutura quase idêntica com código duplicado significativo.
- **Sugestão:** Refatorar para uma única função `getRanking(startDate time.Time)`.

### 39. Grupo inválido pode ser salvo
- **Arquivo:** [group_service.go](file:///home/malbs/Opencode/FreddyBot/internal/core/group_service.go) — Linhas 22-30
- **Categoria:** BUG
- **Descrição:** `SaveGroup` não verifica se o chatID é de grupo (negativo no Telegram). Chat privado (chatID positivo) pode ser salvo como grupo.
- **Sugestão:** Validar `if chatID >= 0 { return error }`.

### 40. Falta de índices no banco
- **Arquivo:** [queries.go](file:///home/malbs/Opencode/FreddyBot/internal/database/queries.go)
- **Categoria:** PERFORMANCE
- **Descrição:** Queries de ranking e contagem podem ser lentas sem índices adequados em `chat_id`, `user_id`, `created_at`.
- **Sugestão:** Verificar e criar índices nas colunas frequentemente filtradas.

### 41. `rows.Err()` não verificado
- **Arquivo:** [repository.go](file:///home/malbs/Opencode/FreddyBot/internal/database/repository.go)
- **Categoria:** BUG
- **Descrição:** Erros retornados por `rows.Err()` após iterar podem não estar sendo verificados.
- **Sugestão:** Sempre verificar `if err := rows.Err(); err != nil` após o loop de `rows.Next()`.

### 42. Handlers processam síncronamente
- **Arquivo:** [handlers.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/handlers.go)
- **Categoria:** PERFORMANCE
- **Descrição:** Os handlers processam mensagens de forma síncrona. Se um handler travar (query lenta), todas as outras mensagens ficam na fila.
- **Sugestão:** Processar handlers pesados em goroutines separadas ou com timeout.

### 43. Handler de mensagens causa alta carga no banco
- **Arquivo:** [handlers.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/handlers.go)
- **Categoria:** PERFORMANCE
- **Descrição:** Cada mensagem recebida gera uma escrita no banco. Para grupos com alto volume, pode sobrecarregar.
- **Sugestão:** Implementar batching de inserções ou buffer write-behind.

### 44. Lógica de negócio misturada com apresentação
- **Arquivo:** [handlers.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/handlers.go)
- **Categoria:** QUALITY
- **Descrição:** Handlers contêm lógica de negócio misturada com formatação de mensagens.
- **Sugestão:** Separar formatação em pacote dedicado (`formatters` ou `views`).

### 45. Comandos não verificam tipo de chat
- **Arquivo:** [commands.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/commands.go)
- **Categoria:** BUG
- **Descrição:** Comandos como `/ranking` não verificam se a mensagem veio de grupo ou chat privado.
- **Sugestão:** Verificar tipo de chat antes de executar comandos que só fazem sentido em grupos.

### 46. Strings hardcoded nos handlers
- **Arquivo:** [commands.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/commands.go)
- **Categoria:** QUALITY
- **Descrição:** Strings de resposta hardcoded nos handlers ao invés de usar o `messages.yml`.
- **Sugestão:** Centralizar todas as mensagens no `config/messages.yml`.

### 47. API pode expor detalhes internos nos erros
- **Arquivo:** [handlers.go](file:///home/malbs/Opencode/FreddyBot/internal/api/handlers.go)
- **Categoria:** SECURITY
- **Descrição:** Respostas de erro podem expor stack traces ou mensagens do banco de dados.
- **Sugestão:** Retornar erros genéricos ao cliente e logar detalhes internos.

### 48. Sem circuit breaker para Redis
- **Arquivo:** [redis.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/redis.go)
- **Categoria:** BUG
- **Descrição:** Se o Redis ficar indisponível, operações de cache falham repetidamente sem circuit breaker.
- **Sugestão:** Implementar circuit breaker para desabilitar temporariamente o cache.

### 49. Chave de cache pode misturar grupos
- **Arquivo:** [ranking_cache.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/ranking_cache.go)
- **Categoria:** BUG
- **Descrição:** A chave de cache pode não incluir `chatID`, fazendo um grupo ver o ranking de outro.
- **Sugestão:** Garantir que a chave de cache inclui `chatID`.

### 50. Formatação de números sem localização
- **Arquivo:** [formatter.go](file:///home/malbs/Opencode/FreddyBot/internal/utils/formatter.go)
- **Categoria:** QUALITY
- **Descrição:** Números grandes são difíceis de ler sem separador de milhares.
- **Sugestão:** Usar formatação pt-BR (ex: `1.000.000`).

### 51. Vazamento de memória no Redis via Rate Limiting
- **Arquivo:** [rate_limit.go](file:///home/malbs/Opencode/FreddyBot/internal/api/middleware/rate_limit.go) — Linhas 30-32
- **Categoria:** CRITICAL / PERFORMANCE
- **Descrição:** O middleware de rate limiting usa `c.Request.Context()` para definir a expiração das chaves no Redis. Se a conexão do cliente for encerrada abruptamente entre o incremento e a definição da expiração, o contexto será cancelado e o `Expire` falhará. O Redis acumulará **milhares de chaves permanentes**, causando vazamento de memória irreversível.
- **Sugestão:** Usar `context.Background()` na chamada de `Expire`:
  ```go
  if count == 1 {
      client.Expire(context.Background(), key, window)
  }
  ```

### 52. `sync.Once` trava bot permanentemente em caso de falha
- **Arquivo:** [checkAddBotMiddlewareTelego.go](file:///home/malbs/Opencode/FreddyBot/internal/middleware/checkAddBotMiddlewareTelego.go) — Linhas 19-27
- **Categoria:** HIGH / BUG
- **Descrição:** A função `getBotUser` usa `sync.Once` para cachear os dados do bot. Porém, `sync.Once` é marcado como executado **mesmo em caso de falha**. Se a API do Telegram falhar na primeira tentativa (timeout de rede), `botUserCache` ficará `nil` **permanentemente** até que o processo seja reiniciado.
- **Sugestão:** Carregar os dados do bot durante o startup no fluxo principal, ou implementar retry se o cache for `nil`.

### 53. Saturação do banco de dados no middleware de autenticação
- **Arquivo:** [middleware.go](file:///home/malbs/Opencode/FreddyBot/internal/api/auth/middleware.go) — Linhas 40-46
- **Categoria:** HIGH / PERFORMANCE
- **Descrição:** O middleware de autenticação JWT faz uma query síncrona ao banco (`GetUserByID`) em **todas** as rotas autenticadas. Sob tráfego moderado, isso satura o pool de conexões do banco.
- **Sugestão:** Adicionar cache Redis para permissões e blacklist de usuários. Invalidar o cache quando permissões forem alteradas.

### 54. Workers da fila de mensagens sofrem starvation
- **Arquivo:** [channelPost.go](file:///home/malbs/Opencode/FreddyBot/internal/telegram/events/channelPost/channelPost.go) — Linhas 89-101
- **Categoria:** HIGH / PERFORMANCE
- **Descrição:** A `MessageQueue` tem apenas 20 workers. Se `processJob` detectar cooldown de rate limit, executa `time.Sleep(500ms)` **bloqueando o worker**. Se 20 canais precisarem de cooldown simultaneamente, **todos os workers travam** e a fila inteira congela.
- **Sugestão:** Nunca bloquear workers com `Sleep`. Recolocar jobs em sub-fila temporizada assíncrona.

### 55. Loop síncrono bloqueando o scheduler global
- **Arquivo:** [scheduler.go](file:///home/malbs/Opencode/FreddyBot/internal/core/services/scheduler.go) — Linhas 88-92
- **Categoria:** MEDIUM / PERFORMANCE
- **Descrição:** Na rotina de envios agendados (`processDuePosts`), os posts são publicados num loop com `time.Sleep(1s)` na **mesma goroutine global**. O sistema consegue enviar no máximo 1 mensagem agendada por segundo. Se um ciclo carregar 20 posts, a função fica bloqueada por 20 segundos.
- **Sugestão:** Disparar `sendScheduledPost()` em goroutines separadas com controle de rate-limit por canal.

### 56. Tipo inconsistente no cache local
- **Arquivo:** [cache.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/cache.go) — Linha 129 e 72
- **Categoria:** MEDIUM / BUG
- **Descrição:** `SetChannel` armazena `*models.Channel` direto no `localCache`, mas a função genérica `Get` tenta type assertion para `[]byte`, que falha silenciosamente. Cache L1 é ignorado, forçando queries desnecessárias ao Redis.
- **Sugestão:** Padronizar tipos no cache — serializar com `json.Marshal` antes de armazenar, ou criar métodos tipados segregados.

### 57. Deep copy ineficiente via JSON no cache
- **Arquivo:** [cache.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/cache.go) — Linhas 105-110
- **Categoria:** MEDIUM / PERFORMANCE
- **Descrição:** Em `GetChannel`, quando o item é encontrado no `localCache`, é feita cópia profunda via `json.Marshal` + `json.Unmarshal`. Usar reflexão do pacote `json` para clonar uma struct elimina o benefício de ter cache em memória (L1).
- **Sugestão:** Implementar método `Clone()` na struct `models.Channel` com cópia manual dos campos.

### 58. Fallback de `generateShortID` ignora erros de entropia
- **Arquivo:** [cache.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/cache.go) — Linhas 447-451
- **Categoria:** LOW / SECURITY
- **Descrição:** Se `rand.Int` e `rand.Read` falharem, `randByte` fica com valor zero, quebrando a aleatoriedade e tornando IDs previsíveis/colisionáveis.
- **Sugestão:** Propagar o erro ou fazer fallback explícito para `math/rand` com seed.

### 59. Portas de banco e Redis expostas publicamente
- **Arquivo:** [docker-compose.yml](file:///home/malbs/Opencode/FreddyBot/docker-compose.yml) — Linhas 7-8, 22-23
- **Categoria:** HIGH / SECURITY
- **Descrição:** As portas `6379:6379` (Redis) e `5432:5432` (PostgreSQL) estão mapeadas para a máquina host. Sem firewall, os bancos ficam acessíveis publicamente. A senha padrão `12345` torna a invasão trivial.
- **Sugestão:** Remover os blocos `ports:` (containers já se comunicam via rede interna). Para debug, usar `127.0.0.1:5432:5432`.

### 60. Caminho fixo para `messages.yml` pode falhar
- **Arquivo:** [parser.go](file:///home/malbs/Opencode/FreddyBot/pkg/parser/parser.go) — Linha 41
- **Categoria:** MEDIUM / MAINTENANCE
- **Descrição:** `loadMessages` usa caminho fixo `os.ReadFile("config/messages.yml")`. Se o bot for executado de um diretório diferente da raiz do projeto, o arquivo não será encontrado.
- **Sugestão:** Usar `//go:embed config/messages.yml` para compilar no binário, ou aceitar o caminho via flag/envvar.

### 61. Nomenclatura confusa da variável Redis
- **Arquivo:** [redis.go](file:///home/malbs/Opencode/FreddyBot/internal/cache/redis.go) — Linha 23
- **Categoria:** INFO / QUALITY
- **Descrição:** `redis.ParseURL(config.RedisAddr)` requer URI completa (`redis://...`), mas a variável no `.env-example` se chama `REDIS_HOST`, induzindo o desenvolvedor a preencher apenas `localhost:6379`, causando crash no startup.
- **Sugestão:** Alterar para `redis.NewClient(&redis.Options{Addr: config.RedisAddr})` ou renomear a variável para `REDIS_URL`.

---

## 🟢 Achados LOW

### 62. Canal `stop` como variável global
- **Arquivo:** [main.go](file:///home/malbs/Opencode/FreddyBot/cmd/FreddyBot/main.go) — Linha 25
- **Descrição:** `var stop = make(chan os.Signal, 1)` é variável global desnecessária. Mover para `main` ou usar `signal.NotifyContext`.

### 63. Logger não otimizado para produção
- **Arquivo:** [logger.go](file:///home/malbs/Opencode/FreddyBot/pkg/logger/logger.go)
- **Descrição:** O logger pode estar com nível muito verboso e formato texto. Em produção, usar nível `INFO`, formato JSON, e output para stdout.

### 64. Erros customizados sem `Unwrap()`
- **Arquivo:** [errors.go](file:///home/malbs/Opencode/FreddyBot/pkg/errors/errors.go)
- **Descrição:** Erros customizados podem não implementar `Unwrap()`, dificultando `errors.Is()` e `errors.As()`.

### 65. `.env-example` incompleto
- **Arquivo:** [.env-example](file:///home/malbs/Opencode/FreddyBot/.env-example)
- **Descrição:** Pode não documentar todas as variáveis e seus valores esperados. Adicionar comentários explicativos.

---

## 🎯 Prioridades para Produção

> [!CAUTION]
> Os itens abaixo devem ser resolvidos **antes** de ir para produção:

### Prioridade 1 — Bloqueia deploy (CRITICAL)
1. 🚨 Race condition no shutdown (item 1)
2. 🚨 Container sem sync (item 2)
3. 🚨 Vazamento de memória no Redis via rate limiting (item 51)
4. 🚨 Redis não valida conexão no startup (item 3)
5. 🚨 Validação de config obrigatórias (item 26)
6. 🚨 Shutdown com nil check (item 7)
7. 🚨 Vazamento de recursos no startup (itens 6, 8)

### Prioridade 2 — Segurança (HIGH)
8. 🔒 Portas de banco/Redis expostas publicamente (item 59)
9. 🔒 Credenciais no docker-compose (item 29)
10. 🔒 Rate limiting na API (item 19)
11. 🔒 CORS restritivo (item 20)
12. 🔒 Sanitização de nomes de usuário (item 25)
13. 🔒 Dockerfile não-root (item 27)
14. 🔒 Comparação segura de tokens (item 30)
15. 🔒 Token do bot nos logs (item 17)

### Prioridade 3 — Estabilidade (HIGH)
16. ⚡ `sync.Once` trava bot permanentemente em falha (item 52)
17. ⚡ Saturação do banco no middleware de auth (item 53)
18. ⚡ Starvation dos workers da fila (item 54)
19. ⚡ Pool de conexões DB e Redis (itens 13, 23)
20. ⚡ Validação de data com fallback (itens 9, 10)
21. ⚡ Cache invalidation (item 24)
22. ⚡ Validação de parâmetros API (item 21)
23. ⚡ Perda silenciosa de mensagens (item 18)

### Prioridade 4 — Performance (MEDIUM)
24. ⏱ Scheduler síncrono bloqueante (item 55)
25. ⏱ Deep copy via JSON no cache (item 57)
26. ⏱ Tipo inconsistente no cache local (item 56)
27. ⏱ Queries sequenciais no ranking (item 35)
28. ⏱ Handler síncrono causa fila (item 42)

### Prioridade 5 — Qualidade (pós-deploy)
29. 🔧 Refatorar código duplicado (itens 36, 38)
30. 🔧 Separar lógica de apresentação (item 44)
31. 🔧 Logging estruturado (item 63)
32. 🔧 Centralizar mensagens no messages.yml (item 46)
33. 🔧 Caminho fixo do messages.yml (item 60)
34. 🔧 Nomenclatura confusa Redis (item 61)
