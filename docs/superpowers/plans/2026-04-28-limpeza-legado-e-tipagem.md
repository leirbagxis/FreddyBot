# Task 1: Limpeza de Código Legado e Tipagem Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remover funcionalidades de economia (Shop, Ranking) e preparar definições de tipos TypeScript baseadas nos modelos Go.

**Architecture:** Mapeamento de structs Go para interfaces TypeScript e limpeza de rotas/componentes não utilizados.

**Tech Stack:** TypeScript, React, React Router.

---

### Task 1: Criar Definições de Tipos

**Files:**
- Create: `webapp/src/types/index.ts`

- [ ] **Step 1: Criar o arquivo de tipos com as interfaces correspondentes aos structs Go**

```typescript
export interface User {
  id: number;
  first_name: string;
  username: string;
  is_admin: boolean;
  channels: Channel[];
}

export interface Channel {
  id: number;
  title: string;
  newPackCaption: string;
  inviteUrl: string;
  ownerId: number;
  defaultCaption?: DefaultCaption;
  buttons: Button[];
}

export interface DefaultCaption {
  captionId: string;
  caption: string;
  messagePermission?: MessagePermission;
  buttonsPermission?: ButtonsPermission;
}

export interface Button {
  buttonId: string;
  nameButton: string;
  buttonUrl: string;
  positionX: number;
  positionY: number;
}

// Interfaces auxiliares baseadas no models.go (DefaultCaption)
export interface MessagePermission {
  // Adicione campos se necessário ou deixe como interface vazia/any se não especificado
  [key: string]: any;
}

export interface ButtonsPermission {
  // Adicione campos se necessário ou deixe como interface vazia/any se não especificado
  [key: string]: any;
}
```

### Task 2: Remover Arquivos Legados

**Files:**
- Delete: `webapp/src/pages/Shop.tsx`
- Delete: `webapp/src/pages/Ranking.tsx`
- Delete: `webapp/src/components/shop/ShopItemCard.tsx`
- Delete: `webapp/src/components/ranking/RankingTable.tsx`

- [ ] **Step 1: Remover arquivos de páginas e componentes de Shop e Ranking**

Execute:
```bash
rm webapp/src/pages/Shop.tsx webapp/src/pages/Ranking.tsx
rm -rf webapp/src/components/shop webapp/src/components/ranking
```

### Task 3: Limpar Rotas e UI

**Files:**
- Modify: `webapp/src/App.tsx`
- Modify: `webapp/src/components/Layout.tsx`
- Modify: `webapp/src/components/dashboard/UserOverview.tsx`
- Modify: `webapp/src/components/dashboard/Inventory.tsx`
- Modify: `webapp/src/pages/Dashboard.tsx`
- Modify: `webapp/src/api.ts`

- [ ] **Step 1: Remover importações e rotas legadas no App.tsx**

- [ ] **Step 2: Remover links de Loja e Rankings no Layout.tsx**

- [ ] **Step 3: Remover botões de Loja e Rankings no UserOverview.tsx**

- [ ] **Step 4: Remover funcionalidade de venda na Loja no Inventory.tsx e Dashboard.tsx**

- [ ] **Step 5: Remover funções de API legadas no api.ts**

### Task 4: Verificação

- [ ] **Step 1: Verificar se não há erros de compilação ou importações quebradas**

Execute:
```bash
cd webapp && npx tsc --noEmit
```
