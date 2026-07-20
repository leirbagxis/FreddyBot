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
import { Users, Hash, Search, ArrowLeft, ChevronRight, User as UserIcon, ShieldCheck, UserX, UserCheck, MessageSquare, Radio, Activity, BarChart3, TrendingUp, Crown, Ban, Mail } from 'lucide-react';
import { useToast } from './Toast';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';

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

// ───── Metric Card ─────

function MetricCard({ icon, label, value, sub, color, delay }: {
  icon: React.ReactNode;
  label: string;
  value: string | number;
  sub?: string;
  color: 'accent' | 'success' | 'warning' | 'danger' | 'info';
  delay?: string;
}) {
  const colors: Record<string, { bg: string; text: string }> = {
    accent: { bg: 'var(--accent-soft)', text: 'var(--accent)' },
    success: { bg: 'var(--success-soft)', text: 'var(--success)' },
    warning: { bg: 'var(--warning-soft)', text: 'var(--warning)' },
    danger: { bg: 'var(--danger-soft)', text: 'var(--danger)' },
    info: { bg: 'rgba(99, 102, 241, 0.06)', text: 'var(--text-secondary)' },
  };
  const c = colors[color];

  return (
    <div
      className="rounded-xl border border-border p-4 animate-stagger-in"
      style={{ animationDelay: delay || '0s' }}
    >
      <div className="flex items-center gap-3 mb-2">
        <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: c.bg, color: c.text }}>
          {icon}
        </div>
        <span className="text-[11px] font-semibold text-muted-foreground uppercase tracking-[0.08em]">{label}</span>
      </div>
      <p className="text-2xl font-extrabold tracking-tight">{value}</p>
      {sub && <p className="text-[11px] text-muted-foreground mt-0.5">{sub}</p>}
    </div>
  );
}

// ───── Distribution Bar ─────

function DistributionBar({ label, count, total, color }: {
  label: string;
  count: number;
  total: number;
  color: string;
}) {
  const pct = total > 0 ? (count / total) * 100 : 0;
  return (
    <div className="flex items-center gap-3">
      <span className="text-[12px] font-medium text-muted-foreground w-16 shrink-0 text-right">{label}</span>
      <div className="flex-1 h-5 rounded-md bg-muted/30 overflow-hidden">
        <div
          className="h-full rounded-md transition-all duration-700"
          style={{ width: `${pct}%`, background: color }}
        />
      </div>
      <span className="text-[13px] font-bold w-8 text-right">{count}</span>
    </div>
  );
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
  const [adminSearch, setAdminSearch] = useState('');
  const [adminChannelCountFilter, setAdminChannelCountFilter] = useState('');
  const [visibleUsersCount, setVisibleUsersCount] = useState(40);
  const [visibleChannelsCount, setVisibleChannelsCount] = useState(40);
  const [channelSearch, setChannelSearch] = useState('');
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
    const dist: Record<number, number> = { 0: 0, 1: 0, 2: 0, 3: 0, '4+': 0 };
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

  const filteredUsers = useMemo(() => {
    const minChannelCount = parseInt(adminChannelCountFilter, 10);
    const hasChannelCountFilter = !Number.isNaN(minChannelCount);

    const filtered = usersList.filter(u => {
      const name = (u.firstName || (u as any).first_name || '').toLowerCase();
      const matchesSearch = name.includes(adminSearch.toLowerCase()) || u.id.toString().includes(adminSearch);
      const matchesCount = hasChannelCountFilter ? (u.channels?.length || 0) >= minChannelCount : true;
      return matchesSearch && matchesCount;
    });

    if (!hasChannelCountFilter) return filtered;

    return [...filtered].sort((a, b) => {
      const channelDiff = (a.channels?.length || 0) - (b.channels?.length || 0);
      if (channelDiff !== 0) return channelDiff;
      const aName = (a.firstName || (a as any).first_name || '').toLowerCase();
      const bName = (b.firstName || (b as any).first_name || '').toLowerCase();
      const nameDiff = aName.localeCompare(bName);
      if (nameDiff !== 0) return nameDiff;
      return a.id - b.id;
    });
  }, [usersList, adminSearch, adminChannelCountFilter]);

  const filteredChannels = useMemo(() => {
    return channelsList.filter(c => {
      return c.title.toLowerCase().includes(channelSearch.toLowerCase()) || c.id.toString().includes(channelSearch);
    });
  }, [channelsList, channelSearch]);

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

    return (
      <div className="space-y-5">
        {/* Metric Grid */}
        <div className="grid grid-cols-2 gap-2.5 sm:grid-cols-3">
          <MetricCard
            icon={<Users size={16} />}
            label="Usuários"
            value={formatNum(totalUsers)}
            sub={`${withChannels} ativos (${activeRate}%)`}
            color="accent"
            delay="0.02s"
          />
          <MetricCard
            icon={<Hash size={16} />}
            label="Canais"
            value={formatNum(totalChannels)}
            sub={`${avgChannels.toFixed(1)} por usuário`}
            color="success"
            delay="0.04s"
          />
          <MetricCard
            icon={<Crown size={16} />}
            label="Admins"
            value={admins}
            sub={totalUsers > 0 ? `${((admins / totalUsers) * 100).toFixed(1)}% dos usuários` : '—'}
            color="warning"
            delay="0.06s"
          />
          <MetricCard
            icon={<Ban size={16} />}
            label="Blacklist"
            value={blacklisted}
            sub={totalUsers > 0 ? `${((blacklisted / totalUsers) * 100).toFixed(1)}% dos usuários` : '—'}
            color="danger"
            delay="0.08s"
          />
          <MetricCard
            icon={<Activity size={16} />}
            label="Taxa de Ativação"
            value={`${activeRate}%`}
            sub={`${withChannels} de ${totalUsers} usam canais`}
            color="info"
            delay="0.1s"
          />
        </div>

        {/* Distribution + Top Users */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          {/* Distribution */}
          <div className="rounded-xl border border-border p-4 animate-stagger-in" style={{ animationDelay: '0.12s' }}>
            <div className="flex items-center gap-2 mb-4">
              <BarChart3 size={16} className="text-accent" />
              <h3 className="text-[13px] font-bold">Distribuição de Canais</h3>
            </div>
            <div className="space-y-2">
              {Object.entries(dist).map(([key, count]) => {
                const colors: Record<string, string> = {
                  '0': 'var(--hint)',
                  '1': 'var(--accent)',
                  '2': 'var(--success)',
                  '3': 'var(--warning)',
                  '4+': 'var(--danger)',
                };
                return (
                  <DistributionBar
                    key={key}
                    label={key === '4+' ? '4+' : key}
                    count={count}
                    total={totalUsers}
                    color={colors[key] || 'var(--accent)'}
                  />
                );
              })}
            </div>
          </div>

          {/* Top Users */}
          <div className="rounded-xl border border-border p-4 animate-stagger-in" style={{ animationDelay: '0.14s' }}>
            <div className="flex items-center gap-2 mb-4">
              <TrendingUp size={16} className="text-accent" />
              <h3 className="text-[13px] font-bold">Top 5 — Mais Canais</h3>
            </div>
            <div className="space-y-2">
              {topUsers.length > 0 ? topUsers.map((u, i) => {
                const name = u.firstName || (u as any).first_name || 'Sem nome';
                const chCount = u.channels?.length || 0;
                return (
                  <div key={u.id} className="flex items-center gap-3 py-1.5">
                    <span className="text-[11px] font-bold text-muted-foreground w-5 shrink-0 text-right">
                      {i + 1}
                    </span>
                    <div className="flex items-center justify-center size-7 rounded-full shrink-0 text-[11px] font-bold" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                      {name.charAt(0).toUpperCase()}
                    </div>
                    <div className="min-w-0 flex-1">
                      <span className="text-[13px] font-semibold truncate block">{name}</span>
                    </div>
                    <Badge variant="secondary" className="text-[10px] font-mono">
                      {chCount} {chCount === 1 ? 'canal' : 'canais'}
                    </Badge>
                  </div>
                );
              }) : (
                <p className="text-[13px] text-muted-foreground text-center py-4">Nenhum usuário com canais</p>
              )}
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="flex flex-wrap gap-2 animate-stagger-in" style={{ animationDelay: '0.16s' }}>
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
    const name = adminSelectedUser.firstName || (adminSelectedUser as any).first_name || 'Sem nome';
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
    const visibleUsers = filteredUsers.slice(0, visibleUsersCount);

    return (
      <div className="space-y-4">
        {/* Stats mini-row */}
        <div className="flex items-center gap-2 flex-wrap">
          <Badge variant="secondary" className="text-[11px] gap-1.5">
            <Users size={12} /> {usersList.length} total
          </Badge>
          <Badge variant="default" className="text-[11px] gap-1.5">
            <Crown size={12} /> {analytics.admins} admins
          </Badge>
          <Badge variant="destructive" className="text-[11px] gap-1.5">
            <Ban size={12} /> {analytics.blacklisted} blacklist
          </Badge>
        </div>

        {/* Search */}
        <div className="flex flex-col gap-2">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
            <Input
              type="text"
              placeholder="Buscar usuário por nome ou ID..."
              className="pl-9 h-10 rounded-xl"
              value={adminSearch}
              onChange={(e) => {
                setAdminSearch(e.target.value);
                setVisibleUsersCount(40);
              }}
            />
          </div>
          <div className="relative">
            <Hash className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
            <Input
              type="number"
              placeholder="Filtrar por mínimo de canais"
              className="pl-9 h-10 rounded-xl"
              value={adminChannelCountFilter}
              onChange={(e) => {
                setAdminChannelCountFilter(e.target.value);
                setVisibleUsersCount(40);
              }}
            />
          </div>
        </div>

        {/* Users list */}
        <div className="space-y-1.5">
          {visibleUsers.length > 0 ? (
            <>
              {visibleUsers.map((u) => (
                <button
                  key={u.id}
                  className="flex items-center w-full text-left gap-3 rounded-xl border border-border p-3 hover:bg-muted/30 transition-colors"
                  onClick={() => setAdminSelectedUser(u)}
                >
                  <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                    <UserIcon size={16} />
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className="text-[13px] font-semibold truncate">{u.firstName || (u as any).first_name || 'Sem nome'}</span>
                      {u.is_admin && <Badge variant="default" className="text-[9px] h-[18px]">Admin</Badge>}
                      {u.is_blacklisted && <Badge variant="destructive" className="text-[9px] h-[18px]">Bloqueado</Badge>}
                    </div>
                    <p className="text-[11px] text-muted-foreground mt-0.5">
                      ID: {u.id} • {u.channels?.length || 0} canais
                    </p>
                  </div>
                  <ChevronRight size={16} className="shrink-0 text-muted-foreground/30" />
                </button>
              ))}
              {filteredUsers.length > visibleUsersCount && (
                <Button
                  variant="secondary"
                  className="w-full mt-2"
                  onClick={() => setVisibleUsersCount(prev => prev + 40)}
                >
                  Carregar mais usuários...
                </Button>
              )}
            </>
          ) : (
            <div className="flex flex-col items-center py-8 text-muted-foreground rounded-xl border border-border">
              <UserIcon size={28} className="opacity-30 mb-2" />
              <p className="text-[13px] font-medium">Nenhum usuário encontrado</p>
            </div>
          )}
        </div>
      </div>
    );
  };

  // ── Channels Tab ──

  const renderChannelsTab = () => {
    const visibleChannels = filteredChannels.slice(0, visibleChannelsCount);

    return (
      <div className="space-y-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
          <Input
            type="text"
            placeholder="Buscar canal por título ou ID..."
            className="pl-9 h-10 rounded-xl"
            value={channelSearch}
            onChange={(e) => {
              setChannelSearch(e.target.value);
              setVisibleChannelsCount(40);
            }}
          />
        </div>

        <div className="space-y-1.5">
          {visibleChannels.length > 0 ? (
            <>
              {visibleChannels.map((c) => (
                <button
                  key={c.id}
                  className="flex items-center w-full text-left gap-3 rounded-xl border border-border p-3 hover:bg-muted/30 transition-colors"
                  onClick={() => navigateToChannel(c.id)}
                >
                  <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--success-soft)', color: 'var(--success)' }}>
                    <Hash size={16} />
                  </div>
                  <div className="min-w-0 flex-1">
                    <span className="text-[13px] font-semibold truncate block">{c.title}</span>
                    <p className="text-[11px] text-muted-foreground mt-0.5">ID: {c.id} • Dono: {c.ownerId}</p>
                  </div>
                  <ChevronRight size={16} className="shrink-0 text-muted-foreground/30" />
                </button>
              ))}
              {filteredChannels.length > visibleChannelsCount && (
                <Button
                  variant="secondary"
                  className="w-full mt-2"
                  onClick={() => setVisibleChannelsCount(prev => prev + 40)}
                >
                  Carregar mais canais...
                </Button>
              )}
            </>
          ) : (
            <div className="flex flex-col items-center py-8 text-muted-foreground rounded-xl border border-border">
              <Hash size={28} className="opacity-30 mb-2" />
              <p className="text-[13px] font-medium">Nenhum canal encontrado</p>
            </div>
          )}
        </div>
      </div>
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
