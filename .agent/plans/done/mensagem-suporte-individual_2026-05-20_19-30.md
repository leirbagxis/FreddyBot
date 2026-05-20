# Plano: Mensagem Direta do Suporte via Dashboard Admin

## Pedido do usuário
Adicionar suporte na dashboard admin para mandar uma mensagem para um usuário específico pelo ID, incluindo um cabeçalho fixo tipo "# Mensagem do Suporte".

## Objetivo técnico
1. Atualizar o payload de broadcast no backend para aceitar um `target_id`.
2. Modificar o controlador de notice para suportar o alvo `single` (envio individual).
3. Adicionar o cabeçalho formatado `# 📨 Mensagem do Suporte` ao texto da mensagem quando enviada individualmente.
4. Atualizar a interface do Dashboard Admin (aba Broadcast) para permitir a seleção de "Usuário Único" e a entrada do ID correspondente.

## Contexto atual
Atualmente o sistema de "Notice" (Broadcast) suporta apenas disparos em massa para "Todos", "Canais" ou "Usuários". Não há uma forma de direcionar uma mensagem via dashboard para um único ID de usuário.

## Arquivos analisados
- `internal/api/controllers/adminController/getAllUserAdminController.go`: Contém a lógica de processamento do broadcast.
- `dashboard/src/components/AdminNoticeTab.tsx`: Componente UI do broadcast.
- `dashboard/src/api.ts`: Interface de chamadas API.

## Arquivos que poderão ser modificados
- `internal/api/controllers/adminController/getAllUserAdminController.go`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/App.tsx`
- `dashboard/src/api.ts`

## Estratégia de implementação

### 1. Backend (Go)
- Expandir a struct `NoticeRequest` para incluir o campo `TargetID (int64)`.
- No método `dispatchNotice`, adicionar o caso `single`.
- Se o alvo for `single`, validar se o `TargetID` foi fornecido.
- Prender o cabeçalho `# 📨 <b>MENSAGEM DO SUPORTE</b>\n\n` ao texto original (convertido de Markdown).
- Colocar o job na `BroadcastQueue` direcionado apenas para aquele ID.

### 2. Frontend (React)
- Atualizar a interface `NoticeRequest` no `api.ts`.
- No `AdminNoticeTab.tsx`:
    - Adicionar a opção "Individual" no grupo de botões de Público-Alvo.
    - Exibir um campo de entrada de texto para o ID do usuário apenas quando "Individual" estiver selecionado.
    - Validar se o ID foi preenchido antes de liberar o botão de disparo.
    - Atualizar a pré-visualização para mostrar como o cabeçalho do suporte ficará.

## Passos detalhados

1. **Modificar `internal/api/controllers/adminController/getAllUserAdminController.go`**:
    - Atualizar `NoticeRequest` struct.
    - Implementar a lógica de alvo `single`.

2. **Modificar `dashboard/src/api.ts`**:
    - Atualizar as interfaces e a função `sendAdminNotice`.

3. **Modificar `dashboard/src/App.tsx`**:
    - Adicionar estado para `noticeTargetId` e repassar para o componente.

4. **Modificar `dashboard/src/components/AdminNoticeTab.tsx`**:
    - Adicionar o novo botão de alvo.
    - Adicionar o input do ID.
    - Atualizar `renderPreview`.

## Riscos
- **IDs Inválidos**: O envio pode falhar se o administrador digitar um ID que nunca interagiu com o bot. O log capturará o erro, mas a UI deve avisar sobre o sucesso da "postagem na fila".

## Impactos esperados
- Administradores poderão responder ou contatar usuários diretamente via Dashboard com uma formatação oficial de suporte.

## Como testar

### Build
```bash
go build -v ./cmd/FreddyBot/... && cd dashboard && npm run build
```

### Testes
1. Acessar Dashboard Admin -> Broadcast.
2. Selecionar "Individual".
3. Digitar um ID de teste (o seu próprio ID).
4. Escrever uma mensagem e clicar em disparar.
5. Verificar se no Telegram a mensagem chegou com o cabeçalho "# 📨 MENSAGEM DO SUPORTE".

## Rollback
Reverter as alterações nos controladores e componentes.
