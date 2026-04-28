# Especificação: Reconstrução do Webapp (Channel Manager)

Este documento detalha a reconstrução total do frontend React para alinhar as funcionalidades ao backend Go (FreddyBot), focando no gerenciamento de canais do Telegram, legendas, botões e permissões.

O usuario entra na dahboard pela rota /dashboard/:channelID, vai pra rota de autenticacao com telegram e ai pega os dados do canal com aquele Id.

## 1. Objetivos
- Substituir a interface de "Loja/Economia" por um sistema de **Gerenciamento de Canais**.
- Implementar integração completa com as rotas de API existentes no Go.
- Garantir que usuários admin tenham visibilidade global de usuários e canais.
- Manter a estética conforme a identidade do projeto.

## 2. Arquitetura do Sistema

### 2.1 Fluxo de Autenticação
1. O Webapp extrai o `initData` do Telegram WebApp.
2. Envia para `POST /api/auth`.
3. O servidor valida os dados e retorna um **JWT via Cookie Seguro (HttpOnly)**.
4. Todas as requisições subsequentes utilizam esse cookie para autenticação.

### 2.2 Rotas da API (Mapeamento)
O arquivo `webapp/src/api.ts` será reconstruído para refletir:
- **Auth:** `POST /api/auth`, `POST /api/me/channels`.
- **Canais:** `GET /api/channel/:channelId`.
- **Legendas:** `PUT /api/channel/:channelId/caption`, `PUT /api/channel/:channelId/newpackcaption`.
- **Permissões:** `PUT /api/channel/:channelId/caption/permissions`, `PUT /api/channel/:channelId/buttons/permissions`.
- **Botões:** `POST/PUT/DELETE` em `/api/channel/:channelId/buttons` e `/api/channel/:channelId/buttons/layout`.
- **Admin:** `GET /admin/api/users`.

## 3. Componentes da Interface

### 3.1 Dashboard (`Dashboard.tsx`)
- Canal do usuario solicitado atravez da rota /dashboard/:channelID
- Cards informativos com Título do Canal e link de convite.
- Botão "Configurar" que redireciona para o Editor do Canal.

### 3.2 Editor de Canal (`ChannelEditor.tsx` - Novo)
Interface com abas para:
- **Legendas:** Edição da legenda padrão e legenda de novos packs.
- **Botões:** Interface para adicionar, editar e organizar o layout dos botões.
- **Permissões:** Toggles para habilitar/desabilitar tipos de mídia (Foto, Vídeo, GIF, etc.) e preview de links.

### 3.3 Painel Admin (`Admin.tsx`)
- Lista global de usuários e seus respectivos canais.
- Ações administrativas conforme definido nos controladores Go (ex: `SendNoticeAdminController`).

## 4. Tecnologias
- **Frontend:** React + TypeScript + Vite.
- **Estilização:** Tailwind CSS (Mantendo o tema escuro/neon).
- **Comunicação:** Axios (Configurado para `withCredentials: true`).

## 5. Próximos Passos
1. Limpeza do código legado (Remoção de Shop, Inventory, Ranking).
2. Atualização dos tipos TypeScript para bater com os modelos Go (`models.go`).
3. Implementação da nova camada de API.
4. Construção das novas telas de Dashboard e Editor.
