# Plano: Refatoração dos Painéis Admin e Abas Restantes (Anti-AI-Slop)

## Pedido do usuário
Reaplicar integralmente as diretrizes do Anti-AI-Slop UI Design Prompt nos componentes restantes do Admin (`ActivityFeed.tsx`, `SystemHealth.tsx`, `FinanceSummary.tsx`, `AlertsPanel.tsx`, `AdminAuditTab.tsx`, `AdminNoticeTab.tsx`, `AdminConfigTab.tsx`, `AdminPremiumFeaturesTab.tsx`, `AdminSubscriptionsTab.tsx`, `AdminMTProtoAccountsTab.tsx`, `ContaTelegramTab.tsx`, `ScheduleTab.tsx`).

## Objetivo
Polir 100% dos painéis administrativos e abas para garantir que todas as caixas de ícones decorativos, bordas duplas, badges excessivos, subtítulos desnecessários e envoltórios de cards tenham sido completamente eliminados.

## Contexto atual
- `OperationsOverview` e `MetricCard` já foram refatorados.
- Painéis secundários do admin (`ActivityFeed`, `SystemHealth`, `FinanceSummary`, `AlertsPanel`) ainda utilizam envoltórios de ícones coloridos e estilos `.admin-card` legados com sombras.

## Arquivos analisados
- `dashboard/src/components/admin/ActivityFeed.tsx`
- `dashboard/src/components/admin/SystemHealth.tsx`
- `dashboard/src/components/admin/FinanceSummary.tsx`
- `dashboard/src/components/admin/AlertsPanel.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/AdminPremiumFeaturesTab.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/components/AdminMTProtoAccountsTab.tsx`
- `dashboard/src/components/ContaTelegramTab.tsx`
- `dashboard/src/components/ScheduleTab.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/admin/ActivityFeed.tsx`
- `dashboard/src/components/admin/SystemHealth.tsx`
- `dashboard/src/components/admin/FinanceSummary.tsx`
- `dashboard/src/components/admin/AlertsPanel.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/components/ContaTelegramTab.tsx`

## Estratégia de implementação

1. **Painéis do Dashboard Admin (`ActivityFeed`, `SystemHealth`, `FinanceSummary`, `AlertsPanel`)**:
   - Substituir a classe `.admin-card` por `border border-border bg-card rounded-md p-4`.
   - Remover ícones dentro de círculos/quadrados coloridos (`.activity-item-icon`, `.finance-stat-icon`, `.alert-item-icon`).
   - Usar ícones inline sutis de 14px/15px e listas com `divide-y divide-border`.

2. **Abas Administrativas (`AdminNoticeTab`, `AdminAuditTab`, `AdminConfigTab`, `ContaTelegramTab`)**:
   - Remover subtítulos óbvios e duplicados.
   - Ajustar raios de curvatura para `rounded-md` e garantir que o formulário/tabela siga o grid limpo.

3. **Validação**:
   - Compilar o frontend (`npm run build`) para verificar integridade total do TypeScript e bundle.

## Passos detalhados
1. Atualizar `ActivityFeed.tsx`.
2. Atualizar `SystemHealth.tsx`.
3. Atualizar `FinanceSummary.tsx`.
4. Atualizar `AlertsPanel.tsx`.
5. Atualizar `AdminNoticeTab.tsx`, `AdminAuditTab.tsx`, `AdminConfigTab.tsx` e `ContaTelegramTab.tsx`.
6. Executar `npm run build` para garantir zero regressões.

## Riscos
- **Baixo**: Alteração estritamente estética de React/Tailwind.

## Impactos esperados
- Coerência visual de 100% da aplicação com as regras do Anti-AI-Slop.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/components/admin/
```
