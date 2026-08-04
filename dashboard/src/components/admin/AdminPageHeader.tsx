import { ReactNode } from 'react';

interface AdminPageHeaderProps {
  eyebrow?: string;
  title: string;
  description?: string;
  actions?: ReactNode;
}

export function AdminPageHeader({ eyebrow = 'Painel administrativo', title, description, actions }: AdminPageHeaderProps) {
  return (
    <header className="minimal-page-header">
      <div className="minimal-page-header-copy">
        <span className="minimal-page-eyebrow">{eyebrow}</span>
        <div className="minimal-page-title">{title}</div>
        {description && <p>{description}</p>}
      </div>
      {actions && <div className="minimal-page-header-actions">{actions}</div>}
    </header>
  );
}
