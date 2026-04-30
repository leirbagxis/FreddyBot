import { type ReactNode } from 'react';

export interface Tab {
  id: string;
  label: string;
  icon: ReactNode;
}

interface Props {
  tabs: Tab[];
  activeTab: string;
  onTabChange: (id: string) => void;
}

export function TabBar({ tabs, activeTab, onTabChange }: Props) {
  return (
    <nav className="bottom-nav">
      {tabs.map((t) => (
        <button
          key={t.id}
          className={`nav-item ${activeTab === t.id ? 'active' : ''}`}
          onClick={() => onTabChange(t.id)}
        >
          {activeTab === t.id && <span className="nav-dot" />}
          {t.icon}
          <span>{t.label}</span>
        </button>
      ))}
    </nav>
  );
}
