# Plano: conexao-conta-telegram-usuario

## Pedido do usuário
Criar funcionalidade para um usuário conectar sua conta pessoal do Telegram ao bot usando a biblioteca `github.com/gotd/td` (MTProto). Se a lib permitir login nativo via Telegram (QR Code), usar esse fluxo. Caso contrário, criar um botão que leva a uma página do Mini App pedindo número, código e 2FA se necessário. Backend deve enviar mensagem de confirmação no console e na conta do usuário. Garantir que nenhum arquivo de sessão seja salvo (apenas banco criptografado) e que nada vaze.

## Objetivo
Permitir que usuários conectem sua conta Telegram pessoal ao bot via MTProto (gotd/td), com:
- Fluxo único: Mini App com phone + code + 2FA (sem QR Code — usuário já está no celular)
- Armazenamento criptografado (AES-256-GCM) da sessão no banco
- Notificação ao usuário via bot quando conectado
- Sem arquivos de sessão em disco

## Contexto atual
- Bot usa `github.com/mymmrac/telego` (Bot API HTTP)
- `github.com/gotd/td` NÃO está no go.mod — precisa ser adicionado
- Usuários são salvos via middleware `SaveUserMiddlewareTelego`
- Dashboard React SPA integrada como Mini App Telegram
- JWT + initData Telegram para autenticação web
- SECRET_KEY já existe no env para criptografia JWT
- AppContainer injeta dependências (services, repos)

## Arquivos analisados
- cmd/FreddyBot/main.go
- internal/container/appContainer.go
- internal/database/models/models.go
- internal/database/database.go
- internal/database/repositories/user.go
- internal/core/services/user.go
- internal/api/api.go
- internal/api/routes/routes.go
- internal/api/controllers/authController.go
- internal/api/auth/jwt.go
- internal/api/auth/signature.go
- internal/api/types/response.go
- internal/api/types/webAppAuthReceive.go
- internal/telegram/client.go
- internal/telegram/loader_telego.go
- pkg/config/config.go
- pkg/logger/logger.go
- pkg/errors/errors.go
- go.mod
- .env-example
- dashboard/src/api.ts
- dashboard/src/App.tsx
- dashboard/src/types.ts
- dashboard/package.json
- Makefile

## Arquivos que poderão ser modificados
- go.mod (adicionar gotd/td e dependências)
- go.sum
- cmd/FreddyBot/main.go (inicializar serviço de clientes)
- internal/container/appContainer.go (adicionar TelegramClientService)
- internal/database/models/models.go (adicionar UserTelegramSession model)
- internal/database/database.go (adicionar AutoMigrate)
- internal/database/repositories/telegram_session.go (NOVO)
- internal/core/services/telegram_client.go (NOVO - serviço gotd/td)
- internal/core/crypto/crypto.go (NOVO - AES-256-GCM)
- internal/api/types/connect.go (NOVO - request/response types)
- internal/api/controllers/connectController.go (NOVO)
- internal/api/routes/routes.go (adicionar rotas)
- internal/telegram/loader_telego.go (adicionar comando /connect)
- internal/telegram/handlers/commands/connect/ (NOVO - handler do comando)
- pkg/config/config.go (adicionar ENC_KEY do env)
- .env-example (adicionar ENC_KEY)
- dashboard/src/api.ts (adicionar funções connect)
- dashboard/src/types.ts (adicionar tipos)
- dashboard/src/App.tsx (adicionar página de connect)
- dashboard/src/components/TelegramConnect.tsx (NOVO)

## Estratégia de implementação

### Camada de Criptografia
Criar `internal/core/crypto/crypto.go` com funções `Encrypt(plaintext, key)`, `Decrypt(ciphertext, key)`. Usar AES-256-GCM com nonce de 12 bytes aleatório. Derivar chave por usuário via HKDF(SECRET_KEY, userID).

### Model + Repositório
Adicionar model `UserTelegramSession` ao banco com:
- `UserID` (PK, FK → User)
- `EncryptedSession` (text, encrypted session data)
- `EncryptedPhoneHash` (text, encrypted hash, opcional)
- `IsActive` (bool)
- `CreatedAt`, `UpdatedAt`

Criar `internal/database/repositories/telegram_session.go` com operações CRUD.

### Serviço Telegram Client (gotd/td)
Criar `internal/core/services/telegram_client.go`:
- `TelegramClientService` struct que gerencia:
  - `GetClient(userID)` - obtém/recria client da sessão criptografada no banco
- `StartPhoneFlow(userID, phone)` - inicia autenticação via phone
  - `SendCode(userID, code)` - envia código OTP
  - `Send2FA(userID, password)` - envia senha 2FA
  - `ConnectUser(userID, client)` - salva sessão criptografada
  - `DisconnectUser(userID)` - remove sessão
  - `SendMessage(userID, chatID, text)` - envia mensagem via user client
- Implementar `gotd/td` `SessionStorage` interface para salvar no banco criptografado
- Usar `session` package do gotd/td para `SessionStorage`

### Phone Auth Flow (Fluxo único)
- Mini App coleta phone number
- Backend inicia phone auth via gotd/td
- Telegram envia código OTP para o número
- Mini App pede código, usuário digita
- Se 2FA necessário, backend detecta (SessionCodeNeededError) e Mini App pede senha 2FA
- Não há QR Code — usuário já está no celular com o Mini App aberto

### Bot Handler
- Comando `/connect` no bot
- Envia mensagem com botão inline → Mini App
- Callback após conexão bem-sucedida

### Notificação
- Após conexão: `logger.Info("TGCONNECT", ...)` no console
- Mensagem via telego (bot) para o usuário informando sucesso

### Dashboard (Mini App)
- Nova página/rota para connect
- Componente React `TelegramConnect.tsx`
- Fluxo: Phone → Code → 2FA (se necessário) → Success
- Comunicação com backend via API REST

## Passos detalhados

1. Adicionar `github.com/gotd/td` e dependências ao go.mod
2. Criar `internal/core/crypto/crypto.go` (AES-256-GCM encrypt/decrypt)
3. Adicionar model `UserTelegramSession` em `internal/database/models/models.go`
4. Adicionar `AutoMigrate` para o novo model em `internal/database/database.go`
5. Criar `internal/database/repositories/telegram_session.go`
6. Adicionar `ENC_KEY` em `pkg/config/config.go` e `.env-example`
7. Criar `internal/core/services/telegram_client.go` com serviço gotd/td
8. Adicionar `TelegramClientService` ao `AppContainer`
9. Criar `internal/api/types/connect.go` com request/response types
10. Criar `internal/api/controllers/connectController.go` (endpoints REST)
11. Adicionar rotas em `internal/api/routes/routes.go`
12. Criar `internal/telegram/handlers/commands/connect/connect.go` (comando /connect)
13. Registrar handler no `loader_telego.go`
14. Criar `dashboard/src/components/TelegramConnect.tsx`
15. Adicionar tipos em `dashboard/src/types.ts`
16. Adicionar funções API em `dashboard/src/api.ts`
17. Adicionar rota no `dashboard/src/App.tsx`
18. Adicionar rota de dashboard no `internal/api/api.go`
19. Rodar `go mod tidy` e `make build` para verificar

## Multi-login (vários usuários simultâneos)
- Cada `UserTelegramSession` é identificada por `UserID` (PK)
- Cada usuário mantém seu próprio client MTProto independente com sua própria `auth_key`
- Clients são gerenciados em um `sync.Map` interno no serviço, chaveado por `UserID`
- Conexões MTProto são independentes (cada uma tem seu próprio DC, auth_key, sessão)
- Não há compartilhamento de estado entre usuários
- Escalabilidade horizontal fica prejudicada (clients MTProto são stateful na instância), mas documentaremos isso como limitação

## Riscos
- gotd/td usa MTProto que requer `dc_id`, `server_address`, `auth_key` — dados sensíveis que exigem criptografia forte
- MTProto pode ser bloqueado em algumas redes (Rússia, China, etc.)
- Sessões expiram se o usuário mudar a senha ou revogar acesso no Telegram
- Múltiplas conexões MTProto simultâneas consomem recursos (cada usuário = 1 client)
- gotd/td requer Go 1.21+ (projeto já usa 1.25, sem problema)
- Dependência pesada (gotd/td + crypto libraries)

## Impactos esperados
- +1 model, +1 repository, +1 service, +1 controller, +1 command handler
- +1 página React no dashboard
- Consumo de memória: cada sessão MTProto ativa mantém conexão WebSocket/TCP
- Consumo de CPU: criptografia AES-256 em cada save/load de sessão
- ENC_KEY necessária no .env (diferente da SECRET_KEY JWT para separação de responsabilidades)

## Compatibilidade
- Linux ✅
- macOS ✅
- Windows ✅ (com ressalvas: MTProto WebSocket pode ter issues)
- Docker ✅
- CI/CD ✅ (sem alterações em infra)

## Como testar

### Build
```bash
go mod tidy
make build-server
```

### Testes
```bash
go build ./...
```

### Execução
```bash
# Configurar ENC_KEY no .env
echo 'ENC_KEY=sua-chave-aes-256-aqui-com-32-caracteres' >> .env
make run
```

## Rollback
```bash
git checkout -- go.mod go.sum
git checkout -- cmd/ internal/ pkg/ dashboard/
git clean -fd
go mod tidy
```

## Observações
- `gotd/td` SessionStorage interface: `gotd/td/session` package fornece interface para persistência. Implementaremos uma versão que salva no banco criptografado.

- A chave ENC_KEY deve ter exatamente 32 bytes (AES-256)
- Cada cliente MTProto mantém uma conexão ativa → considerar pool/garbage collection de clientes inativos
