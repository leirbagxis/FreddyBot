import { useState, memo } from 'react';
import { Menu, Search, Bell, Sun, Moon, User, ChevronDown } from 'lucide-react';
import { useTheme } from '../../hooks/useTheme';

interface AdminTopbarProps {
  onMenuToggle: () => void;
  adminName?: string;
  adminAvatar?: string;
}

export const AdminTopbar = memo(function AdminTopbar({
  onMenuToggle, adminName, adminAvatar
}: AdminTopbarProps) {
  const { theme, toggleTheme } = useTheme();
  const [searchFocused, setSearchFocused] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  return (
    <header className="admin-topbar">
      {/* Left: Menu toggle + Search */}
      <div className="admin-topbar-left">
        <button className="admin-topbar-menu-btn" onClick={onMenuToggle}>
          <Menu size={18} />
        </button>

        <div className={`admin-topbar-search ${searchFocused ? 'focused' : ''}`}>
          <Search size={15} className="admin-topbar-search-icon" />
          <input
            type="text"
            placeholder="Buscar usuários, canais, logs..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onFocus={() => setSearchFocused(true)}
            onBlur={() => setSearchFocused(false)}
            className="admin-topbar-search-input"
          />
          <kbd className="admin-topbar-search-kbd">⌘K</kbd>
        </div>
      </div>

      {/* Right: Actions */}
      <div className="admin-topbar-right">
        {/* Notifications */}
        <button className="admin-topbar-action" title="Notificações">
          <Bell size={17} />
          <span className="admin-topbar-notification-dot" />
        </button>

        {/* Theme toggle */}
        <button className="admin-topbar-action" onClick={toggleTheme} title="Alternar tema">
          {theme === 'dark' ? <Sun size={17} /> : <Moon size={17} />}
        </button>

        {/* Divider */}
        <div className="admin-topbar-divider" />

        {/* Profile */}
        <button className="admin-topbar-profile">
          <div className="admin-topbar-avatar">
            {adminAvatar ? (
              <img src={adminAvatar} alt="" className="admin-topbar-avatar-img" />
            ) : (
              <User size={15} />
            )}
          </div>
          <div className="admin-topbar-profile-info">
            <span className="admin-topbar-profile-name">{adminName || 'Admin'}</span>
            <span className="admin-topbar-profile-role">Owner</span>
          </div>
          <ChevronDown size={14} className="admin-topbar-profile-chevron" />
        </button>
      </div>
    </header>
  );
});
