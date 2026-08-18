# Plano: Refinar UI Mobile Final

## Pedido do usuário
Aplicar apenas ajustes finais à interface mobile atual do FreddyBot, preservando identidade, cores, tipografia, estrutura, cards, ícones e aparência Telegram. A prioridade é eliminar definitivamente a sobreposição da navegação inferior e, depois, refinar slots de reação, espaçamentos, proporções e hierarquia de raios.

## Objetivo
Corrigir o espaço reservado para a navegação fixa no nível do layout e tornar os componentes de configuração ligeiramente mais compactos e claros, sem redesenhar a tela, mudar a arquitetura de informação ou alterar funcionalidades.

## Contexto atual
- A interface já usa tema navy/Telegram escopado em `.channel-dashboard`, cards `rounded-2xl`, previews `rounded-xl`, ícones coloridos e navegação fixa em cápsula.
- A navegação usa `bottom: calc(12px + env(safe-area-inset-bottom))`, mas o espaço compensatório está dividido entre `.main-content` e apenas `.channel-config-screen`. Essa duplicação não estabelece uma única relação explícita com a altura real da barra e não protege uniformemente todas as abas.
- A correção recomendada é definir medidas compartilhadas da navegação e aplicar o padding ao container rolável da página (`.channel-dashboard .main-content`), removendo a compensação específica da aba Legendas.
- Os slots de reação já têm cinco posições funcionais, porém usam borda de 2px e estado selecionado com escala/sombra, deixando-os visualmente mais pesados que o refinamento solicitado.
- Caption Padrão e New Pack Caption usam `p-4`, `space-y-3.5` e previews com `p-3`; há espaço para compactação leve sem alterar conteúdo ou formatação.
- O worktree contém alterações locais já existentes nesses arquivos. Elas serão preservadas e a implementação será limitada a pequenos hunks sobre o estado atual.
- O ambiente continua sem `node`/`npm`; build e renderização só poderão ser executados se o runtime estiver disponível após a aprovação.

## Arquivos analisados
- `.agent/context.md`
- `.agent/memory/memory.md`
- `.agent/skills/fable-method/SKILL.md`
- `.agent/skills/fable-method/references/domains/design-ux.md`
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`
- `dashboard/src/components/CaptionCard.tsx`
- `dashboard/src/components/NewPackCaptionCard.tsx`
- `dashboard/src/components/ReactionsCard.tsx`
- `dashboard/src/components/TabBar.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`
- `dashboard/src/components/CaptionCard.tsx`
- `dashboard/src/components/NewPackCaptionCard.tsx`
- `dashboard/src/components/ReactionsCard.tsx`
- `.agent/memory/memory.md`, somente para registrar a convenção de layout reutilizável

## Estratégia de implementação
Centralizar as dimensões da navegação em variáveis CSS semânticas e reservar no container principal a soma de altura da barra, afastamento inferior, safe area e uma pequena folga de rolagem. Depois, reduzir discretamente padding/gaps dos dois cards de caption e seus previews; suavizar a cápsula ativa e os slots de reação usando bordas finas, background sutil e estados sem sombras grandes. A estrutura JSX, textos, ícones, handlers, formatação das captions e componentes funcionais permanecerão intactos.

## Passos detalhados

1. Definir tokens de layout para altura, afastamento e folga da bottom navigation dentro de `.channel-dashboard`.
2. Aplicar `padding-bottom` e `scroll-padding-bottom` à `.channel-dashboard .main-content`, calculados com os tokens e `env(safe-area-inset-bottom, 0px)`.
3. Remover a compensação específica de `.channel-config-screen`, evitando margens arbitrárias por card ou por seção.
4. Ajustar a bottom navigation para usar os mesmos tokens e reduzir levemente padding, escala do ícone e sombra da cápsula ativa.
5. Refinar os cinco slots de reação com borda de 1px, fundo sutil, alvo de toque confortável e estados empty/pressed/selected/focus claramente distintos.
6. Reduzir discretamente o espaçamento vertical de Caption Padrão e New Pack Caption, mantendo cards `rounded-2xl`, previews `rounded-xl`, controles menores e ações “Editar” intactas.
7. Reduzir apenas o padding vertical dos previews, preservando integralmente Markdown, Unicode, emojis customizados, links e destaque de variáveis.
8. Conferir que “OUTRAS OPÇÕES” permanece como divisor simples após Reações, sem ondas ou containers adicionais.
9. Executar `git diff --check`, procurar regressões de estrutura/estados e, se houver runtime, rodar build e observar a tela em mais de uma largura mobile.
10. Executar `fable-judge`, registrar eventual ressalva de ambiente e concluir o plano.

## Riscos
- Um padding calculado somente para a aba Legendas pode reincidir o problema nas demais abas; por isso a correção será no layout compartilhado.
- Alterar a altura visual da barra sem atualizar sua reserva de espaço pode recriar a sobreposição.
- Compactação excessiva pode reduzir alvos de toque; controles permanecerão com dimensões móveis confortáveis.
- Captions podem conter conteúdo longo e formatação especial; o refinamento não alterará o parser nem o conteúdo do preview.
- As mudanças locais existentes nos mesmos arquivos exigem patches pequenos para evitar sobrescrita.

## Impactos esperados
- Última seção completamente rolável acima da navegação, inclusive com safe area.
- Cápsula ativa mais equilibrada sem perder destaque.
- Slots de reação mais claramente acionáveis e menos pesados.
- Cards de caption ligeiramente mais compactos, mantendo conforto e personalidade.
- Nenhuma mudança de API, dados, conteúdo, navegação ou funcionalidade.

## Compatibilidade
- Linux: CSS moderno e build Vite quando Node/npm estiverem disponíveis.
- macOS: suporte a `env(safe-area-inset-bottom)` para dispositivos com área segura.
- Windows: fallback de safe area para `0px`.
- Docker: nenhum arquivo ou contrato de container será alterado.
- CI/CD: nenhuma dependência nova; mantém o comando de build existente.

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Testes
```bash
cd dashboard && npm run build
```

O projeto não define suíte de testes ou lint no `package.json`. A validação visual deve incluir rolagem até o fim, safe area, slots vazios/selecionados/pressionados, caption curta/longa e bottom nav em larguras mobile.

### Execução
```bash
cd dashboard && npm run dev -- --host 127.0.0.1
```

## Rollback
Reverter somente os hunks desta execução nos arquivos listados, comparando com o diff capturado após a aprovação. Não usar `git checkout`, `git reset` ou sobrescrita integral, pois há alterações locais preexistentes nos mesmos arquivos.

## Observações
- Este plano não redesenha a interface e não altera NativeReactionsCard, CaptionPreview, APIs ou backend.
- A aprovação autoriza somente mudanças locais e verificações; não autoriza commit, push, deploy ou publicação.
