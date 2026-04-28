import axios from 'axios';
import { User, Channel } from './types';

const tg = window.Telegram?.WebApp;

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
});

// Add a request interceptor to include the Telegram initData
api.interceptors.request.use((config) => {
  const initData = window.Telegram?.WebApp?.initData;
  if (initData) {
    config.headers['x-telegram-init-data'] = initData;
  }
  return config;
});

// Add a response interceptor to handle authentication errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && (error.response.status === 401 || error.response.status === 403)) {
      // We can handle global auth errors here if needed
      console.error('Authentication error:', error.response.status);
    }
    return Promise.reject(error);
  }
);

export const login = async (channelId: number, user: any) => {
  const { data } = await api.post('/auth', { channelID: channelId, user });
  return data;
};

export const fetchMyChannels = async (user: any) => {
  const { data } = await api.post('/me/channels', { user });
  return data;
};

export const fetchChannel = async (channelId: number) => {
  const { data } = await api.get(`/channel/${channelId}`);
  return data;
};

export const updateDefaultCaption = async (channelId: number, caption: string) => {
  const { data } = await api.put(`/channel/${channelId}/caption`, { caption });
  return data;
};

export const updateNewPackCaption = async (channelId: number, caption: string) => {
  const { data } = await api.put(`/channel/${channelId}/newpackcaption`, { caption });
  return data;
};

// --- Buttons API ---

export const createButton = async (channelId: number, nameButton: string, buttonUrl: string) => {
  const { data } = await api.post(`/channel/${channelId}/buttons`, { nameButton, buttonUrl });
  return data;
};

export const updateButton = async (channelId: number, buttonId: string, nameButton: string, buttonUrl: string) => {
  const { data } = await api.put(`/channel/${channelId}/buttons/${buttonId}`, { nameButton, buttonUrl });
  return data;
};

export const deleteButton = async (channelId: number, buttonId: string) => {
  const { data } = await api.delete(`/channel/${channelId}/buttons/${buttonId}`);
  return data;
};

export const updateButtonsLayout = async (channelId: number, layout: any) => {
  const { data } = await api.put(`/channel/${channelId}/buttons/layout`, layout);
  return data;
};

// --- Permissions API ---

export const updateMessagePermissions = async (channelId: number, permissions: any) => {
  const { data } = await api.put(`/channel/${channelId}/caption/permissions`, permissions);
  return data;
};

export const updateButtonsPermissions = async (channelId: number, permissions: any) => {
  const { data } = await api.put(`/channel/${channelId}/buttons/permissions`, permissions);
  return data;
};

export const adminFetchUsers = async () => {
  // Use absolute path because it's under /admin/api instead of /api
  const { data } = await api.get('/admin/api/users', { baseURL: '/' });
  return data;
};

export default api;
