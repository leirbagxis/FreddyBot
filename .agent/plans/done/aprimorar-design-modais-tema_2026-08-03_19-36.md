# Plano: Aprimorar Design e Estética dos Modais conforme o Tema / Design System

## Pedido do usuário
"Mas aí tá feio também né, coloque de acordo com o designer/theme"

## Objetivo
Refinar o design e acabamento de todos os modais da aplicação (`ConfirmModal`, `DialogContent`, `DialogOverlay`) para harmonizar perfeitamente com o Design System e com cada tema ativo (Light/Editorial, Dark/Deep Space Neon, Telegram e Minimal):
- Adicionar `backdrop-blur-md` suave ao overlay (`DialogOverlay`) para desfocar o fundo da página em profundidade.
- Aplicar superfície elevada com bordas de vidro sutil (`border: 1px solid var(--border)`), cantos arredondados premium (`18px-20px`), sombras profundas com projeção suave e suporte a `backdrop-filter: blur(20px)` na caixa do modal.
- Refinar a hierarquia visual do `ConfirmModal` (ícone com badge contornado, tipografia ajustada `font-semibold text-[19px]`, botões com `rounded-xl h-10 font-medium`).

## Contexto atual
- A alteração anterior garantiu opacidade total, porém o fundo opaco chapado e a ausência de blur no overlay / refinamento nos botões deixaram o modal com um visual simples/básico sem o acabamento refinado do design system.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/components/ui/dialog.tsx`
- `dashboard/src/components/ConfirmModal.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/ui/dialog.tsx`
- `dashboard/src/components/ConfirmModal.tsx`

## Estratégia de implementação
1. **Em `dashboard/src/components/ui/dialog.tsx`**:
   - Adicionar `backdrop-blur-md` e `bg-black/60` ao `DialogOverlay` para criar desfoque de fundo elegante.
2. **Em `dashboard/src/index.css`**:
   - Ajustar o seletor `[data-slot='dialog-content']` para incluir desfoque de fundo em profundidade (`backdrop-filter: blur(20px)`), raio de borda de `20px`, sombra refinada (`box-shadow: var(--shadow-lg), 0 24px 48px -12px rgba(0, 0, 0, 0.4)`) e ajuste fino da borda conforme o tema.
   - Refinar `--popover` em cada tema para casar harmonicamente com o canvas:
     - `[data-theme="light"]`: `rgba(255, 255, 255, 0.98)` com borda suave `rgba(99, 102, 241, 0.12)`.
     - `[data-theme="dark"]`: `rgba(16, 20, 34, 0.95)` (navy escuro elevado) com borda sutil `rgba(255, 255, 255, 0.1)`.
     - `[data-theme="telegram"]`: `rgba(255, 255, 255, 0.98)`.
3. **Em `dashboard/src/components/ConfirmModal.tsx`**:
   - Refinar os estilos do badge do ícone (borda sutil e fundo `soft`), tipografia do título e descrição, e botões do footer (`rounded-xl`, `h-10`, `font-medium`).

## Passos detalhados
1. Editar `dashboard/src/components/ui/dialog.tsx`:
   - Incluir `backdrop-blur-md` na classe do `DialogOverlay`.
2. Editar `dashboard/src/index.css`:
   - Atualizar a regra `[data-slot='dialog-content']` com as propriedades de acabamento do design system.
3. Editar `dashboard/src/components/ConfirmModal.tsx`:
   - Atualizar estrutura visual interna do modal de confirmação.
4. Rodar o build (`npm --prefix dashboard run build`) para verificar a integridade do código.

## Riscos
- Risco muito baixo. Ajustes exclusivamente estéticos em CSS e classes de componentes de diálogo.

## Impactos esperados
- Modais visualmente deslumbrantes, integrados com a linguagem do design system, com fundo desfocado por overlay e superfície de cartão elevada elegante.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
npm --prefix dashboard run build
```

### Execução
Abrir o modal de "Confirmar Disparo em massa" e "Desconectar Bot" no dashboard para validar a aparência refinada.

## Rollback
Executar `git checkout` nos arquivos alterados.
