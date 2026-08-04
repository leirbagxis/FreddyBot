# Plano: corrigir-modais-confirmacao-opacos

## Pedido do usuário
Corrigir modais de confirmação transparentes e visualmente ruins.

## Objetivo
Garantir que todos os diálogos de confirmação usem uma superfície opaca, borda e sombra consistentes com o painel minimalista, preservando conteúdo, acessibilidade e ações existentes.

## Contexto atual
- `DialogContent` possui superfície padrão `bg-popover`.
- Os diálogos de assinatura e premium sobrescrevem essa superfície com `bg-[var(--bg)]`, variável que não representa o fundo opaco do painel.
- `ConfirmModal` é reutilizado em auditoria, avisos e tela inicial; ajustes no componente base afetam confirmações comuns.

## Arquivos analisados
- `dashboard/src/components/ui/dialog.tsx`
- `dashboard/src/components/ConfirmModal.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/components/PremiumTab.tsx`
- `dashboard/src/index.css`

## Arquivos que poderão ser modificados
- `dashboard/src/components/ui/dialog.tsx`
- `dashboard/src/components/ConfirmModal.tsx`
- `dashboard/src/components/AdminSubscriptionsTab.tsx`
- `dashboard/src/components/PremiumTab.tsx`
- `dashboard/src/index.css`
- Arquivos de teste do dashboard, se existirem e forem necessários
- `.agent/memory/memory.md` e histórico deste plano

## Estratégia de implementação
Manter o componente de diálogo como fonte de verdade para fundo e overlay opacos; remover sobrescritas transparentes das confirmações específicas e aplicar acabamento neutro apenas no escopo administrativo.

## Passos detalhados

1. Reforçar no `DialogContent` uma superfície `bg-popover` opaca, borda e sombra legíveis.
2. Ajustar o overlay para escurecer a tela sem deixar o conteúdo do modal translúcido.
3. Remover `bg-[var(--bg)]` dos diálogos de assinaturas e premium, usando o token de superfície padrão.
4. Aplicar ao `ConfirmModal` espaçamento e rodapé consistentes, sem modificar callbacks ou textos.
5. Verificar visualmente classes e executar type-check/build do dashboard.

## Riscos
- Mudanças globais no diálogo podem afetar modais não administrativos; a alteração preservará tokens existentes e será limitada a superfície, borda, sombra e overlay.
- Não serão alteradas regras de cancelamento, reembolso ou assinatura.

## Impactos esperados
- Confirmações passam a ter fundo sólido e contraste suficiente em desktop e mobile.
- O visual fica coerente entre auditoria, avisos, assinatura e premium.

## Compatibilidade
- Desktop
- Mobile
- Chrome
- Safari/iOS
- Docker
- CI/CD

## Como testar

### Build
```bash
cd dashboard && PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npm run build
```

### Testes
```bash
cd dashboard && PATH=/home/gabriel/.local/share/mise/installs/node/26.5.1/bin:$PATH npx tsc --noEmit
```

### Execução
```bash
make run
```

## Rollback
Reverter apenas as classes e regras visuais modificadas, mantendo os handlers e a estrutura dos diálogos.

## Observações
- Nenhum fluxo de confirmação será removido ou automatizado.
- Nenhuma configuração de backend, banco ou compose será alterada.
