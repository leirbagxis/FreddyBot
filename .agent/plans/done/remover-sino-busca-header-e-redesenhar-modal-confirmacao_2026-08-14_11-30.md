# Plano: Remover Sino do Header, Manter Campo de Busca Visível e Redesenhar Modal de Confirmação do Broadcast

## Pedido do Usuário
1. **Remover Sino do Header**: Remover o ícone e menu de notificações por "sino" (`Bell`) na `.admin-topbar`.
2. **Campo de Busca Sempre Visível**: Garantir o campo de busca (`Search`) visível no header em todas as telas da dashboard admin.
3. **Redesenhar Modal de Confirmação do Broadcast**: Reconstruir a modal de confirmação de envio (`ConfirmModal.tsx`) com a estética da dashboard admin (shadcn UI, ícone com iluminação suave, tipografia nítida, botões modernos e zero textos escuros/transparentes).

---

## Análise Técnica

### 1. Header (`AdminTopbar.tsx`)
- Remover o estado `notificationsOpen`, a ref `notificationsRef`, a função `handleAlert`, o ícone `Bell` e a `section.admin-notifications-panel`.
- Manter o campo de busca `<Input>` visível continuamente, ajustando o placeholder de acordo com a aba ativa (ex.: "Buscar no sistema..." ou "Buscar usuários/canais...").

### 2. Modal de Confirmação (`ConfirmModal.tsx` & `AdminNoticeTab.tsx`)
- Redesenhar `ConfirmModal.tsx` utilizando shadcn `<Dialog>`, `<DialogContent>`, `<DialogHeader>`, `<DialogTitle>`, `<DialogDescription>` e `<DialogFooter>`.
- Adicionar um container de ícone elegante:
  - Para envios perigosos (`danger`): iluminação vermelha suave `bg-red-500/10 text-red-500 border border-red-500/20`.
  - Para confirmação padrão/broadcast: iluminação suave no tom de destaque `bg-accent/15 text-accent border border-accent/30`.
- Tipografia: `DialogTitle` com `text-foreground text-lg font-bold`, `DialogDescription` com `text-muted-foreground text-xs leading-relaxed`.
- Botões: `<Button variant="outline">` para Cancelar e `<Button variant="default">` (ou `destructive`) para Confirmar envio, com `h-11 rounded-xl font-semibold`.

---

## Arquivos que serão modificados

- `dashboard/src/components/admin/AdminTopbar.tsx` (remover sino, manter search contínuo)
- `dashboard/src/components/ConfirmModal.tsx` (redesenhar modal de confirmação)
- `dashboard/src/components/AdminNoticeTab.tsx` (atualizar chamada do ConfirmModal no Broadcast)

---

## Passos Detalhados

### Passo 1 — Atualizar `AdminTopbar.tsx`
1. Remover o import do ícone `Bell` e a estrutura `<div className="admin-notifications">...</div>`.
2. Tornar o container `<div className="admin-topbar-search">` visível em todas as abas.

### Passo 2 — Redesenhar `ConfirmModal.tsx`
1. Reconstruir a modal com design minimalista shadcn:
   - Header centralizado com ícone destacado.
   - Tipografia limpa com cores de texto garantidas `text-foreground` e `text-muted-foreground`.
   - Botões de ação bem espaçados com cantos arredondados `rounded-xl`.

### Passo 3 — Atualizar `AdminNoticeTab.tsx`
1. Passar um texto claro e amigável na abertura da modal no Broadcast.

---

## Como Testar
```bash
cd dashboard && npm run build
```
- Verificar a compilação limpa sem avisos de código não utilizado.
- Testar o header sem o sino, com o campo de busca ativo e a nova modal de confirmação do broadcast.

## Rollback
```bash
git checkout dashboard/src/components/admin/AdminTopbar.tsx dashboard/src/components/ConfirmModal.tsx dashboard/src/components/AdminNoticeTab.tsx
```
