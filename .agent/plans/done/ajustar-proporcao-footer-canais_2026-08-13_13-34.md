# Plano: Ajustar Proporções da Barra Inferior / Footer dos Canais

## Pedido do usuário
Ajustar o footer/barra de navegação dos canais (`.bottom-nav`). Atualmente está "fino e largo" — encontrar o meio-termo ideal (altura mais confortável e largura mais compacta e proporcional).

## Objetivo
Refatorar os estilos CSS de `.bottom-nav` e `.nav-item` em `dashboard/src/index.css` para:
1. **Largura mais compacta**: Reduzir a largura máxima de `400px` para `330px` (`width: calc(100% - 40px); max-width: 330px;`), evitando que a barra flutuante estique demais horizontalmente.
2. **Altura e padding mais confortáveis**:
   - Aumentar o padding de `.bottom-nav` de `4px 3px` para `6px 6px`.
   - Aumentar o padding interno de cada item (`.nav-item`) de `5px 2px 4px` para `8px 4px 7px`.
3. **Escala de ícones e texto**:
   - Aumentar o tamanho dos ícones SVG de `18px` para `20px`.
   - Aumentar o tamanho do texto das abas de `9px` para `10px` (`font-weight: 650`).
   - Arredondar os cantos do item ativo (`border-radius: 18px`).

## Contexto atual
- `.bottom-nav` possui `max-width: 400px` e padding vertical muito reduzido (`4px`), gerando uma aparência achatada e muito esticada na tela.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/components/TabBar.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`

## Estratégia de implementação

1. **`dashboard/src/index.css`**:
   - Atualizar a regra `.bottom-nav`:
     - `max-width: 330px;` (era 400px)
     - `width: calc(100% - 40px);` (era calc(100% - 32px))
     - `padding: 6px 6px;` (era 4px 3px)
     - `gap: 3px;` (era 1px)
     - `border-radius: 24px;` (era 100px)
   - Atualizar a regra `.nav-item`:
     - `padding: 8px 4px 7px;` (era 5px 2px 4px)
     - `gap: 2px;`
     - `font-size: 10px; font-weight: 650;`
     - `border-radius: 18px;`
   - Atualizar `.nav-item svg`:
     - `width: 20px; height: 20px;` (era 18px)

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/index.css`.
2. Executar `npm run build`.

## Riscos
- **Baixo**: Ajuste de CSS puramente proporcional na barra de navegação inferior.

## Impactos esperados
- Barra de navegação flutuante com meio-termo ideal: mais compacta horizontalmente e com altura/ícones perfeitamente proporcionais e confortáveis ao toque.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/index.css
```
