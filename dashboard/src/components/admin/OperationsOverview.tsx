import { memo, useMemo } from 'react';
import { ArrowRight, BellRing, ChartNoAxesColumnIncreasing, Hash, ShieldCheck, Users } from 'lucide-react';
import { Channel, User } from '../../types';
import { getOperationalAlerts, getOverviewMetrics, OperationalAlertKind } from './crmSelectors';

interface OperationsOverviewProps {
  users: User[];
  channels: Channel[];
  onOpenUser: (id: number) => void;
  onViewUsers: () => void;
  onReviewAlert: (kind: OperationalAlertKind) => void;
}

function compact(value: number) {
  return new Intl.NumberFormat('pt-BR', { notation: 'compact', maximumFractionDigits: 1 }).format(value);
}

export const OperationsOverview = memo(function OperationsOverview({ users, channels, onOpenUser, onViewUsers, onReviewAlert }: OperationsOverviewProps) {
  const metrics = useMemo(() => getOverviewMetrics(users, channels), [channels, users]);
  const alerts = useMemo(() => getOperationalAlerts(users, channels), [channels, users]);
  const reviewUsers = useMemo(
    () => [...users].filter((user) => user.is_blacklisted || user.is_admin || !(user.channels?.length)).slice(0, 6),
    [users],
  );

  const items = [
    { label: 'Usuários', value: compact(metrics.totalUsers), note: 'base cadastrada', icon: Users },
    { label: 'Canais', value: compact(metrics.totalChannels), note: 'conectados', icon: Hash },
    { label: 'Administradores', value: compact(metrics.admins), note: 'com acesso', icon: ShieldCheck },
    { label: 'Ativação', value: `${metrics.activationRate}%`, note: `${compact(metrics.activatedUsers)} com canais`, icon: ChartNoAxesColumnIncreasing },
  ];

  return (
    <section className="operations-overview" aria-label="Resumo operacional">
      <header className="operations-overview-heading">
        <div>
          <span>Administração</span>
          <div className="operations-title">Visão geral</div>
          <p>Estado atual da base e acessos que merecem revisão.</p>
        </div>
        <button type="button" className="operations-overview-link" onClick={onViewUsers}>
          Ver usuários <ArrowRight size={15} aria-hidden="true" />
        </button>
      </header>

      <div className="operations-metrics">
        {items.map(({ label, value, note, icon: Icon }) => (
          <article className="operations-metric" key={label}>
            <Icon size={16} aria-hidden="true" />
            <span>{label}</span>
            <strong>{value}</strong>
            <small>{note}</small>
          </article>
        ))}
      </div>

      <div className="operations-insights" aria-label="Indicadores operacionais">
        <section className="operations-health">
          <div className="operations-section-heading">
            <div>
              <span>Saúde da base</span>
              <h2>Ativação e acompanhamento</h2>
            </div>
            <strong>{metrics.activationRate}%</strong>
          </div>
          <div className="operations-progress" role="progressbar" aria-label="Taxa de ativação" aria-valuemin={0} aria-valuemax={100} aria-valuenow={metrics.activationRate}>
            <span style={{ width: `${metrics.activationRate}%` }} />
          </div>
          <div className="operations-health-grid">
            <div><strong>{compact(metrics.activatedUsers)}</strong><span>usuários com canal</span></div>
            <div><strong>{compact(metrics.withoutChannels)}</strong><span>aguardando ativação</span></div>
            <div><strong>{compact(metrics.recentUsers)}</strong><span>novos em 7 dias</span></div>
          </div>
        </section>

        <section className="operations-alerts" aria-labelledby="operations-alerts-title">
          <div className="operations-section-heading">
            <div>
              <span>Notificações</span>
              <h2 id="operations-alerts-title">Pontos que pedem atenção</h2>
            </div>
            <BellRing size={17} aria-hidden="true" />
          </div>
          {alerts.length ? (
            <div className="operations-alert-list">
              {alerts.map((alert) => (
                <button type="button" key={alert.id} className={`operations-alert is-${alert.severity}`} onClick={() => onReviewAlert(alert.id)}>
                  <strong>{alert.count}</strong>
                  <span><b>{alert.title}</b><small>{alert.description}</small></span>
                  <ArrowRight size={15} aria-hidden="true" />
                </button>
              ))}
            </div>
          ) : (
            <p className="operations-alert-empty">Nenhuma notificação operacional no momento.</p>
          )}
        </section>
      </div>

      <section className="operations-queue" aria-labelledby="operations-queue-title">
        <header>
          <div>
            <div id="operations-queue-title" className="operations-queue-title">Fila de revisão</div>
            <p>Usuários com acesso administrativo, bloqueio ou sem canais vinculados.</p>
          </div>
          <span>{reviewUsers.length} exibidos</span>
        </header>
        {reviewUsers.length ? (
          <div className="operations-queue-list">
            {reviewUsers.map((user) => (
              <button type="button" key={user.id} onClick={() => onOpenUser(user.id)}>
                <span className="operations-avatar">{(user.first_name || '?')[0].toUpperCase()}</span>
                <span className="operations-user">
                  <strong>{user.first_name || 'Sem nome'}</strong>
                  <small>ID {user.id} · {user.channels?.length || 0} canais</small>
                </span>
                <span className={user.is_blacklisted ? 'is-risk' : user.is_admin ? 'is-admin' : ''}>
                  {user.is_blacklisted ? 'Bloqueado' : user.is_admin ? 'Admin' : 'Sem canais'}
                </span>
                <ArrowRight size={15} aria-hidden="true" />
              </button>
            ))}
          </div>
        ) : (
          <p className="operations-empty">Nenhuma revisão pendente na base atual.</p>
        )}
      </section>
    </section>
  );
});
