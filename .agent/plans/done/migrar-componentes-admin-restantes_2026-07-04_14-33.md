# Plano: migrar-componentes-admin-restantes

## Pedido do usuário
Completar todos os próximos passos: migrar componentes admin restantes para shadcn/ui, remover CSS antigo, limpar variáveis CSS, build final.

## Objetivo
Remover 100% das classes CSS customizadas antigas (`.btn`, `.card`, `.input`, `.badge`, `.form-area`, `.btn-grid-wrapper`) substituindo por componentes shadcn/ui ou Tailwind, e limpar o CSS e variáveis não utilizadas.

## Contexto atual
- Já migramos 14 componentes para shadcn/ui (ConfirmModal, PermissionsCard, CaptionCard, NewPackCaptionCard, ReactionsCard, DashboardInicioTab, ButtonGrid, ConnectedAccountCard, AuthDisclaimer, AdminConfigTab, AdminDashboard parcial, ContaTelegramTab parcial, AdminAuditTab parcial, App.tsx parcial)
- Restam **37 ocorrências** de classes CSS antigas em **8 arquivos**
- CSS antigo ainda presente: `.btn-*`, `.card`, `.input`, `.badge-*`, `.form-area`, `.btn-grid-wrapper`
- Variáveis CSS potencialmente não utilizadas

## Arquivos analisados
- dashboard/src/index.css
- dashboard/src/App.tsx
- dashboard/src/components/AdminLogsTab.tsx
- dashboard/src/components/AdminDashboard.tsx
- dashboard/src/components/AdminAuditTab.tsx
- dashboard/src/components/RichTextEditor.tsx
- dashboard/src/components/AdminNoticeTab.tsx
- dashboard/src/components/ContaTelegramTab.tsx
- dashboard/src/components/ButtonGrid.tsx

## Arquivos que poderão ser modificados
- dashboard/src/components/AdminLogsTab.tsx
- dashboard/src/components/AdminDashboard.tsx
- dashboard/src/components/AdminAuditTab.tsx
- dashboard/src/components/RichTextEditor.tsx
- dashboard/src/components/AdminNoticeTab.tsx
- dashboard/src/components/ContaTelegramTab.tsx
- dashboard/src/App.tsx
- dashboard/src/components/ButtonGrid.tsx
- dashboard/src/index.css

## Estratégia de implementação
Migração ordenada por arquivo, do mais isolado ao mais complexo. Cada arquivo é migrado completamente (substituindo todas as classes antigas por shadcn/Tailwind). Ao final, remove-se o CSS morto e variáveis não utilizadas. Build verificado após cada fase.

## Passos detalhados

### Fase 1 — AdminLogsTab
- Substituir `<select className="input">` por `<Select>` shadcn
- Substituir `<button className="btn btn-primary">` por `<Button variant="default">`
- Substituir `<button className="btn btn-secondary">` por `<Button variant="secondary">`
- Substituir `<span className="badge badge-muted">` por `<Badge variant="secondary">`
- Substituir `<div className="card">` por `<Card>`
- Build check

### Fase 2 — AdminAuditTab
- Substituir `<div className="card">` por `<Card>`
- Substituir `<button className="btn btn-danger">` por `<Button variant="destructive">`
- Substituir `<button className="btn btn-ghost">` por `<Button variant="ghost">`
- Build check

### Fase 3 — AdminNoticeTab
- Substituir `<div className="card">` por `<Card>`
- Substituir `<button className="btn ...">` por `<Button variant="default">`
- Build check

### Fase 4 — ContaTelegramTab
- Substituir `<div className="card">` por `<Card>`
- Substituir `<button className="btn btn-primary">` por `<Button variant="default">`
- Build check

### Fase 5 — RichTextEditor
- Substituir `<button className="btn btn-primary btn-sm">` por `<Button variant="default" size="sm">`
- Substituir `<button className="btn btn-secondary btn-sm">` por `<Button variant="secondary" size="sm">`
- Build check

### Fase 6 — AdminDashboard (restante)
- Substituir últimos `<button className="btn btn-*">` por `<Button>`
- Substituir últimos `<div className="card">` por `<Card>`
- Build check

### Fase 7 — App.tsx (restante)
- Substituir `<div className="card">` por `<Card>`
- Substituir `<button className="btn btn-primary">` por `<Button>`
- Substituir `<span className="badge badge-*">` por `<Badge>`
- Build check

### Fase 8 — ButtonGrid (form-area + btn-grid-wrapper)
- Substituir `className="form-area"` por Tailwind classes equivalentes
- Substituir `className="btn-grid-wrapper"` por Tailwind classes equivalentes
- Remover CSS `.form-area`, `.btn-grid-wrapper`, `.btn-grid-item` do index.css
- Build check

### Fase 9 — Limpeza do CSS antigo do index.css
- Remover blocos: `.btn`, `.btn-primary`, `.btn-secondary`, `.btn-ghost`, `.btn-danger`, `.btn-sm`, `.btn-success`
- Remover bloco: `.card`
- Remover bloco: `.input`
- Remover bloco: `.badge`, `.badge-accent`, `.badge-success`, `.badge-info`, `.badge-danger`, `.badge-warning`, `.badge-muted`, `.badge-ghost`
- Remover classes auxiliares: `.admin-welcome-card`, `.admin-stat-card`, `.admin-list-item`, `.stat-card-clickable` (se não utilizadas)
- Build check

### Fase 10 — Limpeza de variáveis CSS não utilizadas
- Revisar `:root` em index.css e identificar variáveis não referenciadas em nenhum arquivo .tsx/.css
- Remover variáveis mortas
- Build check

### Fase 11 — Build final e verificação
- `npm run build` sem erros
- Verificar bundle size
- Relatório final de classes removidas

## Riscos
- Quebrar layout de AdminLogsTab (migração mais complexa com Select e filtros)
- Quebrar App.tsx (arquivo grande com muitos cards canais)
- Esquecer de importar componentes shadcn em arquivos que ainda não os importam
- Perder estilos visuais específicos (ex: botão do AdminNoticeTab com shadow e gradiente customizado)

## Impactos esperados
- 0 classes CSS antigas restantes nos componentes
- index.css reduzido significativamente
- UI consistente 100% shadcn/ui + Tailwind
- Manutenção futura simplificada

## Compatibilidade
- Linux ✓
- macOS ✓
- Windows ✓
- CI/CD ✓

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Testes
```bash
# Verificar se não há erros de tipo
npx tsc --noEmit
```

### Preview
```bash
cd dashboard && npm run dev
```

## Rollback
```bash
git checkout -- dashboard/src/components/ dashboard/src/App.tsx dashboard/src/index.css
```

## Observações
- Fases 1-8 podem ser executadas em paralelo entre arquivos diferentes (são independentes)
- A Fase 9 (limpeza CSS) só pode ocorrer após todas as fases 1-8 estarem completas
- Alguns estilos visuais específicos (ex: botão de destaque no AdminNoticeTab) devem ser preservados via Tailwind inline ou shadcn variants
