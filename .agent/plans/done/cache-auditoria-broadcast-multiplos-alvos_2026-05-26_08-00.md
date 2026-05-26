# Plano: cache-auditoria-broadcast-multiplos-alvos

## Pedido do usuário
Adicionar cache na aba de auditoria do painel admin para manter o resultado da varredura do XavolaBot ao trocar de página/aba ou abrir um usuário, até rodar nova varredura ou recarregar a página. Também adicionar suporte para enviar mensagem para canais específicos por ID, sem o título `# MENSAGEM DO SUPORTE`, e permitir múltiplos IDs de canais ou usuários.

## Objetivo
Persistir o resultado da auditoria durante a sessão atual do React sem usar storage persistente, e ampliar o broadcast admin para alvos específicos múltiplos:
- usuários individuais ou múltiplos com cabeçalho de suporte;
- canais individuais ou múltiplos sem cabeçalho de suporte;
- manter os alvos atuais `all`, `users` e `channels`.

## Contexto atual
A aba `AdminAuditTab` mantém `results` em estado local. Quando o admin abre outro painel, usuário ou canal, esse componente pode ser desmontado e o resultado da varredura é perdido.

O broadcast admin usa `noticeTarget` com `channels`, `users`, `all` e `single`. O alvo `single` envia para um único ID e o backend sempre injeta o cabeçalho `# 📨 <b>MENSAGEM DO SUPORTE</b>`. Não existe diferenciação entre usuário específico e canal específico, nem envio para vários IDs.

## Arquivos analisados
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/App.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`
- `internal/api/controllers/adminController/getAllUserAdminController.go`
- `internal/container/appContainer.go`
- `internal/api/routes/routes.go`

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/api.ts`
- `internal/api/controllers/adminController/getAllUserAdminController.go`

## Estratégia de implementação
1. Elevar o estado da auditoria para `DashboardContent`/`AdminDashboard`, mantendo os dados em memória enquanto a página não for recarregada.
2. Fazer `AdminAuditTab` receber `results`, `setResults`, estado de loading e handler de varredura por props, preservando também a limpeza local após exclusão em massa.
3. Expandir o tipo de alvo do broadcast no frontend para separar alvo específico de usuário e canal, com suporte a múltiplos IDs em um campo de texto.
4. Atualizar o payload da API para enviar `targetIds: number[]`, mantendo compatibilidade com `targetId`.
5. Alterar o backend para aceitar:
   - `single` / usuário específico com cabeçalho de suporte;
   - `user_ids` para múltiplos usuários com cabeçalho de suporte;
   - `channel_ids` para múltiplos canais sem cabeçalho de suporte;
   - os alvos globais existentes sem mudança.
6. Ajustar preview, validação e mensagem de confirmação no painel para refletir canal/usuário e múltiplos IDs.

## Passos detalhados
1. Criar estado `auditResults` e `auditLoading` no componente pai da dashboard admin.
2. Mover `fetchAuditCheckBot` para um handler no pai ou passar estado controlado para `AdminAuditTab`.
3. Atualizar `AdminAuditTab` para não inicializar `results` internamente e continuar mostrando os dados ao voltar para a aba.
4. Atualizar `NoticeRequest` em `dashboard/src/api.ts` para incluir `targetIds?: number[]` e novos valores de alvo.
5. Atualizar `App.tsx` para montar `targetIds` a partir de texto com IDs separados por vírgula, espaço ou quebra de linha.
6. Atualizar `AdminNoticeTab` para oferecer opções claras: todos, canais, usuários, usuário(s) por ID, canal(is) por ID.
7. Remover o cabeçalho de suporte no preview quando o alvo for canal específico.
8. Atualizar `NoticeRequest` no backend e `dispatchNotice` para montar listas de IDs e cabeçalho conforme o tipo de alvo.
9. Rodar `npm run build`, `gofmt`, `git diff --check`, `go test ./...` e `go build ./cmd/FreddyBot/main.go` quando possível.

## Riscos
- Múltiplos IDs inválidos precisam ser filtrados/recusados no frontend para evitar disparo parcial acidental.
- Envio para canais por ID pode falhar se o bot não estiver no canal ou não tiver permissão, mas o worker atual já registra erro por job.
- Como o endpoint responde antes do processamento assíncrono, a UI continuará indicando que o disparo foi iniciado, não que todos os alvos receberam.
- A toolchain Go local pode continuar sem `vet`/`compile`, bloqueando `go test` e `go build`.

## Impactos esperados
- Resultado da auditoria permanece visível ao alternar abas ou abrir/voltar de usuário enquanto a página não for recarregada.
- Clicar em nova varredura substitui o cache em memória.
- Recarregar a página limpa o cache da auditoria.
- Broadcast admin aceita múltiplos IDs de usuários ou canais.
- Mensagem para canais específicos não recebe `# MENSAGEM DO SUPORTE`.
- Mensagem para usuários específicos mantém o cabeçalho de suporte.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
npm run build
go build ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./...
```

### Execução
```bash
# Dashboard admin
/admin
```

Validar:
- Rodar auditoria, abrir usuário do resultado, voltar para auditoria e confirmar que os resultados continuam.
- Rodar nova auditoria e confirmar substituição dos resultados.
- Recarregar a página e confirmar que o cache sumiu.
- Enviar broadcast para múltiplos usuários por ID e confirmar cabeçalho de suporte.
- Enviar broadcast para múltiplos canais por ID e confirmar ausência do cabeçalho.

## Rollback
Reverter os arquivos modificados listados acima. O sistema voltará a perder o estado local da auditoria ao desmontar a aba e o broadcast voltará a aceitar apenas os alvos antigos.

## Observações
Há mudanças pendentes anteriores no worktree, incluindo normalização de URLs de botões, alteração do `/info` com miniapp e o arquivo não rastreado `Release`. Essas alterações não devem ser revertidas.
