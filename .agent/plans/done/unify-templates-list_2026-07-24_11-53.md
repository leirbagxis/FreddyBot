# Plano: unify-templates-list

## Pedido do usuário
Unificar a visualização de templates feitos via Dashboard (Legendas) e os feitos via Bot (Post Builder) em uma única aba/lista chamada "Meus Templates", pois o usuário os enxerga como a mesma coisa.

## Objetivo
Criar uma lista única em `UserCaptionTemplateCard.tsx` que mescla `UserCaptionTemplate` e `UserPostTemplate`.

## Contexto atual
Temos dois componentes separados (`UserCaptionTemplateCard` e `UserPostTemplateCard`) que carregam seus dados separadamente e estão separados na UI.

## Arquivos que poderão ser modificados
- `dashboard/src/components/UserCaptionTemplateCard.tsx`
- `dashboard/src/App.tsx` (para remover a importação/renderização do UserPostTemplateCard)

## Estratégia de implementação
1. Em `UserCaptionTemplateCard.tsx`, importar `listUserPostTemplates` e `deleteUserPostTemplate`.
2. Criar um tipo unificado `type MixedTemplate = { type: 'caption', data: UserCaptionTemplate } | { type: 'post', data: UserPostTemplate }`.
3. O `load()` buscará ambas as listas via `Promise.all` e as mesclará em um único estado `items: MixedTemplate[]`.
4. Renderizar a lista unificada.
   - Para os de tipo `caption`, manter o comportamento atual (expande para `TemplateEditor`).
   - Para os de tipo `post`, mostrar um ícone diferente (ex: `LayoutTemplate` em vez de `Hash`) e o título correspondente (`name` em vez de `code`).
   - Quando expandir um tipo `post`, mostrar um aviso temporário: "Este template foi criado pelo Post Builder. A edição pela web estará disponível em breve." (Já que não temos o PostBuilder inteiro no frontend ainda).
5. Remover `UserPostTemplateCard.tsx` do `App.tsx` e deixar apenas o `UserCaptionTemplateCard`.

## Passos detalhados
1. Atualizar `UserCaptionTemplateCard.tsx` para gerenciar as duas fontes de dados.
2. Atualizar a renderização da lista para lidar com o `MixedTemplate`.
3. Ajustar o `App.tsx` para limpar o componente `UserPostTemplateCard`.

## Riscos
- O usuário pode querer editar o `UserPostTemplate` pela web. Como a estrutura de dados é diferente (PostBuilderState JSON vs Caption+Buttons), será necessário deixar claro que a edição web desse tipo será implementada no futuro, ou permitir apenas a exclusão.

## Compatibilidade
- React/Vite
- TailwindCSS

## Rollback
Restaurar `UserCaptionTemplateCard.tsx` e `App.tsx` usando git.
