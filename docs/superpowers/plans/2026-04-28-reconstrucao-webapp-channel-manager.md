# Reconstrução do Webapp (Channel Manager) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transformar o Webapp de uma economia de RPG para um gerenciador de canais do Telegram completo, alinhado ao backend Go.

**Architecture:** Frontend React (Vite) consumindo API Go com autenticação JWT via Cookie. Foco em edição de legendas, botões e permissões de canais.

**Tech Stack:** React, TypeScript, Tailwind CSS, Axios, Lucide React.

---

### Task 1: Limpeza de Código Legado e Tipagem

**Files:**
- Create: `webapp/src/types/index.ts`
- Modify: `webapp/src/App.tsx`
- Delete: `webapp/src/pages/Shop.tsx`, `webapp/src/pages/Ranking.tsx`, `webapp/src/components/shop/*`, `webapp/src/components/ranking/*`

- [ ] **Step 1: Criar definições de tipos baseadas no models.go**
- [ ] **Step 2: Remover arquivos de páginas e componentes legados**
- [ ] **Step 3: Limpar rotas no App.tsx para evitar erros de importação**
- [ ] **Step 4: Commit**

```bash
git add webapp/src/types/index.ts webapp/src/App.tsx
git rm webapp/src/pages/Shop.tsx webapp/src/pages/Ranking.tsx
git commit -m "chore: limpar codigo legado e definir tipos base"
```

### Task 2: Refatoração da Camada de API e Auth

**Files:**
- Modify: `webapp/src/api.ts`, `webapp/src/App.tsx`

- [ ] **Step 1: Reescrever api.ts para mapear as rotas do Go**
- [ ] **Step 2: Atualizar o fluxo de autenticação no App.tsx para usar /api/auth**
- [ ] **Step 3: Implementar tratamento de erro de autenticação (401/403)**
- [ ] **Step 4: Commit**

### Task 3: Novo Dashboard (Lista de Canais)

**Files:**
- Modify: `webapp/src/pages/Dashboard.tsx`
- Create: `webapp/src/components/dashboard/ChannelCard.tsx`

- [ ] **Step 1: Implementar busca de canais do usuário via /api/me/channels**
- [ ] **Step 2: Criar componente ChannelCard com estilo Cyberpunk**
- [ ] **Step 3: Renderizar lista de canais no Dashboard**
- [ ] **Step 4: Commit**

### Task 4: Editor de Canal - Estrutura e Legendas

**Files:**
- Create: `webapp/src/pages/ChannelEditor.tsx`, `webapp/src/components/editor/CaptionTab.tsx`

- [ ] **Step 1: Criar ChannelEditor com sistema de abas (Legendas, Botões, Permissões)**
- [ ] **Step 2: Implementar busca de dados do canal via /api/channel/:id**
- [ ] **Step 3: Criar aba de Legendas com editores para Default e NewPack**
- [ ] **Step 4: Implementar salvamento de legendas**
- [ ] **Step 5: Commit**

### Task 5: Editor de Canal - Botões e Layout

**Files:**
- Create: `webapp/src/components/editor/ButtonsTab.tsx`

- [ ] **Step 1: Criar interface para listar botões existentes**
- [ ] **Step 2: Adicionar modal para novo botão (Nome, URL)**
- [ ] **Step 3: Implementar edição e exclusão de botões**
- [ ] **Step 4: Criar interface de grid para reordenar botões (Layout)**
- [ ] **Step 5: Commit**

### Task 6: Editor de Canal - Permissões de Mídia

**Files:**
- Create: `webapp/src/components/editor/PermissionsTab.tsx`

- [ ] **Step 1: Criar lista de toggles para permissões de mensagem e botões**
- [ ] **Step 2: Mapear os campos do models.go (LinkPreview, Photo, Video, GIF, etc.)**
- [ ] **Step 3: Implementar atualização em tempo real via API**
- [ ] **Step 4: Commit**

### Task 7: Refatoração do Painel Admin

**Files:**
- Modify: `webapp/src/pages/Admin.tsx`

- [ ] **Step 1: Atualizar Admin.tsx para buscar todos os usuários e canais via /admin/api/users**
- [ ] **Step 2: Adaptar a tabela de usuários para o novo modelo de dados**
- [ ] **Step 3: Adicionar busca por ID/Username**
- [ ] **Step 4: Commit**

### Task 8: Polimento Visual e Verificação

**Files:**
- Modify: `webapp/src/index.css`, `webapp/src/App.tsx`

- [ ] **Step 1: Garantir que as animações de terminal e cores neon estão consistentes**
- [ ] **Step 2: Verificar responsividade em dispositivos móveis (Telegram WebApp)**
- [ ] **Step 3: Teste final de todos os fluxos (CRUD de canal, botões e admin)**
- [ ] **Step 4: Commit Final**
