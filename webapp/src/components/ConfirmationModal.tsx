import { X, AlertTriangle, Check, XCircle } from 'lucide-react';import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: (string | undefined | null | false)[]) {
  return twMerge(clsx(inputs));
}

interface ConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  type?: 'danger' | 'primary' | 'success';
  isLoading?: boolean;
}

export default function ConfirmationModal({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmLabel = 'Confirmar',
  cancelLabel = 'Cancelar',
  type = 'primary',
  isLoading = false
}: ConfirmationModalProps) {
  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop - Separate from flex to ensure it covers everything */}
      <div 
        className="fixed inset-0 bg-black/60 backdrop-blur-sm animate-fade-in z-[1000] pointer-events-auto"
        onClick={!isLoading ? onClose : undefined}
      />
      
      {/* Modal Container - Centered Flex */}
      <div className="fixed inset-0 flex items-center justify-center z-[1010] p-md pointer-events-none">
        {/* Modal Card */}
        <div className="relative w-full max-w-[340px] bg-white rounded-[2rem] shadow-[0_30px_100px_-10px_rgba(0,0,0,0.4)] overflow-hidden animate-slide-up border border-border pointer-events-auto">
          {/* Header Icon Decoration */}
          <div className={cn(
            "h-1.5 w-full",
            type === 'danger' ? "bg-red-500" : 
            type === 'success' ? "bg-green-500" : "bg-primary"
          )} />
          
          <div className="p-lg pt-md">
            {/* Close Button */}
            {!isLoading && (
              <button 
                onClick={onClose}
                className="absolute top-4 right-4 p-1.5 text-muted hover:text-black hover:bg-surface rounded-full transition-all"
              >
                <X size={16} />
              </button>
            )}

            <div className="flex flex-col items-center text-center">
              {/* Visual Indicator */}
              <div className={cn(
                "w-12 h-12 rounded-2xl flex items-center justify-center mb-4 shadow-md",
                type === 'danger' ? "bg-red-50 text-red-500" : 
                type === 'success' ? "bg-green-50 text-green-500" : "bg-primary/10 text-primary"
              )}>
                {type === 'danger' && <AlertTriangle size={24} />}
                {type === 'success' && <Check size={24} />}
                {type === 'primary' && <Check size={24} />}
              </div>

              <h3 className="text-xl font-black uppercase tracking-tighter text-black mb-1">
                {title}
              </h3>
              
              <p className="text-muted font-bold text-[11px] uppercase tracking-widest leading-relaxed mb-8">
                {message}
              </p>

              {/* Actions */}
              <div className="grid grid-cols-1 w-full gap-2">
                <button
                  disabled={isLoading}
                  onClick={onConfirm}
                  className={cn(
                    "w-full py-4 rounded-2xl font-black uppercase text-[10px] tracking-[0.2em] transition-all active:scale-95 shadow-lg flex items-center justify-center gap-2",
                    type === 'danger' ? "bg-red-500 text-white hover:bg-red-600 shadow-red-500/10" : 
                    type === 'success' ? "bg-green-500 text-black hover:bg-green-400 shadow-green-500/10" : 
                    "bg-primary text-black hover:bg-[#c9eb00] shadow-primary/10"
                  )}
                >
                  {isLoading ? (
                    <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  ) : (
                    <>
                      <Check size={14} />
                      {confirmLabel}
                    </>
                  )}
                </button>

                <button
                  disabled={isLoading}
                  onClick={onClose}
                  className="w-full py-3 rounded-2xl font-black uppercase text-[9px] tracking-[0.2em] text-muted hover:text-black transition-all flex items-center justify-center gap-2"
                >
                  {!isLoading && (
                    <>
                      <XCircle size={12} />
                      {cancelLabel}
                    </>
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
