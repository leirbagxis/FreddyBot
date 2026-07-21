import { memo } from 'react';
import { DollarSign, TrendingUp, Clock, AlertTriangle } from 'lucide-react';

interface FinanceData {
  dailyRevenue: number;
  monthlyRevenue: number;
  pendingPayments: number;
  overduePayments: number;
  totalRevenue?: number;
  currency?: string;
}

interface FinanceSummaryProps {
  data: FinanceData;
  loading?: boolean;
}

export const FinanceSummary = memo(function FinanceSummary({ data, loading }: FinanceSummaryProps) {
  const currency = data.currency || '⭐';

  return (
    <div className="admin-card">
      <div className="admin-card-header">
        <h3 className="admin-card-title">
          <DollarSign size={15} />
          Resumo Financeiro
        </h3>
      </div>
      <div className="admin-card-body">
        {loading ? (
          <div className="finance-loading">
            {[1, 2, 3, 4].map(i => (
              <div key={i} className="finance-stat skeleton" />
            ))}
          </div>
        ) : (
          <div className="finance-grid">
            <div className="finance-stat">
              <div className="finance-stat-icon" style={{ color: 'var(--success)' }}>
                <TrendingUp size={14} />
              </div>
              <div className="finance-stat-info">
                <span className="finance-stat-label">Receita Diária</span>
                <span className="finance-stat-value">{currency} {data.dailyRevenue}</span>
              </div>
            </div>

            <div className="finance-stat">
              <div className="finance-stat-icon" style={{ color: 'var(--accent)' }}>
                <DollarSign size={14} />
              </div>
              <div className="finance-stat-info">
                <span className="finance-stat-label">Receita Mensal</span>
                <span className="finance-stat-value">{currency} {data.monthlyRevenue}</span>
              </div>
            </div>

            <div className="finance-stat">
              <div className="finance-stat-icon" style={{ color: 'var(--warning)' }}>
                <Clock size={14} />
              </div>
              <div className="finance-stat-info">
                <span className="finance-stat-label">Pagamentos Pendentes</span>
                <span className="finance-stat-value">{data.pendingPayments}</span>
              </div>
            </div>

            <div className="finance-stat">
              <div className="finance-stat-icon" style={{ color: 'var(--danger)' }}>
                <AlertTriangle size={14} />
              </div>
              <div className="finance-stat-info">
                <span className="finance-stat-label">Pagamentos Vencidos</span>
                <span className="finance-stat-value">{data.overduePayments}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
});
