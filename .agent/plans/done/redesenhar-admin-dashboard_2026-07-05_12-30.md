# Plano: Redesenhar Admin Dashboard com Analytics

## Pedido do usuário
Redesenhar completamente a dashboard administrativa, transformando-a em uma dashboard real com analytics, métricas e visualizações.

## Objetivo
Criar uma experiência de dashboard administrativa moderna com cards de métricas, indicadores visuais e listas otimizadas, mantendo todas as funcionalidades existentes (gerenciar usuários, canais, broadcast, auditoria, logs, config).

## Contexto atual
O AdminDashboard atual está todo em um único componente (`AdminDashboard.tsx`) com abas e renderização condicional. Os dados incluem usuários (com info de admin, blacklist, canais) e canais. Atualmente as abas são:
- Users (lista + detalhe do usuário)
- Channels (lista)
- Notice (broadcast)
- Audit (auditoria de bot)
- Logs (logs de eventos)
- Config (configurações)

Cards com `bg-card` e `ring` foram removidos anteriormente (tudo transparente agora).

## Arquivos analisados
- src/components/AdminDashboard.tsx
- src/types.ts (AdminDashboardData, User, Channel)
- src/App.tsx (como o AdminDashboard é chamado)
- src/components/AdminNoticeTab.tsx
- src/components/AdminAuditTab.tsx
- src/components/AdminLogsTab.tsx
- src/components/AdminConfigTab.tsx

## Arquivos que poderão ser modificados
- src/components/AdminDashboard.tsx (reescrita completa)
- src/index.css (se precisar de novos estilos)

## Estratégia de implementação

Vou reescrever o AdminDashboard.tsx com:

1. **Hero Section** com 4-5 métricas em destaque:
   - Total de Usuários
   - Total de Canais
   - Média de Canais por Usuário
   - Admins
   - Blacklist

2. **Barra de distribuição** visual (quantos usuários têm 0, 1, 2, 3+ canais) usando divs como barras horizontais (CSS puro, sem lib)

3. **Top Users** — tabela compacta dos usuários com mais canais

4. **Lista de Usuários** — melhorada com buscador, filtro por canais, indicadores visuais

5. **Lista de Canais** — melhorada com buscador, informações mais ricas

6. Manter todas as abas existentes (users, channels, notice, audit, logs, config)

7. Layout responsivo, sem bordas de cards (consistente com o novo design transparente)

## Passos detalhados

1. Criar métricas computadas a partir dos dados existentes (users, channels)
2. Criar componente de bar chart horizontal (distribuição de canais)
3. Criar seção de Top Users
4. Reescrever renderUsersTab com layout mais limpo e métricas
5. Manter renderChannelsTab, renderNoticeTab, renderAuditTab, renderLogsTab, renderConfigTab funcionando
6. Ajustar header e navegação entre abas
7. Testar build

## Riscos
- Pode ficar pesado se tentar fazer gráficos elaborados — vou usar CSS puro
- Manter compatibilidade com todas as props existentes
- Não quebrar a navegação entre abas

## Impactos esperados
- AdminDashboard.tsx será significativamente maior (~600+ linhas)
- Experiência muito mais rica e informativa
- Nenhuma funcionalidade existente será removida

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

### Execução
Navegar para `/admin/dash` no navegador.

## Rollback
Reverter alterações no AdminDashboard.tsx via git:
```bash
git checkout -- src/components/AdminDashboard.tsx
```

## Observações
- Aproveitar dados que já existem sem chamar APIs extras
- Manter tema consistente (sem fundo nos cards, sem rings)
- Foco em informação relevante, não em decoração vazia
