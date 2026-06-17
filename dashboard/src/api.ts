import { DashboardData, Button, Permission, ChannelsResponse, AdminDashboardData, AdminLogsFilters, AdminLogsResponse, ConnectStatus } from './types';

export interface AuthRequestBody {
    channelID: number;
    user: {
        id: number;
        first_name: string;
        last_name: string;
        username: string;
        photo_url: string;
        auth_date: string;
        hash: string;
    };
}

const apiFetch = async (url: string, options: RequestInit = {}) => {
    const response = await fetch(url, {
        ...options,
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers || {}),
        },
    });

    if (!response.ok) {
        const errBody = await response.text().catch(() => '');
        throw new Error(errBody || `API Error (${response.status})`);
    }

    if (response.status !== 204) {
        try {
            return await response.json();
        } catch {
            return null;
        }
    }
    return null;
};

export const login = async (initData: string, userID: number) => {
    return apiFetch('/api/login', {
        method: 'POST',
        headers: {
            'x-telegram-init-data': initData,
        },
        body: JSON.stringify({ userID }),
    });
};

export const fetchDashboardData = async (channelId: string): Promise<DashboardData> => {
    return apiFetch(`/api/channel/${channelId}`, {
        method: 'GET',
    });
};

export const fetchUserChannels = async (): Promise<ChannelsResponse> => {
    return apiFetch(`/api/me/channels`, {
        method: 'GET',
    });
};

export const fetchAdminDashboard = async (): Promise<AdminDashboardData> => {
    const response = await apiFetch(`/api/admin/overview`);
    const data = response?.data || {};
    return { 
        success: true, // Assuming success if apiFetch didn't throw
        users: data.users || [], 
        channels: data.channels || [] 
    };
};

export const updateDefaultCaption = async (channelId: number, caption: string) => {
    return apiFetch(`/api/channel/${channelId}/caption`, {
        method: 'PUT',
        body: JSON.stringify({ caption }),
    });
};

export const updateNewPackCaption = async (channelId: number, payload: {
    newPackCaption: string;
    newPackMessageButtons?: boolean;
    newPackStickerButtons?: boolean;
    newPackMessagePosition?: 'above' | 'below';
    newPackReplyToSticker?: boolean;
}) => {
    return apiFetch(`/api/channel/${channelId}/newpackcaption`, {
        method: 'PUT',
        body: JSON.stringify(payload),
    });
};

export const updateReactions = async (channelId: number, reactions: string) => {
    return apiFetch(`/api/channel/${channelId}/reactions`, {
        method: 'PUT',
        body: JSON.stringify({ reactions }),
    });
};

export const updateReactionsActive = async (channelId: number, active: boolean) => {
    return apiFetch(`/api/channel/${channelId}/reactions/active`, {
        method: 'PUT',
        body: JSON.stringify({ active }),
    });
};

export const updateReactionPosition = async (channelId: number, reactionPosition: number) => {
    return apiFetch(`/api/channel/${channelId}/reactions/position`, {
        method: 'PUT',
        body: JSON.stringify({ reactionPosition }),
    });
};

export const updateDynamicLinks = async (channelId: number, settings: {
    dynamicLinks: boolean;
    dlBotButtons: boolean;
    dlBotCaptions: boolean;
    dlBotReactions: boolean;
}) => {
    return apiFetch(`/api/channel/${channelId}/dynamic-links`, {
        method: 'PUT',
        body: JSON.stringify(settings),
    });
};

export const updateMessagePermission = async (channelId: number, perms: Permission) => {
    const payload = {
        linkPreview: Boolean(perms.linkPreview),
        message: Boolean(perms.message),
        audio: Boolean(perms.audio),
        video: Boolean(perms.video),
        photo: Boolean(perms.photo),
        document: Boolean(perms.document),
        sticker: Boolean(perms.sticker),
        gif: Boolean(perms.gif),
        reactions: Boolean(perms.reactions),
    };
    return apiFetch(`/api/channel/${channelId}/caption/permissions`, {
        method: 'PUT',
        body: JSON.stringify(payload),
    });
};

export const updateButtonsPermission = async (channelId: number, perms: Permission) => {
    const payload = {
        message: Boolean(perms.message),
        audio: Boolean(perms.audio),
        video: Boolean(perms.video),
        photo: Boolean(perms.photo),
        document: Boolean(perms.document),
        sticker: Boolean(perms.sticker),
        gif: Boolean(perms.gif),
    };
    return apiFetch(`/api/channel/${channelId}/buttons/permissions`, {
        method: 'PUT',
        body: JSON.stringify(payload),
    });
};

export const createButton = async (channelId: number, button: Partial<Button>) => {
    return apiFetch(`/api/channel/${channelId}/buttons`, {
        method: 'POST',
        body: JSON.stringify({
            nameButton: button.nameButton,
            buttonUrl: button.buttonUrl || undefined,
        }),
    });
};

export const updateButton = async (channelId: number, buttonId: string, button: Partial<Button>) => {
    return apiFetch(`/api/channel/${channelId}/buttons/${buttonId}`, {
        method: 'PUT',
        body: JSON.stringify(button),
    });
};

export const deleteButton = async (channelId: number, buttonId: string) => {
    return apiFetch(`/api/channel/${channelId}/buttons/${buttonId}`, {
        method: 'DELETE',
    });
};

export const updateLayoutButtons = async (channelId: number, layout: any[][]) => {
    return apiFetch(`/api/channel/${channelId}/buttons/layout`, {
        method: 'PUT',
        body: JSON.stringify({ layout }),
    });
};

export const transferChannel = async (oldOwnerId: number, newOwnerId: number, channelId: number) => {
    return apiFetch(`/api/channel/transfer`, {
        method: 'POST',
        body: JSON.stringify({ oldOwnerId, newOwnerId, channelId }),
    });
};

export const fetchUserInfo = async (usernameOrId: string) => {
    return apiFetch(`/api/user/info/${usernameOrId}`, {
        method: 'GET',
    });
};

export interface NoticeButton {
    text: string;
    type: string;
    value: string;
}

export type NoticeTarget = 'channels' | 'users' | 'all' | 'single' | 'user_ids' | 'channel_ids';

export interface NoticeRequest {
    message: string;
    target: NoticeTarget;
    targetId?: number;
    targetIds?: number[];
    imageUrl: string;
    buttons: NoticeButton[];
}

export const sendAdminNotice = async (initData: string, payload: NoticeRequest) => {
    return apiFetch(`/api/admin/notice`, {
        method: 'POST',
        headers: {
            'Authorization': `tma ${initData}`,
            'x-telegram-init-data': initData,
        },
        body: JSON.stringify(payload),
    });
};

export const fetchServerConfig = async () => {
    return apiFetch(`/api/admin/config`, {
        method: 'GET',
    });
};

export const updateServerConfig = async (payload: {
    maintence: boolean;
    forceJoin: boolean;
    globalDefaultCaption: string;
    globalNewPackCaption: string;
    fixedPostBuilderEnabled: boolean;
    fixedPostBuilderKey: string;
    fixedPostBuilderPayload: string;
}) => {
    return apiFetch(`/api/admin/config`, {
        method: 'PUT',
        body: JSON.stringify(payload),
    });
};

export const disconnectChannel = async (channelId: number) => {
    // Retorna a promessa Response inteira para podermos conferir o status 204
    // Agora usando a nova rota RESTful: DELETE /api/channel/:id
    return fetch(`/api/channel/${channelId}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
    });
};

export const updateUserAdmin = async (userId: number) => {
    return apiFetch(`/api/admin/users/${userId}/admin`, {
        method: 'POST',
    });
};

export const updateUserBlacklist = async (userId: number) => {
    return apiFetch(`/api/admin/users/${userId}/blacklist`, {
        method: 'POST',
    });
};

export const fetchAuditCheckBot = async () => {
    return apiFetch('/api/admin/audit/checkbot');
};

export const bulkDeleteChannels = async (userId: number, channelIds: number[]) => {
    return apiFetch('/api/admin/audit/bulk-delete', {
        method: 'POST',
        body: JSON.stringify({ userId, channelIds }),
    });
};


export const connectStart = async (phone: string) => {
  return apiFetch('/api/connect/start', {
    method: 'POST',
    body: JSON.stringify({ phone }),
  });
};

export const connectVerify = async (code: string) => {
  return apiFetch('/api/connect/verify', {
    method: 'POST',
    body: JSON.stringify({ code }),
  });
};

export const connect2FA = async (password: string) => {
  return apiFetch('/api/connect/2fa', {
    method: 'POST',
    body: JSON.stringify({ password }),
  });
};

export const connectStatus = async (): Promise<ConnectStatus> => {
  const response = await apiFetch('/api/connect/status', {
    method: 'GET',
  });
  return response?.data || { connected: false, userId: 0 };
};

export const connectDisconnect = async () => {
  return apiFetch('/api/connect/disconnect', {
    method: 'POST',
  });
};

export const fetchAdminLogs = async (filters: AdminLogsFilters = {}): Promise<AdminLogsResponse> => {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && String(value).trim() !== '') {
            params.set(key, String(value));
        }
    });
    const response = await apiFetch(`/api/admin/logs?${params.toString()}`, {
        method: 'GET',
    });
    return response?.data || { events: [], total: 0, limit: filters.limit || 50, offset: filters.offset || 0 };
};
