# Plano: Seletor de Cor para Botões Inline (primary, success, danger)

## Pedido do usuário
"Na parte de colocar botoes/editar botoes, eu queria que tivesse tipo 3 circulos para setar as cores que o telegram suporte, primary, green, danger. Se quiser pode conferir ai. Tbm queria que quando o usuario selecionasse a cor o botao ficaria na cor correspondente e o telegram tbm colocaria a cor no botao"

## Objetivo
Adicionar um seletor de cor com 3 círculos (Primary/azul, Success/verde, Danger/vermelho) nos formulários de criar e editar botões. A cor selecionada será persistida no backend e aplicada tanto visualmente no grid do dashboard quanto enviada ao Telegram via o campo `style` do `InlineKeyboardButton` (Bot API 9.4).

## Contexto atual
- O Telegram Bot API 9.4 suporta o campo `style` em `InlineKeyboardButton` com valores: `"primary"` (azul), `"success"` (verde), `"danger"` (vermelho).
- A biblioteca `telego v1.9.0` já tem o campo `Style string` na struct `InlineKeyboardButton`.
- Atualmente, nem o model Go (`models.Button`), nem o frontend (`types.ts > Button`), nem a API (`ButtonCreateRequest`) possuem campo de cor/estilo.

## Arquivos analisados
- `internal/database/models/models.go` — struct `Button`
- `internal/database/repositories/button.go` — `UpdateButton`
- `internal/core/services/buttons.go` — `CreateButton`, `UpdateButton`
- `internal/api/types/buttons.go` — `ButtonCreateRequest`
- `internal/api/dto/dto.go` — `ButtonDTO`
- `internal/api/dto/mapper.go` — `ToButtonDTO`
- `internal/api/controllers/ButtonsController.go`
- `internal/telegram/events/channelPost/keyboard.go` — `CreateInlineKeyboardTelego`
- `internal/telegram/executor/executor.go` — `InlineKeyboardButton`
- `internal/telegram/executor/botapi.go` — `toTelegoKeyboard`
- `internal/telegram/events/channelPost/dispatch_telego.go` — `toExecutorKeyboard`
- `dashboard/src/types.ts` — `Button`
- `dashboard/src/api.ts` — `updateButton`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/App.tsx` — `handleEditButton`

## Arquivos que poderão ser modificados

### Backend (Go)
1. `internal/database/models/models.go` — adicionar campo `Style` ao struct `Button` e `CustomCaptionButton`
2. `internal/api/types/buttons.go` — adicionar campo `Style` ao `ButtonCreateRequest`
3. `internal/api/dto/dto.go` — adicionar campo `Style` ao `ButtonDTO`
4. `internal/api/dto/mapper.go` — mapear `Style` em `ToButtonDTO` e `ToCustomCaptionButtonDTO`
5. `internal/core/services/buttons.go` — passar `Style` na criação e atualização de botões
6. `internal/database/repositories/button.go` — incluir `style` no update
7. `internal/telegram/events/channelPost/keyboard.go` — setar `Style` no `telego.InlineKeyboardButton`
8. `internal/telegram/executor/executor.go` — adicionar `Style` ao `InlineKeyboardButton`
9. `internal/telegram/executor/botapi.go` — mapear `Style` no `toTelegoKeyboard`
10. `internal/telegram/events/channelPost/dispatch_telego.go` — mapear `Style` no `toExecutorKeyboard`

### Frontend (Dashboard)
11. `dashboard/src/types.ts` — adicionar campo `style` ao `Button`
12. `dashboard/src/components/ButtonGrid.tsx` — adicionar seletor de 3 círculos coloridos nos forms de add/edit, exibir cor no grid

## Estratégia de implementação

### Passo 1 — Backend: Model + DB (AutoMigrate)
Adicionar `Style string` ao `models.Button` e `models.CustomCaptionButton`. GORM AutoMigrate criará a coluna automaticamente.

### Passo 2 — Backend: API Types + DTO
Adicionar `Style string` em `ButtonCreateRequest`, `ButtonDTO`, e nos mappers `ToButtonDTO`/`ToCustomCaptionButtonDTO`.

### Passo 3 — Backend: Service + Repository
Atualizar `CreateButton` e `UpdateButton` para persistir o campo `Style`. Atualizar `UpdateButton` no repositório para incluir `style` no map de updates.

### Passo 4 — Backend: Keyboard Builder
Em `keyboard.go`, setar `btn.Style = b.Style` ao construir o `telego.InlineKeyboardButton`.

### Passo 5 — Backend: Executor
Adicionar `Style string` ao `executor.InlineKeyboardButton`. Mapear nos conversores `toTelegoKeyboard` e `toExecutorKeyboard`.

### Passo 6 — Frontend: Types
Adicionar `style?: string` ao tipo `Button` em `types.ts`.

### Passo 7 — Frontend: ButtonGrid UI
- Nos forms de "Novo botão" e "Editando":
  - Adicionar 3 círculos coloridos (default/cinza, primary/azul, success/verde, danger/vermelho) com seleção visual (anel de seleção).
- No grid, pintar a célula com a cor correspondente ao `style` do botão.
- Ao salvar, enviar `style` no payload junto de `nameButton` e `buttonUrl`.

## Riscos
- **Baixo**: A coluna `style` será adicionada via AutoMigrate (nullable text, sem breaking change). Botões existentes sem `style` simplesmente não terão cor customizada (comportamento padrão do Telegram).

## Impactos esperados
- Usuários poderão escolher a cor de cada botão inline diretamente no dashboard.
- Os botões aparecerão coloridos tanto no grid do dashboard quanto no Telegram.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD
- Telegram Bot API 9.4+ (campo `style` ignorado por clientes mais antigos)

## Como testar

### Build Backend
```bash
cd /home/malbs/Opencode/FreddyBot && go build ./...
```

### Build Frontend
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout internal/database/models/models.go internal/api/types/buttons.go internal/api/dto/dto.go internal/api/dto/mapper.go internal/core/services/buttons.go internal/database/repositories/button.go internal/telegram/events/channelPost/keyboard.go internal/telegram/executor/executor.go internal/telegram/executor/botapi.go internal/telegram/events/channelPost/dispatch_telego.go dashboard/src/types.ts dashboard/src/components/ButtonGrid.tsx
```
