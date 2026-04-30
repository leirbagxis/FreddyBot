import { useState, Dispatch, SetStateAction } from 'react';
import { AdminDashboardData, User, Channel } from '../types';
import { AdminNoticeTab } from './AdminNoticeTab';
import { AdminConfigTab } from './AdminConfigTab';
import { NoticeButton } from '../api';
import { Users, Hash, Search, ArrowLeft, ChevronRight, User as UserIcon, Settings } from 'lucide-react';

interface AdminDashboardProps {
  adminData: AdminDashboardData;
  activeTab: 'users' | 'channels' | 'notice' | 'config';
  navigateToChannel: (id: number) => void;
  selectedUserId: number | null;
  onSelectUser: (id: number | null) => void;
  // Notice tab props
  noticeMessage: string;
  setNoticeMessage: Dispatch<SetStateAction<string>>;
  noticeTarget: 'channels' | 'users' | 'all';
  setNoticeTarget: Dispatch<SetStateAction<'channels' | 'users' | 'all'>>;
  noticeButtons: NoticeButton[];
  handleAddNoticeButton: () => void;
  updateNoticeButton: (index: number, field: keyof NoticeButton, value: string) => void;
  removeNoticeButton: (index: number) => void;
  handleSendNotice: () => void;
  isSendingNotice: boolean;
}

export function AdminDashboard({
  adminData,
  activeTab,
  navigateToChannel,
  selectedUserId,
  onSelectUser,
  noticeMessage, setNoticeMessage,
  noticeTarget, setNoticeTarget,
  noticeButtons, handleAddNoticeButton,
  updateNoticeButton, removeNoticeButton,
  handleSendNotice,
  isSendingNotice
}: AdminDashboardProps) {
  const [adminSearch, setAdminSearch] = useState('');
  const [adminChannelCountFilter, setAdminChannelCountFilter] = useState('');
  const [visibleChannelsCount, setVisibleChannelsCount] = useState(40);
  const [visibleUsersCount, setVisibleUsersCount] = useState(40);

  // Derivar o usuário selecionado das props
  const usersList = adminData.users || [];
  const adminSelectedUser = selectedUserId ? usersList.find(u => u.id === selectedUserId) : null;
  const setAdminSelectedUser = (user: User | null) => onSelectUser(user ? user.id : null);

  const [channelSearch, setChannelSearch] = useState('');

  const renderUserDetail = () => {
    if (!adminSelectedUser) return null;
    const name = adminSelectedUser.firstName || (adminSelectedUser as any).first_name || 'Sem nome';
    return (
      <>
        <button
          onClick={() => setAdminSelectedUser(null)}
          className="mb-2 flex items-center text-sm font-medium transition-colors"
          style={{ color: 'var(--hint)' }}
        >
          <ArrowLeft size={16} className="mr-1.5" /> Voltar para usuários
        </button>

        <div className="admin-welcome-card mb-4">
          <div className="flex items-center">
            <div className="section-icon purple mr-3"><UserIcon size={24} /></div>
            <div className="min-w-0 flex-1">
              <h2 className="text-[16px] font-bold truncate text-[var(--text)]">{name}</h2>
              <p className="text-xs truncate" style={{ color: 'var(--hint)' }}>ID: {adminSelectedUser.id}</p>
            </div>
          </div>
        </div>

        <div className="space-y-3">
          <h3 className="text-sm font-semibold mb-2" style={{ color: 'var(--hint)' }}>Canais do Usuário</h3>
          {adminSelectedUser.channels && adminSelectedUser.channels.length > 0 ? (
            adminSelectedUser.channels.map((c: Channel) => (
            <button key={c.id} className="admin-list-item flex items-center w-full text-left p-4" onClick={() => navigateToChannel(c.id)}>
              <div className="section-icon purple mr-3"><Hash size={20} /></div>
              <div className="min-w-0 flex-1">
                <h3 className="text-[15px] font-semibold truncate">{c.title}</h3>
                <p className="text-xs truncate mt-0.5" style={{ color: 'var(--hint)' }}>ID: {c.id}</p>
              </div>
              <ChevronRight size={18} className="stat-arrow" />
            </button>
          ))
          ) : (
            <div className="card" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', padding: '32px 16px', color: 'var(--hint)' }}>
              <Hash size={32} style={{ opacity: 0.5, marginBottom: 12 }} />
              <p className="text-sm font-medium">Este usuário não possui canais cadastrados</p>
            </div>
          )}
        </div>
      </>
    );
  };

  const renderUsersTab = () => {
    const filteredUsers = usersList.filter(u => {
      const name = (u.firstName || (u as any).first_name || '').toLowerCase();
      const matchesSearch = name.includes(adminSearch.toLowerCase()) || u.id.toString().includes(adminSearch);
      const matchesCount = adminChannelCountFilter ? (u.channels?.length || 0) === parseInt(adminChannelCountFilter, 10) : true;
      return matchesSearch && matchesCount;
    });

    const visibleUsers = filteredUsers.slice(0, visibleUsersCount);

    return (
      <>
        <div className="grid grid-cols-2 gap-3 mb-4">
          <div className="admin-stat-card">
            <div className="admin-stat-icon-glow" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
              <Users size={24} />
            </div>
            <span className="admin-stat-value">{usersList.length}</span>
            <span className="admin-stat-label">Usuários</span>
          </div>
          <div className="admin-stat-card" style={{ animationDelay: '0.1s' }}>
            <div className="admin-stat-icon-glow" style={{ background: 'var(--success-soft)', color: 'var(--success)' }}>
               <Hash size={24} />
            </div>
            <span className="admin-stat-value">
              {adminData.channels?.length || usersList.reduce((acc, u) => acc + (u.channels?.length || 0), 0) || 0}
            </span>
            <span className="admin-stat-label">Canais Totais</span>
          </div>
        </div>

        <div className="flex flex-col gap-2 mb-4">
          <div className="search-bar-container relative">
            <Search className="absolute left-4 top-1/2 transform -translate-y-1/2" size={18} style={{ color: 'var(--hint)' }} />
            <input
              type="text"
              placeholder="Buscar usuário por nome ou ID..."
              className="admin-search-input input"
              value={adminSearch}
              onChange={(e) => {
                setAdminSearch(e.target.value);
                setVisibleUsersCount(40);
              }}
            />
          </div>
          <div className="search-bar-container relative">
            <Hash className="absolute left-4 top-1/2 transform -translate-y-1/2" size={18} style={{ color: 'var(--hint)' }} />
            <input
              type="number"
              placeholder="Filtrar por qtd. de canais"
              className="admin-search-input input"
              value={adminChannelCountFilter}
              onChange={(e) => {
                setAdminChannelCountFilter(e.target.value);
                setVisibleUsersCount(40);
              }}
            />
          </div>
        </div>

        <div className="space-y-3">
          <h3 className="text-sm font-semibold mb-2" style={{ color: 'var(--hint)' }}>Usuários ({filteredUsers.length})</h3>
          {visibleUsers.length > 0 ? (
             <>
             {visibleUsers.map((u) => (
              <button key={u.id} className="admin-list-item flex items-center w-full text-left p-4" onClick={() => setAdminSelectedUser(u)}>
                <div className="section-icon purple mr-3"><UserIcon size={20} /></div>
                <div className="min-w-0 flex-1">
                  <h3 className="text-[15px] font-semibold truncate">{u.firstName || (u as any).first_name || 'Sem nome'}</h3>
                  <p className="text-xs truncate mt-0.5" style={{ color: 'var(--hint)' }}>ID: {u.id} • {u.channels?.length || 0} canais</p>
                </div>
                <ChevronRight size={18} className="stat-arrow" />
              </button>
            ))}
            {filteredUsers.length > visibleUsersCount && (
              <button 
                className="btn btn-secondary w-full py-3" 
                onClick={() => setVisibleUsersCount(prev => prev + 40)}
              >
                Carregar mais usuários...
              </button>
            )}
            </>
          ) : (
            <div className="card" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', padding: '32px 16px', color: 'var(--hint)' }}>
              <UserIcon size={32} style={{ opacity: 0.5, marginBottom: 12 }} />
              <p className="text-sm font-medium">Nenhum usuário encontrado</p>
            </div>
          )}
        </div>
      </>
    );
  };

  const renderChannelsTab = () => {
    const filteredChannels = (adminData.channels || []).filter(c => {
      const matchSearch = c.title.toLowerCase().includes(channelSearch.toLowerCase()) || c.id.toString().includes(channelSearch);
      return matchSearch;
    });

    const visibleChannels = filteredChannels.slice(0, visibleChannelsCount);

    return (
      <>
        <div className="flex flex-col gap-2 mb-4">
          <div className="search-bar-container relative">
            <Search className="absolute left-4 top-1/2 transform -translate-y-1/2" size={18} style={{ color: 'var(--hint)' }} />
            <input
              type="text"
              placeholder="Buscar canal por título ou ID..."
              className="admin-search-input input"
              value={channelSearch}
              onChange={(e) => {
                setChannelSearch(e.target.value);
                setVisibleChannelsCount(40);
              }}
            />
          </div>
        </div>

        <div className="space-y-3">
          <h3 className="text-sm font-semibold mb-2" style={{ color: 'var(--hint)' }}>Todos os Canais ({filteredChannels.length})</h3>
          {visibleChannels.length > 0 ? (
            <>
            {visibleChannels.map((c) => (
              <button key={c.id} className="admin-list-item flex items-center w-full text-left p-4" onClick={() => navigateToChannel(c.id)}>
                <div className="section-icon purple mr-3"><Hash size={20} /></div>
                <div className="min-w-0 flex-1">
                  <h3 className="text-[15px] font-semibold truncate">{c.title}</h3>
                  <p className="text-xs truncate mt-0.5" style={{ color: 'var(--hint)' }}>ID: {c.id}</p>
                </div>
                <ChevronRight size={18} className="stat-arrow" />
              </button>
            ))}
            {filteredChannels.length > visibleChannelsCount && (
              <button 
                className="btn btn-secondary w-full py-3" 
                onClick={() => setVisibleChannelsCount(prev => prev + 40)}
              >
                Carregar mais canais...
              </button>
            )}
            </>
          ) : (
             <div className="card" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', padding: '32px 16px', color: 'var(--hint)' }}>
              <Hash size={32} style={{ opacity: 0.5, marginBottom: 12 }} />
              <p className="text-sm font-medium">Nenhum canal encontrado</p>
            </div>
          )}
        </div>
      </>
    );
  };

  const renderNoticeTab = () => {
    return (
      <AdminNoticeTab
        noticeMessage={noticeMessage}
        setNoticeMessage={setNoticeMessage}
        noticeTarget={noticeTarget}
        setNoticeTarget={setNoticeTarget}
        noticeButtons={noticeButtons}
        handleAddNoticeButton={handleAddNoticeButton}
        updateNoticeButton={updateNoticeButton}
        removeNoticeButton={removeNoticeButton}
        handleSendNotice={handleSendNotice}
        isSendingNotice={isSendingNotice}
      />
    );
  }

  return (
    <div className="space-y-4">
      <div className="admin-welcome-card z-10">
        <div className="welcome-greeting relative z-10">
          <span className="welcome-emoji text-3xl">⚙️</span>
          <div className="min-w-0 flex-1">
            <h2 className="welcome-title text-xl mb-1">Painel <span style={{ color: 'var(--accent)' }}>Administrativo</span></h2>
            <p className="welcome-sub text-sm opacity-80">Gerencie usuários, canais e mensagens globais.</p>
          </div>
        </div>
      </div>

      {activeTab === 'users' && !adminSelectedUser && renderUsersTab()}
      {activeTab === 'users' && adminSelectedUser && renderUserDetail()}
      {activeTab === 'channels' && renderChannelsTab()}
      {activeTab === 'notice' && renderNoticeTab()}
      {activeTab === 'config' && <AdminConfigTab />}
    </div>
  );
}
