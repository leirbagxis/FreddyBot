import React, { useState, useEffect } from 'react';
import { useParams, useLocation, useNavigate } from 'react-router-dom';
import { Activity } from 'lucide-react';
import { adminFetchUsers } from '../api';
import { terminal } from '../components/Terminal';
import { showToast } from '../components/Toast';
import { NAV_ITEMS } from '../components/admin/constants';
import { User } from '../types';

// Sub-components
import Sidebar from '../components/admin/Sidebar';
import UsersTab from '../components/admin/UsersTab';

export default function Admin() {
  const { botId } = useParams<{ botId: string }>();
  const location = useLocation();
  const navigate = useNavigate();
  
  // Data State
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [forbidden, setForbidden] = useState(false);
  
  // UI State
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedUserId, setExpandedUserId] = useState<string | null>(null);
  
  const userRole = localStorage.getItem(`bot_${botId}_role`) || 'admin';
  const query = new URLSearchParams(location.search);
  const activeTab = (query.get('tab') as 'users') || 'users';

  useEffect(() => {
    const handleToggle = () => setIsSidebarOpen(prev => !prev);
    window.addEventListener('toggle-admin-sidebar', handleToggle);
    return () => window.removeEventListener('toggle-admin-sidebar', handleToggle);
  }, []);

  const loadData = async () => {
    setIsLoading(true);
    setForbidden(false);
    try {
      const usersData = await adminFetchUsers();
      setUsers(usersData || []);
      terminal.log(`Admin: Sincronização OK`);
    } catch (err: any) {
      console.error("Fetch error", err);
      if (err.response?.status === 403) {
        setForbidden(true);
      } else {
        showToast.error("ERRO DE SINCRONIZAÇÃO");
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => { 
    loadData(); 
  }, []);

  if (forbidden) {
    return (
      <div className="p-8 flex flex-col items-center justify-center min-h-[60vh]">
        <div className="bg-panel border border-red-500/20 rounded-[2rem] shadow-2xl p-12 text-center max-w-md animate-fade-in duration-700">
          <Activity size={40} className="text-accent mx-auto mb-6 animate-pulse" />
          <h2 className="text-2xl font-black uppercase mb-2 tracking-tighter text-white">ACESSO_NEGADO</h2>
          <p className="text-muted text-[10px] font-bold uppercase tracking-widest mb-8">Nível de autorização insuficiente para Super Admin.</p>
          <button 
            onClick={() => navigate(`/${botId}`)} 
            className="btn-primary w-full py-4 bg-accent text-white shadow-accent-neon hover:shadow-accent-neon hover:scale-[1.02]"
          >
            VOLTAR_PARA_BASE
          </button>
        </div>
      </div>
    );
  }

  const filteredUsers = users.filter(u => 
    u.username?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    u.first_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    u.id?.toString().includes(searchQuery)
  );

  return (
    <div className="min-h-screen bg-surface font-sans selection:bg-primary selection:text-black relative overflow-x-hidden">
      <Sidebar
        botId={botId!}
        activeTab={activeTab}
        isSidebarOpen={isSidebarOpen}
        setIsSidebarOpen={setIsSidebarOpen}
        userRole={userRole}
      />
      {isSidebarOpen && (
        <div className="fixed inset-0 z-[60] bg-black/60 backdrop-blur-sm transition-opacity duration-300" onClick={() => setIsSidebarOpen(false)} />
      )}

      <main className="w-full max-w-[100vw] overflow-x-hidden p-4 lg:p-10">
        <div className="max-w-[1300px] mx-auto space-y-8 lg:space-y-12">
          <div className="hidden lg:flex items-center justify-between mb-6 animate-in fade-in duration-700">
            <header>
              <div className="flex items-center gap-2 mb-1">
                <span className="text-[10px] font-bold uppercase tracking-[0.2em] text-primary mb-xs block mb-0">System Control v2.0</span>
                <div className="w-1.5 h-1.5 rounded-full bg-primary animate-pulse" />
              </div>
              <h2 className="text-3xl font-black uppercase tracking-tighter text-white">
                {NAV_ITEMS.find(n => n.id === activeTab)?.label || 'Painel Admin'}
              </h2>
            </header>
            <div className="w-12 h-12 rounded-2xl bg-primary flex items-center justify-center text-black shadow-neon transition-all duration-500 hover:rotate-12 hover:scale-110 group cursor-pointer">
              <Activity size={24} className="group-hover:animate-pulse" />
            </div>
          </div>

          <div className="animate-in fade-in slide-in-from-bottom-2 duration-500">
            {isLoading && users.length === 0 ? (
               <div className="flex flex-col items-center justify-center py-20">
                 <div className="w-10 h-10 border-4 border-white/5 border-t-primary rounded-full animate-spin mb-4"></div>
                 <div className="text-[10px] font-black uppercase tracking-widest text-primary animate-pulse">Scanning_Network...</div>
               </div>
            ) : (
              <UsersTab 
                searchQuery={searchQuery} 
                setSearchQuery={setSearchQuery} 
                filteredUsers={filteredUsers} 
                expandedUserId={expandedUserId} 
                setExpandedUserId={setExpandedUserId} 
              />
            )}
          </div>
        </div>
      </main>

      <style>{`
        .animate-fade-in { animation: fade-in 0.5s ease-out; }
        @keyframes fade-in { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
        .no-scrollbar::-webkit-scrollbar { display: none; }
        .no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
      `}</style>
    </div>
  );
}
