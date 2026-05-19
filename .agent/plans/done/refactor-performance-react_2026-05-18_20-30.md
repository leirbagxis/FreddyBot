# Plano: Auditoria e Refatoração de Performance (React Engine)

## Pedido do usuário
A dashboard continua "pesada" mesmo após otimizações de CSS. Investigar a fundo o que causa o travamento.

## Diagnóstico Técnico (O que realmente pesa)
1. **Re-renders Massivos:** O componente `DashboardContent` no `App.tsx` possui mais de 1000 linhas e gerencia TODO o estado da aplicação. Qualquer mudança (como digitar um caractere na busca ou clicar em um toggle) faz o React re-processar toda a árvore da dashboard.
2. **Prop Drilling de Funções Instáveis:** Funções como `handleAddNoticeButton`, `updateNoticeButton`, etc., são passadas para o `AdminDashboard` sem `useCallback`. Isso quebra a memoização do componente, forçando-o a renderizar novamente sempre.
3. **Complexidade de Renderização:** O React está gastando muito tempo calculando o "diff" de componentes complexos que não mudaram (ex: renderizar a aba de Botões enquanto você está na aba Início).
4. **Estado "Gordo":** O objeto `data` contém toda a configuração do canal. Atualizá-lo exige clones profundos frequentes, o que gera pressão no coletor de lixo (GC) do JavaScript, causando pequenos engasgos.

## Estratégia de Implementação

**1. Memoização Estrita de Componentes Filhos:**
- Aplicar `React.memo` em: `PermissionsCard`, `CaptionCard`, `NewPackCaptionCard`, `ReactionsCard`, `DashboardInicioTab`, `TabBar`.
- Garantir que as props passadas sejam estáveis.

**2. Estabilização de Callbacks no App.tsx:**
- Envolver TODAS as funções de manipulação de estado (Admin e User) em `useCallback`.
- Remover funções anônimas passadas diretamente nas props (ex: `onSelectUser={(id) => ...}`).

**3. Otimização de Renderização Condicional:**
- Garantir que abas que não estão visíveis não executem lógica pesada durante o render do pai.

**4. Refatoração de Estado Local:**
- Mover estados que não precisam ser globais (como `transferInput`, `isSidebarOpen`, `adminSearch`) para dentro de seus respectivos componentes ou hooks específicos.

## Passos detalhados

1. **Memoizar Componentes**: Editar `PermissionsCard.tsx`, `CaptionCard.tsx`, etc., para usar `memo`.
2. **Estabilizar `App.tsx`**:
   - Envolver handlers de Admin em `useCallback`.
   - Limpar o render principal para ser mais "declarativo" e menos "lógico".
3. **Remover Re-renders Desnecessários**: Ajustar o fluxo de dados para que o `AdminDashboard` não renderize quando você estiver na parte de usuário e vice-versa.

## Riscos
- Erros de dependência no `useCallback` podem causar bugs onde a função usa um estado antigo ("stale closure"). Farei uma revisão rigorosa das dependências.

## Como testar
- Usar a dashboard e observar se a resposta aos cliques é instantânea.
- No modo Admin, a busca de usuários deve ser fluida sem atraso na digitação.
