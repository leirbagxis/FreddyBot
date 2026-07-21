import { useState, useCallback, ReactNode } from 'react';
import { AdminSidebar } from './AdminSidebar';
import { AdminTopbar } from './AdminTopbar';
import { AdminTabId } from '../../App';

interface AdminLayoutProps {
  activeTab: AdminTabId;
  onTabChange: (tab: AdminTabId) => void;
  children: ReactNode;
  adminName?: string;
  adminAvatar?: string;
}

export function AdminLayout({ activeTab, onTabChange, children, adminName, adminAvatar }: AdminLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [mobileOpen, setMobileOpen] = useState(false);

  const toggleSidebar = useCallback(() => {
    setMobileOpen(prev => !prev);
  }, []);

  return (
    <div className="admin-layout-v2" data-theme={document.documentElement.getAttribute('data-theme') || 'dark'}>
      {/* Sidebar */}
      <AdminSidebar
        activeTab={activeTab}
        onTabChange={(tab) => {
          onTabChange(tab);
          setMobileOpen(false);
        }}
        isOpen={sidebarOpen}
        onToggle={() => setSidebarOpen(prev => !prev)}
        mobileOpen={mobileOpen}
        onMobileClose={() => setMobileOpen(false)}
      />

      {/* Mobile overlay */}
      {mobileOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm lg:hidden"
          onClick={() => setMobileOpen(false)}
        />
      )}

      {/* Main content */}
      <div className={`admin-main-content ${sidebarOpen ? '' : 'sidebar-collapsed'}`}>
        <AdminTopbar
          onMenuToggle={toggleSidebar}
          adminName={adminName}
          adminAvatar={adminAvatar}
        />
        <main className="admin-page-content">
          {children}
        </main>
      </div>
    </div>
  );
}
