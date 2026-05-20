# Plano: Agrupamento de Resultados /checkbot e Nova Aba de Auditoria

## Pedido do usuário
1. Melhorar o comando `/checkbot` para que os resultados sejam agrupados por usuário, enviando uma lista de canais por dono.
2. Criar uma nova aba na Dashboard Admin ("Auditoria") que realize a mesma checagem e mostre os resultados agrupados por usuário na interface.

## Objetivo técnico
1. Refatorar o handler `CheckBotAdminHandlerTelego` para agrupar canais encontrados por `OwnerID`.
2. Implementar um novo endpoint de API `/api/admin/audit/checkbot` que performa a varredura e retorna um JSON estruturado por usuário.
3. Adicionar uma aba "Auditoria" no Dashboard React com visualização clara dos canais "clonados" (onde o @XavolaBot é admin).

## Contexto atual
Atualmente o comando `/checkbot` retorna uma lista única e extensa de canais. No dashboard, não existe essa funcionalidade de varredura ativa. O bot gerencia canais de diversos donos e o administrador precisa saber rapidamente quais canais de qual usuário ainda possuem o bot antigo ativo.

## Arquivos analisados
- `internal/telegram/handlers/commands/admin/admin_utils.go`: Contém a lógica do `/checkbot`.
- `internal/api/routes/routes.go`: Registro de rotas.
- `dashboard/src/App.tsx`: Gerenciamento de abas.
- `dashboard/src/components/AdminDashboard.tsx`: Renderização do painel admin.

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/commands/admin/admin_utils.go`
- `internal/api/routes/routes.go`
- `internal/api/controllers/adminController/auditController.go` (Novo)
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminAuditTab.tsx` (Novo)
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`

## Estratégia de implementação

### 1. Refatoração do Comando Bot
- No handler do `/checkbot`, após a varredura concorrente, os canais encontrados serão inseridos em um `map[int64][]Channel`.
- Iterar sobre o mapa e para cada usuário, buscar seu nome e enviar uma mensagem formatada: "Canais de [Nome]: \n - Canal A \n - Canal B".

### 2. Backend (API)
- Criar `AuditController`.
- O método `GetCheckBotAudit` fará a mesma varredura (usando workers para ser rápido).
- Retornará uma lista de objetos: `{ userId, firstName, channels: [...] }`.

### 3. Frontend (Dashboard)
- Adicionar a aba `audit` no `App.tsx`.
- Criar o componente `AdminAuditTab` que:
  - Exibe um botão "Iniciar Varredura".
  - Mostra um loading enquanto a API processa (pode demorar alguns segundos).
  - Exibe os resultados em "cards" ou "acordeões" agrupados por usuário.

## Passos detalhados

1. **Modificar `internal/telegram/handlers/commands/admin/admin_utils.go`**:
    - Alterar `CheckBotAdminHandlerTelego` para agrupar por `OwnerID`.
    - Ajustar o envio de mensagens para ser por usuário.

2. **Criar `internal/api/controllers/adminController/auditController.go`**:
    - Implementar a lógica de varredura e retorno JSON.

3. **Registrar Rota**:
    - Em `internal/api/routes/routes.go`, adicionar `adminRoute.GET("/audit/checkbot", auditController.GetCheckBotAudit)`.

4. **Atualizar Dashboard**:
    - `types.ts`: Adicionar interfaces para os resultados da auditoria.
    - `api.ts`: Adicionar função `fetchAuditCheckBot`.
    - `App.tsx`: Adicionar a aba e o estado necessário.
    - `AdminDashboard.tsx`: Renderizar a nova aba.
    - `AdminAuditTab.tsx`: Implementar a UI da auditoria.

## Riscos
- **Timeout na API**: Se houver milhares de canais, o request HTTP pode dar timeout. O Telegram API é o gargalo. 
  - *Mitigação:* Usar um número maior de workers ou considerar um processo em background com WebSockets/Polling (mas para começar, workers rápidos devem bastar).

## Impactos esperados
- Melhor visibilidade para o administrador sobre o uso do bot antigo por cada cliente/usuário.
- Facilidade de auditoria via interface gráfica.

## Como testar

### Build
```bash
go build -v ./cmd/FreddyBot/... && cd dashboard && npm run build
```

### Testes
1. Rodar `/checkbot` no bot e verificar se as mensagens chegam agrupadas por dono.
2. Abrir a Dashboard -> Aba Auditoria -> Clicar em Varrer.
3. Verificar se a lista de usuários e seus canais com o XavolaBot aparece corretamente.

## Rollback
Reverter as mudanças nos arquivos de admin e dashboard.
