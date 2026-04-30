import { Tab } from './TabBar';

interface Props {
  tabs: Tab[];
  activeTab: string;
  onTabChange: (id: string) => void;
  isCollapsed: boolean; // Aqui usaremos isCollapsed para controlar se está escondida ou não
}

export function AdminSidebar({ tabs, activeTab, onTabChange, isCollapsed }: Props) {
  return (
    <aside className={`admin-sidebar ${isCollapsed ? 'hidden' : 'active'}`}>
      <div className="sidebar-nav">
        {tabs.map((t) => (
          <button
            key={t.id}
            className={`sidebar-nav-item ${activeTab === t.id ? 'active' : ''}`}
            onClick={() => onTabChange(t.id)}
            title={t.label}
          >
            <div className="sidebar-nav-icon">{t.icon}</div>
          </button>
        ))}
      </div>
    </aside>
  );
}
