import { Users, Search, User as UserIcon, Shield, ExternalLink } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { User } from '../../types';

function cn(...inputs: any[]) {
  return twMerge(clsx(inputs));
}

interface UsersTabProps {
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  filteredUsers: User[];
  expandedUserId: string | null;
  setExpandedUserId: (id: string | null) => void;
}

export default function UsersTab({
  searchQuery,
  setSearchQuery,
  filteredUsers,
  expandedUserId,
  setExpandedUserId,
}: UsersTabProps) {
  return (
    <section className="space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-primary/10 rounded-lg text-primary shadow-neon"><Users size={22} className="animate-pulse" /></div>
          <h3 className="text-xl font-black uppercase tracking-tighter text-white">Índice de Usuários</h3>
        </div>
        <div className="relative w-full md:w-80 group">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-muted group-focus-within:text-primary transition-colors" size={16} />
          <input 
            type="text" 
            placeholder="Procurar por Username ou ID..." 
            className="w-full bg-panel border-2 border-white/5 rounded-2xl pl-11 pr-4 py-3 text-[10px] font-black outline-none focus:border-primary focus:shadow-neon text-white transition-all duration-500 placeholder:text-muted" 
            value={searchQuery} 
            onChange={e => setSearchQuery(e.target.value)} 
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredUsers.map((user) => (
          <div 
            key={user.id} 
            onClick={() => setExpandedUserId(expandedUserId === user.id.toString() ? null : user.id.toString())} 
            className={cn(
              "bg-panel rounded-[2.5rem] shadow-xl p-6 lg:p-8 cursor-pointer border-2 transition-all duration-500 relative group/user",
              expandedUserId === user.id.toString() 
                ? "border-primary shadow-neon scale-[1.02] z-10" 
                : "border-white/5 hover:border-primary/20 hover:shadow-2xl hover:-translate-y-1"
            )}
          >
            <div className="flex justify-between items-center gap-3">
              <div className="flex items-center gap-3 min-w-0">
                <div className="w-10 h-10 bg-surface rounded-xl flex items-center justify-center text-muted group-hover/user:bg-primary group-hover/user:text-black transition-all duration-500 shadow-inner">
                  <UserIcon size={20} />
                </div>
                <div className="min-w-0">
                  <div className="text-[11px] font-black uppercase truncate text-white leading-none">
                    {user.first_name} {user.is_admin && <Shield size={10} className="inline text-primary ml-1 animate-pulse" />}
                  </div>
                  <div className="text-[8px] font-bold text-muted flex flex-col uppercase mt-1">
                    <span>ID: {user.id}</span>
                    {user.username && <span className="text-primary/70">@{user.username}</span>}
                  </div>
                </div>
              </div>
              <div className="flex flex-col items-end gap-1">
                <div className="bg-surface px-2 py-1 rounded-lg border border-white/5 text-[9px] font-black tabular-nums shadow-inner shrink-0 text-white/60">
                  {user.channels?.length || 0} Canais
                </div>
              </div>
            </div>

            {expandedUserId === user.id.toString() && (
              <div className="mt-4 pt-4 border-t border-white/5 space-y-3 animate-in fade-in zoom-in-95 duration-300">
                <h4 className="text-[9px] font-black uppercase tracking-widest text-muted">Canais Gerenciados</h4>
                {user.channels && user.channels.length > 0 ? (
                  <div className="space-y-2">
                    {user.channels.map(channel => (
                      <div key={channel.id} className="bg-surface p-3 rounded-xl border border-white/5 flex justify-between items-center group/channel hover:border-primary/30 transition-colors">
                        <div className="min-w-0">
                          <div className="text-[10px] font-black truncate text-white/90">{channel.title}</div>
                          <div className="text-[8px] text-muted truncate">ID: {channel.id}</div>
                        </div>
                        {channel.inviteUrl && (
                          <a 
                            href={channel.inviteUrl} 
                            target="_blank" 
                            rel="noopener noreferrer"
                            onClick={(e) => e.stopPropagation()}
                            className="p-1.5 hover:bg-primary hover:text-black rounded-lg transition-all text-muted"
                          >
                            <ExternalLink size={12} />
                          </a>
                        )}
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-[9px] text-muted italic p-3 text-center bg-surface rounded-xl border border-dashed border-white/5">
                    Nenhum canal encontrado
                  </div>
                )}
              </div>
            )}
          </div>
        ))}
      </div>
      {filteredUsers.length === 0 && (
        <div className="p-12 text-center bg-panel rounded-[2.5rem] border-2 border-dashed border-white/10">
          <Users size={40} className="mx-auto text-muted mb-4 opacity-20" />
          <h3 className="text-sm font-black uppercase tracking-tighter text-muted">Nenhum usuário encontrado</h3>
          <p className="text-[10px] text-muted/60 mt-1 uppercase tracking-widest">SCAN_COMPLETE // 0_RESULTS_FOUND</p>
        </div>
      )}
    </section>
  );
}
