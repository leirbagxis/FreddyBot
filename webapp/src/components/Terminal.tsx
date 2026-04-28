import { useState, useEffect, useCallback } from 'react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: (string | undefined | null | false)[]) {
  return twMerge(clsx(inputs));
}

type LogType = 'info' | 'warn' | 'error' | 'success';

interface LogEntry {
  id: string;
  message: string;
  type: LogType;
  timestamp: string;
}

export function useTerminal() {
  const [logs, setLogs] = useState<LogEntry[]>([]);

  const addLog = useCallback((message: string, type: LogType = 'info') => {
    const newLog: LogEntry = {
      id: Math.random().toString(36).substring(7),
      message: `>> ${message}`,
      type,
      timestamp: new Date().toLocaleTimeString('pt-BR', { hour12: false }),
    };
    setLogs(prev => [...prev.slice(-49), newLog]);
    console.log(`[TERM] ${message}`);
  }, []);

  return { logs, addLog, clearLogs: () => setLogs([]) };
}

// Global context simulation for the terminal
let globalAddLog: (msg: string, type?: LogType) => void = () => {};

export const terminal = {
  log: (msg: string) => globalAddLog(msg, 'info'),
  warn: (msg: string) => globalAddLog(msg, 'warn'),
  error: (msg: string) => globalAddLog(msg, 'error'),
  success: (msg: string) => globalAddLog(msg, 'success'),
};

export default function Terminal({ logs, setAddLog, onClose }: { logs: LogEntry[], setAddLog: any, onClose?: () => void }) {
  useEffect(() => {
    globalAddLog = setAddLog;
  }, [setAddLog]);

  return (
    <div className="bg-[#111111] border-t border-white/10 h-64 overflow-hidden flex flex-col font-mono text-[11px] shadow-2xl">
      <div className="bg-black px-lg py-sm flex justify-between items-center border-b border-white/5">
        <div className="flex items-center gap-md">
          <div className="flex gap-1.5">
            <div className="w-2.5 h-2.5 rounded-full bg-red-500/80"></div>
            <div className="w-2.5 h-2.5 rounded-full bg-yellow-500/80"></div>
            <div className="w-2.5 h-2.5 rounded-full bg-green-500/80"></div>
          </div>
          <span className="text-white/40 font-black tracking-[0.2em] uppercase text-[9px] ml-2">Log de Eventos</span>
        </div>
        
        {onClose && (
          <button 
            onClick={onClose}
            className="text-white/40 hover:text-white transition-colors uppercase font-black text-[9px] tracking-widest px-3 py-1 border border-white/10 rounded-full hover:bg-white/5"
          >
            Fechar Log
          </button>
        )}
      </div>
      
      <div className="flex-1 overflow-y-auto p-md space-y-1.5 scrollbar-thin scrollbar-thumb-white/10">
        {logs.map(log => (
          <div key={log.id} className="flex gap-md animate-fade-in group">
            <span className="text-white/20 shrink-0 select-none">[{log.timestamp}]</span>
            <span className={cn(
              "font-medium",
              log.type === 'success' ? 'text-primary' : 
              log.type === 'error' ? 'text-red-400' : 
              log.type === 'warn' ? 'text-yellow-400' : 'text-white/70'
            )}>
              {log.message}
            </span>
          </div>
        ))}
        {logs.length === 0 && (
          <div className="text-white/20 italic opacity-50 flex items-center gap-2">
            <span className="animate-pulse">_</span>
            Aguardando eventos...
          </div>
        )}
      </div>
    </div>
  );
}
