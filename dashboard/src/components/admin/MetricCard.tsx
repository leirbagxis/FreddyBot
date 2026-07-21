import { memo, ReactNode } from 'react';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';

interface MetricCardProps {
  title: string;
  value: string | number;
  change?: number; // percentage change from previous period
  changeLabel?: string;
  icon: ReactNode;
  iconColor?: string;
  sparkline?: number[]; // simple sparkline data points
  loading?: boolean;
}

export const MetricCard = memo(function MetricCard({
  title, value, change, changeLabel, icon, iconColor, sparkline, loading
}: MetricCardProps) {
  const trend = change != null ? (change > 0 ? 'up' : change < 0 ? 'down' : 'neutral') : null;

  return (
    <div className="metric-card">
      <div className="metric-card-header">
        <div className="metric-card-icon" style={{ color: iconColor || 'var(--accent)' }}>
          {icon}
        </div>
        {trend && (
          <div className={`metric-card-trend ${trend}`}>
            {trend === 'up' && <TrendingUp size={12} />}
            {trend === 'down' && <TrendingDown size={12} />}
            {trend === 'neutral' && <Minus size={12} />}
            <span>{change != null ? `${change > 0 ? '+' : ''}${change}%` : '—'}</span>
          </div>
        )}
      </div>

      <div className="metric-card-value">{loading ? '—' : value}</div>
      <div className="metric-card-title">{title}</div>

      {sparkline && sparkline.length > 0 && (
        <div className="metric-card-sparkline">
          <svg viewBox="0 0 100 30" preserveAspectRatio="none">
            <defs>
              <linearGradient id={`spark-${title.replace(/\s/g, '')}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="var(--accent)" stopOpacity="0.3" />
                <stop offset="100%" stopColor="var(--accent)" stopOpacity="0" />
              </linearGradient>
            </defs>
            <path
              d={`M ${sparkline.map((v, i) => `${(i / (sparkline.length - 1)) * 100},${30 - (v / Math.max(...sparkline)) * 28}`).join(' L ')} L 100,30 L 0,30 Z`}
              fill={`url(#spark-${title.replace(/\s/g, '')})`}
            />
            <polyline
              points={sparkline.map((v, i) => `${(i / (sparkline.length - 1)) * 100},${30 - (v / Math.max(...sparkline)) * 28}`).join(' ')}
              fill="none"
              stroke="var(--accent)"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </div>
      )}

      {changeLabel && (
        <div className="metric-card-change-label">{changeLabel}</div>
      )}
    </div>
  );
});
