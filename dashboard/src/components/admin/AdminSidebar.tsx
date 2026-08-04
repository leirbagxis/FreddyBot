import { memo, ReactNode } from 'react';
import {
  LayoutDashboard, Users, Hash, MessageSquare, FileClock,
  Settings, Zap, Crown, Star, Smartphone, ChevronLeft,
  ChevronRight
} from 'lucide-react';
import { AdminTabId } from '../../App';

interface SidebarItem {
  id: AdminTabId;
  label: string;
  icon: ReactNode;
  section?: string;
}

const SIDEBAR_ITEMS: SidebarItem[] = [
  { id: 'overview', label: 'Visão Geral', icon: <LayoutDashboard size={18} />, section: 'Principal' },
  { id: 'users', label: 'Usuários', icon: <Users size={18} />, section: 'Principal' },
  { id: 'channels', label: 'Canais', icon: <Hash size={18} />, section: 'Principal' },
  { id: 'notice', label: 'Broadcast', icon: <MessageSquare size={18} />, section: 'Operações' },
  { id: 'audit', label: 'Auditoria', icon: <Zap size={18} />, section: 'Operações' },
  { id: 'logs', label: 'Logs', icon: <FileClock size={18} />, section: 'Operações' },
  { id: 'accounts', label: 'Contas MTProto', icon: <Smartphone size={18} />, section: 'Premium' },
  { id: 'premium-features', label: 'Features', icon: <Crown size={18} />, section: 'Premium' },
  { id: 'subscriptions', label: 'Assinaturas', icon: <Star size={18} />, section: 'Premium' },
  { id: 'config', label: 'Configurações', icon: <Settings size={18} />, section: 'Sistema' },
];

interface AdminSidebarProps {
  activeTab: AdminTabId;
  onTabChange: (tab: AdminTabId) => void;
  isOpen: boolean;
  onToggle: () => void;
  mobileOpen: boolean;
}

export const AdminSidebar = memo(function AdminSidebar({
  activeTab, onTabChange, isOpen, onToggle, mobileOpen
}: AdminSidebarProps) {
  // Group items by section
  const sections = SIDEBAR_ITEMS.reduce((acc, item) => {
    const section = item.section || 'Outros';
    if (!acc[section]) acc[section] = [];
    acc[section].push(item);
    return acc;
  }, {} as Record<string, SidebarItem[]>);

  return (
    <aside className={`admin-sidebar ${isOpen ? 'open' : 'collapsed'} ${mobileOpen ? 'mobile-open' : ''}`}>
      <div className="sidebar-logo">
        {isOpen && (
          <span className="sidebar-logo-name">FreddyBot</span>
        )}
        {!isOpen && <span className="sidebar-logo-mark">F</span>}
      </div>

      <nav className="sidebar-nav">
        {Object.entries(sections).map(([sectionName, items]) => (
          <div key={sectionName} className="sidebar-section">
            {isOpen && sectionName !== 'Principal' && (
              <div className="sidebar-section-label">{sectionName}</div>
            )}
            {items.map((item) => (
              <button
                type="button"
                key={item.id}
                onClick={() => onTabChange(item.id)}
                className={`sidebar-nav-item ${activeTab === item.id ? 'active' : ''}`}
                title={!isOpen ? item.label : undefined}
                aria-current={activeTab === item.id ? 'page' : undefined}
              >
                <span className="sidebar-nav-icon">{item.icon}</span>
                {isOpen && (
                  <span className="sidebar-nav-label">{item.label}</span>
                )}
              </button>
            ))}
          </div>
        ))}
      </nav>

      <button
        type="button"
        className="sidebar-collapse-btn"
        onClick={onToggle}
        aria-label={isOpen ? 'Recolher menu' : 'Expandir menu'}
      >
        {isOpen ? <ChevronLeft size={16} /> : <ChevronRight size={16} />}
      </button>
    </aside>
  );
});
