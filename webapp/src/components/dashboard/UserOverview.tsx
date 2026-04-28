import { Link } from 'react-router-dom';
import { ShoppingBag, Trophy, ArrowRight, Shield } from 'lucide-react';

interface UserOverviewProps {
  botId: string | undefined;
  profile: any;
  botConfig: any;
}

export default function UserOverview({ botId, profile, botConfig }: UserOverviewProps) {
  return (
    <div className="space-y-lg animate-fade-in">
      <div className="grid grid-cols-2 gap-md">
        <Link to={`/${botId}/shop`} className="refined-card group bg-primary border-transparent hover:bg-white transition-all shadow-neon">
          <div className="flex flex-col gap-sm">
            <div className="w-10 h-10 rounded-xl bg-black/10 flex items-center justify-center text-black"><ShoppingBag size={20} /></div>
            <div>
              <div className="text-[10px] font-black uppercase tracking-widest text-black/60">ACESSAR_ESTOQUE</div>
              <div className="text-xl font-black uppercase tracking-tighter flex items-center gap-2 text-black">Loja <ArrowRight size={18} /></div>
            </div>
          </div>
        </Link>
        <Link to={`/${botId}/ranking`} className="refined-card group border-white/5 hover:border-primary/50 transition-all">
          <div className="flex flex-col gap-sm">
            <div className="w-10 h-10 rounded-xl bg-white/5 flex items-center justify-center text-primary group-hover:scale-110 transition-transform"><Trophy size={20} /></div>
            <div>
              <div className="text-[10px] font-black uppercase tracking-widest text-white/40">VISUALIZAR</div>
              <div className="text-xl font-black uppercase tracking-tighter flex items-center gap-2 text-white">Rankings <ArrowRight size={18} /></div>
            </div>
          </div>
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
        <div className="refined-card bg-black border-primary/20 shadow-neon overflow-hidden">
          <div className="text-[10px] font-bold uppercase tracking-[0.2em] text-primary mb-xs block opacity-60">SALDO_ATUAL</div>
          <div className="flex items-baseline gap-sm mt-2 overflow-hidden">
            <span className="text-primary font-black text-2xl shrink-0">{botConfig.currency_symbol}</span>
            <h2 className="text-5xl md:text-6xl font-black text-primary tracking-tighter truncate leading-tight drop-shadow-[0_0_10px_rgba(217,255,0,0.3)]">
              {profile.balance.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
            </h2>
          </div>
          <div className="mt-lg pt-md border-t border-white/5 flex justify-between font-black text-[9px] text-white/20 uppercase tracking-[0.3em]">
            <span>{botConfig.currency_name}</span>
            <span className="flex items-center gap-2">
              <span className="w-1 h-1 rounded-full bg-primary animate-pulse" />
              CONTA_VERIFICADA
            </span>
          </div>
        </div>

        <div className="refined-card col-span-1 border-white/5">
          <header className="flex justify-between items-start mb-lg">
            <div className="overflow-hidden">
              <div className="text-[10px] font-bold uppercase tracking-[0.2em] text-white/40 mb-xs block">IDENTIDADE_NEURAL</div>
              <h3 className="text-2xl font-black text-white leading-none truncate glitch-text" data-text={`${profile.first_name} ${profile.last_name || ''}`}>
                {profile.first_name} {profile.last_name}
              </h3>
            </div>
            <div className="w-12 h-12 bg-white/5 border border-white/5 rounded-2xl flex items-center justify-center text-primary shrink-0 shadow-inner"><Shield size={24} /></div>
          </header>
          <div className="font-bold text-[10px] space-y-1 text-white/40 uppercase tracking-[0.2em]">
            <div>NEURAL_ID: <span className="text-white/60">{profile.telegram_user_id}</span></div>
            <div>STATUS: <span className="text-primary animate-pulse">CONNECTED_ENCRYPTED</span></div>
          </div>
        </div>
      </div>
    </div>
  );
}
