# Plano: User Post Templates na Dashboard

## Pedido do usuário
Ligar o `UserPostTemplate` (templates do post builder criados no bot pelo /template) na dashboard para que o usuário possa listá-los e excluí-los.

## Objetivo
Criar a infraestrutura de API e UI para listar e gerenciar `UserPostTemplate`s do usuário na interface web.

## Contexto atual
- O bot já permite criar `UserPostTemplate` usando `/template`.
- O modelo e o service (`UserPostTemplateService`) já existem.
- Falta o Controller (`UserPostTemplateController`) para servir os dados na API `/api/me/post-templates`.
- Falta a UI (`UserPostTemplateCard.tsx`) na dashboard (`/me/channels`).

## Arquivos que poderão ser modificados
- `internal/api/routes/routes.go`
- `internal/container/appContainer.go`
- `dashboard/src/App.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`

## Arquivos que serão criados
- `internal/api/controllers/userPostTemplateController.go`
- `dashboard/src/components/UserPostTemplateCard.tsx`

## Estratégia de implementação

### Backend
1. Criar `UserPostTemplateController` com injenção de `UserPostTemplateService`.
2. Adicionar métodos:
   - `ListUserPostTemplates(c *gin.Context)`
   - `DeleteUserPostTemplate(c *gin.Context)`
3. Registrar as rotas no `routes.go`:
   - `GET /api/me/post-templates`
   - `DELETE /api/me/post-templates/:id`
4. Registrar o Controller no `appContainer.go`.

### Frontend
1. Atualizar `types.ts` para incluir a interface `UserPostTemplate`.
2. Adicionar as chamadas em `api.ts`:
   - `listUserPostTemplates()`
   - `deleteUserPostTemplate(id: string)`
3. Criar `UserPostTemplateCard.tsx` que fará:
   - Uma lista simples dos templates (nome do template).
   - Um botão de lixeira para excluir (com confirmação).
4. Em `App.tsx`, adicionar um card novo para "Templates de Post" logo abaixo de "Meus Templates" (que são os de legenda), que abrirá esse novo componente em tela cheia (como feito anteriormente).

## Passos detalhados
1. Criar controller no backend e expor rotas CRUD simples (listar e excluir).
2. Atualizar o frontend (API, tipos) para consumir as novas rotas.
3. Criar o componente React.
4. Integrar o componente na dashboard.
5. Testar backend via curl/dashboard e frontend verificando se lista e apaga.

## Riscos
- O JSON armazenado em `UserPostTemplate` não deve ser alterado diretamente sem validação de estrutura, por isso, manteremos a edição restrita (apenas listar e excluir por enquanto).

## Rollback
- Reverter as alterações em `routes.go` e remover os arquivos recém criados.
