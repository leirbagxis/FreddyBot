# Plano: Reajustar Footer Mobile

## Pedido do usuário
Corrigir o footer de navegação para que permaneça sempre fixo à frente do conteúdo, deixe de parecer transparente e fique um pouco maior e visualmente mais espesso.

## Objetivo
Reforçar a bottom navigation existente sem alterar sua estrutura, itens ou identidade: posição fixa confiável, fundo sólido, camada frontal explícita e dimensões ligeiramente maiores, mantendo a compensação do conteúdo sincronizada.

## Contexto atual
- `.bottom-nav` já declara `position: fixed` e `z-index: 1000`.
- Seu fundo usa `var(--card)`, que é transparente nos temas light e Telegram e semitransparente no dark global; isso explica a aparência transparente em alguns estados.
- A geometria compartilhada usa `--bottom-nav-height: 52px`, e o padding do conteúdo depende desse token.
- A barra tem apenas `padding: 4px`, borda muito sutil e sombra leve, o que reduz sua espessura visual.
- Não existem outros overrides de `.bottom-nav` alterando sua posição.
- Há alterações locais preexistentes em `dashboard/src/index.css`, que serão preservadas.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/components/TabBar.tsx`
- `.agent/memory/memory.md`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `.agent/memory/memory.md`, somente para atualizar a convenção de layout

## Estratégia de implementação
Criar um token sólido específico para o fundo da navegação, aumentar moderadamente a altura compartilhada e o padding, reforçar borda/sombra sem glassmorphism e garantir a camada fixa com isolamento e z-index explícitos. Como o padding inferior da página usa o mesmo token de altura, ele crescerá automaticamente junto com o footer.

## Passos detalhados

1. Aumentar `--bottom-nav-height` de forma moderada.
2. Definir `--bottom-nav-bg` sólido para o dashboard de canal e uma variação navy no tema dark.
3. Aplicar o novo fundo sólido em `.bottom-nav`.
4. Reforçar `position: fixed`, isolamento, z-index, padding, borda e sombra de forma contida.
5. Confirmar que o cálculo de `padding-bottom` do layout continua usando a nova altura.
6. Rodar `git diff --check` e executar `fable-judge`; tentar o build caso Node/npm esteja disponível.

## Riscos
- Aumentar a barra sem manter o token compartilhado causaria nova sobreposição; altura e reserva continuarão ligadas à mesma variável.
- Fundo excessivamente contrastante poderia descaracterizar o tema; será usado navy sólido próximo aos cards atuais.
- Sombra excessiva criaria efeito flutuante pesado; o ajuste será discreto.

## Impactos esperados
- Footer sempre visualmente à frente do conteúdo.
- Fundo opaco e legível.
- Barra ligeiramente maior e mais robusta.
- Conteúdo ainda totalmente rolável acima da navegação.
- Nenhuma mudança nos itens, ações ou arquitetura.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker: sem alterações
- CI/CD: sem dependências ou scripts novos

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Testes
```bash
cd dashboard && npm run build
```

### Execução
```bash
cd dashboard && npm run dev -- --host 127.0.0.1
```

## Rollback
Reverter apenas as declarações adicionadas ou ajustadas no bloco `.channel-dashboard` e `.bottom-nav` desta execução, sem usar reset ou checkout nos arquivos com alterações locais.

## Observações
- Nenhuma alteração em `TabBar.tsx` é necessária; a estrutura atual já atende ao pedido.
- A aprovação não autoriza commit, push ou deploy.
