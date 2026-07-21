export interface ServerConfig {
  id: number;
  maintence: boolean;
  forceJoin: boolean;
  globalDefaultCaption: string;
  globalNewPackCaption: string;
  fixedPostBuilderEnabled: boolean;
  fixedPostBuilderKey: string;
  fixedPostBuilderPayload: string;
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
  newPackMessageButtons: boolean;
  newPackStickerButtons: boolean;
  newPackMessagePosition: 'above' | 'below';
  newPackReplyToSticker: boolean;
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


export interface ChannelEvent {
  id: string;
  channelId: number;
  channelTitle: string;
  ownerId: number;
  actorId: number;
  source: string;
  eventType: string;
  status: string;
  messageType: string;
  telegramMessageId: number;
  sessionId: string;
  errorMessage: string;
  metadata: string;
  created_at: string;
}

export interface AdminLogsResponse {
  events: ChannelEvent[];
  total: number;
  limit: number;
  offset: number;
}

/* ===== Premium Features (Admin) ===== */
export interface PremiumFeature {
  key: string;
  name: string;
  description: string;
  enabled: boolean;
  price: number;
  createdAt: string;
  updatedAt: string;
}

export interface AdminLogsFilters {
  channelId?: string;
  ownerId?: string;
  actorId?: string;
  source?: string;
  eventType?: string;
  status?: string;
  sessionId?: string;
  q?: string;
  dateFrom?: string;
  dateTo?: string;
  limit?: number;
  offset?: number;
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

/* ===== Subscription / Premium ===== */
export interface Subscription {
  id: string;
  userId: number;
  status: 'active' | 'canceled' | 'expired';
  currentPeriodStart: string;
  currentPeriodEnd: string;
  extraChannels: number;
  cancelAtPeriodEnd: boolean;
  telegramPaymentId: string;
  extraChannelPayments: string;
  createdAt: string;
  updatedAt: string;
}

export interface UserFeatures {
  managedPremiumAccount?: boolean;
  customEmojis?: boolean;
  extraChannels?: number;
}

export interface SubscriptionStatus {
  hasSubscription: boolean;
  subscription?: Subscription;
  features?: UserFeatures;
  basePrice: number;
  extraChannelPrice: number;
  starsTestMode?: boolean;
  hasAccount?: boolean;
  premiumEnabled?: boolean;
  connectedAccountEnabled?: boolean;
}

/* ===== Connected Account (MTProto) ===== */
export interface AccountStatus {
  status: 'connected' | 'disconnected';
  telegramId?: number;
  username?: string;
  avatarUrl?: string;
  connectedAt?: string;
  lastUsedAt?: string;
}

export interface AuthStatus {
  step: 'phone' | 'code' | 'password' | 'done' | 'error';
  error?: string;
  hasPassword?: boolean;
}

/* ===== Admin MTProto Accounts ===== */
export interface AdminMTProtoAccount {
  id: string;
  label: string;
  phoneNumber: string;
  telegramUserId: number;
  username: string;
  firstName: string;
  enabled: boolean;
  status: string; // "connected", "disconnected", "error"
  lastUsedAt: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface AdminAuthStep {
  step: 'phone' | 'code' | 'password' | 'done' | 'error';
  sessionId?: string;
  hasPassword?: boolean;
  error?: string;
}

/* ===== Invoice / Payment ===== */

export type InvoiceStatus = 'paid' | 'cancelled' | 'failed' | 'pending';

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
        showPopup: (params: { title?: string; message: string; buttons?: { type?: string; text: string; id?: string }[] }, callback?: (buttonId: string) => void) => void;
        setHeaderColor: (color: string) => void;
        setBackgroundColor: (color: string) => void;
        openInvoice: (url: string, callback: (status: InvoiceStatus) => void) => void;
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
