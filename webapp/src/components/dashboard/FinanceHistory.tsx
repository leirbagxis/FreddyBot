import { Coins } from 'lucide-react';
import { useState, useMemo } from 'react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: any[]) {
  return twMerge(clsx(inputs));
}

interface FinanceHistoryProps {
  logs: any[];
  hasMore: boolean;
  loadingMore: boolean;
  handleLoadMore: () => void;
  getLogDescription: (log: any) => string;
  getOperationLabel: (type: string) => string;
}

type FilterType = 'all' | 'admin' | 'shop' | 'consume' | 'transfer' | 'reward';

export default function FinanceHistory({
  logs,
  hasMore,
  loadingMore,
  handleLoadMore,
  getLogDescription,
  getOperationLabel
}: FinanceHistoryProps) {
  const [activeFilter, setActiveFilter] = useState<FilterType>('all');

  const filteredLogs = useMemo(() => {
    if (activeFilter === 'all') return logs;
    return logs.filter(log => {
      switch (activeFilter) {
        case 'admin': return log.operation_type === 'admin_add' || log.operation_type === 'admin_reduce';
        case 'shop': return log.operation_type === 'purchase' || log.operation_type === 'purchase_web' || log.operation_type === 'sell_shop' || log.operation_type === 'sell_shop_web';
        case 'consume': return log.operation_type === 'consume';
        case 'transfer': return log.operation_type === 'sell_player' || log.operation_type === 'buy_player';
        case 'reward': return log.operation_type === 'reward';
        default: return true;
      }
    });
  }, [logs, activeFilter]);

  const filters: { id: FilterType; label: string }[] = [
    { id: 'all', label: 'Todos' },
    { id: 'admin', label: 'Admin' },
    { id: 'shop', label: 'Loja' },
    { id: 'consume', label: 'Usados' },
    { id: 'transfer', label: 'Trocas' },
    { id: 'reward', label: 'Bônus' },
  ];

  return (
    <div className="refined-card p-0 overflow-hidden border-2 border-black/5 animate-fade-in mb-32">
      <header className="p-6 border-b border-border bg-black text-white flex flex-col gap-4">
        <div className="flex justify-between items-center">
          <div className="flex items-center gap-3">
            <Coins size={20} className="text-primary animate-spin" style={{ animationDuration: '3s' }} />
            <h3 className="text-xl font-black uppercase tracking-tighter">Extrato Financeiro</h3>
          </div>
          <span className="text-[9px] font-black uppercase tracking-widest opacity-60 text-right">Últimas Movimentações</span>
        </div>

        <div className="flex items-center gap-2 overflow-x-auto no-scrollbar py-1">
          {filters.map(filter => (
            <button
              key={filter.id}
              onClick={() => setActiveFilter(filter.id)}
              className={cn(
                "px-4 py-1.5 rounded-xl text-[8px] font-black uppercase tracking-widest transition-all whitespace-nowrap",
                activeFilter === filter.id 
                  ? "bg-primary text-black" 
                  : "bg-white/10 text-white/40 hover:text-white hover:bg-white/20"
              )}
            >
              {filter.label}
            </button>
          ))}
        </div>
      </header>
      
      <div className="bg-black text-white p-4 flex justify-between items-center px-6 border-b border-white/5">
        <span className="text-[9px] font-black uppercase tracking-widest opacity-60">Operação / Data</span>
        <span className="text-[9px] font-black uppercase tracking-widest opacity-60 text-right">Valor</span>
      </div>
      <div className="divide-y divide-border">
        {filteredLogs.length === 0 ? (
          <div className="py-20 text-center opacity-20 font-black uppercase tracking-widest text-xs italic">Nenhuma transação encontrada</div>
        ) : (
          filteredLogs.map((log, index) => (
            <div 
              key={log.id} 
              className="flex items-center justify-between p-4 sm:p-6 hover:bg-black/[0.02] transition-all group stagger-item"
              style={{ animationDelay: `${index * 40}ms` }}
            >
              <div className="flex items-center gap-3 sm:gap-5">
                <div className={cn(
                  "w-10 h-10 sm:w-12 sm:h-12 rounded-xl sm:rounded-2xl flex items-center justify-center border-2 shrink-0 transition-transform group-hover:scale-110",
                  log.amount > 0 ? "bg-green-500/10 border-green-500/20 text-green-600" : "bg-red-500/10 border-red-500/20 text-red-600"
                )}>
                  <Coins size={20} className="sm:w-6 sm:h-6 animate-spin-slow" style={{ animationDuration: '3s' }} />
                </div>
                <div>
                  <div className="text-xs sm:text-sm font-black uppercase tracking-tight text-black">{getLogDescription(log)}</div>
                  <div className="text-[8px] sm:text-[10px] font-bold text-muted uppercase tracking-widest mt-0.5">
                    {getOperationLabel(log.operation_type)} • {new Date(log.created_at).toLocaleString('pt-BR')}
                  </div>
                </div>
              </div>
              <div className={cn("text-base sm:text-xl font-black tracking-tighter tabular-nums shrink-0 ml-2", log.amount > 0 ? "text-green-600" : "text-red-600")}>
                {log.amount > 0 ? '+' : ''}{log.amount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
              </div>
            </div>
          ))
        )}
      </div>

      {hasMore && (
        <div className="p-6 border-t border-border bg-surface/50 flex justify-center">
          <button
            disabled={loadingMore}
            onClick={handleLoadMore}
            className="group flex items-center gap-3 px-8 py-3 rounded-2xl border-2 border-black/5 hover:border-black transition-all bg-white shadow-sm hover:shadow-md"
          >
            <Coins size={16} className={cn("text-muted group-hover:text-black transition-colors animate-spin", loadingMore ? "" : "animate-none")} />
            <span className="text-[10px] font-black uppercase tracking-widest text-muted group-hover:text-black transition-colors">
              {loadingMore ? 'Processando...' : 'Carregar mais transações'}
            </span>
          </button>
        </div>
      )}
    </div>
  );
}
