import { useState, useEffect } from 'react';
import { X, CheckCircle, AlertCircle, Info } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: (string | undefined | null | false)[]) {
  return twMerge(clsx(inputs));
}

type ToastType = 'success' | 'error' | 'info';

interface Toast {
  id: string;
  message: string;
  type: ToastType;
}

// Global hook/state simulation for toasts
let globalAddToast: (msg: string, type: ToastType) => void = () => {};

export const showToast = {
  success: (msg: string) => globalAddToast(msg, 'success'),
  error: (msg: string) => globalAddToast(msg, 'error'),
  info: (msg: string) => globalAddToast(msg, 'info'),
};

export default function ToastContainer() {
  const [toasts, setToasts] = useState<Toast[]>([]);

  useEffect(() => {
    globalAddToast = (message: string, type: ToastType) => {
      const id = Math.random().toString(36).substring(7);
      setToasts(prev => [...prev, { id, message, type }]);
      setTimeout(() => {
        setToasts(prev => prev.filter(t => t.id !== id));
      }, 5000);
    };
  }, []);

  return (
    <div className="fixed bottom-24 md:bottom-auto md:top-xl right-md md:right-xl z-[200] flex flex-col gap-sm max-w-[calc(100vw-2rem)]">
      {toasts.map(toast => (
        <div 
          key={toast.id}
          className={cn(
            "w-full md:min-w-[320px] p-lg flex items-center gap-md border shadow-2xl animate-slide-in rounded-3xl",
            toast.type === 'success' ? 'bg-black text-white border-primary/20' : 
            toast.type === 'error' ? 'bg-red-50 text-red-900 border-red-100' : 
            'bg-white text-black border-border'
          )}
        >
          <div className={cn(
            "w-10 h-10 rounded-2xl flex items-center justify-center shrink-0",
            toast.type === 'success' ? 'bg-primary text-black' : 
            toast.type === 'error' ? 'bg-red-500 text-white' : 
            'bg-black text-white'
          )}>
            {toast.type === 'success' && <CheckCircle size={20} />}
            {toast.type === 'error' && <AlertCircle size={20} />}
            {toast.type === 'info' && <Info size={20} />}
          </div>
          
          <div className="flex-1">
            <div className={cn(
              "text-[9px] font-black uppercase tracking-[0.2em] mb-1 opacity-40",
              toast.type === 'error' ? 'text-red-900' : ''
            )}>
              {toast.type}
            </div>
            <div className="text-xs font-black uppercase tracking-tight leading-tight">{toast.message}</div>
          </div>
          
          <button 
            onClick={() => setToasts(prev => prev.filter(t => t.id !== toast.id))}
            className="text-muted hover:text-black transition-colors p-1"
          >
            <X size={16} />
          </button>
        </div>
      ))}
    </div>
  );
}
