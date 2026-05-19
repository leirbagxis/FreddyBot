# Plano: atualizar-esquema-cores-dashboard_2026-05-18_16-30.md

## Pedido do usuário
O usuário deseja que as cores da dashboard sejam "tipo um creme com um verde abacatado para algumas coisas".

## Objetivo técnico
Substituir a paleta de cores atual (baseada no design Sentri/Violeta) por uma paleta personalizada com tons de creme (fundo) e verde abacate (acentos).

## Contexto atual
A dashboard utiliza variáveis CSS para temas Light e Dark no arquivo `dashboard/src/index.css`. Atualmente, o tema light é predominantemente branco/violeta escuro e o tema dark é violeta profundo/limão.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/hooks/useTheme.ts`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/hooks/useTheme.ts` (ajuste de cores de cabeçalho do Telegram)

## Estratégia de implementação
1.  Definir uma paleta de cores harmoniosa:
    - **Creme (Light):** `#FCF8F0` (Cream) para o fundo, `#FFFFFF` para cards.
    - **Verde Abacate (Light):** `#568203` (Avocado) para acentos primários.
    - **Verde Abacate Suave (Light):** `#7EA04D` ou variantes translúcidas para hovers e estados secundários.
2.  Atualizar as variáveis `:root` no `index.css` para refletir essas mudanças no tema `light`.
3.  Ajustar o tema `dark` para manter a consistência, possivelmente usando tons de "Olive/Forest" para o fundo e mantendo o verde limão/abacate atual como acento.

## Passos detalhados

1.  **Backup:** (Já verifiquei que existem arquivos `.bak`, mas farei alterações cirúrgicas).
2.  **Modificar `dashboard/src/index.css`:**
    - Atualizar `[data-theme="light"]`:
        - `--bg`: `#FCF8F0`
        - `--card`: `#FFFFFF`
        - `--accent`: `#568203`
        - `--accent-soft`: `rgba(86, 130, 3, 0.1)`
        - `--accent-hover`: `rgba(86, 130, 3, 0.15)`
        - `--link`: `#568203`
        - `--border-active`: `#568203`
    - Atualizar `[data-theme="dark"]` para tons de oliva/abacate escuro se necessário (para não destoar).
3.  **Modificar `dashboard/src/hooks/useTheme.ts`:**
    - Atualizar as cores de `setHeaderColor` e `setBackgroundColor` no `useEffect` para combinar com o novo creme.
4.  **Validar:** Abrir o dashboard (se possível num ambiente de preview) ou verificar via código se todas as referências a `--accent` e `--bg` foram contempladas.

## Riscos
- **Contraste:** O verde abacate em fundo creme pode ter contraste baixo se não for bem escolhido (WCAG). Usarei um tom de verde mais escuro (`#568203`) para o texto/acentos principais no tema light.
- **Estética:** A mudança pode afetar ícones e gradientes que usam cores fixas ou outras variáveis (ex: `--success`).

## Impactos esperados
- Dashboard com aparência mais orgânica e "soft".
- Alinhamento com a preferência visual do usuário.

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Execução
Manual (visual) após a aplicação das mudanças.

## Rollback
Restaurar `dashboard/src/index.css` a partir do `index.css.bak` ou via `git checkout`.

## Observações
O tema Dark atual já usa um "accent-lime" (`#c2ef4e`) que é próximo ao "verde abacatado". Vou mantê-lo ou ajustá-lo levemente para ser mais "abacate" (menos neon).
