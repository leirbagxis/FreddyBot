import { Link, useNavigate } from 'react-router-dom';
import { Shield } from 'lucide-react';
import { NAV_ITEMS } from './constants';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: any[]) {
  return twMerge(clsx(inputs));
}

interface SidebarProps {
  botId: string;
  activeTab: string;
  isSidebarOpen: boolean;
  setIsSidebarOpen: (open: boolean) => void;
  userRole: string;
}

export default function Sidebar({
  botId,
  activeTab,
  isSidebarOpen,
  setIsSidebarOpen,
  userRole,
}: SidebarProps) {
  const navigate = useNavigate();

  return (
    <aside className={cn(
      "fixed inset-y-0 left-0 z-[70] w-64 bg-black text-white will-change-transform transition-transform duration-500 ease-out",
      isSidebarOpen ? "translate-x-0 shadow-2xl" : "-translate-x-full"
    )}>
      <div className="h-full flex flex-col p-6 overflow-y-auto no-scrollbar">
        <div className="mb-10">
          <Link to={`/${botId}`} className="flex items-center gap-3 mb-8 group" onClick={() => setIsSidebarOpen(false)}>
            <div className="w-10 h-10 bg-primary text-black rounded-xl flex items-center justify-center font-black text-xl transition-all duration-500 group-hover:rotate-6 group-hover:scale-110 shadow-lg">F</div>
            <div className="transition-transform duration-500 group-hover:translate-x-1">
              <h1 className="font-black text-base leading-none tracking-tighter normal-case">Freddy Bot</h1>
              <span className="text-[9px] font-bold text-white/40 uppercase tracking-widest">Painel Admin</span>
            </div>
          </Link>

          <nav className="space-y-1">
            {NAV_ITEMS.map((item) => (
              <button
                key={item.id}
                onClick={() => {
                  navigate(`/${botId}/admin?tab=${item.id}`);
                  setIsSidebarOpen(false);
                }}
                className={cn(
                  "w-full flex items-center gap-3 px-5 py-3.5 rounded-xl text-[10px] font-black tracking-widest transition-all duration-300 group relative",
                  activeTab === item.id 
                    ? "bg-primary text-black shadow-lg" 
                    : "text-white/40 hover:text-white hover:bg-white/5"
                )}
              >
                <item.icon size={16} className={cn("transition-all duration-500", activeTab === item.id ? "animate-pulse" : "group-hover:scale-110 group-hover:rotate-3")} />
                {item.label}
              </button>
            ))}
          </nav>
        </div>

        <div className="mt-auto pt-6 border-t border-white/10 space-y-3">
          <div className="bg-white/5 p-4 rounded-2xl border border-white/5 transition-colors duration-300 hover:bg-white/10">
            <span className="status-label text-white/40 mb-1 text-[8px] font-bold uppercase tracking-[0.15em] block">Privilégio</span>
            <div className="text-[10px] font-black text-white flex items-center gap-2 uppercase tracking-tighter">
              <Shield size={12} className={cn("transition-transform duration-700", userRole === 'owner' ? "text-primary animate-bounce" : "text-white/60")} />
              {userRole === 'owner' ? 'Proprietário' : 'Administrador'}
            </div>
          </div>
        </div>
      </div>
    </aside>
  );
}
