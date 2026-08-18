# Plano: Corrigir acesso negado para usuário novo (cookie-only)

## Pedido do usuário
Usuários novos recebem a tela "Ops! Acesso negado / Acesso não autorizado" ao abrir o dashboard. O usuário pediu para manter autenticação 100% via cookie (sem token em localStorage).

## Objetivo
Garantir que o cookie `token` definido em `/api/login` funcione dentro do WebView do Telegram (contexto de iframe/terceiros) e que usuários novos sejam aceitos, mantendo o modo cookie-only.

## Contexto atual
- `AuthMiddlewareJWT` (`internal/api/auth/middleware.go:22-32`) só autentica via cookie `token` ou header `Authorization: Bearer`. A mensagem "Acesso não autorizado" (linha 29) ocorre quando não há token algum.
- `Login` (`internal/api/controllers/authController.go:86-94`) define cookie com `SameSite: http.SameSiteStrictMode` e `Secure: config.AppEnv != "dev"`.
- Cookies `SameSite=Strict` são descartados/não enviados em contexto de iframe/WebView de terceiros do Telegram Mini App (tdesktop#30889, Chromium#41352976), quebrando a primeira autenticação de usuários novos.
- Usuário que nunca interagiu com o bot não existe no banco; o middleware retorna "Usuário não encontrado ou inativo" (`middleware.go:44`) e o `Login` não faz upsert.
- No frontend, `App.tsx:262` lê `authRes.isBlacklisted` mas o campo correto é `authRes.data.isBlacklisted`.

## Arquivos analisados
- internal/api/auth/middleware.go
- internal/api/controllers/authController.go
- dashboard/src/App.tsx
- dashboard/src/api.ts
- internal/api/auth/jwt.go
- pkg/config

## Arquivos que poderão ser modificados
- internal/api/controllers/authController.go
- dashboard/src/App.tsx

## Estratégia de implementação
1. No backend, ajustar atributos do cookie para serem compatíveis com o WebView do Telegram:
   - Prod: `SameSite=None; Secure=true`.
   - Dev: `SameSite=Lax; Secure=false`.
2. No `Login`, fazer upsert do usuário quando não existir, usando dados do initData do Telegram (nome, username), mantendo blacklist/role existentes.
3. No frontend, corrigir leitura de blacklist e adicionar auto-re-login com retry único em falha 401 usando initData salvo em CloudStorage.

## Passos detalhados

1. Criar o plano em `.agent/plans/pending/` e mover para `approved/` após aprovação.
2. `internal/api/controllers/authController.go`:
   - Ajustar o `http.Cookie` do token:
     - `Secure: config.AppEnv != "dev"`.
     - `SameSite: http.SameSiteNoneMode` quando prod, `http.SameSiteLaxMode` quando dev.
   - No bloco de determinação de role, quando `GetUserByID` não encontrar o usuário, fazer `UpsertUser` com `UserId`, `FirstName` (limpo de HTML) e `Username` extraídos do initData (`userDataRaw` é JSON; extrair via `json.Unmarshal` em struct com `id`, `first_name`, `username`).
3. `dashboard/src/App.tsx`:
   - Corrigir `authRes.isBlacklisted` → `authRes.data?.isBlacklisted`.
   - Criar função de auto-re-login: guardar initData em `sessionStorage`/`CloudStorage` já existente; ao detectar 401 nas chamadas autenticadas, reexecutar `login(initData, userID)` uma única vez e re-tentar a chamada original.
4. Rodar testes Go e `make build`.
5. Mover plano para `done/` e atualizar `.agent/memory/memory.md`.

## Riscos
- `SameSite=None` exige `Secure`; em dev via HTTP o login quebraria — mitigado usando `Lax` em dev.
- Retry automático pode mascarar erros reais — limitado a 1 tentativa.
- Alteração do upsert deve preservar flags existentes (admin/blacklist) — usar `UpsertUser` que só cria/atualiza nome, sem resetar flags.

## Impactos esperados
- Usuários novos conseguem autenticar no WebView do Telegram.
- Sessão continua funcionando para usuários existentes.
- Blacklist e controle de admin preservados.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Build
```bash
make build
```

### Testes
```bash
go test ./internal/api/...
```

### Execução
```bash
./Release
```

## Rollback
- Reverter o commit das alterações (ou restaurar atributos do cookie e remover upsert).

## Observações
- Mantido o modo cookie-only: nenhum token em localStorage.
