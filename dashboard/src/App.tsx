import { useState, useEffect, useCallback, memo } from 'react';
import { DashboardData, Button, TelegramUser, AdminDashboardData, Channel, AuditResult } from './types';
import {
  login, fetchDashboardData, fetchUserChannels, fetchAdminDashboard,
  updateMessagePermission, updateButtonsPermission,
  createButton, deleteButton, updateButton, updateLayoutButtons,
  updateDefaultCaption, updateNewPackCaption, updateReactions, 
  updateReactionPosition, updateDynamicLinks, transferChannel, fetchUserInfo,
  sendAdminNotice, NoticeButton, NoticeRequest, NoticeTarget, disconnectChannel, fetchAuditCheckBot
} from './api';
import { PermissionsCard } from './components/PermissionsCard';
import { ButtonGrid } from './components/ButtonGrid';
import { CaptionCard } from './components/CaptionCard';
import { NewPackCaptionCard } from './components/NewPackCaptionCard';
import { ReactionsCard } from './components/ReactionsCard';
import { AdminDashboard } from './components/AdminDashboard';
import { DashboardInicioTab } from './components/DashboardInicioTab';
import { TabBar, Tab } from './components/TabBar';
import { AdminSidebar } from './components/AdminSidebar';
import { ToastProvider, useToast } from './components/Toast';
import { useTheme } from './hooks/useTheme';
import {
  Users, Hash, Sun, Moon, Send, ExternalLink, MousePointerClick, Link2,
  LayoutDashboard, Type, Grid3X3, Shield, MessageCircle,
  AlertTriangle, ChevronRight, MessageSquare, Menu, ArrowLeft, Zap, Settings, FileClock
} from 'lucide-react';

const tabs: Tab[] = [
  { id: 'geral', label: 'Início', icon: <LayoutDashboard size={22} /> },
  { id: 'legendas', label: 'Legendas', icon: <Type size={22} /> },
  { id: 'botoes', label: 'Botões', icon: <Grid3X3 size={22} /> },
  { id: 'permissoes', label: 'Permissões', icon: <Shield size={22} /> },
];

const adminTabs: Tab[] = [
  { id: 'users', label: 'Usuários', icon: <Users size={22} /> },
  { id: 'channels', label: 'Canais', icon: <Hash size={22} /> },
  { id: 'audit', label: 'Auditoria', icon: <Zap size={22} /> },
  { id: 'logs', label: 'Logs', icon: <FileClock size={22} /> },
  { id: 'notice', label: 'Broadcast', icon: <MessageSquare size={22} /> },
  { id: 'config', label: 'Configurações', icon: <Settings size={22} /> },
];

const permLabels: Record<string, string> = {
  message: 'Mensagem', audio: 'Áudio', video: 'Vídeo',
  photo: 'Foto', document: 'Arquivo', sticker: 'Sticker', gif: 'GIF', linkPreview: 'Link Preview',
  reactions: 'Reações',
};

function getChannelIdFromUrl(): string | null {
  const match = window.location.pathname.match(/\/dashboard\/(-?\d+)/);
  return match ? match[1] : null;
}

function isChannelsRoute(): boolean {
  return window.location.pathname.startsWith('/me/channels');
}

function isRootRoute(): boolean {
  return window.location.pathname === '/' || window.location.pathname === '';
}

function isAdminDashRoute(): boolean {
  return window.location.pathname.startsWith('/admin/dash');
}

type AdminTabId = 'users' | 'channels' | 'notice' | 'config' | 'audit' | 'logs';

function getInitialAdminTabFromUrl(): AdminTabId {
  const tab = new URLSearchParams(window.location.search).get('tab');
  return tab === 'logs' ? 'logs' : 'users';
}

function getInitialLogsChannelIdFromUrl(): string {
  return new URLSearchParams(window.location.search).get('channelId') || '';
}

type AuthState = 'idle' | 'authenticating' | 'authenticated' | 'error';

const MemoizedAdminDashboard = memo(AdminDashboard);

const DashboardContent = memo(function DashboardContent() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('geral');
  const [adminActiveTab, setAdminActiveTab] = useState<AdminTabId>(() => getInitialAdminTabFromUrl());
  const [adminSelectedUserId, setAdminSelectedUserId] = useState<number | null>(null);
  const [tgUser, setTgUser] = useState<TelegramUser | null>(null);
  const [authState, setAuthState] = useState<AuthState>('idle');
  const [authError, setAuthError] = useState<string>('');
  const toast = useToast();
  const { theme, toggleTheme } = useTheme();

  const [adminData, setAdminData] = useState<AdminDashboardData | null>(null);
  const [noticeMessage, setNoticeMessage] = useState('');
  const [noticeImageUrl, setNoticeImageUrl] = useState('');
  const [noticeTarget, setNoticeTarget] = useState<NoticeTarget>('all');
  const [noticeTargetId, setNoticeTargetId] = useState<string>('');
  const [noticeButtons, setNoticeButtons] = useState<NoticeButton[]>([]);
  const [isSendingNotice, setIsSendingNotice] = useState(false);
  const [auditResults, setAuditResults] = useState<AuditResult[] | null>(null);
  const [auditLoading, setAuditLoading] = useState(false);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [initialLogsChannelId] = useState(() => getInitialLogsChannelIdFromUrl());

  useEffect(() => {
    const savedUid = sessionStorage.getItem('lastAdminUserId');
    if (savedUid) {
      setAdminSelectedUserId(parseInt(savedUid, 10));
    }
  }, []);

  const channelId = getChannelIdFromUrl();
  const isAdmin = isAdminDashRoute();
  const isRoot = isRootRoute();
  const isChannels = isChannelsRoute() || isRoot;
  const isSpecificChannel = !isAdmin && !isChannels && !!channelId;

  const handleBack = useCallback(() => {
    const source = sessionStorage.getItem('navSource');
    if (source === 'admin') {
      window.location.href = '/admin/dash';
    } else {
      window.location.href = '/me/channels';
    }
  }, []);

  useEffect(() => {
    const tg = window.Telegram?.WebApp;
    if (tg) {
      tg.ready();
      tg.expand();
      if (tg.initDataUnsafe?.user) setTgUser(tg.initDataUnsafe.user);

      if (isSpecificChannel) {
        tg.BackButton.show();
        tg.BackButton.onClick(handleBack);
      } else {
        tg.BackButton.hide();
      }
    }
    return () => {
      window.Telegram?.WebApp.BackButton.offClick(handleBack);
    };
  }, [isSpecificChannel, handleBack]);

  const handleBlacklist = useCallback(() => {
    const tg = window.Telegram?.WebApp;
    if (tg) {
      tg.showConfirm("🚫 Você está na blacklist e seu acesso ao painel foi bloqueado. Em caso de dúvidas, acione a /ouvidoria no bot.", () => {
        tg.close();
      });
    } else {
      alert("🚫 Você está na blacklist e seu acesso ao painel foi bloqueado.");
    }
  }, []);

  useEffect(() => {
    (async () => {
      setLoading(true);
      setAuthState('authenticating');

      const tg = window.Telegram?.WebApp;
      const initData = tg?.initData || '';
      const userID = tg?.initDataUnsafe?.user?.id || 0;

      try {
        if (isRootRoute()) {
          setAuthState('authenticated');
          setData({
            channel: null as any,
            user: {
              id: 0,
              first_name: 'Convidado',
              username: '',
              is_admin: false,
              is_blacklisted: false,
              isContribute: false,
              created_at: '',
              updated_at: '',
              channels: []
            }
          });
          return;
        }

        const authRes = await login(initData, userID);
        if (!authRes.success) throw new Error(authRes.message || 'Falha no login');

        if (authRes.isBlacklisted) {
          handleBlacklist();
          return;
        }

        setAuthState('authenticated');

        if (isAdminDashRoute()) {
          const response = await fetchAdminDashboard();
          setAdminData(response);
        } else if (isChannelsRoute()) {
          const response = await fetchUserChannels();
          const channelsData = Array.isArray(response?.data) ? response.data : (response?.data?.channels || response?.channels || []);
          setData({
            channel: null as any,
            user: {
              id: userID,
              first_name: tg?.initDataUnsafe?.user?.first_name || 'Usuário',
              is_admin: false,
              is_blacklisted: false,
              isContribute: false,
              created_at: '',
              updated_at: '',
              channels: channelsData,
              username: tg?.initDataUnsafe?.user?.username || ''
            }
          });
        } else if (channelId) {
          const response = await fetchDashboardData(channelId);
          const dashRes = response?.data || response;
          
          if (dashRes.user?.is_blacklisted) {
            handleBlacklist();
            return;
          }

          // Resolver conflito de posição entre botões e reações (se houver)
          if (dashRes.channel && dashRes.channel.buttons) {
            const { buttons, reactionPosition } = dashRes.channel;
            const hasConflict = buttons.some(b => b.positionY === reactionPosition);
            if (hasConflict) {
              const maxBtnY = buttons.reduce((max, b) => Math.max(max, b.positionY), -1);
              dashRes.channel.reactionPosition = maxBtnY + 1;
            }
          }

          setData(dashRes);
        }

        if (tg?.CloudStorage && initData) {
          tg.CloudStorage.setItem('initData', initData);
        }

      } catch (err: any) {
        console.warn('Auth/fetch failed, checking error:', err);
        const errMsg = err?.message || 'Erro na autenticação';

        // @ts-ignore
        if (import.meta.env.DEV) {
          const { fallbackData, mockAdminData } = await import('./mockData');
          setAuthState('authenticated');
          if (isAdminDashRoute()) setAdminData(mockAdminData);
          else if (isChannelsRoute()) setData({ channel: null as any, user: { ...fallbackData.user, channels: mockAdminData.channels } });
          else setData(fallbackData);
        } else {
          setAuthState('error');
          setAuthError(errMsg);
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [channelId, handleBlacklist]);

  const handleMsgPerm = useCallback(async (field: string, value: boolean) => {
    if (!data?.channel?.defaultCaption) return;
    const cid = data.channel.id;
    
    setData(p => {
      if (!p?.channel?.defaultCaption) return p;
      return {
        ...p, channel: {
          ...p.channel, defaultCaption: {
            ...p.channel.defaultCaption,
            messagePermission: { ...(p.channel.defaultCaption.messagePermission || {}), [field]: value }
          }
        }
      };
    });

    try {
      const currentPerms = data.channel.defaultCaption.messagePermission || {};
      const newPerms = { ...currentPerms, [field]: value };
      await updateMessagePermission(cid, newPerms);
      toast(`${permLabels[field] || field} ${value ? 'ativado' : 'desativado'}`, value ? 'success' : 'info');
    } catch {
      setData(data); 
      toast(`Erro ao atualizar permissão`, 'error');
    }
  }, [toast, data]);

  const handleBtnPerm = useCallback(async (field: string, value: boolean) => {
    if (!data?.channel?.defaultCaption) return;
    const cid = data.channel.id;

    setData(p => {
      if (!p?.channel?.defaultCaption) return p;
      return {
        ...p, channel: {
          ...p.channel, defaultCaption: {
            ...p.channel.defaultCaption,
            buttonsPermission: { ...(p.channel.defaultCaption.buttonsPermission || {}), [field]: value }
          }
        }
      };
    });

    try {
      const currentPerms = data.channel.defaultCaption.buttonsPermission || {};
      const newPerms = { ...currentPerms, [field]: value };
      await updateButtonsPermission(cid, newPerms);
      toast(`${permLabels[field] || field} ${value ? 'ativado' : 'desativado'}`, value ? 'success' : 'info');
    } catch {
      setData(data);
      toast(`Erro ao atualizar permissão`, 'error');
    }
  }, [toast, data]);

  const handleDynamicLinks = useCallback(async (field: string, value: boolean) => {
    if (!data) return;
    const cid = parseInt(String(channelId), 10);

    const newSettings = {
      dynamicLinks: field === 'dynamicLinks' ? value : data.channel.dynamicLinks,
      dlBotButtons: field === 'dlBotButtons' ? value : data.channel.dlBotButtons,
      dlBotCaptions: field === 'dlBotCaptions' ? value : data.channel.dlBotCaptions,
      dlBotReactions: field === 'dlBotReactions' ? value : data.channel.dlBotReactions,
    };

    setData(p => {
      if (!p) return p;
      return { ...p, channel: { ...p.channel, ...newSettings } };
    });

    try {
      await updateDynamicLinks(cid, newSettings);
      const labels: Record<string, string> = {
        dynamicLinks: 'Links Dinâmicos',
        dlBotButtons: 'Botões do Bot',
        dlBotCaptions: 'Legendas do Bot',
        dlBotReactions: 'Reações do Bot'
      };
      toast(`${labels[field] || field} ${value ? 'ativados' : 'desativados'}`, value ? 'success' : 'info');
    } catch {
      setData(data);
      toast(`Erro ao atualizar configuração`, 'error');
    }
  }, [toast, channelId, data]);

  const handleAddButton = useCallback(async (button: Button) => {
    const cid = parseInt(String(channelId), 10);
    try {
      const resp = await createButton(cid, button);
      const newButtonData = resp?.data || resp?.Data || resp;
      const realId = newButtonData?.buttonId || newButtonData?.ButtonID || newButtonData?.id;

      if (!realId) throw new Error("ID not returned from API");

      const finalButton = { ...button, buttonId: realId };

      setData(p => {
        if (!p) return p;
        return { ...p, channel: { ...p.channel, buttons: [...p.channel.buttons, finalButton] } };
      });
      toast(`"${button.nameButton}" adicionado`, 'success');

      setData(p => {
        if (!p) return p;
        const allButtons = p.channel.buttons;
        const layout: any[][] = [];
        const maxRow = allButtons.reduce((max, b) => Math.max(max, b.positionY), 0);
        for (let currentY = 0; currentY <= maxRow; currentY++) {
          const rowButtons = allButtons
            .filter(b => b.positionY === currentY)
            .sort((a, b) => a.positionX - b.positionX)
            .map(b => ({
              buttonId: b.buttonId,
              nameButton: b.nameButton,
              buttonUrl: b.buttonUrl,
              positionX: b.positionX,
              positionY: b.positionY,
            }));
          layout.push(rowButtons);
        }
        updateLayoutButtons(cid, layout).catch(console.error);
        return p;
      });
    } catch (err) {
      console.error(err);
      toast(`Erro ao adicionar botão`, 'error');
    }
  }, [toast, channelId]);

  const handleDeleteButton = useCallback(async (buttonId: string) => {
    const cid = parseInt(String(channelId), 10);
    try {
      await deleteButton(cid, buttonId);
      setData(p => {
        if (!p) return p;
        const btn = p.channel.buttons.find(b => b.buttonId === buttonId);
        if (btn) toast(`"${btn.nameButton}" excluído`, 'error');
        return { ...p, channel: { ...p.channel, buttons: p.channel.buttons.filter(b => b.buttonId !== buttonId) } };
      });
    } catch {
      toast(`Erro ao excluir botão`, 'error');
    }
  }, [toast, channelId]);

  const handleEditButton = useCallback(async (buttonId: string, updates: Partial<Button>) => {
    const cid = parseInt(String(channelId), 10);
    try {
      await updateButton(cid, buttonId, updates);
      setData(p => {
        if (!p) return p;
        return {
          ...p, channel: {
            ...p.channel, buttons: p.channel.buttons.map(b =>
              b.buttonId === buttonId ? { ...b, ...updates, updated_at: new Date().toISOString() } : b
            )
          }
        };
      });
      toast('Botão atualizado', 'success');
    } catch {
      toast('Erro ao atualizar botão', 'error');
    }
  }, [toast, channelId]);

  const handleMoveButton = useCallback(async (buttonId: string, x: number, y: number) => {
    const cid = parseInt(String(channelId), 10);
    if (!data) return;

    const updatedButtons = data.channel.buttons.map(b =>
      b.buttonId === buttonId ? { ...b, positionX: x, positionY: y, updated_at: new Date().toISOString() } : b
    );

    const currentReactionPos = data.channel.reactionPosition;
    let desiredReactionPos = currentReactionPos;
    
    const hasConflict = updatedButtons.some(b => b.positionY === currentReactionPos);
    if (hasConflict) {
      const maxBtnY = updatedButtons.reduce((max, b) => Math.max(max, b.positionY), -1);
      desiredReactionPos = maxBtnY + 1;
    }

    try {
      if (desiredReactionPos !== currentReactionPos) {
        await updateReactionPosition(cid, 99);
      }

      const layout: any[][] = [];
      const maxRow = Math.max(...updatedButtons.map(b => b.positionY));
      for (let currentY = 0; currentY <= maxRow; currentY++) {
        const rowButtons = updatedButtons
          .filter(b => b.positionY === currentY)
          .sort((a, b) => a.positionX - b.positionX)
          .map(b => ({
            buttonId: b.buttonId,
            nameButton: b.nameButton,
            buttonUrl: b.buttonUrl,
            positionX: b.positionX,
            positionY: b.positionY,
          }));
        layout.push(rowButtons);
      }

      await updateLayoutButtons(cid, layout);
      
      if (desiredReactionPos !== currentReactionPos) {
        await updateReactionPosition(cid, desiredReactionPos);
      }

      setData(p => {
        if (!p) return p;
        return {
          ...p, channel: {
            ...p.channel, 
            buttons: updatedButtons,
            reactionPosition: desiredReactionPos
          }
        };
      });
      toast('Botão movido com sucesso', 'info');
    } catch (err: any) {
      toast(err.message || 'Erro ao mover botão', 'error');
      fetchDashboardData(String(channelId)).then(setData);
    }
  }, [toast, channelId, data]);

  const handleMoveReactions = useCallback(async (y: number) => {
    const cid = parseInt(String(channelId), 10);
    if (!data) return;

    const conflictingButtons = data.channel.buttons.filter(b => b.positionY === y);
    if (conflictingButtons.length > 0) {
      const names = conflictingButtons.map(b => `"${b.nameButton}"`).join(', ');
      toast(`Não é possível mover: a linha ${y + 1} possui botões (${names})`, 'error');
      return;
    }

    try {
      await updateReactionPosition(cid, y);
      setData(p => {
        if (!p) return p;
        return { ...p, channel: { ...p.channel, reactionPosition: y } };
      });
      toast('Posição das reações atualizada', 'info');
    } catch (err: any) {
      toast(err.message || 'Erro ao mover reações', 'error');
    }
  }, [toast, channelId, data]);

  const handleUpdateCaption = useCallback(async (text: string) => {
    const cid = parseInt(String(channelId), 10);
    try {
      await updateDefaultCaption(cid, text);
      setData(p => {
        if (!p) return p;
        return { ...p, channel: { ...p.channel, defaultCaption: { ...p.channel.defaultCaption, caption: text } } };
      });
      toast('Caption atualizada', 'success');
    } catch {
      toast('Erro ao atualizar caption', 'error');
    }
  }, [toast, channelId]);

  const handleUpdateNewPack = useCallback(async (settings: {
    caption: string;
    messageButtons: boolean;
    stickerButtons: boolean;
    messagePosition: 'above' | 'below';
    replyToSticker: boolean;
  }) => {
    const cid = parseInt(String(channelId), 10);
    try {
      await updateNewPackCaption(cid, {
        newPackCaption: settings.caption,
        newPackMessageButtons: settings.messageButtons,
        newPackStickerButtons: settings.stickerButtons,
        newPackMessagePosition: settings.messagePosition,
        newPackReplyToSticker: settings.replyToSticker,
      });
      setData(p => {
        if (!p) return p;
        return {
          ...p,
          channel: {
            ...p.channel,
            newPackCaption: settings.caption,
            newPackMessageButtons: settings.messageButtons,
            newPackStickerButtons: settings.stickerButtons,
            newPackMessagePosition: settings.messagePosition,
            newPackReplyToSticker: settings.replyToSticker,
          }
        };
      });
      toast('New Pack atualizada', 'success');
    } catch {
      toast('Erro ao atualizar New Pack', 'error');
    }
  }, [toast, channelId]);

  const handleUpdateReactions = useCallback(async (text: string) => {
    const cid = parseInt(String(channelId), 10);
    if (!data) return;

    try {
      await updateReactions(cid, text);
      
      let newPos = data.channel.reactionPosition;
      const hasConflict = data.channel.buttons.some(b => b.positionY === newPos);
      
      if (text.trim() !== '' && (hasConflict || (newPos === 0 && data.channel.buttons.some(b => b.positionY === 0)))) {
        const maxBtnY = data.channel.buttons.reduce((max, b) => Math.max(max, b.positionY), -1);
        newPos = maxBtnY + 1;
        await updateReactionPosition(cid, newPos);
      }

      setData(p => {
        if (!p) return p;
        return { ...p, channel: { ...p.channel, reactions: text, reactionPosition: newPos } };
      });
      toast('Reações atualizadas', 'success');
    } catch (err: any) {
      toast(err.message || 'Erro ao atualizar reações', 'error');
    }
  }, [toast, channelId, data]);

  const getGreeting = useCallback(() => {
    const h = new Date().getHours();
    if (h < 12) return 'Bom dia';
    if (h < 18) return 'Boa tarde';
    return 'Boa noite';
  }, []);

  const getGreetingEmoji = useCallback(() => {
    const h = new Date().getHours();
    if (h < 6) return '🌙';
    if (h < 12) return '☀️';
    if (h < 18) return '🌤️';
    return '🌙';
  }, []);

  const [showDisconnect, setShowDisconnect] = useState(false);
  const [showDisconnectSuccess, setShowDisconnectSuccess] = useState(false);

  const handleDisconnect = useCallback(() => {
    setShowDisconnect(true);
  }, []);

  const [isDisconnecting, setIsDisconnecting] = useState(false);

  const confirmDisconnect = useCallback(async () => {
    setIsDisconnecting(true);
    try {
      if (!channelId) throw new Error("ID do canal não encontrado");

      const cid = parseInt(String(channelId), 10);
      const res = await disconnectChannel(cid);

      if (res.status === 204) {
        setShowDisconnectSuccess(true);
      } else {
        const errText = await res.text().catch(() => '');
        throw new Error(errText || `Erro na API (${res.status})`);
      }
    } catch (err: any) {
      toast(err.message || 'Erro ao desconectar o bot', 'error');
    } finally {
      setIsDisconnecting(false);
      setShowDisconnect(false);
    }
  }, [channelId, toast]);

  const parseNoticeTargetIds = useCallback((raw: string) => {
    return raw
      .split(/[\s,;]+/)
      .map(v => v.trim())
      .filter(Boolean)
      .map(v => Number.parseInt(v, 10))
      .filter(v => Number.isFinite(v) && v !== 0);
  }, []);

  const handleSendNotice = useCallback(async () => {
    if (!noticeMessage.trim()) {
      toast('A mensagem não pode estar vazia', 'error');
      return;
    }

    const specificTarget = noticeTarget === 'single' || noticeTarget === 'user_ids' || noticeTarget === 'channel_ids';
    const targetIds = specificTarget ? parseNoticeTargetIds(noticeTargetId) : [];
    if (specificTarget && targetIds.length === 0) {
      toast('Informe pelo menos um ID válido', 'error');
      return;
    }

    setIsSendingNotice(true);
    try {
      const tg = window.Telegram?.WebApp;
      const initData = tg?.initData || '';

      const payload: NoticeRequest = {
        message: noticeMessage,
        imageUrl: noticeImageUrl,
        target: noticeTarget,
        targetId: noticeTarget === 'single' ? targetIds[0] : undefined,
        targetIds: noticeTarget === 'user_ids' || noticeTarget === 'channel_ids' ? targetIds : undefined,
        buttons: noticeButtons
      };

      await sendAdminNotice(initData, payload);
      toast('Mensagem enviada com sucesso!', 'success');
      setNoticeMessage('');
      setNoticeImageUrl('');
      setNoticeTargetId('');
      setNoticeButtons([]);
    } catch (err: any) {
      toast(err.message || 'Erro ao enviar mensagem', 'error');
    } finally {
      setIsSendingNotice(false);
    }
  }, [noticeMessage, noticeImageUrl, noticeTarget, noticeTargetId, noticeButtons, parseNoticeTargetIds, toast]);

  const handleRunAudit = useCallback(async () => {
    setAuditLoading(true);
    try {
      const res = await fetchAuditCheckBot();
      if (res.success) {
        const auditData = Array.isArray(res.data) ? res.data : [];
        setAuditResults(auditData);
        if (auditData.length === 0) {
          toast("Varredura concluída: nenhum canal com @XavolaBot.", "success");
        } else {
          toast(`Auditoria concluída: ${auditData.length} usuários afetados`, "info");
        }
      } else {
        throw new Error(res.message || 'Erro na auditoria');
      }
    } catch (err: any) {
      toast(err.message || 'Erro ao realizar auditoria', 'error');
    } finally {
      setAuditLoading(false);
    }
  }, [toast]);

  const handleAddNoticeButton = useCallback(() => {
    setNoticeButtons(prev => [...prev, { text: '', type: 'url', value: '' }]);
  }, []);

  const updateNoticeButton = useCallback((index: number, field: keyof NoticeButton, value: string) => {
    setNoticeButtons(prev => {
      const newBtns = [...prev];
      newBtns[index] = { ...newBtns[index], [field]: value };
      return newBtns;
    });
  }, []);

  const removeNoticeButton = useCallback((index: number) => {
    setNoticeButtons(prev => prev.filter((_, i) => i !== index));
  }, []);

  const navigateToChannel = useCallback((id: number) => {
    if (isAdminDashRoute()) {
      sessionStorage.setItem('navSource', 'admin');
    } else {
      sessionStorage.removeItem('navSource');
    }
    window.location.href = `/dashboard/${id}`;
  }, []);

  const onSelectAdminUser = useCallback((id: number | null) => {
    setAdminSelectedUserId(id);
    if (id) sessionStorage.setItem('lastAdminUserId', id.toString());
    else sessionStorage.removeItem('lastAdminUserId');
  }, []);

  const openAdminUserDetail = useCallback((id: number) => {
    setAdminSelectedUserId(id);
    sessionStorage.setItem('lastAdminUserId', id.toString());
    setAdminActiveTab('users');
    setIsSidebarOpen(false);
  }, []);

  const openSupportNoticeForUser = useCallback((id: number) => {
    setAdminSelectedUserId(id);
    sessionStorage.setItem('lastAdminUserId', id.toString());
    setNoticeTarget('single');
    setNoticeTargetId(id.toString());
    setAdminActiveTab('notice');
    setIsSidebarOpen(false);
  }, []);

  useEffect(() => {
    // Tab switch effect handled via CSS entrance animations
  }, [activeTab, adminActiveTab, loading]);

  if (authState === 'error') {
    let displayMessage = authError || 'Não foi possível autenticar. Tente novamente pelo Telegram.';
    
    try {
      if (displayMessage.startsWith('{')) {
        const parsed = JSON.parse(displayMessage);
        displayMessage = parsed.message || parsed.error || displayMessage;
      }
    } catch (e) {}

    return (
      <div className="app-layout">
        <div className="main-content" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh', gap: 16, textAlign: 'center' }}>
          <div style={{ width: 64, height: 64, borderRadius: 20, background: 'var(--danger-soft)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 8 }}>
            <AlertTriangle size={32} style={{ color: 'var(--danger)' }} />
          </div>
          <h2 style={{ fontSize: 20, fontWeight: 800 }}>Ops! Acesso negado</h2>
          <p style={{ fontSize: 15, color: 'var(--hint)', maxWidth: 320, lineHeight: 1.6 }}>{displayMessage}</p>
          <button className="btn btn-primary" onClick={() => window.location.href = '/me/channels'} style={{ marginTop: 12, minWidth: 200 }}>
            <ArrowLeft size={18} /> Voltar para Meus Canais
          </button>
        </div>
      </div>
    );
  }

  if (loading || (!data && !adminData)) {
    return (
      <div className="app-layout">
        <div className="main-content space-y-4" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh', gap: 12 }}>
          <div className="auth-spinner" />
          <p style={{ fontSize: 14, color: 'var(--hint)', fontWeight: 500 }}>
            {authState === 'authenticating' ? 'Autenticando...' : 'Carregando dados...'}
          </p>
        </div>
      </div>
    );
  }

  const channel = data?.channel;
  const user = data?.user;
  const displayName = tgUser?.first_name || user?.firstName || user?.first_name || 'Administrador';
  const initials = displayName[0]?.toUpperCase() || '?';

  return (
    <div className={`app-layout ${isAdmin ? 'admin-layout' : ''}`}>
      {isAdmin && adminData && (
        <>
          {isSidebarOpen && (
            <div 
              className="admin-overlay" 
              onClick={() => setIsSidebarOpen(false)} 
            />
          )}
          <AdminSidebar
            tabs={adminTabs}
            activeTab={adminActiveTab}
            onTabChange={(id) => {
              setAdminActiveTab(id as any);
              setIsSidebarOpen(false);
            }}
            isCollapsed={!isSidebarOpen}
          />
        </>
      )}
      
      <div className={isAdmin ? 'app-main' : 'w-full flex flex-col min-h-screen'}>
        <div className="top-bar">
          {isAdmin && (
            <button 
              className="sidebar-trigger-btn mr-2" 
              onClick={() => setIsSidebarOpen(!isSidebarOpen)}
              title={isSidebarOpen ? "Desativar menu" : "Ativar menu"}
            >
              <Menu size={22} />
            </button>
          )}

          {isSpecificChannel && (
            <button 
              className="sidebar-trigger-btn mr-2" 
              onClick={handleBack}
              title="Voltar"
            >
              <ArrowLeft size={22} />
            </button>
          )}

          <div className="top-avatar">
            {tgUser?.photo_url ? (
              <img src={tgUser.photo_url} alt="" />
            ) : (
              initials
            )}
          </div>
          <div className="min-w-0 flex-1">
            <h1 className="text-[15px] font-bold truncate">{displayName}</h1>
            <p className="text-xs truncate" style={{ color: 'var(--hint)' }}>{isChannels ? 'Meus Canais' : (isAdmin ? 'Painel Admin' : channel?.title)}</p>
          </div>
          <button className="theme-switch" onClick={toggleTheme} title={`Tema atual: ${theme === 'telegram' ? 'Telegram' : theme === 'dark' ? 'Escuro' : 'Claro'}`}>
            {theme === 'telegram' ? <Send size={17} /> : theme === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
          </button>
        </div>

        <div className="main-content">
          {isAdmin && adminData && (
            <div className="tab-content-wrapper">
              <MemoizedAdminDashboard
                adminData={adminData}
                activeTab={adminActiveTab}
                navigateToChannel={navigateToChannel}
                selectedUserId={adminSelectedUserId}
                onSelectUser={onSelectAdminUser}
                onOpenUserDetail={openAdminUserDetail}
                onMessageUser={openSupportNoticeForUser}
                noticeMessage={noticeMessage}
                setNoticeMessage={setNoticeMessage}
                noticeImageUrl={noticeImageUrl}
                setNoticeImageUrl={setNoticeImageUrl}
                noticeTarget={noticeTarget}
                setNoticeTarget={setNoticeTarget}
                noticeTargetId={noticeTargetId}
                setNoticeTargetId={setNoticeTargetId}
                noticeButtons={noticeButtons}
                handleAddNoticeButton={handleAddNoticeButton}
                updateNoticeButton={updateNoticeButton}
                removeNoticeButton={removeNoticeButton}
                handleSendNotice={handleSendNotice}
                isSendingNotice={isSendingNotice}
                auditResults={auditResults}
                setAuditResults={setAuditResults}
                auditLoading={auditLoading}
                handleRunAudit={handleRunAudit}
                initialLogsChannelId={initialLogsChannelId}
              />
            </div>
          )}

          {(!isAdmin && (!channel || isChannelsRoute())) && (
            <div className="space-y-4">
              {isChannels && (
                <div className="card" style={{ padding: '20px' }}>
                  <div className="flex items-center gap-3">
                    <span className="text-2xl">{getGreetingEmoji()}</span>
                    <div className="min-w-0 flex-1">
                      <h2 className="text-[16px] font-bold">{getGreeting()}, <span style={{ color: 'var(--accent)' }}>{displayName}</span></h2>
                      <p className="text-[12px]" style={{ color: 'var(--hint)' }}>Selecione um canal para gerenciar suas configurações.</p>
                    </div>
                  </div>
                </div>
              )}

              <div className="space-y-3">
                {isChannels && <h3 className="text-sm font-semibold mb-2" style={{ color: 'var(--hint)' }}>Canais Encontrados</h3>}

                {user?.channels && user?.channels.length > 0 && isChannels ? (
                  user?.channels.map((c: Channel) => (
                    <button key={c.id} className="card stat-card-clickable" style={{ display: 'flex', alignItems: 'center', width: '100%', textAlign: 'left', padding: '16px' }} onClick={() => navigateToChannel(c.id)}>
                      <div className="section-icon purple mr-3"><Hash size={20} /></div>
                      <div className="min-w-0 flex-1">
                        <h3 className="text-[15px] font-semibold truncate">{c.title}</h3>
                        <p className="text-xs truncate mt-0.5" style={{ color: 'var(--hint)' }}>ID: {c.id}</p>
                      </div>
                      <ChevronRight size={18} className="stat-arrow" />
                    </button>
                  ))
                ) : (!channel ? (
                  <div className="card" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', padding: '40px 24px', color: 'var(--hint)' }}>
                    <div style={{ width: 64, height: 64, borderRadius: '50%', background: 'var(--accent-soft)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 16 }}>
                      <Shield size={32} style={{ color: 'var(--accent)' }} />
                    </div>
                    <h3 className="text-[18px] font-bold" style={{ color: 'var(--text)', marginBottom: 8 }}>Mantenha seu canal organizado</h3>
                    <p className="text-[14px]" style={{ opacity: 0.8, marginBottom: 8, lineHeight: 1.5 }}>
                      O LegendasBOT ajuda a gerenciar botões, legendas automáticas e permissões de forma simples e rápida.
                    </p>
                    <p className="text-[13px]" style={{ opacity: 0.7, marginBottom: 24 }}>
                      Para começar, adicione este bot como <strong style={{ color: 'var(--text)' }}>administrador</strong> no seu canal do Telegram.
                    </p>
                    <div style={{ width: '100%', height: 1, background: 'var(--border)', marginBottom: 24 }}></div>
                    <h4 className="text-[14px] font-semibold" style={{ color: 'var(--text)', marginBottom: 12 }}>Fique por dentro das novidades!</h4>
                    <p className="text-[13px]" style={{ opacity: 0.8, marginBottom: 16 }}>
                      Entre no nosso canal oficial para acompanhar atualizações, dicas e novos recursos.
                    </p>
                    <a
                      href="https://t.me/LegendasBOTTopic"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="btn btn-primary"
                      style={{ width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 8 }}
                    >
                      <ExternalLink size={18} />
                      Entrar no Canal de Atualizações
                    </a>
                  </div>
                ) : null)}
              </div>
            </div>
          )}

          {!isChannels && !isAdmin && activeTab === 'geral' && channel && (
            <div className="tab-content-wrapper">
              <DashboardInicioTab
                channel={channel}
                displayName={displayName}
                getGreeting={getGreeting}
                getGreetingEmoji={getGreetingEmoji}
                handleDisconnect={handleDisconnect}
                showDisconnect={showDisconnect}
                setShowDisconnect={setShowDisconnect}
                isDisconnecting={isDisconnecting}
                confirmDisconnect={confirmDisconnect}
                showDisconnectSuccess={showDisconnectSuccess}
                setShowDisconnectSuccess={setShowDisconnectSuccess}
              />
            </div>
          )}

          {!isChannels && !isAdmin && activeTab === 'legendas' && channel && (
            <div className="space-y-4 tab-content-wrapper">
              <CaptionCard caption={channel.defaultCaption} onUpdate={handleUpdateCaption} />
              <NewPackCaptionCard
                caption={channel.newPackCaption}
                messageButtons={channel.newPackMessageButtons ?? true}
                stickerButtons={channel.newPackStickerButtons ?? true}
                messagePosition={channel.newPackMessagePosition ?? 'above'}
                replyToSticker={channel.newPackReplyToSticker ?? false}
                onUpdate={handleUpdateNewPack}
              />
              <ReactionsCard reactions={channel.reactions} onUpdate={handleUpdateReactions} />
            </div>
          )}

          {!isChannels && !isAdmin && activeTab === 'botoes' && channel && (
            <ButtonGrid
              buttons={channel.buttons}
              reactions={channel.reactions}
              reactionPosition={channel.reactionPosition}
              channelId={channel.id}
              onAdd={handleAddButton}
              onDelete={handleDeleteButton}
              onEdit={handleEditButton}
              onMove={handleMoveButton}
              onMoveReactions={handleMoveReactions}
            />
          )}

          {!isChannels && !isAdmin && activeTab === 'permissoes' && channel && (
            <div className="space-y-4 tab-content-wrapper">
              <div className="card">
                <div className="section-header">
                  <div className="section-icon purple">
                    <Zap size={18} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="text-[15px] font-semibold truncate">Configurações de Reações</h3>
                    <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
                      {channel.defaultCaption?.messagePermission?.reactions ? 'Ativadas' : 'Desativadas'}
                    </p>
                  </div>
                  <span className={`badge ${channel.defaultCaption?.messagePermission?.reactions ? 'badge-accent' : 'badge-ghost'}`}>
                    {channel.defaultCaption?.messagePermission?.reactions ? 'ON' : 'OFF'}
                  </span>
                </div>
                <div className="space-y-2">
                  <div
                    className={`perm-row ${channel.defaultCaption?.messagePermission?.reactions ? 'on' : ''}`}
                    onClick={() => handleMsgPerm('reactions', !channel.defaultCaption?.messagePermission?.reactions)}
                  >
                    <div className="flex items-center gap-3 min-w-0">
                      <span
                        className="flex-shrink-0"
                        style={{
                          color: channel.defaultCaption?.messagePermission?.reactions ? 'var(--accent)' : 'var(--hint)',
                          opacity: channel.defaultCaption?.messagePermission?.reactions ? 1 : 0.4
                        }}
                      >
                        <Zap size={16} />
                      </span>
                      <span className="text-[13px] font-medium">Ativar Reações em Posts</span>
                    </div>
                    <div className={`toggle ${channel.defaultCaption?.messagePermission?.reactions ? 'on' : ''}`} />
                  </div>
                </div>
              </div>

              {/* Links Dinâmicos */}
              <div className="card">
                <div className="section-header">
                  <div className="section-icon purple">
                    <Link2 size={18} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="text-[15px] font-semibold truncate">Links Dinâmicos</h3>
                    <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
                      Transforma links em botões automaticamente
                    </p>
                  </div>
                  <span className={`badge ${channel.dynamicLinks ? 'badge-accent' : 'badge-ghost'}`}>
                    {channel.dynamicLinks ? 'ON' : 'OFF'}
                  </span>
                </div>
                <div className="space-y-2">
                  <div
                    className={`perm-row ${channel.dynamicLinks ? 'on' : ''}`}
                    onClick={() => handleDynamicLinks('dynamicLinks', !channel.dynamicLinks)}
                  >
                    <div className="flex items-center gap-3 min-w-0">
                      <span
                        className="flex-shrink-0"
                        style={{
                          color: channel.dynamicLinks ? 'var(--accent)' : 'var(--hint)',
                          opacity: channel.dynamicLinks ? 1 : 0.4
                        }}
                      >
                        <ExternalLink size={16} />
                      </span>
                      <span className="text-[13px] font-medium">Ativar Links Dinâmicos</span>
                    </div>
                    <div className={`toggle ${channel.dynamicLinks ? 'on' : ''}`} />
                  </div>

                  {/* Sub-toggles (só aparecem se o principal estiver ON) */}
                  {channel.dynamicLinks && (
                    <div className="pl-6 space-y-2 mt-2 border-l-2 border-[var(--border)] ml-4 animate-in fade-in slide-in-from-left-2 duration-300">
                      <div className="text-[10px] font-bold opacity-40 uppercase mb-2 tracking-wider">Regras de Exceção</div>
                      
                      <div
                        className={`perm-row ${channel.dlBotButtons ? 'on' : ''}`}
                        onClick={() => handleDynamicLinks('dlBotButtons', !channel.dlBotButtons)}
                      >
                        <div className="flex items-center gap-3">
                          <MousePointerClick size={14} className={channel.dlBotButtons ? 'text-[var(--accent)]' : 'text-[var(--hint)]'} />
                          <span className="text-[12px]">Manter Botões do Bot</span>
                        </div>
                        <div className={`toggle sm ${channel.dlBotButtons ? 'on' : ''}`} />
                      </div>

                      <div
                        className={`perm-row ${channel.dlBotCaptions ? 'on' : ''}`}
                        onClick={() => handleDynamicLinks('dlBotCaptions', !channel.dlBotCaptions)}
                      >
                        <div className="flex items-center gap-3">
                          <Type size={14} className={channel.dlBotCaptions ? 'text-[var(--accent)]' : 'text-[var(--hint)]'} />
                          <span className="text-[12px]">Manter Legendas do Bot</span>
                        </div>
                        <div className={`toggle sm ${channel.dlBotCaptions ? 'on' : ''}`} />
                      </div>

                      <div
                        className={`perm-row ${channel.dlBotReactions ? 'on' : ''}`}
                        onClick={() => handleDynamicLinks('dlBotReactions', !channel.dlBotReactions)}
                      >
                        <div className="flex items-center gap-3">
                          <Zap size={14} className={channel.dlBotReactions ? 'text-[var(--accent)]' : 'text-[var(--hint)]'} />
                          <span className="text-[12px]">Manter Reações do Bot</span>
                        </div>
                        <div className={`toggle sm ${channel.dlBotReactions ? 'on' : ''}`} />
                      </div>
                      
                      <p className="text-[10px] opacity-40 italic mt-2">
                        * Estas regras só se aplicam se um link dinâmico for detectado na postagem.
                      </p>
                    </div>
                  )}
                </div>
              </div>

              <PermissionsCard
                title="Permissões de Mensagem"
                icon={<MessageCircle size={18} />}
                permission={channel.defaultCaption?.messagePermission}
                onToggle={handleMsgPerm}
              />
              <PermissionsCard
                title="Permissões de Botões"
                icon={<MousePointerClick size={18} />}
                permission={channel.defaultCaption?.buttonsPermission}
                onToggle={handleBtnPerm}
              />
            </div>
          )}
        </div>
      </div>
      
      {!isChannels && !isAdmin && (
        <TabBar tabs={tabs} activeTab={activeTab} onTabChange={setActiveTab} />
      )}
    </div>
  );
});

export default function App() {
  return (
    <ToastProvider>
      <DashboardContent />
    </ToastProvider>
  );
}
