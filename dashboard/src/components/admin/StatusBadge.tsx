import { memo } from 'react';

type BadgeVariant = 'success' | 'warning' | 'danger' | 'accent' | 'default';

interface StatusBadgeProps {
  label: string;
  variant?: BadgeVariant;
  dot?: boolean;
  size?: 'sm' | 'md';
}

const VARIANT_CLASSES: Record<BadgeVariant, string> = {
  success: 'badge-success',
  warning: 'badge-warning',
  danger: 'badge-danger',
  accent: 'badge-accent',
  default: 'badge-default',
};

export const StatusBadge = memo(function StatusBadge({
  label,
  variant = 'default',
  dot = false,
  size = 'sm',
}: StatusBadgeProps) {
  return (
    <span className={`status-badge ${VARIANT_CLASSES[variant]} badge-${size}`}>
      {dot && <span className="status-badge-dot" />}
      {label}
    </span>
  );
});
