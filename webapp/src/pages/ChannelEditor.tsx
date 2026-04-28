import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { fetchChannel } from '../api';
import { Channel } from '../types';
import { Terminal, Type, MousePointer2, Shield, Settings, ChevronLeft, Loader2 } from 'lucide-react';
import CaptionTab from '../components/editor/CaptionTab';
import ButtonsTab from '../components/editor/ButtonsTab';
import PermissionsTab from '../components/editor/PermissionsTab';

type Tab = 'captions' | 'buttons' | 'permissions' | 'settings';

export default function ChannelEditor() {
  const { channelId } = useParams<{ channelId: string }>();
  const navigate = useNavigate();
  const [channel, setChannel] = useState<Channel | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<Tab>('captions');

  useEffect(() => {
    const loadChannel = async () => {
      if (!channelId) return;
      try {
        const data = await fetchChannel(parseInt(channelId, 10));
        setChannel(data);
      } catch (err) {
        console.error("Failed to load channel", err);
      } finally {
        setLoading(false);
      }
    };
    loadChannel();
  }, [channelId]);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-20 font-mono text-primary">
        <Loader2 className="animate-spin mb-4" size={32} />
        <div className="tracking-[0.2em] uppercase text-xs">EXECUTING_SYSTEM_INIT...</div>
      </div>
    );
  }

  if (!channel) {
    return (
      <div className="p-6 font-mono border-2 border-red-500 bg-black text-red-500">
        <div className="text-xl font-bold mb-4">ERROR: CHANNEL_NOT_FOUND</div>
        <button 
          onClick={() => navigate(-1)}
          className="px-4 py-2 border border-red-500 hover:bg-red-500 hover:text-black transition-colors"
        >
          &gt; RETURN_TO_BASE
        </button>
      </div>
    );
  }

  return (
    <div className="min-h-screen font-mono text-primary bg-surface p-4 md:p-8 selection:bg-primary selection:text-black">
      {/* Terminal Header */}
      <div className="mb-8 border-b border-white/5 pb-6 flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 text-[10px] text-primary/50 mb-2 font-black tracking-widest">
            <Terminal size={14} className="animate-pulse" />
            <span>ENCRYPTED_CHANNEL_ID: {channelId}</span>
          </div>
          <h1 className="text-3xl font-black tracking-tighter uppercase flex items-center gap-3 text-white glitch-text" data-text={channel.title}>
            <span className="w-2 h-8 bg-primary animate-pulse" />
            {channel.title}
          </h1>
        </div>
        <button 
          onClick={() => navigate(-1)}
          className="btn-secondary group flex items-center gap-2"
        >
          <ChevronLeft size={14} className="group-hover:-translate-x-1 transition-transform" />
          BACK_TO_DASHBOARD
        </button>
      </div>

      {/* Navigation Tabs */}
      <div className="flex flex-wrap gap-2 mb-8 border-b border-white/5 pb-4">
        <button
          onClick={() => setActiveTab('captions')}
          className={`flex items-center gap-2 px-6 py-2 border rounded-full transition-all text-[10px] font-black tracking-widest ${activeTab === 'captions' ? 'border-primary bg-primary text-black shadow-neon' : 'border-white/5 text-muted hover:text-white hover:bg-white/5'}`}
        >
          <Type size={16} />
          LEGENDAS
        </button>
        <button
          onClick={() => setActiveTab('buttons')}
          className={`flex items-center gap-2 px-6 py-2 border rounded-full transition-all text-[10px] font-black tracking-widest ${activeTab === 'buttons' ? 'border-primary bg-primary text-black shadow-neon' : 'border-white/5 text-muted hover:text-white hover:bg-white/5'}`}
        >
          <MousePointer2 size={16} />
          BOTÕES
        </button>
        <button
          onClick={() => setActiveTab('permissions')}
          className={`flex items-center gap-2 px-6 py-2 border rounded-full transition-all text-[10px] font-black tracking-widest ${activeTab === 'permissions' ? 'border-primary bg-primary text-black shadow-neon' : 'border-white/5 text-muted hover:text-white hover:bg-white/5'}`}
        >
          <Shield size={16} />
          PERMISSÕES
        </button>
        <button
          className="flex items-center gap-2 px-6 py-2 border rounded-full transition-all text-[10px] font-black tracking-widest border-white/5 text-muted/20 cursor-not-allowed"
          disabled
        >
          <Settings size={16} />
          CONFIG
        </button>
      </div>

      {/* Tab Content */}
      <div className="animate-in fade-in slide-in-from-bottom-2 duration-500">
        {activeTab === 'captions' && channel && (
          <CaptionTab channel={channel} onUpdate={setChannel} />
        )}
        {activeTab === 'buttons' && channel && (
          <ButtonsTab channel={channel} onUpdate={setChannel} />
        )}
        {activeTab === 'permissions' && channel && (
          <PermissionsTab channel={channel} onUpdate={setChannel} />
        )}
      </div>

      {/* Terminal Footer Info */}
      <div className="mt-12 pt-6 border-t border-white/5 text-[9px] text-primary/30 flex justify-between uppercase tracking-[0.2em] font-black">
        <div>System: FreddyBot // Neural_Core v2.0</div>
        <div className="flex items-center gap-2">
          <span className="w-1.5 h-1.5 rounded-full bg-primary animate-pulse" />
          Status: Online // Connection: Secure_Encrypted
        </div>
      </div>

      <style>{`
        @keyframes pulse {
          0%, 100% { opacity: 1; }
          50% { opacity: 0; }
        }
      `}</style>
    </div>
  );
}
