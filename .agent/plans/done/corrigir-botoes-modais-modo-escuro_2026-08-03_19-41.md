# Plano: Corrigir Botões Finos e Modal Branco no Modo Escuro

## Pedido do usuário
"Continua feio para um caralho. Os botoes estao finos e quando ta no modo escuro fica branco a modal de confirmacao. pqp"

## Objetivo
1. **Corrigir Fundo Branco no Modo Escuro**: Garantir que no modo escuro (`[data-theme="dark"]` e `.dark`), os modais e diálogos tenham fundo escuro elegante (`#121724` / `#0f1422`) com texto claro (`#f1f5f9`), corrigindo a regra de `.admin-layout-v2` que estava forçando `--minimal-paper: #ffffff` e `--popover: #ffffff`.
2. **Aumentar Espessura/Tamanho dos Botões**: Substituir os botões finos (`h-8` / `32px`) no `ConfirmModal` e em modais de confirmação por botões robustos, com altura `h-12` (48px), padding amplo (`px-5 py-3`), cantos arredondados `rounded-xl`, tipografia `font-semibold text-sm` e estados de hover/foco premium.

## Contexto atual
- A classe `.admin-layout-v2` definia variáveis locais `--minimal-paper: #ffffff` e `--popover: var(--minimal-paper)`. Como essa regra aparecia ao final do `index.css`, ela sobrescrevia `--popover` para branco mesmo quando o tema escuro estava ativo (`[data-theme="dark"]`).
- O `ConfirmModal` usava a variação padrão de tamanho do botão (que possui apenas 32px-40px de altura), dando visual "fino" e frágil.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/src/components/ConfirmModal.tsx`
- `dashboard/src/components/ui/button.tsx`
- `dashboard/src/components/ui/dialog.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/ConfirmModal.tsx`

## Estratégia de implementação
1. **Em `dashboard/src/index.css`**:
   - Adicionar a regra de override para modo escuro em `.admin-layout-v2`:
     ```css
     [data-theme="dark"] .admin-layout-v2,
     .dark .admin-layout-v2,
     [data-theme="dark"] {
       --minimal-paper: #121724;
       --minimal-text: #e8edf5;
       --minimal-muted: #94a3b8;
       --minimal-line: rgba(255, 255, 255, 0.08);
       --popover: #121724;
       --popover-foreground: #e8edf5;
     }
     ```
   - Forçar explicitamente `[data-theme="dark"] [data-slot='dialog-content']` e `.dark [data-slot='dialog-content']` a usarem `background-color: #121724 !important; color: #e8edf5 !important; border-color: rgba(255, 255, 255, 0.12) !important;`.
2. **Em `dashboard/src/components/ConfirmModal.tsx`**:
   - Atualizar os botões para altura robusta `h-12` (48px), `px-5`, `font-semibold text-[15px]`, `rounded-xl`, garantindo área de toque ampla, visual imponente e excelente ergonomia.

## Passos detalhados
1. Editar `dashboard/src/index.css`:
   - Adicionar overrides escuros para `.admin-layout-v2` e `[data-slot='dialog-content']`.
2. Editar `dashboard/src/components/ConfirmModal.tsx`:
   - Atualizar a altura dos botões para `h-12` (48px) e tipografia `font-semibold text-[15px]`.
3. Compilar o dashboard (`npm --prefix dashboard run build`) para verificar.

## Riscos
- Risco muito baixo. Ajustes estéticos direcionados a cores no escuro e dimensão de botões.

## Impactos esperados
- Modal no modo escuro totalmente escuro (`#121724`), sem qualquer trecho branco sobressalente.
- Botões encorpados, altos (48px) e elegantes, fáceis de clicar e visualmente marcantes.

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
Abrir modais em modo escuro e modo claro para verificar botões encorpados e cor de fundo escura correta.

## Rollback
`git checkout` nos arquivos alterados.
