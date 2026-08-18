# Plano: postbuilder-uses-caption-templates

## Pedido do usuário
O usuário quer que os templates do PostBuilder e do Dashboard sejam a MESMA entidade real no banco de dados. Os templates criados na dashboard devem poder ser abertos pelo PostBuilder e vice-versa. O usuário odiou a mensagem fictícia do comando "/template". O menu de templates no post builder está uma bagunça e precisa ser limpo para carregar esses templates unificados.

## Objetivo
Unificar os dois domínios usando `UserCaptionTemplate` para ambos. Eliminar o `UserPostTemplate`. Adicionar suporte a reações no `UserCaptionTemplate`.

## Contexto atual
- `UserCaptionTemplate` armazena Nome (code), Legenda e Botões com posições (X, Y).
- `UserPostTemplate` salvava o JSON inteiro do estado do PostBuilder.
- Na UI da dashboard, nós falseamos uma "união" (chamando ambos), o que não era o que o usuário queria.

## Arquivos que serão modificados
- `internal/database/models/user_caption_template.go`: Adicionar campo `Reactions string`.
- `internal/core/services/user_caption_template.go`: Incluir `Reactions` no `Create` e `Update`.
- `internal/api/controllers/userCaptionTemplateController.go`: Parsear o campo `Reactions`.
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`: 
  - `pb-save-template` vai criar um `UserCaptionTemplate` preenchendo Code, Legenda, Reações e mapeando os botões para `UserCaptionTemplateButton`.
  - `pb-list-templates` vai listar os `UserCaptionTemplate`.
  - `pb-load-template` vai buscar o `UserCaptionTemplate`, sobrescrever a Legenda, Botões e Reações no `PostBuilderState`.
- `internal/database/models/user_post_template.go`: Remover.
- `dashboard/src/components/UserTemplatesManager.tsx`: Simplificar para exibir apenas `UserCaptionTemplate`, já que a divisão não existe mais. Remover os textos sobre "/template".

## Estratégia de implementação
1. **Backend - DB & Services**:
   - Adicionar `Reactions` em `UserCaptionTemplate`.
   - Modificar construtores e SQL queries (Gorm trata automaticamente) em `UserCaptionTemplateService`.
2. **Backend - PostBuilder**:
   - Refatorar a listagem, salvamento e carregamento para injetar o `UserCaptionTemplateService` e usar os métodos `List`, `Create`, `CreateButton`, etc.
   - O salvamento agrupará Title+Body+Footer na `Caption` do template, e mapeará a lista flat de botões do post builder para a grid do template (Y = index, X = 0).
3. **Frontend**:
   - No `UserTemplatesManager.tsx`, voltar a usar apenas `listUserCaptionTemplates`.
   - Remover as lógicas de tipo "MixedTemplate" e o if falso do `/template`.

## Passos detalhados
- Atualizar modelo no Go.
- Atualizar Service e Controller.
- Refatorar bot handler.
- Limpar Frontend.
- Verificar compilação Go e NPM run build.

## Como testar
- Criar um template na dashboard.
- Abrir o bot, usar o post builder, listar templates. O template da dashboard deve aparecer.
- Carregar o template.
- Salvar um template no bot, ir na dashboard. O template deve aparecer e permitir edição.

## Compatibilidade
- React, Go, Gorm.

## Rollback
Desfazer commits caso ocorra falhas massivas.
