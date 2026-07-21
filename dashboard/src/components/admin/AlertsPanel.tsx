import { memo } from 'react';
import { AlertTriangle, Shield, CreditCard, Ban, AlertOctagon } from 'lucide-react';

interface Alert {
  id: string;
  type: 'invalid_group' | 'private_channel' | 'expired_payment' | 'critical_error' | 'blocked_user';
  message: string;
  count: number;
}

interface AlertsPanelProps {
  alerts: Alert[];
  loading?: boolean;
}

const ALERT_CONFIG = {
  invalid_group: { icon: AlertTriangle, color: 'var(--warning)', label: 'Grupos Inválidos' },
  private_channel: { icon: Shield, color: 'var(--accent)', label: 'Canais Privados' },
  expired_payment: { icon: CreditCard, color: 'var(--danger)', label: 'Pagamentos Expirados' },
  critical_error: { icon: AlertOctagon, color: 'var(--danger)', label: 'Erros Críticos' },
  blocked_user: { icon: Ban, color: 'var(--warning)', label: 'Usuários Bloqueados' },
};

export const AlertsPanel = memo(function AlertsPanel({ alerts, loading }: AlertsPanelProps) {
  if (loading) {
    return (
      <div className="admin-card">
        <div className="admin-card-header">
          <h3 className="admin-card-title">Alertas</h3>
        </div>
        <div className="admin-card-body">
          {[1, 2, 3].map(i => (
            <div key={i} className="alert-item skeleton" />
          ))}
        </div>
      </div>
    );
  }

  const totalAlerts = alerts.reduce((sum, a) => sum + a.count, 0);

  return (
    <div className="admin-card">
      <div className="admin-card-header">
        <h3 className="admin-card-title">
          <AlertTriangle size={15} />
          Alertas
        </h3>
        {totalAlerts > 0 && (
          <span className="admin-card-badge warning">{totalAlerts}</span>
        )}
      </div>
      <div className="admin-card-body alerts-list">
        {alerts.length === 0 && (
          <div className="alerts-empty">
            <Shield size={24} />
            <span>Nenhum alerta no momento</span>
          </div>
        )}
        {alerts.map((alert) => {
          const config = ALERT_CONFIG[alert.type];
          const Icon = config.icon;
          return (
            <div key={alert.id} className="alert-item">
              <div className="alert-item-icon" style={{ color: config.color }}>
                <Icon size={14} />
              </div>
              <div className="alert-item-content">
                <span className="alert-item-label">{config.label}</span>
                <span className="alert-item-message">{alert.message}</span>
              </div>
              <span className="alert-item-count">{alert.count}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
});
