# Admin Panel Refactoring Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the Admin panel to focus on monitoring users and their channels, removing legacy store bot features.

**Architecture:** Simplified single-tab layout with a robust searchable user list and expandable channel details.

**Tech Stack:** React, TypeScript, Tailwind CSS, Lucide React, Axios.

---

### Task 1: Update Constants and Sidebar

**Files:**
- Modify: `webapp/src/components/admin/constants.ts`
- Modify: `webapp/src/components/admin/Sidebar.tsx`

- [ ] **Step 1: Simplify NAV_ITEMS in `constants.ts`**
  Remove all legacy items and keep only "USUÁRIOS".

```typescript
import { Users } from 'lucide-react';

export const NAV_ITEMS = [
  { id: 'users', label: 'USUÁRIOS', icon: Users },
];
```

- [ ] **Step 2: Clean up `Sidebar.tsx`**
  Remove the "Sincronizar" button (loadData will be handled by the main Admin page) and keep it simple.

- [ ] **Step 3: Commit changes**
  `git commit -m "refactor: simplify admin navigation and sidebar"`

### Task 2: Refactor UsersTab to UsersList

**Files:**
- Modify: `webapp/src/components/admin/UsersTab.tsx`

- [ ] **Step 1: Update Types and Component Signature**
  Adapt `UsersTabProps` to use the new `User` model from `../types`.

- [ ] **Step 2: Update Layout to show Channels**
  Modify the card to show channel count and list channels when expanded.

- [ ] **Step 3: Remove legacy Balance/Economy buttons**

- [ ] **Step 4: Commit changes**
  `git commit -m "refactor: update UsersTab to new data model and remove legacy logic"`

### Task 3: Rewrite Admin.tsx

**Files:**
- Modify: `webapp/src/pages/Admin.tsx`

- [ ] **Step 1: Clean up imports**
  Remove all legacy API imports and tab components.

- [ ] **Step 2: Simplify State**
  Only keep `users`, `isLoading`, `searchQuery`, and `expandedUserId`.

- [ ] **Step 3: Update loadData**
  Use `adminFetchUsers()` and handle errors.

- [ ] **Step 4: Implement Filter Logic**
  Filter by Username or ID.

- [ ] **Step 5: Clean up Render Logic**
  Only render the header and `UsersTab`.

- [ ] **Step 6: Commit changes**
  `git commit -m "refactor: complete rewrite of Admin.tsx for user monitoring"`

### Task 4: Final Cleanup

**Files:**
- Delete: `webapp/src/components/admin/FactoryTab.tsx`
- Delete: `webapp/src/components/admin/EconomyTab.tsx`
- Delete: `webapp/src/components/admin/LogsTab.tsx`
- Delete: `webapp/src/components/admin/ConfigTab.tsx`
- Delete: `webapp/src/components/admin/BotTab.tsx`

- [ ] **Step 1: Delete all unused component files**

- [ ] **Step 2: Commit changes**
  `git commit -m "cleanup: remove legacy admin tab components"`
