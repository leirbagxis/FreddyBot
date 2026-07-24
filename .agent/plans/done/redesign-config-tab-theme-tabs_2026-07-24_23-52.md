# Plano: redesign-config-tab-theme-tabs

## Pedido do usuário
Três tarefas: (1) refazer a aba de Configurações da dashboard admin (está muito poluída), (2) corrigir o toggle de mudança de tema que não está funcionando, (3) redesenhar o menu tabs.

## Objetivo
- Migrar AdminConfigTab de `.cfg-*` classes CSS customizadas para componentes shadcn (Card, CardHeader, CardContent, etc.)
- Corrigir o Theme Toggle no AdminTopbar: ícone não reflete estado telegram, data-theme estático no AdminLayout, duas instâncias de useTheme() dessincronizadas
- Redesenhar navegação de abas (AdminSidebar): melhorar visual e usabilidade

## Contexto atual
- AdminConfigTab.tsx usa 24 classes `.cfg-*` customizadas definidas em index.css, mistura shadcn (Switch, Button, Input) com HTML puro (textarea)
- Tema usa hook `useTheme()` sem context — App.tsx e AdminTopbar.tsx chamam `useTheme()` separadamente, criando duas instâncias de estado que podem dessincronizar
- AdminLayout lê `data-theme` de `<html>` uma vez no render e nunca atualiza
- AdminTopbar mostra apenas Sun/Moon, sem indicar estado `telegram`
- AdminSidebar (em `components/admin/`) já tem design seccionado, mas a versão antiga em `components/AdminSidebar.tsx` ainda existe e pode causar confusão
- AdminDashboard renderiza tabs condicionalmente com subtítulos descritivos

## Arquivos analisados
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminSidebar.tsx` (antigo, em `components/`)
- `dashboard/src/components/admin/AdminSidebar.tsx` (novo, em `components/admin/`)
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/components/ui/card.tsx`
- `dashboard/src/components/ui/textarea.tsx`
- `dashboard/src/index.css` (regras `.cfg-*` e `.admin-*`)

## Arquivos que poderão ser modificados
- `dashboard/src/components/AdminConfigTab.tsx` — redesign completo
- `dashboard/src/components/admin/AdminTopbar.tsx` — fix ícone do toggle
- `dashboard/src/components/admin/AdminLayout.tsx` — remover data-theme estático
- `dashboard/src/components/AdminSidebar.tsx` — remover arquivo obsoleto
- `dashboard/src/components/admin/AdminSidebar.tsx` — redesign visual
- `dashboard/src/index.css` — remover classes `.cfg-*` após migração

## Estratégia de implementação

### 1. AdminConfigTab redesign
Substituir toda a estrutura de `.cfg-*` classes por Cards shadcn:
- Cada seção (Sistema, Legendas, PostBuilder) vira um Card com CardHeader + CardContent
- Loading state mantém spinner, mas usando shadcn Skeleton se disponível
- Toggle rows usam flex row com Switch + label/descrição
- Textarea customizado vira shadcn Textarea
- Hint com inline code mantido mas estilizado com classes Tailwind
- Save button no footer do último Card ou num CardFooter

### 2. Theme toggle fix
- AdminTopbar: corrigir ícone para mostrar Sun (dark), Moon (light), Send (telegram) — igual App.tsx
- AdminLayout: remover `data-theme` estático (redundante, CSS vars globais já funcionam), ou usar `useTheme()` reativamente

### 3. Menu tabs redesign
- Remover `components/AdminSidebar.tsx` obsoleto (versão antiga de 27 linhas)
- Melhorar `components/admin/AdminSidebar.tsx`: ajustar padding, hover states, transições
- Adicionar indicadores visuais mais claros para aba ativa
- Garantir que AdminLayout importe do path correto (`./admin/AdminSidebar`)

### 4. CSS cleanup
- Remover todas as regras `.cfg-*` de index.css após confirmar que AdminConfigTab não as usa mais

## Passos detalhados

1. Redesign AdminConfigTab:
   a. Substituir estrutura por Cards shadcn (Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter)
   b. Sistema section: toggle rows com Switch + label/desc com shadcn Card
   c. Legendas section: RichTextEditor dentro de Card
   d. PostBuilder section: toggle + Input + Textarea + hint, tudo em Card
   e. Adicionar save button no CardFooter do último card
   f. Remover imports de `.cfg-*` CSS (só Tailwind/shadcn)

2. Fix theme toggle:
   a. Atualizar AdminTopbar para mostrar ícone correto para cada estado (telegram → Send, dark → Sun, light → Moon)
   b. Remover `data-theme` estático do AdminLayout (ou torná-lo reativo)

3. Redesign tabs:
   a. Remover `src/components/AdminSidebar.tsx` (arquivo obsoleto)
   b. Melhorar visual do `admin/AdminSidebar.tsx` com animações e hover states
   c. Verificar que AdminLayout importa do caminho correto

4. CSS cleanup:
   a. Remover blocos `.cfg-*` de index.css
   b. Verificar build

5. Verificar build (tsc --noEmit)

## Riscos
- AdminConfigTab usa RichTextEditor que é componente custom — manter integração
- AdminTopbar importa useTheme de hooks/useTheme — se criarmos ThemeProvider depois, precisaremos migrar
- AdminDashboard chama `AdminConfigTab` sem props — manter essa interface
- Alguns toggles disparam save imediato via handleSave() — preservar esse comportamento

## Impactos esperados
- AdminConfigTab visualmente mais limpo e consistente com o resto do admin CRM
- Theme toggle funcional e com feedback visual correto
- Sidebar mais polida
- Remoção de ~50 linhas de CSS customizado de index.css
- Código mais fácil de manter (shadcn components)

## Compatibilidade
- Linux ✓
- macOS ✓
- Windows ✓
- Docker ✓
- CI/CD ✓

## Como testar

### Build
```bash
cd dashboard && npx tsc --noEmit
```

### Execução
```bash
cd dashboard && npm run dev
```

## Rollback
```bash
git checkout -- src/components/AdminConfigTab.tsx src/components/admin/AdminTopbar.tsx src/components/admin/AdminLayout.tsx src/index.css src/components/AdminSidebar.tsx src/components/admin/AdminSidebar.tsx
```

## Observações
- Manter mesma API de save (handleSave com overrides) e toggle imediato
- Manter comportamento de loading/saving states
- Não quebrar o RichTextEditor
