# Plano: Ajustar Tema Telegram e Header Transparente/Blur na Admin Dashboard

## Pedido do usuário
"O tema ainda esta branco e nao o thema vindo do teelegram. Tbm coloque o header transparente/bluor e deixe somente a parte de search e o + no header da dashboard admin, pfvr"

## Objetivo
1. **Fundo do Tema Telegram**: Remover os fundos brancos fixos `#ffffff` de `.admin-layout-v2` e `.admin-sidebar`, fazendo com que a Admin respeite o tema dinâmico do Telegram (`var(--bg)` / `var(--tg-theme-bg-color)`).
2. **Header Transparente com Blur**:
   - Estilizar a `.admin-topbar` com fundo transparente e efeito de desfoque (`backdrop-filter: blur(12px)` / `-webkit-backdrop-filter: blur(12px)`).
3. **Simplificação do Header**:
   - Manter no header apenas a barra de pesquisa (`Search`) e o botão de adicionar (`+` / Novo broadcast).
   - Remover contadores, ordenações e perfil do cabeçalho para uma interface minimalista.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/components/admin/AdminTopbar.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/admin/AdminTopbar.tsx`

## Estratégia de implementação

1. **`dashboard/src/index.css`**:
   - Atualizar a classe `.admin-layout-v2` para usar `background: var(--bg)` em vez de `#ffffff`.
   - Atualizar `.admin-sidebar` para usar `background: rgba(255, 255, 255, 0.03)` com `backdrop-filter: blur(12px)` e `border: none`.
   - Atualizar `.admin-topbar` para usar `background: rgba(255, 255, 255, 0.03)`, `backdrop-filter: blur(12px)` e `border: none`.
   - Atualizar `.admin-topbar-search` com fundo `rgba(255, 255, 255, 0.05)`, cantos arredondados (`12px`) e sem bordas duras.

2. **`dashboard/src/components/admin/AdminTopbar.tsx`**:
   - Simplificar o JSX do componente `AdminTopbar`:
     - Lado esquerdo: Botão de menu mobile e barra de busca (`Search input`).
     - Lado direito: Apenas o botão de `+` (Novo broadcast).
     - Remover seletores de ordenação, filtros, notificações e perfil do header.

3. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Riscos
- **Nenhum**: Apenas otimização de layout e CSS do cabeçalho da admin.

## Impactos esperados
- A Dashboard Admin utilizará o tema dinâmico escuro/translúcido do Telegram.
- O cabeçalho terá um visual de vidro (*glassmorphism*) transparente com desfoque, contendo apenas a busca e o botão `+`.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/index.css dashboard/src/components/admin/AdminTopbar.tsx
```
