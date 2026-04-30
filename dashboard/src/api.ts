import { DashboardData, Button, Permission, ChannelsResponse, AdminDashboardData } from './types';

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
    const data = await apiFetch(`/api/admin/overview`);
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

export const updateNewPackCaption = async (channelId: number, newPackCaption: string) => {
    return apiFetch(`/api/channel/${channelId}/newpackcaption`, {
        method: 'PUT',
        body: JSON.stringify({ newPackCaption }),
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

export interface NoticeRequest {
    message: string;
    target: 'channels' | 'users' | 'all';
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

export const updateServerConfig = async (maintence: boolean, forceJoin: boolean) => {
    return apiFetch(`/api/admin/config`, {
        method: 'PUT',
        body: JSON.stringify({ maintence, forceJoin }),
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
