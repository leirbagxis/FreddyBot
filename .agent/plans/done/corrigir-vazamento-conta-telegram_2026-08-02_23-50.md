# Plano: Corrigir vazamento da feature "Conta Telegram" desligada

## Pedido do usuário
Quando o admin desliga a feature "Conta Telegram Pessoal", o card "Conta Telegram" às vezes continua aparecendo na tela `/me/channels` e vaza para o usuário.

## Objetivo
Fazer com que as flags de feature (`connectedAccountEnabled`, `premiumEnabled`) sejam **fail-closed** e carregadas **após a autenticação**, impedindo que features desligadas vazem.

## Contexto atual
- Efeito em `dashboard/src/App.tsx:161-183` roda no mount, ANTES do login (deps `[]`).
- `fetchSubscriptionStatus()`/`fetchAccountStatus()` são endpoints protegidos → 401 sem token → `.catch()` → `s = null`.
- `caEnabled = s?.connectedAccountEnabled !== false` e `useState(true)` são **fail-open** → quando o fetch falha, a flag vira `true` e o card aparece.
- O efeito não re-executa após o login, então o valor errado persiste.
- O card "Conta Telegram" (`App.tsx:1039`), a aba "Conta" (`App.tsx:849`) e o modal (`App.tsx:1369`) são gated por `connectedAccountEnabled`.
- `PremiumTab` (`App.tsx:1053`) é renderizado sem gate e lista "Conta Telegram Gerenciada" como benefício.

## Arquivos analisados
- dashboard/src/App.tsx
- dashboard/src/components/PremiumTab.tsx
- internal/core/services/subscription_service.go

## Arquivos que poderão ser modificados
- dashboard/src/App.tsx

## Estratégia de implementação
1. Fail-closed: `useState(false)` para `premiumEnabled` e `connectedAccountEnabled`.
2. No efeito de features, usar `=== true` em vez de `!== false`.
3. Fazer o efeito re-executar após a autenticação: dependência em `authState`, ignorando quando `authState !== 'authenticated'`.
4. Gate do `PremiumTab` por `premiumEnabled` na tela `/me/channels` (mesma causa raiz; evita vazar o card premium/benefício quando o fetch falha).
5. Rebuild do dashboard.

## Passos detalhados

1. `dashboard/src/App.tsx`:
   - Linhas 121-122: `premiumEnabled` e `connectedAccountEnabled` iniciam `false`.
   - Efeito (linhas 161-183): adicionar `if (authState !== 'authenticated') return;`, mudar para `pEnabled = s?.premiumEnabled === true` e `caEnabled = s?.connectedAccountEnabled === true`, e dep `[authState]`.
   - Linha ~1053: envolver `PremiumTab` com `{premiumEnabled && (...)}`.
2. Rebuild: `make build`.

## Riscos
- Flash breve de cards ocultos até o fetch pós-auth completar (aceitável).
- Se o fetch pós-auth falhar (WebView sem cookie), a feature fica oculta (fail-closed) em vez de vazar — comportamento desejado.

## Impactos esperados
- Feature desligada não vaza mais na tela `/me/channels`.
- Usuários legítimos veem o card/aba apenas quando o backend confirma a feature ativa.

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
- Reverter mudanças de default e gate em `App.tsx`.

## Observações
- Mantida a sessão cookie-only (sem token em localStorage).
