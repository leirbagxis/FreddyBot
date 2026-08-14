import { useState, useCallback, ReactNode } from 'react';
import { AdminSidebar } from './AdminSidebar';
import { AdminTopbar } from './AdminTopbar';
import { AdminTabId } from '../../App';
import { AdminCrmProvider } from './AdminCrmContext';
import { Channel, User } from '../../types';

interface AdminLayoutProps {
  activeTab: AdminTabId;
  onTabChange: (tab: AdminTabId) => void;
  children: ReactNode;
  adminName?: string;
  adminAvatar?: string;
  users: User[];
  channels: Channel[];
}

export function AdminLayout({ activeTab, onTabChange, children, adminName, adminAvatar, users, channels }: AdminLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [mobileOpen, setMobileOpen] = useState(false);

  const toggleSidebar = useCallback(() => {
    if (window.matchMedia('(max-width: 960px)').matches) {
      setMobileOpen(prev => !prev);
    } else {
      setSidebarOpen(prev => !prev);
    }
  }, []);

  const handleTabChange = useCallback((tab: AdminTabId) => {
    onTabChange(tab);
    setMobileOpen(false);

    const url = new URL(window.location.href);
    if (tab === 'overview') url.searchParams.delete('tab');
    else url.searchParams.set('tab', tab);
    if (tab !== 'logs') url.searchParams.delete('channelId');
    window.history.replaceState({}, '', `${url.pathname}${url.search}${url.hash}`);
  }, [onTabChange]);

  return (
    <AdminCrmProvider onNavigate={handleTabChange}>
      <div className="admin-layout-v2">
        <AdminSidebar
          activeTab={activeTab}
          onTabChange={handleTabChange}
          isOpen={sidebarOpen}
          onToggle={() => setSidebarOpen(prev => !prev)}
          mobileOpen={mobileOpen}
          adminName={adminName}
          adminAvatar={adminAvatar}
        />

        {mobileOpen && (
          <button
            type="button"
            className="admin-mobile-overlay"
            onClick={() => setMobileOpen(false)}
            aria-label="Fechar menu administrativo"
          />
        )}

        <div className={`admin-main-content ${sidebarOpen ? '' : 'sidebar-collapsed'}`}>
          <AdminTopbar
            activeTab={activeTab}
            onMenuToggle={toggleSidebar}
            onNavigate={handleTabChange}
            adminName={adminName}
            adminAvatar={adminAvatar}
            users={users}
            channels={channels}
          />
          <main className="admin-page-content">
            {children}
          </main>
        </div>
      </div>
    </AdminCrmProvider>
  );
}
