# Plano: corrigir-issues-review

## Pedido do usuário
Corrigir as issues de CRITICAL e HIGH identificadas no code review da feature de subscriptions/separator/PremiumConfigTab.

## Objetivo
Resolver 7 issues de segurança e qualidade introduzidas pela feature, sem modificar código pré-existente não relacionado.

## Arquivos que serão modificados

1. `internal/api/controllers/channelController.go` — validação EmojiEntitiesJSON + whitelist Type
2. `dashboard/src/components/PremiumConfigTab.tsx` — useCallback + apiFetch + catch logging
3. `dashboard/src/components/AdminSubscriptionsTab.tsx` — importar Subscription de types.ts
4. `dashboard/src/App.tsx` — catch(() => {}) logging
5. `dashboard/src/components/CaptionPreview.tsx` — sanitização HTML no dangerouslySetInnerHTML

## Passos detalhados

### 1. channelController.go — Validar EmojiEntitiesJSON
- Fazer unmarshal do JSON e validar estrutura `[{type, offset, length, emoji_id}]`
- Se inválido, retornar BadRequest

### 2. channelController.go — Whitelist Separator.Type
- Validar que Type é "sticker" ou "custom_emoji"
- Valor default continua "custom_emoji"

### 3. PremiumConfigTab.tsx — Envolver funções em useCallback
- `addEmojiToSeparator`, `removeEmojiFromSeparator`, `saveSeparator`, `deleteSeparator`

### 4. PremiumConfigTab.tsx — Usar apiFetch para emoji history
- Substituir `fetch('/api/emoji/history', ...)` por `apiFetch('/api/emoji/history')`

### 5. AdminSubscriptionsTab.tsx — Importar Subscription de types.ts
- Remover interface local duplicada
- Importar `Subscription` de `../types`

### 6. App.tsx + PremiumConfigTab — Logging em catch
- Substituir `.catch(() => {})` por `.catch(err => console.warn(...))`

### 7. CaptionPreview.tsx — Sanitizar HTML antes de dangerouslySetInnerHTML
- Adicionar função de escape para tags HTML perigosas
- Ou usar regex para remover tags <script> e event handlers

## Riscos
- Modificar CaptionPreview pode afetar outros componentes que o usam
- A validação de EmojiEntitiesJSON pode quebrar clients existentes se o formato esperado for diferente
  - Mitigação: validar apenas estrutura mínima, não rejeitar campos extras

## Como testar

### Build frontend
```bash
npm run build
```

### Build backend
```bash
go build ./cmd/...
```

## Rollback
- Reverter alterações com `git checkout -- <arquivo>` para cada arquivo
