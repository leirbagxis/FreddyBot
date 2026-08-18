# Plano: Refinar Dashboard Mobile Telegram

## Pedido do usuário
Redesenhar a tela de configuração de canais do FreddyBot como uma Telegram Mini App mobile, preservando aproximadamente 80% da identidade, da estrutura de informação, do conteúdo e das funcionalidades existentes, com melhorias de hierarquia, espaçamento, consistência, estados de interação e acabamento visual.

## Objetivo
Consolidar a aba de configurações em uma composição `Página → Card de seção → Conteúdo`, com tema navy inspirado no Telegram, azul reservado a ações importantes, cards coerentes para captions e reações, opções secundárias compactas e navegação inferior segura em telas móveis.

## Contexto atual
- O dashboard usa React 19, TypeScript, Vite 7, Tailwind CSS 4 e componentes shadcn/ui.
- A tela já está dividida nas abas Início, Legendas, Botões e Permissões, além de abas condicionais, e sua lógica deve ser preservada.
- Os componentes atuais já contêm parte do redesign solicitado: ícones por categoria, cinco slots de reação, botão azul de salvamento, divisor “OUTRAS OPÇÕES”, preview de variáveis e cápsula ativa no rodapé.
- Existe um plano concluído anterior com escopo muito semelhante; esta execução será um refinamento baseado no estado real renderizado, evitando reimplementar o que já está correto.
- Não foi encontrado um `brand.md`; a fonte de verdade visual disponível é `dashboard/src/index.css`, os componentes compartilhados e as superfícies vizinhas do dashboard.
- O worktree contém alterações locais extensas, inclusive em todos os arquivos centrais deste redesign. Essas alterações serão preservadas e qualquer edição será feita sobre o estado atual, sem reset ou sobrescrita ampla.

## Arquivos analisados
- `.agent/context.md`
- `.agent/memory/memory.md`
- `.agent/skills/fable-method/SKILL.md`
- `.agent/skills/fable-method/references/domains/design-ux.md`
- `.agent/plans/done/redesign-dashboard-mobile-config-telegram_2026-08-14_11-45.md`
- `dashboard/package.json`
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`
- `dashboard/src/components/CaptionCard.tsx`
- `dashboard/src/components/CaptionPreview.tsx`
- `dashboard/src/components/NewPackCaptionCard.tsx`
- `dashboard/src/components/ReactionsCard.tsx`
- `dashboard/src/components/NativeReactionsCard.tsx`
- `dashboard/src/components/TabBar.tsx`
- `dashboard/src/components/WaveDivider.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`
- `dashboard/src/components/CaptionCard.tsx`
- `dashboard/src/components/CaptionPreview.tsx`
- `dashboard/src/components/NewPackCaptionCard.tsx`
- `dashboard/src/components/ReactionsCard.tsx`
- `dashboard/src/components/NativeReactionsCard.tsx`
- `dashboard/src/components/TabBar.tsx`
- `.agent/memory/memory.md`, somente se surgir uma convenção reutilizável
- `.agent/decisions.md`, somente se houver decisão arquitetural relevante

## Estratégia de implementação
Primeiro será criada uma referência visual do estado atual em viewport mobile. Em seguida, as lacunas entre a interface renderizada e o pedido serão corrigidas de forma cirúrgica, priorizando tokens semânticos no tema dark, componentes já existentes e nenhuma mudança em contratos de API ou arquitetura de informação. Os estados de edição, vazio, selecionado, carregamento, desabilitado, foco, toque e overflow serão tratados junto com cada componente. A navegação fixa será validada com safe area e conteúdo longo. O resultado será observado em mais de uma largura mobile e submetido a build e revisão `fable-judge` antes de ser considerado concluído.

## Passos detalhados

1. Executar o dashboard localmente com dados disponíveis ou mocks existentes e registrar o estado atual da aba Legendas em pelo menos duas larguras móveis.
2. Ajustar os tokens dark/Telegram e a composição geral para obter fundo navy profundo, cards navy ligeiramente mais claros, contraste legível e azul Telegram concentrado em ações e estados ativos.
3. Uniformizar Caption Padrão e New Pack Caption: cabeçalho, ícone, título, descrição, ação de edição, preview, destaque de links/variáveis e estados de edição/cancelamento/salvamento.
4. Refinar Reações / Votos para cinco alvos de toque claros, estados vazio/selecionado/foco/pressionado, remoção acessível, comportamento de entrada preservado e botão primário de largura total com loading.
5. Tornar “Reações Nativas do Telegram” uma opção secundária compacta, mantendo a expansão funcional de modo e emojis quando ativada, sem criar aninhamentos visuais desnecessários.
6. Consolidar o separador “OUTRAS OPÇÕES”, espaçamento vertical, padding interno e comportamento de conteúdo longo, removendo apenas decoração redundante dentro desta superfície.
7. Refinar a navegação inferior mantendo rótulos e ícones, cápsula azul ativa, foco/pressed state, safe area e padding suficiente para nunca cobrir o último conteúdo.
8. Verificar overflow de captions e Unicode, estados vazio/loading/erro aplicáveis, navegação por teclado, rótulos acessíveis e contraste das novas combinações de cor.
9. Rodar build, observar novamente o resultado nas larguras planejadas, comparar com o pedido original e executar o passe `fable-judge`.
10. Registrar somente decisões/memória realmente reutilizáveis e mover este plano para `done/` após a verificação final.

## Riscos
- As alterações locais não commitadas podem se sobrepor ao escopo; edições amplas poderiam apagar trabalho existente.
- O estado atual já implementa boa parte do pedido, portanto um novo redesign indiscriminado poderia reduzir consistência ou alterar mais que os 20% desejados.
- Captions com Markdown, HTML do Telegram, emojis customizados e Unicode decorativo podem quebrar ou transbordar se o preview for simplificado.
- A entrada de emoji por campo invisível pode ter comportamento diferente entre teclado desktop e seletor nativo mobile; isso precisa ser observado, não apenas inferido.
- A barra fixa pode cobrir conteúdo em iOS/Telegram se `safe-area-inset-bottom` e o padding final divergirem.
- Estados assíncronos de reações nativas devem manter feedback e evitar mudanças otimistas inconsistentes em caso de erro.

## Impactos esperados
- Melhor leitura e hierarquia sem alterar conteúdo, abas ou contratos de dados.
- Alvos de toque mais claros e feedback consistente para edição, switches e seleção de emojis.
- Maior coerência com uma Telegram Mini App em telas móveis.
- Menos ruído visual e menor dependência de cards aninhados ou decoração excessiva.
- Nenhum impacto esperado no backend ou nas APIs.

## Compatibilidade
- Linux: desenvolvimento e build via Vite; validação visual em navegador Chromium disponível.
- macOS: layout responsivo e safe area compatíveis com navegadores modernos.
- Windows: sem APIs específicas de sistema operacional.
- Docker: nenhum contrato de container será alterado; o bundle continuará sendo produzido pelo fluxo existente.
- CI/CD: nenhuma dependência nova; `npm run build` permanece como verificação de compilação.

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Testes
```bash
cd dashboard && npm run build
```

O projeto não declara suíte de testes ou lint no `package.json`; além do build, serão feitas verificações manuais dos estados interativos e inspeção visual em pelo menos duas larguras mobile.

### Execução
```bash
cd dashboard && npm run dev -- --host 127.0.0.1
```

## Rollback
Como o worktree já contém alterações do usuário nos mesmos arquivos, não será usado `git checkout`, `git reset` ou sobrescrita integral. O rollback seguro será feito revertendo apenas os hunks introduzidos nesta execução, identificados pelo diff produzido após a aprovação, e restaurando o plano à pasta apropriada se necessário.

## Observações
- O redesign será restrito à dashboard de configuração de canal; superfícies administrativas e backend ficam fora do escopo.
- A implementação preservará funcionalidades e dados existentes, inclusive opções avançadas do New Pack Caption e das reações nativas.
- A aprovação deste plano autoriza apenas alterações locais e verificações; não autoriza commit, push, deploy ou publicação.
