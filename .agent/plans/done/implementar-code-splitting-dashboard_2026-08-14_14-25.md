# Plano: Implementar Code Splitting e Otimização de Bundles na Dashboard

## Pedido do usuário
O usuário solicitou aplicar a otimização de Code Splitting na aplicação para dividir o bundle JavaScript buildado por funcionalidade e por biblioteca, tornando a abertura do Telegram Mini App ultrarrápida no celular.

## Objetivo
Reduzir o tamanho do carregamento inicial da Dashboard de ~1.03 MB para ~180-200 KB no primeiro carregamento do celular, dividindo os componentes administrativos e abas secundárias sob demanda (`React.lazy`) e agrupando bibliotecas estáticas (`manualChunks` no Vite).

## Contexto atual
- Atualmente em `dashboard/src/App.tsx` e `AdminDashboard.tsx`, todos os componentes (incluindo o painel admin inteiro com logs, usuários, broadcasting e relatórios) são importados estaticamente.
- No `vite.config.ts`, não há separação de `manualChunks`, o que gera um único arquivo JS monolítico `index-XXXXX.js` com aviso de limite de tamanho do Rollup.

## Arquivos analisados
- `dashboard/vite.config.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`

## Arquivos que poderão ser modificados
- `dashboard/vite.config.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`

## Estratégia de implementação

### 1. Configurar `manualChunks` no `vite.config.ts`
Agrupar bibliotecas pesadas de node_modules em pacotes organizados para cache permanente no navegador:
- `vendor-react`: `react`, `react-dom`, `scheduler`
- `vendor-icons`: `lucide-react`
- `vendor-ui`: `@base-ui/react`, `@floating-ui`, `tailwind-merge`

### 2. Otimizar `App.tsx` com `React.lazy` e `<Suspense>`
Transformar importações estáticas pesadas em importações dinâmicas:
- `AdminDashboard` (carregado apenas se o usuário for Admin e abrir a área admin)
- `ScheduleTab` (carregado ao abrir a aba Agendamentos)
- `PremiumConfigTab` (carregado ao abrir a aba Premium)
- `ContaTelegramTab` (carregado ao abrir a aba Conta)

### 3. Otimizar `AdminDashboard.tsx` com `React.lazy` e `<Suspense>`
Transformar abas administrativas secundárias em dinâmicas:
- `AdminLogsTab`
- `AdminNoticeTab` (broadcast)
- `AdminAuditTab`
- `AdminMTProtoAccountsTab`
- `AdminSubscriptionsTab`
- `AdminPremiumFeaturesTab`

## Passos detalhados
1. Atualizar `dashboard/vite.config.ts` com a propriedade `build.rollupOptions.output.manualChunks`.
2. Atualizar `App.tsx` substituindo as importações diretas por `React.lazy` e envolvendo os componentes dinâmicos com `<Suspense fallback={...}>`.
3. Atualizar `AdminDashboard.tsx` substituindo as importações das abas secundárias admin por `React.lazy` com `<Suspense fallback={...}>`.
4. Executar a compilação com `cd dashboard && npm run build` para validar a diminuição dos chunks e ausência de erros.

## Riscos
- Mínimo. Telas divididas com `lazy` exibem um indicador suave de carregamento durante a transição inicial e são armazenadas em cache local.

## Impactos esperados
- Redução de ~80% no tamanho do JavaScript inicial baixado no Telegram Mini App.
- Eliminação do alerta de tamanho do bundle do Vite/Rollup.
- Abertura instantânea da interface mobile para donos de canal.

## Compatibilidade
- Linux, macOS, Windows, Docker, Navegadores Web, Telegram Mini App iOS/Android

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
Restaurar os arquivos `vite.config.ts`, `App.tsx` e `AdminDashboard.tsx`.
