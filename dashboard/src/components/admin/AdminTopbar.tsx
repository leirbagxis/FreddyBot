import { useEffect, useRef, useState, memo } from 'react';
import { Menu, Search, User, ArrowUpDown, SlidersHorizontal, Plus } from 'lucide-react';
import { AdminTabId } from '../../App';
import { AdminCrmFilter, AdminCrmSort, useAdminCrmControls } from './AdminCrmContext';
import { Channel, User as UserData } from '../../types';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

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
  activeTab, onMenuToggle, onNavigate, adminName, adminAvatar
}: AdminTopbarProps) {
  const { searchQuery, setSearchQuery, sortBy, setSortBy, filterBy, setFilterBy } = useAdminCrmControls();
  const [searchFocused, setSearchFocused] = useState(false);
  const searchRef = useRef<HTMLInputElement>(null);
  const supportsCustomerControls = activeTab === 'overview' || activeTab === 'users';

  useEffect(() => {
    const handleShortcut = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
        event.preventDefault();
        searchRef.current?.focus();
      }
    };
    window.addEventListener('keydown', handleShortcut);
    return () => window.removeEventListener('keydown', handleShortcut);
  }, []);

  return (
    <header className="admin-topbar">
      <div className="admin-topbar-left">
        <Button variant="ghost" size="icon" type="button" className="admin-topbar-menu-btn" onClick={onMenuToggle} aria-label="Abrir menu administrativo">
          <Menu size={18} />
        </Button>

        <div className={`admin-topbar-search ${searchFocused ? 'focused' : ''}`}>
          <Search size={15} className="admin-topbar-search-icon" aria-hidden="true" />
          <Input
            ref={searchRef}
            type="search"
            placeholder={activeTab === 'channels' ? 'Buscar canal...' : activeTab === 'users' ? 'Buscar usuário...' : 'Buscar no sistema...'}
            value={searchQuery}
            onChange={(event) => setSearchQuery(event.target.value)}
            onFocus={() => setSearchFocused(true)}
            onBlur={() => setSearchFocused(false)}
            className="admin-topbar-search-input"
            aria-label="Buscar no sistema"
          />
        </div>
      </div>

      <div className="admin-topbar-right">
        {supportsCustomerControls && (
          <label className="admin-topbar-select-control">
            <ArrowUpDown size={15} aria-hidden="true" />
            <span>Ordenar</span>
            <select value={sortBy} onChange={(event) => setSortBy(event.target.value as AdminCrmSort)} aria-label="Ordenar resultados">
              <option value="recent">Recentes</option>
              <option value="name">Nome</option>
              <option value="channels">Mais canais</option>
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

        <div className="admin-topbar-profile" aria-label={`Administrador ${adminName || 'Admin'}`} title={adminName || 'Admin'}>
          {adminAvatar ? (
            <img src={adminAvatar} alt="" className="admin-topbar-avatar-img" referrerPolicy="no-referrer" />
          ) : (
            <User size={15} aria-hidden="true" />
          )}
          <span>Eu</span>
        </div>

        <Button type="button" className="admin-topbar-primary" onClick={() => onNavigate('notice')}>
          <Plus size={15} aria-hidden="true" />
          <span>Novo broadcast</span>
        </Button>
      </div>
    </header>
  );
});
