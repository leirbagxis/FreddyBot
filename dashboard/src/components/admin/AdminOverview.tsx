import { memo } from 'react';
import { Users, UserPlus, Hash, MessageSquare, Send, CreditCard, AlertTriangle, Clock } from 'lucide-react';
import { MetricCard } from './MetricCard';
import { ActivityFeed } from './ActivityFeed';
import { SystemHealth } from './SystemHealth';
import { FinanceSummary } from './FinanceSummary';
import { AlertsPanel } from './AlertsPanel';

interface OverviewData {
  metrics: {
    activeUsers: { value: number; change: number; sparkline: number[] };
    newUsers: { value: number; change: number; sparkline: number[] };
    activeChannels: { value: number; change: number; sparkline: number[] };
    postsToday: { value: number; change: number; sparkline: number[] };
    scheduledPosts: { value: number; change: number };
    messagesSent: { value: number; change: number; sparkline: number[] };
    paymentsReceived: { value: number; change: number };
    systemErrors: { value: number; change: number };
  };
  activity: Array<{
    id: string;
    type: any;
    message: string;
    timestamp: string;
  }>;
  systemHealth: Array<{
    name: string;
    status: 'healthy' | 'warning' | 'error';
    latency?: number;
  }>;
  finance: {
    dailyRevenue: number;
    monthlyRevenue: number;
    pendingPayments: number;
    overduePayments: number;
  };
  alerts: Array<{
    id: string;
    type: any;
    message: string;
    count: number;
  }>;
}

interface AdminOverviewProps {
  data: OverviewData;
  loading?: boolean;
}

export const AdminOverview = memo(function AdminOverview({ data, loading }: AdminOverviewProps) {
  const { metrics } = data;

  return (
    <div className="admin-overview">
      {/* Metrics Grid */}
      <div className="admin-metrics-grid">
        <MetricCard
          title="Usuários Ativos"
          value={metrics.activeUsers.value.toLocaleString()}
          change={metrics.activeUsers.change}
          icon={<Users size={18} />}
          iconColor="var(--accent)"
          sparkline={metrics.activeUsers.sparkline}
          loading={loading}
        />
        <MetricCard
          title="Novos Usuários"
          value={metrics.newUsers.value.toLocaleString()}
          change={metrics.newUsers.change}
          changeLabel="vs. período anterior"
          icon={<UserPlus size={18} />}
          iconColor="var(--success)"
          sparkline={metrics.newUsers.sparkline}
          loading={loading}
        />
        <MetricCard
          title="Canais Ativos"
          value={metrics.activeChannels.value.toLocaleString()}
          change={metrics.activeChannels.change}
          icon={<Hash size={18} />}
          iconColor="var(--accent)"
          sparkline={metrics.activeChannels.sparkline}
          loading={loading}
        />
        <MetricCard
          title="Publicações Hoje"
          value={metrics.postsToday.value.toLocaleString()}
          change={metrics.postsToday.change}
          icon={<Send size={18} />}
          iconColor="var(--warning)"
          sparkline={metrics.postsToday.sparkline}
          loading={loading}
        />
        <MetricCard
          title="Agendados"
          value={metrics.scheduledPosts.value.toLocaleString()}
          change={metrics.scheduledPosts.change}
          icon={<Clock size={18} />}
          iconColor="var(--accent)"
          loading={loading}
        />
        <MetricCard
          title="Mensagens Enviadas"
          value={metrics.messagesSent.value.toLocaleString()}
          change={metrics.messagesSent.change}
          icon={<MessageSquare size={18} />}
          iconColor="var(--success)"
          sparkline={metrics.messagesSent.sparkline}
          loading={loading}
        />
        <MetricCard
          title="Pagamentos"
          value={`⭐ ${metrics.paymentsReceived.value}`}
          change={metrics.paymentsReceived.change}
          icon={<CreditCard size={18} />}
          iconColor="var(--success)"
          loading={loading}
        />
        <MetricCard
          title="Erros"
          value={metrics.systemErrors.value.toLocaleString()}
          change={metrics.systemErrors.change}
          icon={<AlertTriangle size={18} />}
          iconColor="var(--danger)"
          loading={loading}
        />
      </div>

      {/* Content Grid */}
      <div className="admin-content-grid">
        {/* Left column */}
        <div className="admin-content-col">
          <ActivityFeed items={data.activity} loading={loading} />
          <FinanceSummary data={data.finance} loading={loading} />
        </div>

        {/* Right column */}
        <div className="admin-content-col">
          <SystemHealth items={data.systemHealth} loading={loading} />
          <AlertsPanel alerts={data.alerts} loading={loading} />
        </div>
      </div>
    </div>
  );
});
