# Plano: implementar-biblioteca-rascunhos-posts_2026-08-15_20-46

## Pedido do usuário
Implementar a funcionalidade de **Biblioteca de Rascunhos / Posts Salvos para Reutilização**.
Permitir que o usuário salve postagens com mídias, textos, botões e reações em sua biblioteca pessoal e reutilize-as a qualquer momento tanto pelo bot do Telegram quanto pelo Dashboard Web para enviar direto ou agendar.

## Objetivo
Conectar o serviço `UserPostTemplateService` com o menu do Telegram Bot e com a API REST/Dashboard Web, oferecendo uma experiência completa de gerenciamento de rascunhos e reutilização de posts.

## Contexto atual
- O banco de dados e o repositório Go já possuem os modelos `UserPostTemplate` e o serviço `UserPostTemplateService` com os métodos `SaveTemplate`, `ListTemplates`, `GetTemplateByID` e `DeleteTemplate`.
- Atualmente essa estrutura ainda não possui interface completa nem no menu principal do bot do Telegram nem na Dashboard Web.

## Arquivos analisados
- `internal/database/models/user_post_template.go`
- `internal/core/services/user_post_template.go`
- `internal/database/repositories/user_post_template.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/loader_telego.go`
- `internal/api/controllers/templateController.go` *(se existir / a criar)*
- `internal/api/routes/routes.go`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`

## Arquivos que poderão ser modificados / criados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_drafts/my_drafts.go` *(novo handler para listar rascunhos no bot)*
- `internal/telegram/loader_telego.go`
- `internal/api/controllers/postTemplateController.go` *(novo)*
- `internal/api/routes/routes.go`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`
- `dashboard/src/components/PostDraftsModal.tsx` *(novo componente de rascunhos no Dashboard)*

## Estratégia de implementação

1. **Integração no Telegram Bot (PostBuilder e Menu de Rascunhos)**:
   - Adicionar o botão `💾 Salvar Rascunho` no menu do PostBuilder (`showMenuTelego`).
   - Criar o handler `my_drafts.go` com o callback `my-drafts` para listar os rascunhos salvos do usuário.
   - Ao selecionar um rascunho na lista do Telegram:
     - 🚀 **Enviar para Canal**: Seleciona o canal e envia na hora.
     - 📅 **Agendar Envio**: Abre o fluxo de agendamento com os dados do rascunho.
     - ✏️ **Editar no PostBuilder**: Carrega os dados do rascunho na sessão ativa do PostBuilder.
     - 🗑️ **Excluir Rascunho**: Deleta o rascunho da biblioteca do usuário.

2. **Endpoints REST API para a Dashboard Web**:
   - Criar `PostTemplateController`:
     - `GET /api/v1/templates` $\rightarrow$ Lista todos os rascunhos/templates salvos pelo usuário autenticado.
     - `POST /api/v1/templates` $\rightarrow$ Salva a sessão atual do PostBuilder ou dados enviados como um novo rascunho.
     - `GET /api/v1/templates/:id` $\rightarrow$ Obtém detalhes do rascunho.
     - `DELETE /api/v1/templates/:id` $\rightarrow$ Exclui um rascunho.
     - `POST /api/v1/templates/:id/load` $\rightarrow$ Carrega o rascunho para a sessão ativa do PostBuilder.

3. **Interface no Dashboard Web (React)**:
   - Adicionar o modal/seção "📂 Meus Rascunhos" no PostBuilder da Dashboard Web.
   - Exibir lista de rascunhos salvos com botões de "Carregar no PostBuilder", "Agendar Envio" e "Excluir".

4. **Validação & Testes**:
   - Testes unitários (`go test ./...`).
   - Compilação do binário (`go build ./cmd/FreddyBot`).
   - Build do Dashboard (`npm run build`).

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/implementar-biblioteca-rascunhos-posts_2026-08-15_20-46.md` e obter aprovação do usuário.
2. Criar `my_drafts.go` e adicionar opções no `postBuilder.go` e `loader_telego.go`.
3. Criar `postTemplateController.go` e registrar as rotas na API REST.
4. Adicionar suporte a rascunhos no Dashboard Web em React.
5. Executar testes e compilações (`go test ./...`, `go build`, `npm run build`).

## Riscos
- Mínimo. A estrutura de banco e serviços de backend já está pronta e testada.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Testes Go
```bash
go test ./...
go build ./cmd/FreddyBot
```

### Build Frontend
```bash
cd dashboard && npm run build
```

## Rollback
`git checkout . && git clean -fd`
