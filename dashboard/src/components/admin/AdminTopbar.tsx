import { useEffect, useMemo, useRef, useState, memo } from 'react';
import { Bell, Menu, Search, User, ArrowRight, ArrowUpDown, SlidersHorizontal, Plus } from 'lucide-react';
import { AdminTabId } from '../../App';
import { AdminCrmFilter, AdminCrmSort, useAdminCrmControls } from './AdminCrmContext';
import { Channel, User as UserData } from '../../types';
import { getOperationalAlerts, OperationalAlertKind } from './crmSelectors';

interface AdminTopbarProps {
  activeTab: AdminTabId;
  onMenuToggle: () => void;
  onNavigate: (tab: AdminTabId) => void;
  adminName?: string;
  adminAvatar?: string;
  users: UserData[];
  channels: Channel[];
}

export const AdminTopbar = memo(function AdminTopbar({
  activeTab, onMenuToggle, onNavigate, adminName, adminAvatar, users, channels
}: AdminTopbarProps) {
  const { searchQuery, setSearchQuery, sortBy, setSortBy, filterBy, setFilterBy } = useAdminCrmControls();
  const [searchFocused, setSearchFocused] = useState(false);
  const [notificationsOpen, setNotificationsOpen] = useState(false);
  const searchRef = useRef<HTMLInputElement>(null);
  const notificationsRef = useRef<HTMLDivElement>(null);
  const supportsCustomerControls = activeTab === 'overview' || activeTab === 'users';
  const supportsSearch = supportsCustomerControls || activeTab === 'channels';
  const alerts = useMemo(() => getOperationalAlerts(users, channels), [channels, users]);

  useEffect(() => {
    const handleShortcut = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k' && supportsSearch) {
        event.preventDefault();
        searchRef.current?.focus();
      }
    };
    window.addEventListener('keydown', handleShortcut);
    return () => window.removeEventListener('keydown', handleShortcut);
  }, [supportsSearch]);

  useEffect(() => {
    const closeOnOutsideInteraction = (event: PointerEvent) => {
      if (notificationsRef.current && !notificationsRef.current.contains(event.target as Node)) setNotificationsOpen(false);
    };
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setNotificationsOpen(false);
    };
    document.addEventListener('pointerdown', closeOnOutsideInteraction);
    document.addEventListener('keydown', closeOnEscape);
    return () => {
      document.removeEventListener('pointerdown', closeOnOutsideInteraction);
      document.removeEventListener('keydown', closeOnEscape);
    };
  }, []);

  const handleAlert = (kind: OperationalAlertKind) => {
    if (kind === 'blacklisted') setFilterBy('blacklisted');
    if (kind === 'without-channels') setFilterBy('without-channels');
    if (kind === 'new-users') setFilterBy('all');
    onNavigate('users');
    setNotificationsOpen(false);
  };

  return (
    <header className="admin-topbar">
      <div className="admin-topbar-left">
        <button type="button" className="admin-topbar-menu-btn" onClick={onMenuToggle} aria-label="Abrir menu administrativo">
          <Menu size={18} />
        </button>

        {supportsSearch ? (
          <label className={`admin-topbar-search ${searchFocused ? 'focused' : ''}`}>
            <Search size={15} className="admin-topbar-search-icon" aria-hidden="true" />
            <input
              ref={searchRef}
              type="search"
              placeholder={activeTab === 'channels' ? 'Buscar canal...' : 'Buscar usuário...'}
              value={searchQuery}
              onChange={(event) => setSearchQuery(event.target.value)}
              onFocus={() => setSearchFocused(true)}
              onBlur={() => setSearchFocused(false)}
              className="admin-topbar-search-input"
              aria-label={activeTab === 'channels' ? 'Buscar canais' : 'Buscar usuários'}
            />
          </label>
        ) : (
          <span className="admin-topbar-context">Administração</span>
        )}
      </div>

      <div className="admin-topbar-right">
        {supportsSearch && (
          <label className="admin-topbar-select-control">
            <ArrowUpDown size={15} aria-hidden="true" />
            <span>Ordenar</span>
            <select value={sortBy} onChange={(event) => setSortBy(event.target.value as AdminCrmSort)} aria-label="Ordenar resultados">
              <option value="recent">Recentes</option>
              <option value="name">Nome</option>
              {supportsCustomerControls && <option value="channels">Mais canais</option>}
            </select>
          </label>
        )}

        {supportsCustomerControls && (
          <label className="admin-topbar-select-control">
            <SlidersHorizontal size={15} aria-hidden="true" />
            <span>Filtros</span>
            <select value={filterBy} onChange={(event) => setFilterBy(event.target.value as AdminCrmFilter)} aria-label="Filtrar usuários">
              <option value="all">Todos</option>
              <option value="admins">Admins</option>
              <option value="blacklisted">Blacklist</option>
              <option value="with-channels">Com canais</option>
              <option value="without-channels">Sem canais</option>
            </select>
          </label>
        )}

        <div className="admin-notifications" ref={notificationsRef}>
          <button
            type="button"
            className="admin-notifications-trigger"
            onClick={() => setNotificationsOpen((open) => !open)}
            aria-label={`${alerts.length} notificações operacionais`}
            aria-expanded={notificationsOpen}
            aria-haspopup="dialog"
          >
            <Bell size={16} aria-hidden="true" />
            {alerts.length > 0 && <span>{alerts.length > 9 ? '9+' : alerts.length}</span>}
          </button>
          {notificationsOpen && (
            <section className="admin-notifications-panel" role="dialog" aria-label="Notificações operacionais">
              <header>
                <div><span>Monitoramento</span><strong>Notificações</strong></div>
                <small>{alerts.length ? `${alerts.length} ativas` : 'Tudo em ordem'}</small>
              </header>
              {alerts.length ? (
                <div>
                  {alerts.map((alert) => (
                    <button type="button" key={alert.id} onClick={() => handleAlert(alert.id)}>
                      <b>{alert.count}</b>
                      <span><strong>{alert.title}</strong><small>{alert.description}</small></span>
                      <ArrowRight size={14} aria-hidden="true" />
                    </button>
                  ))}
                </div>
              ) : <p>Nenhuma ação pendente na base atual.</p>}
            </section>
          )}
        </div>

        <div className="admin-topbar-profile" aria-label={`Administrador ${adminName || 'Admin'}`} title={adminName || 'Admin'}>
          {adminAvatar ? (
            <img src={adminAvatar} alt="" className="admin-topbar-avatar-img" referrerPolicy="no-referrer" />
          ) : (
            <User size={15} aria-hidden="true" />
          )}
          <span>Eu</span>
        </div>

        <button type="button" className="admin-topbar-primary" onClick={() => onNavigate('notice')}>
          <Plus size={15} aria-hidden="true" />
          <span>Novo broadcast</span>
        </button>
      </div>
    </header>
  );
});
