# Plano: Polimento Final de UI/UX da Página "Meus Templates"

## Pedido do usuário
Realizar o polimento final de UI/UX na página "Meus Templates" (`UserTemplatesManager.tsx` e `ButtonGrid.tsx`), sem alterar a estrutura conceitual já aprovada:
1. Remover o símbolo/artefato `▦ ` antes do título **"Botões Inline"**.
2. **Manter estritamente o controle compacto de grade `− 4 + × − 1 +`** em sua forma original no cabeçalho do componente `ButtonGrid.tsx`, sem criar seletores em linhas separadas.
3. Melhorar a visibilidade e o feedback de toque dos slots vazios com um ícone `+` nitidamente visível e fundo sutil.
4. Destacar o contraste e legibilidade dos botões já configurados para diferenciação imediata.
5. Tornar o alça de arrasto (`GripVertical`) suavemente mais visível sem competir com o texto do botão.
6. Tornar a altura do container da grade dinâmica conforme o número real de linhas (ex: 1 linha de botões gera um container compacto sem grande espaço em branco abaixo).

## Objetivo
Refinar a clareza tátil, contraste visual e compacidade da grade de botões inline e da página de templates sem realizar redesigns ou quebrar contratos visuais.

## Contexto atual
- `UserTemplatesManager.tsx` exibe um caractere `▦ ` antes do título "Botões Inline".
- `ButtonGrid.tsx` inicializava `rows` com mínimo de 3 linhas fixas, deixando espaço vazio em branco quando havia apenas 1 linha configurada.
- O controle compacto `− 4 + × − 1 +` já existe no `ButtonGrid.tsx` e deve ser o único controle de dimensões mantido.
- Os ícones `+` dos slots vazios estavam com baixa opacidade (`text-muted-foreground/20`).

## Arquivos analisados
- `dashboard/src/components/UserTemplatesManager.tsx`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/index.css`

## Arquivos que poderão ser modificados
- `dashboard/src/components/UserTemplatesManager.tsx`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação

### 1. Remoção de Artefato em `UserTemplatesManager.tsx`
- Remover o símbolo `▦ ` do elemento `<h5>`, deixando apenas o ícone azul de grade e o texto limpo **`Botões Inline`**.
- Remover quaisquer controles duplicados de linhas/colunas do `UserTemplatesManager.tsx`, delegando inteiramente para o controle compacto nativo do `ButtonGrid.tsx`.

### 2. Controle Compacto de Grade `− 4 + × − 1 +` (`ButtonGrid.tsx`)
- Garantir que o controle compacto `− [col] + × − [row] +` no cabeçalho do `ButtonGrid.tsx` esteja perfeitamente alinhado, limpo e com botões táteis confortáveis.

### 3. Altura Dinâmica da Grade (`ButtonGrid.tsx`)
- Ajustar o cálculo inicial de `rows` para `Math.max(1, maxActiveY + 1)` quando `hideReactions` for `true`.
- Se o template tiver 1 linha de botões, o container renderiza 1 linha compacta (56px) sem criar 3 linhas vazias desnecessárias.

### 4. Slots Vazios e Botões Existentes (`ButtonGrid.tsx` e `index.css`)
- Slots vazios: exibir ícone `+` visível em contêiner suave com opacidade ajustada e efeito tátil (`active:scale-95`).
- Botões ocupados: aplicar contraste nítido (`bg-card border-border/90`), texto em negrito (`font-bold text-foreground`) e alça de arrasto (`GripVertical`) suavemente mais nítida (`text-muted-foreground/45 group-hover:text-muted-foreground/80`).

## Passos detalhados
1. Atualizar `UserTemplatesManager.tsx` limpando o título "Botões Inline" (removendo `▦ `).
2. Atualizar `ButtonGrid.tsx` ajustando o cálculo de linhas dinâmicas, os slots de célula vazios com ícone `+` visível e alça de arrasto dos botões configurados.
3. Atualizar estilos `.grid-cell` em `index.css` para legibilidade máxima.
4. Executar `cd dashboard && npm run build` para testar.

## Riscos
- Nenhum risco. Ajuste puramente estético e de usabilidade de frontend.

## Impactos esperados
- Interface extremamente compacta e limpa.
- Slots da grade imediatamente visíveis e fáceis de tocar.
- Controle de dimensões mantido 100% fiel ao formato original compactado `− 4 + × − 1 +`.

## Compatibilidade
- Linux, macOS, Windows, Mobile (iOS/Android)

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
Restaurar `UserTemplatesManager.tsx`, `ButtonGrid.tsx` e `index.css`.
