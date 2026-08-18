# Plano: Importar Fonte Montserrat e Aumentar Escala Tipográfica do SideMenu

## Pedido do usuário
1. **Importar Fonte Montserrat**: Adicionar a fonte `https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&display=swap`.
2. **Aumentar Tamanho de Fonte no SideMenu**: Aumentar a escala dos textos e proporções no menu lateral para melhor legibilidade e presença visual.

## Objetivo
Configurar a fonte Google Fonts Montserrat no projeto e recalibrar a tipografia do `SideMenu.tsx` (nome do usuário, ID, título "MINHA CONTA", contador de canais, itens de menu e ícones) com tamanhos maiores e mais legíveis.

## Contexto atual
- `dashboard/index.html` / `dashboard/src/index.css` utilizam 'Geist Variable' e 'Plus Jakarta Sans'.
- `SideMenu.tsx` utiliza atualmente tamanhos pequenos (`text-xs`, `text-[11px]`, `text-4xl`).

## Arquivos analisados
- `dashboard/index.html`
- `dashboard/src/index.css`
- `dashboard/src/components/SideMenu.tsx`

## Arquivos que poderão ser modificados
- `dashboard/index.html`
- `dashboard/src/index.css`
- `dashboard/src/components/SideMenu.tsx`

## Estratégia de implementação

1. **Importar Montserrat (`index.html` / `index.css`)**:
   - Adicionar `<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&display=swap">` no `index.html`.
   - Adicionar `@import url('https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&display=swap');` no topo de `index.css`.
   - Adicionar `'Montserrat'` no topo da família de fontes.

2. **Aumentar Escala Tipográfica em `SideMenu.tsx`**:
   - Aplicar `font-family: 'Montserrat', sans-serif;` no container do drawer.
   - Nome do usuário: de `text-sm` (14px) para `text-base` (16px) `font-bold`.
   - ID do usuário: de `text-[11px]` para `text-xs` (12px).
   - Título "MINHA CONTA": de `text-[10px]` para `text-xs` (12px) `font-bold uppercase tracking-wider`.
   - Número de canais: de `text-4xl` (36px) para `text-5xl` (48px) `font-black`.
   - Subtítulo "canais cadastrados": de `text-[11px]` para `text-xs` (12px).
   - Itens de menu (Meus Canais, Meus Templates, Meus Agendamentos, Suporte): de `text-xs` (12px) para `text-sm` (14px) `font-semibold`.
   - Ícones esquerdos: tamanho `18` (era 15).
   - Ícones de seta direita `>`: tamanho `16` (era 14).
   - Aumentar padding dos botões (`py-3.5 px-4`).

3. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/index.html` e `dashboard/src/index.css`.
2. Editar `dashboard/src/components/SideMenu.tsx`.
3. Executar `npm run build`.

## Riscos
- **Baixo**: Mudança puramente visual e tipográfica na UI.

## Impactos esperados
- Menu lateral com a nova tipografia Montserrat, textos maiores e perfeitamente legíveis.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/index.html
git checkout dashboard/src/index.css
git checkout dashboard/src/components/SideMenu.tsx
```
