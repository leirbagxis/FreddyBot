import { useState, useCallback, createContext, useContext } from 'react';
import { CheckCircle2, Info, XCircle } from 'lucide-react';

interface ToastItem {
  id: string;
  message: string;
  type: 'success' | 'info' | 'error';
}

type ToastFn = (message: string, type?: 'success' | 'info' | 'error') => void;
const ToastContext = createContext<ToastFn>(() => {});
export function useToast() { return useContext(ToastContext); }

const icons = {
  success: <CheckCircle2 size={16} />,
  info: <Info size={16} />,
  error: <XCircle size={16} />,
};
const colors = {
  success: 'var(--success)',
  info: 'var(--accent)',
  error: 'var(--danger)',
};

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const addToast: ToastFn = useCallback((message, type = 'success') => {
    const id = Math.random().toString(36).slice(2);
    setToasts(prev => [...prev, { id, message, type }]);
    setTimeout(() => setToasts(prev => prev.filter(t => t.id !== id)), 2500);
  }, []);

  return (
    <ToastContext.Provider value={addToast}>
      {children}
      {toasts.length > 0 && (
        <div className="toast-wrap">
          {toasts.map(t => (
            <div key={t.id} className="toast-msg">
              <span style={{ color: colors[t.type], flexShrink: 0 }}>{icons[t.type]}</span>
              <span className="truncate">{t.message}</span>
            </div>
          ))}
        </div>
      )}
    </ToastContext.Provider>
  );
}
