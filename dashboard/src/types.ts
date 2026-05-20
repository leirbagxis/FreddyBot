export interface ServerConfig {
  id: number;
  maintence: boolean;
  forceJoin: boolean;
  globalDefaultCaption: string;
  globalNewPackCaption: string;
  created_at: string;
  updated_at: string;
}

export interface Permission {
  messagePermissionId?: string;
  buttonsPermissionId?: string;
  linkPreview?: boolean;
  message: boolean;
  audio: boolean;
  video: boolean;
  photo: boolean;
  document: boolean;
  sticker: boolean;
  gif: boolean;
  reactions?: boolean;
  ownerCaptionId: string;
  created_at: string;
  updated_at: string;
}

export interface Caption {
  captionId: string;
  caption: string;
  messagePermission: Permission;
  buttonsPermission: Permission;
  ownerChannelId: number;
  created_at: string;
  updated_at: string;
}

export interface Button {
  buttonId: string;
  nameButton: string;
  buttonUrl: string;
  positionX: number;
  positionY: number;
  ownerChannelId: number;
  created_at: string;
  updated_at: string;
}

export interface Channel {
  id: number;
  title: string;
  newPackCaption: string;
  inviteUrl: string;
  ownerId: number;
  reactions: string;
  reactionPosition: number;
  dynamicLinks: boolean;
  dlBotButtons: boolean;
  dlBotCaptions: boolean;
  dlBotReactions: boolean;
  defaultCaption: Caption;
  buttons: Button[];
  customCaptions: Caption[];
  created_at: string;
  updated_at: string;
}

export interface User {
  id: number;
  first_name: string;
  username: string;
  is_admin: boolean;
  is_blacklisted: boolean;
  isContribute: boolean;
  channels: Channel[] | null;
  created_at: string;
  updated_at: string;
}

export interface AdminDashboardData {
  success: boolean;
  users: User[];
  channels: Channel[];
}

export interface AuditResult {
  userId: number;
  firstName: string;
  channels: Channel[];
}

export interface DashboardData {
  channel: Channel;
  user: User;
}

export interface ChannelsResponse {
  channels: Channel[];
  success: boolean;
}

/* ===== Telegram WebApp ===== */
export interface TelegramUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  photo_url?: string;
  language_code?: string;
}

declare global {
  interface Window {
    Telegram?: {
      WebApp: {
        ready: () => void;
        expand: () => void;
        close: () => void;
        colorScheme: 'light' | 'dark';
        themeParams: Record<string, string | undefined>;
        initData: string;
        initDataUnsafe: {
          query_id?: string;
          user?: TelegramUser;
          auth_date?: string | number;
          hash?: string;
        };
        viewportHeight: number;
        viewportStableHeight: number;
        isExpanded: boolean;
        headerColor: string;
        backgroundColor: string;
        showConfirm: (message: string, callback: (ok: boolean) => void) => void;
        setHeaderColor: (color: string) => void;
        setBackgroundColor: (color: string) => void;
        BackButton: {
          isVisible: boolean;
          show: () => void;
          hide: () => void;
          onClick: (callback: () => void) => void;
          offClick: (callback: () => void) => void;
        };
        CloudStorage: {
          setItem: (key: string, value: string, callback?: (error: Error | null, success?: boolean) => void) => void;
          getItem: (key: string, callback: (error: Error | null, value?: string) => void) => void;
          removeItem: (key: string, callback?: (error: Error | null, success?: boolean) => void) => void;
          getKeys: (callback: (error: Error | null, keys?: string[]) => void) => void;
        };
      };
    };
  }
}
