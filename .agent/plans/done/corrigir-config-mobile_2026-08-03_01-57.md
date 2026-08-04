# Plano: corrigir config mobile

## Pedido do usuário

Corrigir a responsividade da aba de Configurações no celular.

## Objetivo

Fazer a aba `Configurações` caber e permanecer utilizável em telas de 390 px, sem rolagem horizontal, controles sobrepostos ou ações inacessíveis. Os campos, payload de configuração, salvamento, refresh do cache e toggles atuais serão preservados.

## Contexto atual

- A tela é renderizada por `AdminConfigTab.tsx` dentro do shell responsivo `AdminLayout.tsx`.
- Os três blocos funcionais são Sistema, Legendas e PostBuilder.
- O CSS atual apenas aplica acabamento genérico a `.admin-config-page`; não há breakpoint específico para o cabeçalho do card, `CardAction`, linhas com `Switch`, editores de texto rico e rodapé de salvamento.
- Em celular, essa ausência deixa a composição de desktop prevalecer, especialmente na seção PostBuilder.

## Arquivos analisados

- `.agent/context.md`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/index.css`
- Componentes atuais de `Card`, `Button`, `Input`, `Switch`, `Textarea` e `RichTextEditor` usados pela aba.

## Arquivos que poderão ser modificados

- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/index.css`
- `.agent/memory/memory.md`
- `.agent/decisions.md` (somente se a correção exigir uma decisão reutilizável)

## Estratégia de implementação

Criar uma composição mobile específica para Configurações, sem alterar a lógica:

- Marcar semanticamente os cards e linhas de configuração para o CSS tratar somente essa tela.
- Em até 640 px, empilhar cabeçalho e ação do PostBuilder, reduzir padding, permitir quebra de textos longos e manter cada `Switch` acessível à direita.
- Fazer inputs, editor rico, textarea JSON e o código de uso inline ocuparem 100% da largura disponível, com quebra segura para conteúdo longo.
- Transformar o rodapé em ação larga e fácil de tocar no celular, mantendo o botão e handler atuais.
- Não mudar endpoints, estados React, payloads ou comportamento de salvamento.

## Passos detalhados

1. Adicionar classes de escopo aos grupos Sistema, Legendas e PostBuilder, às linhas com switch e ao rodapé de salvar.
2. Ajustar os cards para largura mínima zero e espaçamento responsivo, evitando overflow de editor, textarea ou conteúdo `code`.
3. Criar breakpoint de 640 px para cabeçalho/`CardAction`, linhas de controle, editores, bloco de PostBuilder e botão Salvar.
4. Preservar foco visível, tamanho de toque dos toggles e leitura de descrições em telas pequenas.
5. Construir o dashboard e renderizar a aba em 390 px e desktop para verificar ausência de rolagem horizontal, botão acessível e continuidade do layout maior.
6. Atualizar a memória e mover este plano para concluído após a validação.

## Riscos

- `RichTextEditor` possui sua própria toolbar e pode impor largura mínima; o ajuste será escopado ao container de Configurações para não afetar o Broadcast nem editores de usuários.
- O `CardAction` é usado em layout de cabeçalho; a alteração deve preservar a ação de refresh e seu estado de carregamento.
- A tela depende de API autenticada para conteúdo completo. Serão verificados tanto o estado de loading quanto a composição carregada quando disponível, sem salvar nem alterar configurações.

## Impactos esperados

- Configurações utilizável em celular sem corte lateral.
- Campos e editores legíveis, controles com alvos de toque adequados e botão Salvar acessível.
- Nenhuma mudança funcional no backend ou no fluxo de administração.

## Compatibilidade

- Linux
- macOS
- Windows
- Docker
- CI/CD
- Desktop e WebView/mobile do Telegram

## Como testar

### Build

```bash
cd dashboard && npm run build
```

### Testes

```bash
cd dashboard && npx tsc --noEmit
```

Falhas TypeScript preexistentes serão separadas de qualquer regressão introduzida pela correção.

### Execução

```bash
cd dashboard && npm run dev
```

Abrir `/admin/dash?tab=config` em 390 px e desktop; conferir loading, Sistema, Legendas, PostBuilder ativo e inativo, foco de teclado e ausência de rolagem horizontal. Não executar salvar, toggle ou refresh durante a validação.

## Rollback

Reverter somente as classes e regras mobile introduzidas em `AdminConfigTab.tsx` e `index.css`, sem tocar nas alterações administrativas já existentes.

## Observações

- Critério de aceite: em 390 px o usuário consegue ler e editar todos os campos sem zoom, corte ou sobreposição e alcança o botão Salvar.
- A correção será isolada da aba Broadcast, que também usa editor rico mas possui composição diferente.
