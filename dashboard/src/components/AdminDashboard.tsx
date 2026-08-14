import { useState, useMemo, useTransition, useEffect, Dispatch, SetStateAction, Suspense, lazy } from 'react';
import { AdminDashboardData, User, Channel, AuditResult } from '../types';
import { NoticeButton, NoticeTarget, updateUserAdmin, updateUserBlacklist } from '../api';

const AdminNoticeTab = lazy(() => import('./AdminNoticeTab').then(m => ({ default: m.AdminNoticeTab })));
const AdminConfigTab = lazy(() => import('./AdminConfigTab').then(m => ({ default: m.AdminConfigTab })));
const AdminAuditTab = lazy(() => import('./AdminAuditTab').then(m => ({ default: m.AdminAuditTab })));
const AdminLogsTab = lazy(() => import('./AdminLogsTab').then(m => ({ default: m.AdminLogsTab })));
const AdminMTProtoAccountsTab = lazy(() => import('./AdminMTProtoAccountsTab').then(m => ({ default: m.AdminMTProtoAccountsTab })));
const AdminPremiumFeaturesTab = lazy(() => import('./AdminPremiumFeaturesTab').then(m => ({ default: m.AdminPremiumFeaturesTab })));
const AdminSubscriptionsTab = lazy(() => import('./AdminSubscriptionsTab').then(m => ({ default: m.AdminSubscriptionsTab })));
import { Hash, ArrowLeft, ChevronRight, User as UserIcon, ShieldCheck, UserX, UserCheck, MessageSquare } from 'lucide-react';
import { useToast } from './Toast';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { DataTable, Column } from './admin/DataTable';
import { StatusBadge } from './admin/StatusBadge';
import { OperationsOverview } from './admin/OperationsOverview';
import { AdminPageHeader } from './admin/AdminPageHeader';
import { useAdminCrmControls } from './admin/AdminCrmContext';
import { filterAndSortChannels, filterAndSortUsers } from './admin/crmSelectors';

interface AdminDashboardProps {
  adminData: AdminDashboardData;
  activeTab: 'overview' | 'users' | 'channels' | 'notice' | 'config' | 'audit' | 'logs' | 'accounts' | 'premium-features' | 'subscriptions';
  initialLogsChannelId?: string;
  navigateToChannel: (id: number) => void;
  selectedUserId: number | null;
  onSelectUser: (id: number | null) => void;
  onOpenUserDetail: (id: number) => void;
  onMessageUser: (id: number) => void;
  // Notice tab props
  noticeMessage: string;
  setNoticeMessage: Dispatch<SetStateAction<string>>;
  noticeImageUrl: string;
  setNoticeImageUrl: Dispatch<SetStateAction<string>>;
  noticeTarget: NoticeTarget;
  setNoticeTarget: Dispatch<SetStateAction<NoticeTarget>>;
  noticeTargetId: string;
  setNoticeTargetId: Dispatch<SetStateAction<string>>;
  noticeButtons: NoticeButton[];
  handleAddNoticeButton: () => void;
  updateNoticeButton: (index: number, field: keyof NoticeButton, value: string) => void;
  removeNoticeButton: (index: number) => void;
  handleSendNotice: () => void;
  isSendingNotice: boolean;
  auditResults: AuditResult[] | null;
  setAuditResults: Dispatch<SetStateAction<AuditResult[] | null>>;
  auditLoading: boolean;
  handleRunAudit: () => void;
  toast: (message: string, type: 'success' | 'error' | 'info') => void;
}

// ───── Main Component ─────

export function AdminDashboard({
  adminData,
  activeTab,
  navigateToChannel,
  selectedUserId,
  onSelectUser,
  onOpenUserDetail,
  onMessageUser,
  noticeMessage, setNoticeMessage,
  noticeImageUrl, setNoticeImageUrl,
  noticeTarget, setNoticeTarget,
  noticeTargetId, setNoticeTargetId,
  noticeButtons, handleAddNoticeButton,
  updateNoticeButton, removeNoticeButton,
  handleSendNotice,
  isSendingNotice,
  auditResults, setAuditResults, auditLoading, handleRunAudit,
  initialLogsChannelId
}: AdminDashboardProps) {
  const toast = useToast();
  const { searchQuery, sortBy, filterBy, setFilterBy, navigateToTab } = useAdminCrmControls();

  const [localActiveTab, setLocalActiveTab] = useState(activeTab);
  const [isPending, startTransition] = useTransition();

  const [localUsers, setLocalUsers] = useState<User[]>(adminData.users || []);

  useEffect(() => {
    setLocalUsers(adminData.users || []);
  }, [adminData.users]);

  useEffect(() => {
    startTransition(() => {
      setLocalActiveTab(activeTab);
    });
  }, [activeTab]);

  const usersList = localUsers;
  const channelsList = adminData.channels || [];
  const visibleUsers = useMemo(
    () => filterAndSortUsers(usersList, searchQuery, filterBy, sortBy),
    [filterBy, searchQuery, sortBy, usersList],
  );
  const visibleChannels = useMemo(
    () => filterAndSortChannels(channelsList, searchQuery, sortBy),
    [channelsList, searchQuery, sortBy],
  );

  const adminSelectedUser = useMemo(() =>
    selectedUserId ? usersList.find(u => u.id === selectedUserId) : null,
    [selectedUserId, usersList]);

  const setAdminSelectedUser = (user: User | null) => onSelectUser(user ? user.id : null);

  // ── User Actions ──

  const handleToggleAdmin = async (uid: number) => {
    try {
      const res = await updateUserAdmin(uid);
      if (res.success) {
        const isAdmin = res.data?.isAdmin;
        setLocalUsers(prev => prev.map(u => u.id === uid ? { ...u, is_admin: isAdmin } : u));
        toast(isAdmin ? "Usuário promovido a Admin" : "Privilégios de Admin removidos", "success");
      }
    } catch (err: any) {
      toast(err.message || "Erro ao atualizar status de admin", "error");
    }
  };

  const handleToggleBlacklist = async (uid: number) => {
    try {
      const res = await updateUserBlacklist(uid);
      if (res.success) {
        const isBlacklisted = res.data?.isBlacklisted;
        setLocalUsers(prev => prev.map(u => u.id === uid ? { ...u, is_blacklisted: isBlacklisted } : u));
        toast(isBlacklisted ? "Usuário adicionado à Blacklist" : "Usuário removido da Blacklist", isBlacklisted ? "error" : "success");
      }
    } catch (err: any) {
      toast(err.message || "Erro ao atualizar status de blacklist", "error");
    }
  };

  // ── Overview Tab ──

  const renderOverviewTab = () => (
    <div className="admin-overview-page">
      <OperationsOverview
        users={usersList}
        channels={channelsList}
        onOpenUser={(id) => {
          onOpenUserDetail(id);
          navigateToTab('users');
        }}
        onViewUsers={() => navigateToTab('users')}
        onReviewAlert={(kind) => {
          if (kind === 'blacklisted') setFilterBy('blacklisted');
          if (kind === 'without-channels') setFilterBy('without-channels');
          if (kind === 'new-users') setFilterBy('all');
          navigateToTab('users');
        }}
      />
    </div>
  );

  // ── User Detail ──

  const renderUserDetail = () => {
    if (!adminSelectedUser) return null;
    const name = adminSelectedUser.first_name || 'Sem nome';
    return (
      <div className="space-y-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setAdminSelectedUser(null)}
          className="text-muted-foreground"
        >
          <ArrowLeft size={16} className="mr-1.5" /> Voltar para usuários
        </Button>

        <div className="rounded-xl border border-border p-4">
          <div className="flex items-center gap-3 mb-4">
            <div className="flex items-center justify-center size-12 rounded-xl shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
              <UserIcon size={24} />
            </div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <h2 className="text-base font-bold truncate">{name}</h2>
                {adminSelectedUser.is_admin && <Badge variant="default" className="text-[10px]">Admin</Badge>}
                {adminSelectedUser.is_blacklisted && <Badge variant="destructive" className="text-[10px]">Bloqueado</Badge>}
              </div>
              <p className="text-xs text-muted-foreground">ID: {adminSelectedUser.id} • {adminSelectedUser.channels?.length || 0} canais</p>
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <Button
              variant="default"
              size="sm"
              className="w-full"
              onClick={() => {
                onMessageUser(adminSelectedUser.id);
                navigateToTab('notice');
              }}
            >
              <MessageSquare size={16} />
              Mensagem de Suporte
            </Button>
            <Button
              variant={adminSelectedUser.is_admin ? "secondary" : "default"}
              size="sm"
              className="w-full"
              onClick={() => handleToggleAdmin(adminSelectedUser.id)}
            >
              <ShieldCheck size={16} />
              {adminSelectedUser.is_admin ? "Remover Admin" : "Tornar Admin"}
            </Button>
            <Button
              variant={adminSelectedUser.is_blacklisted ? "secondary" : "destructive"}
              size="sm"
              className="w-full"
              onClick={() => handleToggleBlacklist(adminSelectedUser.id)}
            >
              {adminSelectedUser.is_blacklisted ? <UserCheck size={16} /> : <UserX size={16} />}
              {adminSelectedUser.is_blacklisted ? "Remover Blacklist" : "Add Blacklist"}
            </Button>
          </div>
        </div>

        <div className="space-y-2">
          <h3 className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Canais do Usuário</h3>
          {adminSelectedUser.channels && adminSelectedUser.channels.length > 0 ? (
            adminSelectedUser.channels.map((c: Channel) => (
              <button
                key={c.id}
                className="flex items-center w-full text-left gap-3 rounded-xl border border-border p-3 hover:bg-muted/30 transition-colors"
                onClick={() => navigateToChannel(c.id)}
              >
                <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                  <Hash size={16} />
                </div>
                <div className="min-w-0 flex-1">
                  <h3 className="text-[13px] font-semibold truncate">{c.title}</h3>
                  <p className="text-[11px] text-muted-foreground mt-0.5">ID: {c.id}</p>
                </div>
                <ChevronRight size={16} className="shrink-0 text-muted-foreground/30" />
              </button>
            ))
          ) : (
            <div className="flex flex-col items-center py-6 text-muted-foreground rounded-xl border border-border">
              <Hash size={28} className="opacity-30 mb-2" />
              <p className="text-[13px] font-medium">Este usuário não possui canais</p>
            </div>
          )}
        </div>
      </div>
    );
  };

  // ── Users Tab ──

  const renderUsersTab = () => {
    const userColumns: Column<any>[] = [
      { key: 'first_name', label: 'Nome', render: (_: any, row: any) => (
        <div className="flex items-center gap-2">
          <div className="flex items-center justify-center size-7 rounded-full shrink-0 text-[11px] font-bold" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
            {(row.first_name || '?')[0].toUpperCase()}
          </div>
          <div className="min-w-0">
            <span className="text-sm font-semibold truncate block">{row.first_name || 'Sem nome'}</span>
            <span className="text-[10px] text-muted-foreground">ID: {row.id}</span>
          </div>
        </div>
      )},
      { key: 'channels', label: 'Canais', align: 'center', render: (v: any) => (
        <Badge variant="secondary" className="text-[10px] font-mono">{v ? v.length : 0}</Badge>
      )},
      { key: 'is_admin', label: 'Admin', align: 'center', render: (v: boolean) => (
        v ? <StatusBadge label="Admin" variant="accent" dot /> : <span className="text-[11px] text-muted-foreground">—</span>
      )},
      { key: 'is_blacklisted', label: 'Bloqueado', align: 'center', render: (v: boolean) => (
        v ? <StatusBadge label="Bloqueado" variant="danger" dot /> : <span className="text-[11px] text-muted-foreground">—</span>
      )},
    ];

    return (
      <div className="space-y-4">
        <DataTable
          columns={userColumns}
          data={visibleUsers}
          searchable={false}
          pageSize={15}
          emptyMessage="Nenhum usuário encontrado"
          actions={(row: any) => (
            <div className="flex items-center gap-1">
              <button
                className="p-1.5 rounded-lg hover:bg-muted/50 transition-colors"
                onClick={() => setAdminSelectedUser(row)}
                title="Ver detalhes"
              >
                <ChevronRight size={16} className="text-muted-foreground/40" />
              </button>
            </div>
          )}
        />
      </div>
    );
  };

  // ── Channels Tab ──

  const renderChannelsTab = () => {
    const channelColumns: Column<any>[] = [
      { key: 'title', label: 'Canal', render: (_: any, row: any) => (
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center size-8 rounded-lg shrink-0" style={{ background: 'var(--success-soft)', color: 'var(--success)' }}>
            <Hash size={16} />
          </div>
          <div className="min-w-0">
            <span className="text-sm font-semibold truncate block">{row.title}</span>
            <span className="text-[10px] text-muted-foreground">ID: {row.id}</span>
          </div>
        </div>
      )},
      { key: 'ownerId', label: 'Dono', align: 'center' },
      { key: 'subscriberCount', label: 'Inscritos', align: 'center', render: (v: any) => (
        v ? <Badge variant="secondary" className="text-[10px]">{v}</Badge> : <span className="text-[11px] text-muted-foreground">—</span>
      )},
    ];

    return (
      <DataTable
        columns={channelColumns}
        data={visibleChannels}
        searchable={false}
        pageSize={15}
        emptyMessage="Nenhum canal encontrado"
        actions={(row: any) => (
          <button
            className="p-1.5 rounded-lg hover:bg-muted/50 transition-colors"
            onClick={() => navigateToChannel(row.id)}
            title="Abrir canal"
          >
            <ChevronRight size={16} className="text-muted-foreground/40" />
          </button>
        )}
      />
    );
  };

  // ── Notice Tab ──

  const renderNoticeTab = () => {
    return (
      <div className="space-y-4">
        <AdminNoticeTab
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
          users={usersList}
          channels={channelsList}
        />
      </div>
    );
  };

  // ── Render ──

  const tabCopy: Record<typeof localActiveTab, { title: string; description: string }> = {
    overview: { title: 'Visão geral', description: 'Acompanhamento operacional da base' },
    users: { title: 'Usuários', description: 'Gerencie usuários, acessos e canais vinculados' },
    channels: { title: 'Canais', description: 'Todos os canais conectados ao FreddyBot' },
    notice: { title: 'Broadcast', description: 'Envie comunicações para usuários e canais' },
    audit: { title: 'Auditoria', description: 'Verifique a presença e o estado do bot nos canais' },
    logs: { title: 'Logs', description: 'Investigue o histórico operacional do sistema' },
    config: { title: 'Configurações', description: 'Defina o comportamento global do FreddyBot' },
    accounts: { title: 'Contas MTProto', description: 'Gerencie contas usadas na edição de postagens' },
    'premium-features': { title: 'Features premium', description: 'Controle recursos e preços premium' },
    subscriptions: { title: 'Assinaturas', description: 'Gerencie assinaturas e pagamentos dos usuários' },
  };

  return (
    <div className={`admin-crm-page ${localActiveTab === 'overview' ? 'is-overview' : ''} ${isPending ? 'is-pending' : ''}`}>
      {localActiveTab !== 'overview' && (
        <AdminPageHeader
          eyebrow="Painel administrativo"
          title={tabCopy[localActiveTab].title}
          description={tabCopy[localActiveTab].description}
        />
      )}

      {localActiveTab === 'overview' && renderOverviewTab()}
      {localActiveTab === 'users' && !adminSelectedUser && renderUsersTab()}
      {localActiveTab === 'users' && adminSelectedUser && renderUserDetail()}
      {localActiveTab === 'channels' && renderChannelsTab()}
      {localActiveTab === 'audit' && (
        <div className="space-y-4">
          <Suspense fallback={<div className="p-8 text-center text-xs text-muted-foreground animate-pulse">Carregando auditoria...</div>}>
            <AdminAuditTab
              navigateToChannel={navigateToChannel}
              onOpenUser={(id) => {
                onOpenUserDetail(id);
                navigateToTab('users');
              }}
              results={auditResults}
              setResults={setAuditResults}
              loading={auditLoading}
              onRunAudit={handleRunAudit}
            />
          </Suspense>
        </div>
      )}
      {localActiveTab === 'notice' && (
        <Suspense fallback={<div className="p-8 text-center text-xs text-muted-foreground animate-pulse">Carregando aviso...</div>}>
          {renderNoticeTab()}
        </Suspense>
      )}
      {localActiveTab === 'logs' && (
        <div className="space-y-4">
          <Suspense fallback={<div className="p-8 text-center text-xs text-muted-foreground animate-pulse">Carregando logs...</div>}>
            <AdminLogsTab navigateToChannel={navigateToChannel} initialChannelId={initialLogsChannelId} />
          </Suspense>
        </div>
      )}
      {localActiveTab === 'config' && (
        <div className="space-y-4">
          <Suspense fallback={<div className="p-8 text-center text-xs text-muted-foreground animate-pulse">Carregando configurações...</div>}>
            <AdminConfigTab />
          </Suspense>
        </div>
      )}
      {localActiveTab === 'accounts' && (
        <div className="space-y-4">
          <Suspense fallback={<div className="p-8 text-center text-xs text-muted-foreground animate-pulse">Carregando contas MTProto...</div>}>
            <AdminMTProtoAccountsTab />
          </Suspense>
        </div>
      )}
      {localActiveTab === 'premium-features' && (
        <div className="space-y-4">
          <Suspense fallback={<div className="p-8 text-center text-xs text-muted-foreground animate-pulse">Carregando recursos premium...</div>}>
            <AdminPremiumFeaturesTab toast={toast} />
          </Suspense>
        </div>
      )}
      {localActiveTab === 'subscriptions' && (
        <div className="space-y-4">
          <Suspense fallback={<div className="p-8 text-center text-xs text-muted-foreground animate-pulse">Carregando assinaturas...</div>}>
            <AdminSubscriptionsTab toast={toast} />
          </Suspense>
        </div>
      )}
    </div>
  );
}
