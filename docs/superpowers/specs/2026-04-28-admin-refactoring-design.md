# Design Doc: Admin Panel Refactoring

**Date:** 2026-04-28
**Topic:** Admin Panel Refactoring (Task 7)
**Status:** Draft

## 1. Goal
Refactor the Admin panel in `webapp/src/pages/Admin.tsx` to focus on monitoring users and their channels, removing all legacy "Store Bot" features (Economy, Items, Factory, etc.).

## 2. Approach
Simplified Full Rewrite. We will strip away the multi-tab complexity and focus on a robust, searchable user list that displays channel ownership.

## 3. Architecture & Components

### 3.1 Data Source
- API: `adminFetchUsers()`
- Model: `User { id, first_name, username, is_admin, channels: Channel[] }`

### 3.2 UI Components
- **AdminPage (Admin.tsx):** Main container, handles data fetching and search state.
- **Sidebar:** Updated to show only relevant links (Users).
- **UsersList:** A refined version of the previous UsersTab, optimized for the new model.
- **UserCard:** Displays user summary (Name, ID, Admin status, Channel count).
- **ChannelList (Nested):** Shown when a user is expanded, listing their channels with titles and invite links.

## 4. Cleanup Plan
- Remove imports and usage of: `adminFetchItems`, `fetchAdminBalanceLogs`, `adminFetchConfig`, `adminCreateItem`, etc.
- Delete files:
    - `webapp/src/components/admin/FactoryTab.tsx`
    - `webapp/src/components/admin/EconomyTab.tsx`
    - `webapp/src/components/admin/LogsTab.tsx`
    - `webapp/src/components/admin/ConfigTab.tsx`
    - `webapp/src/components/admin/BotTab.tsx`
- Update `webapp/src/components/admin/constants.ts` to remove legacy `NAV_ITEMS`.
- Update `webapp/src/components/admin/Sidebar.tsx` to reflect the new navigation.

## 5. Functional Requirements
- Fetch all users on mount.
- Filter users by Username or ID (case-insensitive).
- Display user details: First Name, Username, Telegram ID, Admin Badge.
- Display "Channels Managed" count.
- Expand user to see channel details: Title, Invite Link.
- Remove all legacy "Balance" and "Economy" buttons/logic.

## 6. Implementation Steps
1. Update `constants.ts` to reflect the new `NAV_ITEMS`.
2. Clean up `Sidebar.tsx`.
3. Refactor `UsersTab.tsx` into a robust `UsersList`.
4. Rewrite `Admin.tsx` to integrate the new `UsersList` and remove legacy logic.
5. Delete unused component files.
