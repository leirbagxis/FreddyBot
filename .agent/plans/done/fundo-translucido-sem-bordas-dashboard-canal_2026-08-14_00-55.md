# Plano: Fundo Translucido (bg-white/5) e Sem Bordas nas Divs do Dashboard do Canal (/dashboard/:channelId)

## Pedido do usuário
"Na dashboard /dashboad/:channelId eu quero que voce coloque em todas as divs o fundo branco leve igual tem no menu lateral e sem as bordas nas divs, pfvr"

## Objetivo
Configurar todas as divs, cards e containers da tela de dashboard de canal (`/dashboard/:channelId`) com o fundo translúcido `bg-white/5` (`rgba(255, 255, 255, 0.05)`), idêntico ao menu lateral, e remover todas as bordas visíveis de cards e elementos internos.

## Contexto atual
- As variáveis CSS `--card`, `--card-elevated`, `--surface` e `--input-bg` estavam com fundo opaco.
- Alguns containers e cards no `App.tsx` e componentes internos (Reações, Links Dinâmicos, Permissões por Tipo, Editor de Legendas) possuem classes com bordas (`border border-border`, `border-l-2`).

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/CaptionCard.tsx`
- `dashboard/src/components/ReactionsCard.tsx`
- `dashboard/src/components/ButtonGrid.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/App.tsx`

## Estratégia de implementação

1. **`dashboard/src/index.css`**:
   - Ajustar as variáveis globais de temas (`[data-theme="telegram"]` e `[data-theme="light"]`):
     - `--card: rgba(255, 255, 255, 0.05);` (fundo translúcido idêntico ao menu lateral)
     - `--card-elevated: rgba(255, 255, 255, 0.05);`
     - `--surface: rgba(255, 255, 255, 0.05);`
     - `--input-bg: rgba(255, 255, 255, 0.05);`
     - `--border: transparent;`
   - Atualizar regras CSS de cards:
     - `.content-card`, `.action-card`, `.caption-preview`, `.rte-container`, `.rte-toolbar`: `border: none;` e `background: rgba(255, 255, 255, 0.05);`.
     - Sobrescrever sub-elementos com `bg-muted/50` ou `border-border` para remover bordas e usar `background: rgba(255, 255, 255, 0.04)`.

2. **`dashboard/src/hooks/useTheme.ts`**:
   - Atualizar `applyTelegramVars()` para atribuir `rgba(255, 255, 255, 0.05)` a `--card`, `--surface`, `--input-bg` e `--card-elevated`, mantendo a transparência de bordas `--border: transparent`.

3. **`dashboard/src/App.tsx`**:
   - Remover as classes de borda em itens da aba de permissões/configurações do canal (e.g. trocar `border border-border` por `border-none` e garantir `bg-white/5` nas sub-divs).

4. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Riscos
- **Nenhum**: Alteração visual solicitada especificamente pelo usuário para ter o mesmo fundo translúcido `bg-white/5` do menu lateral sem bordas nas divs.

## Impactos esperados
- Todas as divs e cards da tela do canal (`/dashboard/:channelId`) exibirão o fundo branco leve translúcido idêntico ao menu lateral (`bg-white/5`), sem bordas visíveis, criando um visual limpo e fluido.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/index.css dashboard/src/hooks/useTheme.ts dashboard/src/App.tsx
```
