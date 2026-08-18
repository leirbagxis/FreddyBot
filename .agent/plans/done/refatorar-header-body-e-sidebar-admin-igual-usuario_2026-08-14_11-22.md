# Plano: Refatorar Header para Pertencer ao Body e Redesenhar Sidebar Admin como a Sidebar de Usuários

## Pedido do Usuário
1. **Header da Admin (`.admin-topbar`)**: Deve ter o efeito de pertencer 100% ao body (fundo transparente/mesma cor do body `var(--bg)`, sem retângulos com cores separadas ou bordas pesadas de divisão).
2. **Menu Lateral da Admin (`AdminSidebar.tsx`)**: Deve ser redesenhado para ficar no mesmo padrão visual elegante da Sidebar de Usuários (`SideMenu.tsx`), com foto/iniciais do perfil no topo, grupos de botões em cards translúcidos `bg-white/5` / `bg-muted/30` com cantos arredondados (`rounded-2xl`), divisores sutis e fonte Montserrat.

---

## Análise do Estado Atual vs. Desejado

### Header (`.admin-topbar`)
- **Atual**: Possuía retângulos de cor e bordas pesadas.
- **Desejado**:
  - `background: transparent;` (ou `var(--bg)` fluído contínuo com a página).
  - `border-bottom: 1px solid transparent;` (sem linha divisória pesada).
  - Os controles (busca, ordenação, filtro, perfil, novo broadcast) flutuam suavemente sobre o fundo do body.

### Menu Lateral Admin (`AdminSidebar.tsx`)
- **Inspiração (`SideMenu.tsx` de Usuários)**:
  - Fundo: `background: var(--bg)` contínuo.
  - Perfil no topo: Foto circular do administrador (ou inicial) + Nome + subtítulo "Administrador".
  - Seções agrupadas em blocos estilizados `rounded-2xl bg-muted/20 border border-border/30` ou `bg-white/5`.
  - Itens de menu com ícone, label e divisor sutil (`border-b border-border/20`).
  - Item Ativo: `bg-accent/15 text-accent font-semibold`.
  - Item Hover: `hover:bg-muted/40 transition-colors`.
  - Estilo minimalista, limpo e integrado.

---

## Arquivos que serão modificados

- `dashboard/src/components/admin/AdminSidebar.tsx` (redesenhar layout no padrão do `SideMenu.tsx`)
- `dashboard/src/index.css` (remover bordas/fundos pesados do `.admin-sidebar` e `.admin-topbar`, tornando o header fluido com o body e a sidebar igual à de usuários)

---

## Passos Detalhados

### Passo 1 — Redesenhar `AdminSidebar.tsx` (estilo `SideMenu.tsx`)
1. Adicionar cabeçalho de perfil no topo da sidebar (foto `adminAvatar` ou inicial circular + nome do administrador `adminName` + tag/subtítulo "Administrador").
2. Agrupar os itens do menu (Principal, Operações, Premium, Sistema) dentro de containers em estilo card arredondado (`rounded-2xl bg-muted/30 dark:bg-white/5 p-1.5 space-y-1`).
3. Estilizar cada item de navegação com a estética do `SideMenu.tsx`: `flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-colors`.
4. Aplicar efeito ativo `bg-accent/15 text-accent font-semibold` e hover `hover:bg-muted/50 dark:hover:bg-white/10`.

### Passo 2 — Atualizar CSS do Header (`.admin-topbar`) no `index.css`
1. Remover fundos retangulares destacados do `.admin-topbar`.
2. Definir:
   ```css
   .admin-layout-v2 .admin-topbar {
     position: sticky;
     top: 0;
     z-index: 40;
     min-height: 64px;
     height: 64px;
     padding: 0 24px;
     border-bottom: 1px solid transparent;
     background: var(--bg);
     color: var(--text);
   }
   ```
3. O header parecerá 100% integrado e pertencente ao corpo da página (`body`), sem separação artificial.

### Passo 3 — Atualizar CSS do Sidebar (`.admin-sidebar`) no `index.css`
1. Configurar `.admin-sidebar` para ter fundo `var(--bg)`, borda sutil `border-r border-border/30` e fonte Montserrat.
2. Garantir que a sidebar recolhida (collapsed) e expandida (open) mantenham os mesmos cantos e fluidez visual.

---

## Como Testar
```bash
cd dashboard && npm run build
```
- Verificar a compilação limpa do Vite.
- Testar visualmente: o header estará integrado ao body e a sidebar admin terá o visual idêntico ao menu lateral de usuários.

## Rollback
```bash
git checkout dashboard/src/components/admin/AdminSidebar.tsx dashboard/src/index.css
```
