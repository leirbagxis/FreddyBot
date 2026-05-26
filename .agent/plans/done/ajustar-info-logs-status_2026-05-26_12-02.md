# Plano: ajustar-info-logs-status

## Pedido do usuário
Adicionar no comando `/info` um botão para abrir os logs do canal e ajustar as cores dos status dos logs na dashboard: sucesso verde, info azul, erro vermelho e ignorado laranja.

## Objetivo
Melhorar o acesso direto aos logs por canal a partir do comando admin `/info` e deixar os badges de status visualmente coerentes na aba Logs da dashboard admin.

## Contexto atual
O comando `/info` em `internal/telegram/handlers/commands/admin/admin_channels.go` já busca o canal, valida informações e envia uma mensagem com botão WebApp "Abrir canal na dashboard".

A dashboard admin já possui a aba `Logs` em `dashboard/src/components/AdminLogsTab.tsx`, com badges de status. Hoje o status `success` usa verde, `error` usa vermelho, mas `info` e `skipped` caem em estilo secundário/mutado.

A rota admin da dashboard usa `/admin/dash`, e a aba Logs busca `/api/admin/logs` com filtros como `channelId`. Para abrir logs filtrados direto pelo botão, será necessário suportar querystring na dashboard, por exemplo `/admin/dash?tab=logs&channelId=<id>`.

## Arquivos analisados
- `internal/telegram/handlers/commands/admin/admin_channels.go`
- `dashboard/src/components/AdminLogsTab.tsx`
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/commands/admin/admin_channels.go`
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminLogsTab.tsx`
- possivelmente `dashboard/src/index.css`

## Estratégia de implementação
1. No `/info`, gerar uma URL WebApp para `/admin/dash?tab=logs&channelId=<channelID>` usando a base já utilizada por `auth.GenerateMiniAppUrl` ou derivando da URL gerada existente.
2. Adicionar um segundo botão inline WebApp abaixo de "Abrir canal na dashboard", com texto "Abrir logs do canal".
3. Na dashboard admin, ler `tab=logs` da querystring ao iniciar e selecionar a aba Logs automaticamente quando existir.
4. Passar o `channelId` da querystring para `AdminLogsTab` como filtro inicial.
5. Em `AdminLogsTab`, iniciar os filtros com `channelId` quando recebido e carregar os logs já filtrados.
6. Ajustar o mapeamento visual dos status:
   - `success` -> verde
   - `info` -> azul/accent
   - `error` -> vermelho
   - `skipped` -> laranja/warning
7. Se não houver classe CSS existente para danger/warning/accent em badges, adicionar classes pequenas e reutilizáveis em `index.css`.

## Passos detalhados
1. Adicionar helper local para construir a URL de logs do canal no comando `/info`, preservando segurança e HTML escaping.
2. Alterar o `ReplyMarkup` do `/info` para conter dois botões WebApp em linhas separadas.
3. Adicionar leitura de `URLSearchParams(window.location.search)` em `dashboard/src/App.tsx` para definir aba inicial `logs` quando `tab=logs`.
4. Criar estado/prop para `initialLogsChannelId` e repassar para `AdminDashboard` e `AdminLogsTab`.
5. Atualizar `AdminLogsTab` para aceitar `initialChannelId?: string`.
6. Ajustar função/classe de status do badge.
7. Atualizar CSS se necessário para `badge-danger`, `badge-warning` e `badge-accent`.
8. Rodar build da dashboard e tentar validações Go, registrando qualquer limitação do toolchain local.

## Riscos
- A dashboard precisa aceitar querystring sem quebrar as abas atuais.
- Se a URL WebApp admin for montada incorretamente, o botão abre a dashboard sem filtro ou falha no Telegram.
- O estado inicial dos logs não deve impedir o admin de limpar ou trocar o filtro manualmente.

## Impactos esperados
- Admin/owner consegue sair do `/info <channelID>` direto para os logs daquele canal.
- A aba Logs fica mais legível: info azul, ignorado laranja, erro vermelho e sucesso verde.
- Nenhuma mudança no schema do banco ou no fluxo de processamento de posts.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD
- Telegram WebApp

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
npm run build
```

### Testes
```bash
go test ./...
```

### Execução
```bash
# No Telegram, como admin/owner:
/info <channel_id>
# Clicar em "Abrir logs do canal" e confirmar a aba Logs filtrada pelo canal.
```

## Rollback
Reverter as alterações no comando `/info`, remover a leitura de querystring para logs e voltar o mapeamento antigo dos badges.

## Observações
O ambiente Go local está com toolchain quebrado em validações recentes (`no such tool "vet"` e `no such tool "compile"`). Se persistir, a validação Go ficará limitada até corrigir o toolchain.
