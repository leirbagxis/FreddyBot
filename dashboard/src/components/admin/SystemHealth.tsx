import { memo } from 'react';
import { Activity, Database, Server, Wifi, CheckCircle, AlertCircle, XCircle } from 'lucide-react';

interface HealthItem {
  name: string;
  status: 'healthy' | 'warning' | 'error';
  latency?: number;
  lastCheck?: string;
}

interface SystemHealthProps {
  items: HealthItem[];
  loading?: boolean;
}

const STATUS_CONFIG = {
  healthy: { icon: CheckCircle, color: 'var(--success)', label: 'Online' },
  warning: { icon: AlertCircle, color: 'var(--warning)', label: 'Lento' },
  error: { icon: XCircle, color: 'var(--danger)', label: 'Offline' },
};

const SERVICE_ICONS: Record<string, any> = {
  'API Telegram': Wifi,
  'Banco de Dados': Database,
  'Filas': Activity,
  'Workers': Server,
  'Serviço Online': Server,
};

export const SystemHealth = memo(function SystemHealth({ items, loading }: SystemHealthProps) {
  if (loading) {
    return (
      <div className="admin-card">
        <div className="admin-card-header">
          <h3 className="admin-card-title">Saúde do Sistema</h3>
        </div>
        <div className="admin-card-body">
          {[1, 2, 3].map(i => (
            <div key={i} className="health-item skeleton" />
          ))}
        </div>
      </div>
    );
  }

  const healthyCount = items.filter(i => i.status === 'healthy').length;
  const overallStatus = items.some(i => i.status === 'error')
    ? 'error'
    : items.some(i => i.status === 'warning')
      ? 'warning'
      : 'healthy';

  return (
    <div className="admin-card">
      <div className="admin-card-header">
        <h3 className="admin-card-title">
          <Activity size={15} />
          Saúde do Sistema
        </h3>
        <div className={`health-overall-badge ${overallStatus}`}>
          {healthyCount}/{items.length} online
        </div>
      </div>
      <div className="admin-card-body health-list">
        {items.map((item, idx) => {
          const config = STATUS_CONFIG[item.status];
          const StatusIcon = config.icon;
          const ServiceIcon = SERVICE_ICONS[item.name] || Server;
          return (
            <div key={idx} className="health-item">
              <div className="health-item-left">
                <ServiceIcon size={14} className="health-item-service-icon" />
                <span className="health-item-name">{item.name}</span>
              </div>
              <div className="health-item-right">
                {item.latency != null && (
                  <span className="health-item-latency">{item.latency}ms</span>
                )}
                <div className="health-item-status" style={{ color: config.color }}>
                  <StatusIcon size={12} />
                  <span>{config.label}</span>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
});
