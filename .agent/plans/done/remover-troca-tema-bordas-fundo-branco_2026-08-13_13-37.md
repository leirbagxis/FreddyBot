# Plano: Remover Troca de Tema e Usar Apenas Tema do Telegram

## Pedido do usuário
Atuar como designer profissional: remover o modo de troca de tema manual e deixar apenas o tema vindo do Telegram. Remover todos os ícones/botões de troca de tema de todo o dashboard. Remover a cor de todas as bordas possíveis do dashboard e aplicar um fundo branco claro.

## Objetivo
1. Simplificar o hook `useTheme` para sempre usar o tema Telegram (sem toggle manual).
2. Remover todos os botões/ícones de troca de tema (top-bar e SideMenu).
3. Remover as props `theme`/`toggleTheme` do SideMenu.
4. Alterar o `--border` no CSS para `transparent` em todos os temas, eliminando visualmente todas as bordas decorativas.
5. Aplicar fundo branco claro (`--bg: #ffffff` ou `#fafafa`), delegando ao `applyTelegramVars()` sobrescrever dinamicamente quando em Telegram.

## Contexto atual
- O hook `useTheme.ts` suporta 3 temas: `light`, `dark`, `telegram`, com toggle cíclico e persistência em localStorage.
- O botão de troca de tema aparece em 2 locais: top-bar (`App.tsx` L1004-1006) e SideMenu (`SideMenu.tsx` L173-191).
- As bordas usam a variável CSS `--border` com cores semi-transparentes de accent em cada tema.

## Arquivos analisados
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/index.css`

## Arquivos que poderão ser modificados
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação

### Passo 1 — `useTheme.ts`: Simplificar para Telegram-only
- O hook sempre inicializa com `'telegram'` como tema.
- Remover a função `toggleTheme` (retornar um no-op para não quebrar tipos).
- Remover a lógica de `localStorage` para `theme-manual`.
- Remover os temas `light`/`dark` do ciclo — sempre aplicar `applyTelegramVars()`.
- Manter a detecção de `tg.colorScheme` (light/dark) dentro do `applyTelegramVars` para que o Telegram controle automaticamente se é claro/escuro.

### Passo 2 — `App.tsx`: Remover botão de tema do top-bar
- Remover as linhas L1004-1006 (o `<button className="theme-switch">` com Sun/Moon/Send).
- Remover os imports de `Sun`, `Moon`, `Send` se não forem usados em outro lugar (Sun e Moon usados no greeting — manter).
- Remover a prop `toggleTheme` passada ao `<SideMenu />` (L1424).
- Remover a prop `theme` passada ao `<SideMenu />` (L1423).

### Passo 3 — `SideMenu.tsx`: Remover botão "Alternar Tema"
- Remover props `theme` e `toggleTheme` da interface e da desestruturação.
- Remover o bloco do botão "Alternar Tema" (L173-191).
- Remover imports `Sun`, `Moon`, `Send` que eram usados apenas para o tema.

### Passo 4 — `index.css`: Bordas transparentes e fundo branco claro
- Em `[data-theme="light"]`, `[data-theme="dark"]`, `[data-theme="telegram"]` e `.dark`:
  - `--border: transparent;`
  - `--border-active: transparent;` (ou manter accent para bordas ativas de foco)
- No tema `[data-theme="telegram"]`:
  - `--bg: #fafafa;` (branco claro suave, será sobrescrito pelo `applyTelegramVars` quando em Telegram real)
- Manter `border-radius` intactos em todos os componentes.
- A variável `--border` sendo `transparent` automaticamente remove a cor de todas as bordas que usam `border: 1px solid var(--border)` sem alterar estrutura.

### Passo 5 — Remover CSS `.theme-switch`
- Remover o bloco CSS `.theme-switch` (L1311-1338) e a animação `@keyframes theme-icon-in` que só serviam ao botão de troca de tema.

## Riscos
- **Baixo**: Componentes que dependiam de `--border` para separação visual perderão a borda, mas o efeito é o desejado pelo usuário.
- **Baixo**: A prop `toggleTheme` no SideMenu e App.tsx será removida — sem impacto externo pois é passada internamente.

## Impactos esperados
- Dashboard visualmente mais limpo: sem bordas visíveis, fundo branco claro, tema exclusivamente controlado pelo Telegram.
- Remoção de toda UI de troca de tema (top-bar e menu lateral).
- Experiência nativa do Telegram: o tema adapta automaticamente ao esquema de cores do app Telegram do usuário.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/hooks/useTheme.ts dashboard/src/App.tsx dashboard/src/components/SideMenu.tsx dashboard/src/index.css
```

## Observações
- Os imports `Sun`, `Sunrise`, `CloudMoon` no `App.tsx` são usados no greeting (bom-dia/boa-tarde/boa-noite) e devem ser mantidos.
- O `Send` import no `App.tsx` pode ser removido se não for usado em outros lugares além do botão de tema.
