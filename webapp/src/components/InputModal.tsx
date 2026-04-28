import { useState } from 'react';
import { X, Zap, type LucideIcon } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: any[]) {
  return twMerge(clsx(inputs));
}

interface InputModalProps {
  onClose: () => void;
  onConfirm: (value: string) => void;
  title: string;
  label: string;
  placeholder: string;
  itemName: string;
  icon: LucideIcon;
  isLoading?: boolean;
  maxLength?: number;
}

export default function InputModal({ 
  onClose, 
  onConfirm, 
  title,
  label,
  placeholder,
  itemName, 
  icon: Icon,
  isLoading,
  maxLength = 32
}: InputModalProps) {
  const [value, setValue] = useState('');

  return (
    <div className="fixed inset-0 z-[110] flex items-center justify-center p-4">
      <div 
        className="absolute inset-0 bg-black/60 backdrop-blur-sm animate-fade-in" 
        onClick={onClose}
      />
      
      <div className="relative w-full max-w-sm bg-white border-[3px] border-black rounded-[1.5rem] shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] overflow-hidden animate-slide-in">
        <header className="bg-black p-4 flex justify-between items-center text-primary">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-primary/20 flex items-center justify-center border border-primary/30">
              <Icon size={16} className="text-primary" />
            </div>
            <div>
              <h3 className="text-sm font-black uppercase tracking-tight text-white normal-case">{title}</h3>
            </div>
          </div>
          <button 
            onClick={onClose}
            className="w-8 h-8 flex items-center justify-center rounded-full hover:bg-white/10 transition-colors text-white"
          >
            <X size={18} />
          </button>
        </header>

        <div className="p-6 space-y-4">
          <div className="text-[10px] font-bold text-muted uppercase tracking-wider leading-relaxed">
            Utilizando item: <span className="text-black font-black">[{itemName}]</span>
          </div>

          <div className="space-y-1.5">
            <label className="text-[9px] font-black uppercase tracking-widest text-black/40 ml-1">{label}</label>
            <div className="relative group">
              <input 
                autoFocus
                type="text"
                maxLength={maxLength}
                value={value}
                onChange={(e) => setValue(e.target.value)}
                placeholder={placeholder}
                className="w-full bg-surface border-2 border-black rounded-xl p-3 pr-10 text-base font-black placeholder:text-black/5 outline-none transition-all"
              />
              <Zap className={cn(
                "absolute right-3 top-1/2 -translate-y-1/2 transition-colors",
                value.length > 0 ? "text-primary fill-primary" : "text-black/5"
              )} size={16} />
            </div>
            <div className="flex justify-between px-1">
              <span className="text-[7px] font-bold text-muted uppercase">Max {maxLength} chars</span>
              <span className={cn(
                "text-[7px] font-black uppercase font-mono",
                value.length > maxLength * 0.8 ? "text-error" : "text-black"
              )}>{value.length}/{maxLength}</span>
            </div>
          </div>

          <button
            disabled={!value.trim() || isLoading}
            onClick={() => onConfirm(value.trim())}
            className="w-full bg-black text-primary py-4 rounded-xl font-black uppercase text-[10px] tracking-[0.2em] flex items-center justify-center gap-2 transition-all active:scale-95 disabled:opacity-50 disabled:grayscale shadow-md"
          >
            {isLoading ? (
              <div className="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin" />
            ) : (
              <>
                Confirmar <Zap size={12} />
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
