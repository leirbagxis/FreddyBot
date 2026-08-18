# Plano: User Caption Templates (nível de usuário)

## Pedido do usuário
Mover custom captions/templates para `/me/channels` (nível de usuário, não por canal). Cada template tem: nome curto (usado em postagens para identificar), legenda (editor igual ao padrão) e botões com drag & drop (igual ao ButtonGrid padrão).

## Objetivo
Criar sistema de templates de legenda no nível do usuário, com:
- Modelo próprio (`UserCaptionTemplate`) com short name, caption, buttons
- API `/api/me/templates` (CRUD)
- Pipeline do bot busca por userID + short name (em vez de channel.CustomCaptions)
- UI na página `/me/channels` com editor completo: short name, caption, drag-drop buttons
- Compatibilidade retroativa com `CustomCaption` por canal (mantido para não quebrar)

## Contexto atual
- `CustomCaption` (per-channel): code, caption, linkPreview, buttons — pipeline busca em `channel.CustomCaptions`
- `UserPostTemplate` (per-user): genérico, guarda PostBuilderState como JSON
- Bot pipeline: `stage_transform_telego.go` + `utils_v2.go` → `findCustomCaption(channel, hashtag)`
- Frontend: `/me/channels` mostra lista de canais + cards (saudação, conta Telegram, Premium)
- Aba "Legendas" no dashboard do canal já tem CustomCaptionsCard e CaptionTemplateCard

## Arquivos analisados
- `internal/database/models/models.go` — CustomCaption, CustomCaptionButton
- `internal/telegram/events/channelPost/stage_transform_telego.go` — pipeline
- `internal/telegram/events/channelPost/stage_decorate_telego.go` — pipeline
- `internal/telegram/events/channelPost/utils_v2.go` — findCustomCaption
- `internal/api/routes/routes.go` — rotas
- `dashboard/src/App.tsx` — renderização da página /me/channels + legendas tab
- `dashboard/src/components/ButtonGrid.tsx` — drag & drop
- `dashboard/src/components/CaptionCard.tsx` — editor de caption padrão

## Estratégia de implementação

### Backend

1. **Novo modelo** `UserCaptionTemplate` (tabela `user_caption_templates`):
   ```go
   type UserCaptionTemplate struct {
       ID        string    `gorm:"type:text;primaryKey"`
       UserID    int64     `gorm:"index:idx_user_code,unique"`
       Code      string    `gorm:"index:idx_user_code,unique"`  // short name
       Caption   string    `gorm:"type:text"`
       Buttons   []UserCaptionTemplateButton `gorm:"foreignKey:OwnerTemplateID;constraint:OnDelete:CASCADE;"`
       CreatedAt time.Time
       UpdatedAt time.Time
   }
   ```
   - `Code` único por usuário (unique constraint on user_id + code)
   - Botões com positionX, positionY (igual CustomCaptionButton)
   - Sem link_preview (simplificado)

2. **Novo service** `UserCaptionTemplateService`:
   - CRUD básico (create, list, get, update, delete)
   - `GetByUserAndCode(userID, code)` — para pipeline

3. **Novo controller** `UserCaptionTemplateController`:
   - `GET /api/me/templates` — lista do usuário
   - `POST /api/me/templates` — criar
   - `GET /api/me/templates/:id` — obter
   - `PUT /api/me/templates/:id` — atualizar (caption + short name)
   - `PUT /api/me/templates/:id/layout` — atualizar layout dos botões
   - `DELETE /api/me/templates/:id` — excluir

4. **Pipeline** (stage_transform_telego.go + utils_v2.go):
   - `findUserCaptionTemplate(ownerID, hashtag)` — busca no novo modelo
   - Prioridade: UserCaptionTemplate > CustomCaption (se existir)
   - Compatibilidade: se achar no user template, usa; senão, fallback pro channel custom caption

5. **Container**: injetar novo service + repo

### Frontend

6. **Componente** `UserCaptionTemplateCard.tsx`:
   - Card na página `/me/channels` (entre o PremiumTab e a lista de canais)
   - Estado: listando templates
   - Ao clicar em um template: expande/mostra editor
   - "Novo template" → inline edit
   - Editor:
     - Input "Nome Curto" (short name, usado em posts: `#shortname`)
     - Textarea/editor de legenda (como CaptionCard)
     - ButtonGrid com drag & drop (reutilizar ButtonGrid ou criar adaptador)
   - Botões "Salvar", "Excluir" (com confirmação)

7. **App.tsx**: adicionar na seção `/me/channels` após PremiumTab/antes da lista de canais

### Pipeline

8. **utils_v2.go**: adicionar `findUserCaptionTemplate(userID, hashtag)` que busca em `UserCaptionTemplate` por userID + code
9. **stage_transform_telego.go**: após extrair hashtag, primeiro tenta user template; se não achar, fallback para channel custom caption

## Arquivos que serão criados
- `internal/database/models/user_caption_template.go` — modelo
- `internal/database/repositories/user_caption_template.go` — repo
- `internal/core/services/user_caption_template.go` — service
- `internal/api/controllers/userCaptionTemplateController.go` — controller
- `dashboard/src/components/UserCaptionTemplateCard.tsx` — componente

## Arquivos que serão modificados
- `internal/database/database.go` — adicionar modelo na migração
- `internal/container/appContainer.go` — injetar service + repo
- `internal/api/routes/routes.go` — rotas
- `internal/telegram/events/channelPost/utils_v2.go` — nova função de lookup
- `internal/telegram/events/channelPost/stage_transform_telego.go` — pipeline
- `dashboard/src/App.tsx` — adicionar na página /me/channels
- `dashboard/src/types.ts` — tipos
- `dashboard/src/api.ts` — funções de API

## Riscos
- Quebrar pipeline existente se fallback não funcionar — manter fallback para `channel.CustomCaptions`
- Unique constraint `(user_id, code)` pode conflitar se usuário tentar criar template duplicado — tratamento de erro adequado
- Novo modelo = nova migration — precisa rodar auto-migrate
- UI na página /me/channels precisa ser responsiva com os cards existentes

## Impactos esperados
- Templates funcionam em QUALQUER canal do usuário (não precisa configurar por canal)
- Editor completo com caption + drag-drop buttons
- Pipeline ainda suporta custom captions antigas (fallback)
- Usuário gerencia tudo de um lugar só

## Como testar

### Build
```bash
go build ./...
cd dashboard && npm run build
```

### Pipeline
1. Criar template com short name "promo" e legenda "Oferta!"
2. Postar `#promo` em qualquer canal → bot substitui pela legenda e botões do template
3. Postar hashtag que não existe → fallback para CustomCaption do canal (se existir)

## Rollback
- Remover rotas e controller
- Remover componente do frontend
- Remover modelo (ou manter, sem uso)
- Reverter pipeline para usar só `channel.CustomCaptions`
