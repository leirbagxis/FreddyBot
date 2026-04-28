import { Channel } from '../../types';
import { Settings, ExternalLink, Cpu } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

interface ChannelCardProps {
  channel: Channel;
}

export default function ChannelCard({ channel }: ChannelCardProps) {
  const navigate = useNavigate();

  return (
    <div className="relative group overflow-hidden bg-[#111111] border-2 border-[#D9FF00]/20 hover:border-[#D9FF00] transition-all duration-300 rounded-none p-md">
      {/* Technical corner elements */}
      <div className="absolute top-0 right-0 w-6 h-6 border-t-2 border-r-2 border-[#D9FF00]/30 group-hover:border-[#D9FF00] transition-colors" />
      <div className="absolute bottom-0 left-0 w-3 h-3 border-b-2 border-l-2 border-[#D9FF00]/30 group-hover:border-[#D9FF00] transition-colors" />
      
      {/* Scanning line effect */}
      <div className="absolute inset-0 pointer-events-none bg-gradient-to-b from-transparent via-[#D9FF00]/5 to-transparent h-1/2 w-full -translate-y-full group-hover:animate-[scan_2s_linear_infinite]" />

      <div className="flex flex-col h-full relative z-10">
        <div className="flex items-start justify-between mb-lg">
          <div className="flex items-center gap-sm">
            <div className="p-xs bg-[#D9FF00]/10 border border-[#D9FF00]/20">
              <Cpu className="text-[#D9FF00]" size={18} />
            </div>
            <div>
              <h3 className="text-white font-mono text-sm tracking-wider uppercase truncate max-w-[140px]">
                {channel.title}
              </h3>
              <div className="flex items-center gap-xs">
                <span className="w-1.5 h-1.5 rounded-full bg-success animate-pulse" />
                <span className="text-[#D9FF00]/60 text-[9px] font-mono uppercase">
                  Active // CH_{channel.id}
                </span>
              </div>
            </div>
          </div>
          <a 
            href={channel.inviteUrl} 
            target="_blank" 
            rel="noopener noreferrer"
            className="text-[#D9FF00]/40 hover:text-[#D9FF00] transition-colors p-xs"
            title="Abrir no Telegram"
          >
            <ExternalLink size={14} />
          </a>
        </div>

        <div className="mt-auto">
          <button 
            onClick={() => navigate(`/channel/${channel.id}`)}
            className="w-full flex items-center justify-center gap-sm bg-transparent hover:bg-[#D9FF00] border border-[#D9FF00] text-[#D9FF00] hover:text-[#111111] font-mono text-[10px] font-black uppercase tracking-[0.2em] py-sm transition-all duration-300"
          >
            <Settings size={12} />
            Configurar
          </button>
        </div>
      </div>

      <style>{`
        @keyframes scan {
          0% { transform: translateY(-100%); }
          100% { transform: translateY(200%); }
        }
      `}</style>
    </div>
  );
}
