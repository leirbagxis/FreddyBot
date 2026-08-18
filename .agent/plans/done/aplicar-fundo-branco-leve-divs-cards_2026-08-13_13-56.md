# Plano: Aplicar Fundo Branco Leve nas Divs e Cards da Dashboard

## Pedido do usuário
"eu pedi para colocar nos fundos das divs um branco leve e nao foi aplicado"

## Objetivo
Aplicar um fundo branco leve (`#ffffff` / `rgba(255, 255, 255, 0.95)`) em todas as divs, cards, containers e superfícies do dashboard (`.action-card`, `.content-card`, modais, inputs e divs do menu), garantindo contraste sobre um fundo de página suave (`#f4f4f6`).

## Contexto atual
- As variáveis CSS `--card`, `--card-elevated`, `--surface` e `--input-bg` nos blocos `[data-theme="telegram"]` e `[data-theme="light"]` do `index.css` estavam definidas como `transparent`, fazendo com que as divs ficassem transparentes em vez de brancas.
- No `SideMenu.tsx`, os containers das divs usavam `bg-white/5`, ficando imperceptíveis em telas claras.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/components/SideMenu.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/components/SideMenu.tsx`

## Estratégia de implementação

1. **`dashboard/src/index.css`**:
   - Em `[data-theme="telegram"]` e `[data-theme="light"]`:
     - `--card: #ffffff;`
     - `--card-elevated: #ffffff;`
     - `--surface: #ffffff;`
     - `--surface-hover: #f4f4f5;`
     - `--input-bg: #ffffff;`
     - `--bg: #f4f4f6;` (fundo cinza/bege muito suave para dar contraste elegante às divs brancas)
     - Manter `--border: transparent;` conforme solicitado anteriormente.

2. **`dashboard/src/hooks/useTheme.ts`**:
   - Em `applyTelegramVars()`:
     - No fallback de modo claro (`scheme === 'light'`), definir:
       - `bg: '#f4f4f6'`
       - `card: '#ffffff'`
       - `surface: '#ffffff'`
     - Garantir que `--card`, `--card-elevated`, `--surface` e `--input-bg` recebam `#ffffff` quando em modo claro, sobressaindo com elegância.

3. **`dashboard/src/components/SideMenu.tsx`**:
   - Atualizar a div do bloco "Minha Conta" e a div do menu de navegação: trocar `bg-white/5` por `bg-white shadow-sm` para que as divs do menu lateral exibam o fundo branco leve.

4. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Riscos
- **Nenhum**: Alteração estética para atender exatamente ao requisito de ter fundo branco leve em todas as divs.

## Impactos esperados
- Todas as divs e cards (`.action-card`, `.content-card`, modais, menu lateral, inputs) com fundo branco leve (`#ffffff`) sob um fundo de página suave (`#f4f4f6`), sem bordas, garantindo visual limpo, moderno e com contraste perfeito.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/index.css dashboard/src/hooks/useTheme.ts dashboard/src/components/SideMenu.tsx
```
