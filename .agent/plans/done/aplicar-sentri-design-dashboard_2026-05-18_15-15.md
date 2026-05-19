# Plano: aplicar-sentri-design-dashboard_2026-05-18_15-15.md

## Pedido do usuário
Aplicar as cores e fontes do `@DESIGN.md` na dashboard, mantendo um backup para reversão.

## Objetivo
Transformar a identidade visual da dashboard para o estilo "Sentri": paleta violeta profunda/meia-noite, acentos em verde lima elétrico e tipografia Rubik/Monaco.

## Contexto atual
- A dashboard usa variáveis CSS no `:root` e classes do Tailwind (importadas no topo do CSS).
- As animações já foram simplificadas para CSS puro.

## Arquivos analisados
- `@DESIGN.md` (especificação de design)
- `dashboard/src/index.css` (estilos atuais)
- `dashboard/src/App.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/App.tsx` (para ajustes de layout se necessário)

## Estratégia de implementação
1.  **Backup:** Criar `index.css.bak` (já solicitado via comando).
2.  **Fontes:** Adicionar o import do Google Fonts para a família `Rubik`.
3.  **Variáveis CSS:** Mapear os tokens do `@DESIGN.md` para as variáveis existentes no `index.css`:
    - `--bg`: `#150f23` (Midnight)
    - `--card`: `#1f1633` (Ink Deep)
    - `--accent`: `#c2ef4e` (Electric Lime)
    - `--text`: `#ffffff` (On Primary)
    - Tipografia: Configurar `Rubik` como fonte principal.
4.  **Ajustes de Componentes:**
    - Atualizar botões para seguirem o padrão `button-primary` (black-violet no light, white no dark).
    - Ajustar `border-radius` conforme a escala `rounded` (xs, sm, md, lg, xl, xxl).
    - Adaptar o layout de "Welcome Card" para usar o gradiente violeta/lime sugerido.

## Passos detalhados

1.  **Atualizar index.css:**
    - Inserir `@import url('https://fonts.googleapis.com/css2?family=Rubik:wght@400;500;600;700&display=swap');`.
    - Sobrescrever as variáveis de cor nos blocos `[data-theme="light"]` e `[data-theme="dark"]` para convergirem para a estética Sentri.
    - Atualizar `font-family` global para `Rubik`.
    - Ajustar estilos de botões, cards e inputs para refletirem os tokens específicos (ex: uppercase nos botões com tracking).

2.  **Refinar Elementos de Marca:**
    - Aplicar o gradiente "Sentri" nos ícones e áreas de destaque.
    - Adicionar o efeito de "lime keyword" em títulos onde fizer sentido.

## Riscos
- O contraste do verde lima elétrico sobre fundo escuro é excelente, mas sobre fundo claro precisa de cuidado (o `@DESIGN.md` recomenda usar o lime apenas como destaque tipográfico).
- A mudança radical de cores pode exigir ajustes em ícones do `lucide-react` que tenham cores fixas (embora a maioria use `currentColor`).

## Impactos esperados
- Dashboard com visual premium, "developer-focused" e moderno.
- Identidade visual coerente com a marca Sentri.

## Como testar
- Verificar legibilidade em ambos os temas (embora o Sentri foque muito no Dark).
- Validar se os botões e estados de hover/active continuam claros.

## Rollback
- `cp dashboard/src/index.css.bak dashboard/src/index.css`
- `cp dashboard/src/App.tsx.bak dashboard/src/App.tsx`
