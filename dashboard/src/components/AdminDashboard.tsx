import { useState, useMemo, useTransition, useEffect, Dispatch, SetStateAction } from 'react';
import { AdminDashboardData, User, Channel, AuditResult } from '../types';
import { AdminNoticeTab } from './AdminNoticeTab';
import { AdminConfigTab } from './AdminConfigTab';
import { AdminAuditTab } from './AdminAuditTab';
import { AdminLogsTab } from './AdminLogsTab';
import { AdminMTProtoAccountsTab } from './AdminMTProtoAccountsTab';
import { AdminPremiumFeaturesTab } from './AdminPremiumFeaturesTab';
import { AdminSubscriptionsTab } from './AdminSubscriptionsTab';
import { NoticeButton, NoticeTarget, updateUserAdmin, updateUserBlacklist } from '../api';
import { Users, Hash, ArrowLeft, ChevronRight, User as UserIcon, ShieldCheck, UserX, UserCheck, MessageSquare, Radio, BarChart3, Crown, Ban, Mail, TrendingUp } from 'lucide-react';
import { useToast } from './Toast';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { MetricCard } from './admin/MetricCard';
import { DataTable, Column } from './admin/DataTable';
import { StatusBadge } from './admin/StatusBadge';

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

// ───── Helpers ─────

function formatNum(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M';
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K';
  return String(n);
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

  // ── Analytics ──
  const analytics = useMemo(() => {
    const totalUsers = usersList.length;
    const totalChannels = channelsList.length;
    const admins = usersList.filter(u => u.is_admin).length;
    const blacklisted = usersList.filter(u => u.is_blacklisted).length;
    const withChannels = usersList.filter(u => (u.channels?.length || 0) > 0).length;
    const avgChannels = totalUsers > 0 ? (totalChannels / totalUsers) : 0;

    // Channel distribution
    const dist: Record<string, number> = { 0: 0, 1: 0, 2: 0, 3: 0, '4+': 0 };
    usersList.forEach(u => {
      const c = u.channels?.length || 0;
      if (c >= 4) dist['4+']++;
      else dist[c] = (dist[c] || 0) + 1;
    });

    // Top users by channel count
    const topUsers = [...usersList]
      .sort((a, b) => (b.channels?.length || 0) - (a.channels?.length || 0))
      .slice(0, 5);

    return { totalUsers, totalChannels, admins, blacklisted, withChannels, avgChannels, dist, topUsers };
  }, [usersList, channelsList]);

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

  const renderOverviewTab = () => {
    const { totalUsers, totalChannels, admins, blacklisted, withChannels, avgChannels, dist, topUsers } = analytics;
    const activeRate = totalUsers > 0 ? Math.round((withChannels / totalUsers) * 100) : 0;

    const distributionColors: Record<string, string> = {
      '0': 'var(--hint)',
      '1': 'var(--accent)',
      '2': 'var(--success)',
      '3': 'var(--warning)',
      '4+': 'var(--danger)',
    };

    const topUserColumns: Column<any>[] = [
      { key: 'rank', label: '#', width: '48px', render: (_: any, row: any) => (
        <span className="text-[11px] font-bold text-muted-foreground">{row.rank}</span>
      )},
      { key: 'initial', label: '', width: '36px', render: (_: any, row: any) => (
        <div className="flex items-center justify-center size-7 rounded-full text-[11px] font-bold" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
          {row.name.charAt(0).toUpperCase()}
        </div>
      )},
      { key: 'name', label: 'Nome', render: (_: any, row: any) => (
        <span className="text-[13px] font-semibold">{row.name}</span>
      )},
      { key: 'channels', label: 'Canais', align: 'right', render: (_: any, row: any) => (
        <Badge variant="secondary" className="text-[10px] font-mono">
          {row.channels} {row.channels === 1 ? 'canal' : 'canais'}
        </Badge>
      )},
    ];

    const topUserData = topUsers.map((u, i) => ({
      rank: i + 1,
      initial: u.first_name?.charAt(0) || '?',
      name: u.first_name || 'Sem nome',
      channels: u.channels?.length || 0,
      id: u.id,
    }));

    return (
      <div className="space-y-5">
        {/* Metric Grid */}
        <div className="admin-metrics-grid">
          <MetricCard
            title="Usuários"
            value={formatNum(totalUsers)}
            changeLabel={`${withChannels} ativos (${activeRate}%)`}
            icon={<Users size={18} />}
            iconColor="var(--accent)"
          />
          <MetricCard
            title="Canais"
            value={formatNum(totalChannels)}
            changeLabel={`${avgChannels.toFixed(1)} por usuário`}
            icon={<Hash size={18} />}
            iconColor="var(--success)"
          />
          <MetricCard
            title="Admins"
            value={admins}
            changeLabel={totalUsers > 0 ? `${((admins / totalUsers) * 100).toFixed(1)}%` : '—'}
            icon={<Crown size={18} />}
            iconColor="var(--warning)"
          />
          <MetricCard
            title="Blacklist"
            value={blacklisted}
            changeLabel={totalUsers > 0 ? `${((blacklisted / totalUsers) * 100).toFixed(1)}%` : '—'}
            icon={<Ban size={18} />}
            iconColor="var(--danger)"
          />
        </div>

        {/* Distribution + Top Users */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-4">
                <BarChart3 size={16} className="text-accent" />
                <h3 className="text-sm font-bold">Distribuição de Canais</h3>
              </div>
              <div className="space-y-2">
                {Object.entries(dist).map(([key, count]) => {
                  const pct = totalUsers > 0 ? (count / totalUsers) * 100 : 0;
                  return (
                    <div key={key} className="flex items-center gap-3">
                      <span className="text-xs font-medium text-muted-foreground w-12 shrink-0 text-right">{key === '4+' ? '4+' : key}</span>
                      <div className="flex-1 h-5 rounded-md bg-muted/30 overflow-hidden">
                        <div
                          className="h-full rounded-md transition-all duration-700"
                          style={{ width: `${pct}%`, background: distributionColors[key] || 'var(--accent)' }}
                        />
                      </div>
                      <span className="text-sm font-bold w-8 text-right">{count}</span>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-4">
                <TrendingUp size={16} className="text-accent" />
                <h3 className="text-sm font-bold">Top 5 — Mais Canais</h3>
              </div>
              <DataTable
                columns={topUserColumns}
                data={topUserData}
                searchable={false}
                pageSize={5}
                emptyMessage="Nenhum usuário com canais"
              />
            </CardContent>
          </Card>
        </div>

        {/* Quick Actions */}
        <div className="flex flex-wrap gap-2">
          <Button variant="secondary" size="sm" onClick={() => window.location.href = '/admin/dash?tab=users'}>
            <Users size={14} /> Gerenciar Usuários
          </Button>
          <Button variant="secondary" size="sm" onClick={() => window.location.href = '/admin/dash?tab=channels'}>
            <Hash size={14} /> Ver Canais
          </Button>
          <Button variant="secondary" size="sm" onClick={() => window.location.href = '/admin/dash?tab=notice'}>
            <Mail size={14} /> Enviar Broadcast
          </Button>
        </div>
      </div>
    );
  };

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
            <Button variant="default" size="sm" className="w-full" onClick={() => onMessageUser(adminSelectedUser.id)}>
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
          data={usersList}
          searchable={true}
          searchPlaceholder="Buscar por nome ou ID..."
          searchKeys={['first_name', 'username', 'id']}
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
        data={channelsList}
        searchable={true}
        searchPlaceholder="Buscar canal por título ou ID..."
        searchKeys={['title', 'id']}
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
        />
      </div>
    );
  };

  // ── Render ──

  return (
    <div className={`space-y-4 ${isPending ? 'opacity-60 pointer-events-none' : ''}`} style={{ transition: 'opacity 0.2s ease' }}>
      {/* Header */}
      <div className="flex items-center justify-between animate-stagger-in" style={{ animationDelay: '0s' }}>
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center size-10 rounded-xl shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
            <Radio size={20} />
          </div>
          <div>
            <h1 className="text-lg font-bold">Painel Administrativo</h1>
            <p className="text-xs text-muted-foreground">
              {localActiveTab === 'overview' && 'Métricas e visão geral do sistema'}
              {localActiveTab === 'users' && 'Gerencie todos os usuários da plataforma'}
              {localActiveTab === 'channels' && 'Todos os canais cadastrados'}
              {localActiveTab === 'notice' && 'Envie mensagens globais para usuários'}
              {localActiveTab === 'audit' && 'Auditoria de bots nos canais'}
              {localActiveTab === 'logs' && 'Histórico de eventos do sistema'}
              {localActiveTab === 'config' && 'Configurações globais do servidor'}
              {localActiveTab === 'accounts' && 'Contas Telegram para edição de postagens'}
              {localActiveTab === 'premium-features' && 'Gerencie as features premium do sistema'}
              {localActiveTab === 'subscriptions' && 'Gerencie assinaturas de todos os usuários'}
            </p>
          </div>
        </div>
      </div>

      <div className="h-px bg-border/50" />

      {localActiveTab === 'overview' && renderOverviewTab()}
      {localActiveTab === 'users' && !adminSelectedUser && renderUsersTab()}
      {localActiveTab === 'users' && adminSelectedUser && renderUserDetail()}
      {localActiveTab === 'channels' && renderChannelsTab()}
      {localActiveTab === 'audit' && (
        <div className="space-y-4">
          <AdminAuditTab
            navigateToChannel={navigateToChannel}
            onOpenUser={onOpenUserDetail}
            results={auditResults}
            setResults={setAuditResults}
            loading={auditLoading}
            onRunAudit={handleRunAudit}
          />
        </div>
      )}
      {localActiveTab === 'notice' && renderNoticeTab()}
      {localActiveTab === 'logs' && (
        <div className="space-y-4">
          <AdminLogsTab navigateToChannel={navigateToChannel} initialChannelId={initialLogsChannelId} />
        </div>
      )}
      {localActiveTab === 'config' && (
        <div className="space-y-4">
          <AdminConfigTab />
        </div>
      )}
      {localActiveTab === 'accounts' && (
        <div className="space-y-4">
          <AdminMTProtoAccountsTab />
        </div>
      )}
      {localActiveTab === 'premium-features' && (
        <div className="space-y-4">
          <AdminPremiumFeaturesTab toast={toast} />
        </div>
      )}
      {localActiveTab === 'subscriptions' && (
        <div className="space-y-4">
          <AdminSubscriptionsTab toast={toast} />
        </div>
      )}
    </div>
  );
}
