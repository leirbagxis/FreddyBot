import { memo } from 'react';
import {
  UserPlus, CreditCard, Hash, Trash2, Send, AlertTriangle,
  Clock, MessageSquare
} from 'lucide-react';

interface ActivityItem {
  id: string;
  type: 'user_joined' | 'payment' | 'channel_added' | 'channel_removed' | 'post_sent' | 'error' | 'broadcast';
  message: string;
  timestamp: string;
  userId?: number;
  channelId?: number;
}

interface ActivityFeedProps {
  items: ActivityItem[];
  loading?: boolean;
}

const ACTIVITY_CONFIG = {
  user_joined: { icon: UserPlus, color: 'var(--success)', label: 'Usuário' },
  payment: { icon: CreditCard, color: 'var(--accent)', label: 'Pagamento' },
  channel_added: { icon: Hash, color: 'var(--success)', label: 'Canal' },
  channel_removed: { icon: Trash2, color: 'var(--danger)', label: 'Canal' },
  post_sent: { icon: Send, color: 'var(--accent)', label: 'Publicação' },
  error: { icon: AlertTriangle, color: 'var(--danger)', label: 'Erro' },
  broadcast: { icon: MessageSquare, color: 'var(--warning)', label: 'Broadcast' },
};

function formatTimeAgo(timestamp: string): string {
  const now = new Date();
  const then = new Date(timestamp);
  const diffMs = now.getTime() - then.getTime();
  const diffMin = Math.floor(diffMs / 60000);
  if (diffMin < 1) return 'agora';
  if (diffMin < 60) return `${diffMin}m`;
  const diffH = Math.floor(diffMin / 60);
  if (diffH < 24) return `${diffH}h`;
  const diffD = Math.floor(diffH / 24);
  return `${diffD}d`;
}

export const ActivityFeed = memo(function ActivityFeed({ items, loading }: ActivityFeedProps) {
  if (loading) {
    return (
      <div className="admin-card">
        <div className="admin-card-header">
          <h3 className="admin-card-title">Atividade Recente</h3>
        </div>
        <div className="admin-card-body">
          {[1, 2, 3, 4, 5].map(i => (
            <div key={i} className="activity-item skeleton" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="admin-card">
      <div className="admin-card-header">
        <h3 className="admin-card-title">
          <Clock size={15} />
          Atividade Recente
        </h3>
        <span className="admin-card-badge">{items.length} eventos</span>
      </div>
      <div className="admin-card-body activity-feed-list">
        {items.length === 0 && (
          <div className="activity-empty">
            <MessageSquare size={24} />
            <span>Nenhuma atividade recente</span>
          </div>
        )}
        {items.map((item) => {
          const config = ACTIVITY_CONFIG[item.type] || ACTIVITY_CONFIG.error;
          const Icon = config.icon;
          return (
            <div key={item.id} className="activity-item">
              <div className="activity-item-icon" style={{ color: config.color }}>
                <Icon size={14} />
              </div>
              <div className="activity-item-content">
                <span className="activity-item-message">{item.message}</span>
                <span className="activity-item-time">{formatTimeAgo(item.timestamp)}</span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
});
