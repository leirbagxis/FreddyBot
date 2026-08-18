# Plano: Correções da Auditoria de Produção

## Pedido do usuário
Corrigir todos os 35 achados identificados na auditoria de produção do FreddyBot, priorizando itens críticos e de alta severidade.

## Objetivo
Tornar o FreddyBot production-ready corrigindo bugs, falhas de segurança, problemas de performance e melhorias de qualidade de código identificados na auditoria.

## Contexto atual
- Backend em Go (Gin + GORM + telego)
- Banco: PostgreSQL (prod) / SQLite (dev)
- Cache: Redis + go-cache (L1 in-memory)
- MTProto via gotd/td para contas conectadas
- API REST + Telegram Bot (polling/webhook)
- Dashboard React embutido
- Pagamentos via Telegram Stars

## Arquivos analisados
- cmd/FreddyBot/main.go
- internal/container/appContainer.go
- internal/database/database.go
- internal/cache/*.go
- internal/core/services/*.go
- internal/api/**/*.go
- internal/telegram/**/*.go
- internal/middleware/*.go
- pkg/config/config.go
- pkg/logger/logger.go
- pkg/errors/errors.go
- Dockerfile, docker-compose.yml, .env-example, .gitignore

## Arquivos que poderão ser modificados

### Fase 1 — Críticos
- docker-compose.yml
- internal/database/database.go
- internal/cache/redis.go
- internal/container/appContainer.go
- cmd/FreddyBot/main.go

### Fase 2 — Segurança & Auth
- internal/api/auth/middleware.go
- internal/api/auth/jwt.go
- internal/api/routes/routes.go
- internal/cache/cache.go
- pkg/config/config.go

### Fase 3 — Performance & Cache
- internal/telegram/executor/factory.go
- internal/cache/cache.go
- internal/cache/local.go
- internal/middleware/checkAddBotMiddlewareTelego.go
- internal/core/services/scheduler.go

### Fase 4 — Qualidade de código
- internal/database/database.go (typo Maintence)
- pkg/config/config.go (typo SecreteKey)
- internal/api/auth/jwt.go (typo secreteKey)
- internal/api/auth/signature.go (typo secreteKey)
- internal/core/services/admin_account.go (logs sensíveis)
- internal/core/services/subscription_service.go (logs sensíveis)

### Fase 5 — Infraestrutura
- .gitignore
- Dockerfile
- docker-compose.yml

### Fase 6 — Melhorias futuras (documentar)
- pkg/logger/logger.go
- testes de integração (novos)

## Estratégia de implementação

O plano está dividido em **6 fases** ordenadas por prioridade e risco. Cada fase é independente e pode ser commitada separadamente. As fases 1-3 são **obrigatórias** antes de produção. As fases 4-5 são **fortemente recomendadas**. A fase 6 é **desejável** mas pode ser feita depois.

---

## Passos detalhados

### ═══ FASE 1: CRÍTICOS (5 itens) ═══

#### Passo 1.1 — Remover credenciais hardcoded do docker-compose.yml
**Achado #1 | CRÍTICO**

**Arquivo:** `docker-compose.yml`

**Alterações:**
```yaml
# ANTES
environment:
  POSTGRES_USER: postgres
  POSTGRES_PASSWORD: 12345
  POSTGRES_DB: freddybot

# DEPOIS
environment:
  POSTGRES_USER: ${POSTGRES_USER:-postgres}
  POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?Defina POSTGRES_PASSWORD}
  POSTGRES_DB: ${POSTGRES_DB:-freddybot}
```

Também adicionar `networks: app-network` ao serviço postgres (achado #15).

**Também:** Atualizar `.env-example` para incluir:
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=
POSTGRES_DB=freddybot
```

---

#### Passo 1.2 — Substituir panic() por retorno de erro na inicialização
**Achado #2 | CRÍTICO**

**Arquivo:** `internal/database/database.go`

**Alterações:**
1. Alterar assinatura de `InitDB()` para retornar `(*gorm.DB, error)`:
```go
func InitDB() (*gorm.DB, error) {
    // ...
    db, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("falha ao conectar no banco: %w", err)
    }
    // ...
    if err = db.AutoMigrate(...); err != nil {
        return nil, fmt.Errorf("falha na migração: %w", err)
    }
    // ...
    return db, nil
}
```

2. Alterar `initServerConfig` e `seedPremiumFeatures` para propagar erro (já retornam error, basta não usar panic no caller)

**Arquivo:** `internal/cache/redis.go`

3. Alterar `GetRedisClient()` para retornar `(*redis.Client, error)`:
```go
func GetRedisClient() (*redis.Client, error) {
```
Isso requer ajustar todos os callers de `GetRedisClient()` — ou manter o singleton mas logar o erro e retornar um client que falha gracefully.

**Decisão pragmática:** Dado que `GetRedisClient()` é chamado em ~30 lugares, a abordagem mais segura é:
- Manter `GetRedisClient()` com signature atual
- Substituir `panic()` por `log.Fatalf()` que ao menos flusha os logs antes de morrer
- Em `InitDB()`, retornar error normalmente

**Arquivo:** `cmd/FreddyBot/main.go`

4. Ajustar `main()` para tratar o erro de `InitDB()`:
```go
db, err := database.InitDB()
if err != nil {
    logger.Error("APP", "Falha ao iniciar banco: %v", err)
    os.Exit(1)
}
```

---

#### Passo 1.3 — Broadcast workers com contexto cancelável
**Achado #3 | CRÍTICO**

**Arquivo:** `internal/container/appContainer.go`

**Alterações:**
1. Alterar `startBroadcastWorkers` para receber `ctx`:
```go
func (c *AppContainer) startBroadcastWorkers(ctx context.Context, workerCount int) {
    for i := 0; i < workerCount; i++ {
        go c.broadcastWorker(ctx)
    }
}
```

2. Alterar `broadcastWorker` para respeitar `ctx.Done()`:
```go
func (c *AppContainer) broadcastWorker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case job, ok := <-c.BroadcastQueue:
            if !ok {
                return
            }
            // ... usar ctx ao invés de context.Background() no SendPhoto/SendMessage
            _, err = c.TelegoBot.SendPhoto(ctx, params)
            // ...
        }
    }
}
```

3. No `StartBackground`, passar `ctx`:
```go
c.startBroadcastWorkers(ctx, 5)
```

---

#### Passo 1.4 — PRAGMA condicionado ao dialector, não ao ambiente
**Achado #4 | CRÍTICO**

**Arquivo:** `internal/database/database.go`

**Alterações:**
```go
// ANTES
if config.AppEnv == "dev" {
    db.Exec("PRAGMA foreign_keys = ON;")
}

// DEPOIS — Detectar o dialector usado
isSQLite := false
switch dbDriver {
case "sqlite":
    isSQLite = true
case "":
    if config.AppEnv == "dev" {
        isSQLite = true
    }
}
if isSQLite {
    db.Exec("PRAGMA foreign_keys = ON;")
}
```

Usar uma variável local `isSQLite` definida no switch de seleção do dialector.

---

#### Passo 1.5 — Migrações PostgreSQL condicionais
**Achado #5 | CRÍTICO**

**Arquivo:** `internal/database/database.go`

**Alterações:**
```go
// Só executar migrações PostgreSQL-specific se o dialector for Postgres
if !isSQLite {
    db.Exec(`DO $$ BEGIN ... END $$;`)
    // ... todas as migrações com sintaxe PostgreSQL
}
```

Reutilizar a variável `isSQLite` do passo anterior.

---

### ═══ FASE 2: SEGURANÇA & AUTH (6 itens) ═══

#### Passo 2.1 — Corrigir comparação de erro Redis inconsistente
**Achado #8 | ALTO**

**Arquivo:** `internal/cache/cache.go`

**Alterações:** Substituir **todas** as instâncias de:
```go
if err.Error() == "redis: nil" {
```
por:
```go
if err == redis.Nil {
```

**Linhas afetadas:** 54, 175, 204, 239, 271, 300, 334, 370

---

#### Passo 2.2 — Remover auth middleware duplicado
**Achado #9 | ALTO**

**Arquivo:** `internal/api/routes/routes.go`

**Alterações:** Remover a linha 172:
```go
// REMOVER ESTA LINHA:
accountRoutes.Use(auth.AuthMiddlewareJWT(c))
```

O middleware já está aplicado no grupo `api` (linha 49).

---

#### Passo 2.3 — Type assertion segura no RequireRole
**Achado #14 | ALTO**

**Arquivo:** `internal/api/auth/middleware.go`

**Alterações:**
```go
// ANTES
role := userRole.(Role)

// DEPOIS
role, ok := userRole.(Role)
if !ok {
    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
        "success": false,
        "message": "Tipo de cargo inválido",
    })
    return
}
```

Fazer o mesmo na `AuthorizeChannel` (linhas 113-114).

---

#### Passo 2.4 — Timing-safe comparison no ValidateTelegramInitData
**Achado #31 | BAIXO (mas segurança)**

**Arquivo:** `internal/api/auth/middleware.go`

**Alterações:**
```go
// ANTES
if expectedHash != hash {

// DEPOIS
if !hmac.Equal([]byte(expectedHash), []byte(hash)) {
```

Já importa `crypto/hmac` no topo do arquivo.

---

#### Passo 2.5 — Proteção extra para StarsTestMode
**Achado #12 | ALTO**

**Arquivo:** `pkg/config/config.go`

**Alterações:** Adicionar validação dupla:
```go
StarsTestMode = os.Getenv("STARS_TEST_MODE") == "true"
if StarsTestMode && os.Getenv("GO_ENV") == "production" {
    logger.Error("CONFIG", "⛔ STARS_TEST_MODE=true em produção! Desativando automaticamente.")
    StarsTestMode = false
}
```

---

#### Passo 2.6 — Tratar erro no generateShortID
**Achado #11 | ALTO**

**Arquivo:** `internal/cache/cache.go`

**Alterações:**
```go
func generateShortID(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            // Fallback seguro: usar crypto/rand bytes direto
            randByte := make([]byte, 1)
            rand.Read(randByte)
            b[i] = charset[int(randByte[0])%len(charset)]
            continue
        }
        b[i] = charset[num.Int64()]
    }
    return string(b)
}
```

---

### ═══ FASE 3: PERFORMANCE & CACHE (5 itens) ═══

#### Passo 3.1 — Cachear botInfo no middleware
**Achado #10 | ALTO**

**Arquivo:** `internal/middleware/checkAddBotMiddlewareTelego.go`

**Alterações:**
```go
// Adicionar variável de pacote para cachear botInfo
var (
    cachedBotInfo   *telego.User
    cachedBotInfoMu sync.Once
)

func getBotInfo(b *telego.Bot) *telego.User {
    cachedBotInfoMu.Do(func() {
        info, err := b.GetMe(context.Background())
        if err == nil {
            cachedBotInfo = info
        }
    })
    return cachedBotInfo
}
```

E usar `getBotInfo(b)` em vez de `b.GetMe(context.Background())` na linha 79.

---

#### Passo 3.2 — Cache L1 retornar deep copy
**Achado #17 | MÉDIO**

**Arquivo:** `internal/cache/cache.go`

**Alterações no GetChannel:**
```go
// Em vez de retornar o ponteiro diretamente do L1:
if val, found := localCache.Get(key); found {
    if channel, ok := val.(*models.Channel); ok {
        // Deep copy via JSON round-trip
        copy := *channel
        return &copy, nil
    }
}
```

Nota: Uma cópia rasa (`*channel`) é suficiente se os campos do Channel são value types. Se tiver slices/maps aninhados, usar JSON marshal/unmarshal.

Verificar struct `models.Channel` — tem slices como `Buttons`, `CustomCaptions`, etc. Então precisa de deep copy real:
```go
if channel, ok := val.(*models.Channel); ok {
    data, _ := json.Marshal(channel)
    var copy models.Channel
    json.Unmarshal(data, &copy)
    return &copy, nil
}
```

---

#### Passo 3.3 — Invalidar cache do ExecutorFactory nos fluxos relevantes
**Achado #6 | ALTO**

**Arquivos a modificar:**
- `internal/core/services/connected_account.go` — em `SaveSession()` e `Disconnect()`
- `internal/core/services/subscription_service.go` — em `activateSubscription()` e `ExpireSubscriptions()`

**Problema:** O `ExecutorFactory` vive no `container`, e os services não têm acesso a ele.

**Solução:** Adicionar um callback/hook no `AppContainer`:
```go
// Em appContainer.go
func (c *AppContainer) OnAccountChanged(userID int64) {
    c.ExecutorFactory.InvalidateCache(userID)
}
```

E chamar este callback nos services relevantes. Para evitar dependência circular, passar a função como parâmetro no construtor dos services, ou usar um evento.

**Abordagem mais simples:** Remover o cache da factory inteiramente (o lookup de HasConnectedAccount/HasPremiumManagedAccount já usa cache L1/Redis e é rápido):
```go
func (f *ExecutorFactory) ForUser(ctx context.Context, userID int64) TelegramExecutor {
    // Sem cache — consulta é rápida via Redis/L1
    if f.mtproto != nil && f.provider.HasConnectedAccount(ctx, userID) {
        return NewUserExecutor(userID, f.botAPI, f.mtproto)
    }
    // ...
}
```

---

#### Passo 3.4 — Limitar batch size no scheduler
**Achado #27 | MÉDIO**

**Arquivo:** `internal/core/services/scheduler.go`

**Alterações:**
```go
func (s *SchedulerService) processDuePosts() {
    ctx := context.Background()
    now := time.Now()
    // ...
    posts, err := s.repo.ClaimDuePosts(ctx, now)
    // ...
    
    // Limitar processamento por ciclo
    maxPerCycle := 20
    if len(posts) > maxPerCycle {
        posts = posts[:maxPerCycle]
    }
    
    for _, post := range posts {
        s.sendScheduledPost(ctx, &post)
        time.Sleep(1 * time.Second)
    }
}
```

---

#### Passo 3.5 — Reduzir TTL do selected_channel
**Achado #16 | MÉDIO**

**Arquivo:** `internal/cache/cache.go`

**Alterações:**
```go
// ANTES
return client.Set(ctx, key, channelID, 43200*time.Minute).Err()

// DEPOIS — 24 horas é suficiente
return client.Set(ctx, key, channelID, 24*time.Hour).Err()
```

---

### ═══ FASE 4: QUALIDADE DE CÓDIGO (5 itens) ═══

#### Passo 4.1 — Corrigir typo "Maintence" → "Maintenance"
**Achado #19 | MÉDIO**

**Atenção:** Este typo está no **nome da coluna** do banco de dados. Corrigir o nome do campo Go não altera a coluna existente automaticamente (GORM usa o nome do campo para mapear). 

**Decisão:** Manter o typo no banco por compatibilidade, mas usar tag `gorm:"column:maintence"` para que o Go use o nome correto:

**Arquivo:** `internal/database/models/models.go` (campo ServerConfig)
```go
Maintenance bool `gorm:"column:maintence" json:"maintenance"`
```

Ou, se o banco for recriável:
1. Renomear campo para `Maintenance`
2. Adicionar migração manual: `ALTER TABLE server_configs RENAME COLUMN maintence TO maintenance;`

**Recomendação:** Abordagem com tag `gorm:"column:maintence"` para não quebrar banco existente.

---

#### Passo 4.2 — Documentar typo "SecreteKey" (NÃO renomear)
**Achado #20 | MÉDIO**

**Decisão:** Renomear `SecreteKey` para `SecretKey` em todo o projeto quebraria a compatibilidade com `.env` existentes (campo `SECRET_KEY` não muda, mas o código Go sim). Como é um export público, isso pode quebrar dependências.

**Ação:** Adicionar alias e deprecation notice:
```go
// SecretKey é o alias correto para SecreteKey (typo histórico mantido por compatibilidade)
var SecretKey = SecreteKey
```

Ou aceitar o débito técnico e documentar em `.agent/decisions.md`.

---

#### Passo 4.3 — Mascarar dados sensíveis nos logs
**Achado #22 | MÉDIO**

**Arquivos:**
- `internal/core/services/admin_account.go` — linha 155
- `internal/core/services/subscription_service.go` — linhas 184, 271

**Alterações:**
```go
// ANTES
logger.Bot("📱 Iniciando auth MTProto admin: label=%s phone=%s", label, phoneNumber)

// DEPOIS
maskedPhone := phoneNumber[:3] + "****" + phoneNumber[len(phoneNumber)-2:]
logger.Bot("📱 Iniciando auth MTProto admin: label=%s phone=%s", label, maskedPhone)
```

Para payloads:
```go
// Truncar payload nos logs
truncated := payload
if len(truncated) > 20 {
    truncated = truncated[:20] + "..."
}
```

---

#### Passo 4.4 — Melhorar fallback do newJTI
**Achado #21 | MÉDIO**

**Arquivo:** `internal/api/auth/jwt.go`

**Alterações:**
```go
func newJTI() string {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        // Fallback com mais entropia que UnixNano
        b2 := make([]byte, 16)
        for i := range b2 {
            b2[i] = byte(time.Now().UnixNano() & 0xff)
            time.Sleep(time.Nanosecond)
        }
        return fmt.Sprintf("%x", b2)
    }
    return fmt.Sprintf("%x", b)
}
```

---

#### Passo 4.5 — Rate limiting na API
**Achado #23 | MÉDIO**

**Arquivo:** `internal/api/routes/routes.go` (ou novo arquivo `internal/api/middleware/rate_limit.go`)

**Alterações:** Adicionar middleware de rate limiting baseado em Redis:

```go
package middleware

func RateLimit(cache *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        key := fmt.Sprintf("rl:%s", ip)
        count, _ := cache.Incr(c.Request.Context(), key).Result()
        if count == 1 {
            cache.Expire(c.Request.Context(), key, window)
        }
        if count > int64(limit) {
            c.AbortWithStatusJSON(429, gin.H{"message": "Muitas requisições"})
            return
        }
        c.Next()
    }
}
```

Aplicar nas rotas:
```go
api.Use(middleware.RateLimit(redisClient, 100, time.Minute))
// Rate limit mais restritivo para login:
api.POST("/login", middleware.RateLimit(redisClient, 10, time.Minute), authController.Login)
```

---

### ═══ FASE 5: INFRAESTRUTURA (3 itens) ═══

#### Passo 5.1 — Limpar binários e imagens do repositório Git
**Achado #29, #30 | BAIXO**

**Comandos:**
```bash
# Remover binários do tracking (já estão no .gitignore)
git rm --cached FreddyBot Release 2>/dev/null || true

# Remover imagens da raiz
mkdir -p docs/images
git mv *.png docs/images/ 2>/dev/null || true
```

**Nota:** Os binários já estão no `.gitignore`, mas foram commitados antes. `git rm --cached` remove do tracking sem deletar localmente.

Para realmente limpar o histórico (opcional, requer rewrite):
```bash
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch FreddyBot Release' HEAD
```

---

#### Passo 5.2 — Remover .env do tracking se existir
**Achado #28 | BAIXO**

```bash
git rm --cached .env 2>/dev/null || true
```

O `.gitignore` já tem `.env`, mas o arquivo pode ter sido commitado antes da regra.

---

#### Passo 5.3 — Shutdown graceful completo no main.go
**Achado #13 | ALTO**

**Arquivo:** `cmd/FreddyBot/main.go`

**Alterações:**
```go
func main() {
    db, err := database.InitDB()
    if err != nil {
        logger.Error("APP", "Erro ao iniciar banco: %v", err)
        os.Exit(1)
    }
    defer func() {
        sqlDB, _ := db.DB()
        if sqlDB != nil {
            sqlDB.Close()
        }
    }()

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    webhookHandler, bot, app, err := telegram.StartBot(db)
    if err != nil {
        logger.Error("APP", "Erro ao iniciar bot: %v", err)
        return
    }
    defer func() {
        if err := cache.CloseRedis(); err != nil {
            logger.Error("APP", "Erro ao fechar Redis: %v", err)
        }
    }()

    app.StartBackground(ctx)

    go func() {
        if err := api.StartApi(ctx, app, webhookHandler); err != nil {
            logger.Error("APP", "Erro ao iniciar API: %v", err)
            stop()
        }
    }()

    <-ctx.Done()
    logger.Info("APP", "🧹 Encerrando app com segurança...")
    
    // Shutdown timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = shutdownCtx // usar nos shutdowns abaixo se necessário
}
```

**Nota:** O shutdown do `BotHandler` requer expor o `bh` do `telegram.StartBot`. Isso precisa de refator leve na assinatura de `StartBot` para retornar também o `bh`, ou encapsular num `Closer`.

---

### ═══ FASE 6: MELHORIAS FUTURAS (documentar) ═══

#### Passo 6.1 — Structured logging (JSON)
**Achado #33 | BAIXO**

Substituir `pkg/logger` por `slog` (stdlib Go 1.21+) ou `zap`. Não será feito nesta iteração, mas documentar a decisão.

#### Passo 6.2 — Migração controlada com goose/migrate
**Achado #25 | MÉDIO**

Substituir `AutoMigrate` por migrações versionadas. Requer planejamento separado. Documentar como trabalho futuro.

#### Passo 6.3 — Testes de integração para pagamento
**Achado #34 | BAIXO**

Adicionar testes para: intent expirada, pagamento duplicado, race conditions. Trabalho futuro.

#### Passo 6.4 — ConnMaxIdleTime no pool do DB
**Achado #26 | MÉDIO**

```go
sqlDB.SetConnMaxIdleTime(5 * time.Minute)
```

Pode ser incluído na Fase 1 (passo 1.2) junto com as alterações em database.go.

#### Passo 6.5 — SCAN do Redis com namespace
**Achado #7 | ALTO**

Refatorar as chaves Redis para usar prefixo `user:{ID}:` em vez de sufixo `:{ID}`. Requer migração de chaves existentes. Trabalho futuro.

---

## Riscos

- **Fase 1.2 (panic→error):** Alterar `GetRedisClient()` propaga mudanças para ~30 callers. Mitigação: usar `log.Fatalf` como intermediário.
- **Fase 1.4-1.5 (PRAGMA/SQL condicional):** Se a detecção do dialector falhar, pode executar SQL incompatível. Mitigação: usar variável explícita `isSQLite` definida no switch.
- **Fase 3.3 (cache factory):** Remover cache pode aumentar latência marginalmente. Mitigação: as consultas subjacentes já usam Redis L1/L2.
- **Fase 4.1 (typo Maintence):** Renomear coluna em banco existente pode causar downtime. Mitigação: usar tag `gorm:"column:maintence"`.
- **Fase 4.5 (rate limiting):** Rate limit muito agressivo pode bloquear usuários legítimos. Mitigação: começar com limites generosos (100 req/min).

## Impactos esperados

- Eliminação de 5 vulnerabilidades críticas
- Shutdown graceful sem perda de dados
- Prevenção de panics em produção
- Integridade referencial garantida em SQLite
- Proteção contra brute force na API
- Cache consistente sem race conditions
- Logs limpos sem dados sensíveis

## Compatibilidade
- ✅ Linux
- ✅ macOS  
- ✅ Docker
- ✅ PostgreSQL (prod)
- ✅ SQLite (dev)
- ✅ CI/CD (sem breaking changes na API pública)

## Como testar

### Build
```bash
go build -o FreddyBot ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./... -v -count=1
```

### Testes específicos de segurança
```bash
go test ./internal/api/auth/... -v -run TestSecurity
go test ./internal/core/services/... -v -run TestSubscription
go test ./internal/telegram/executor/... -v -run TestFactory
```

### Verificação manual
1. Iniciar app em modo dev → verificar que PRAGMA foreign_keys está ativo
2. Iniciar app em modo prod → verificar que migrações PostgreSQL executam
3. Enviar SIGINT → verificar shutdown graceful (sem goroutines pendentes)
4. Testar login com rate limit → verificar bloqueio após 10 tentativas/min
5. Verificar logs → nenhum telefone ou payload completo visível

## Rollback

Cada fase tem seu próprio commit. Em caso de problema:
```bash
# Reverter última fase
git revert HEAD

# Reverter fase específica
git log --oneline  # encontrar commit da fase
git revert <commit-hash>
```

Os binários e .env removidos do tracking podem ser restaurados:
```bash
git checkout HEAD~1 -- FreddyBot Release .env
```

## Observações

1. **Ordem importa:** As fases 1-3 devem ser implementadas nessa ordem. A fase 1 é pré-requisito para as demais.

2. **Commits recomendados:**
   - `fix(security): remove hardcoded credentials from docker-compose`
   - `fix(startup): replace panic with error return in InitDB`
   - `fix(shutdown): add context cancellation to broadcast workers`
   - `fix(db): condition PRAGMA and migrations on dialector type`
   - `fix(auth): use redis.Nil constant, safe type assertions`
   - `feat(api): add rate limiting middleware`
   - `refactor(cache): return deep copy from L1 cache`
   - `chore(repo): remove tracked binaries and images`

3. **Estimativa de tempo:** ~4-6 horas de implementação para as fases 1-5.

4. **A fase 6 (structured logging, migrações controladas, testes de integração) pode ser planejada separadamente** como um projeto de melhoria contínua pós-produção.
