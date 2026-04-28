import { useState } from 'react';
import { updateMessagePermissions, updateButtonsPermissions } from '../../api';
import { Channel, MessagePermission, ButtonsPermission } from '../../types';
import { 
  Shield, 
  MessageSquare, 
  MousePointer2, 
  Link, 
  Type, 
  Music, 
  Video, 
  Image as ImageIcon, 
  StickyNote, 
  PlayCircle,
  Info
} from 'lucide-react';
import { showToast } from '../Toast';

interface PermissionsTabProps {
  channel: Channel;
  onUpdate: (channel: Channel) => void;
}

const defaultMessagePermissions: MessagePermission = {
  linkPreview: true,
  message: true,
  audio: true,
  video: true,
  photo: true,
  sticker: true,
  gif: true,
};

const defaultButtonsPermissions: ButtonsPermission = {
  message: true,
  audio: true,
  video: true,
  photo: true,
  sticker: true,
  gif: true,
};

const Switch = ({ active, label, icon: Icon, onClick }: { active: boolean, label: string, icon: any, onClick: () => void }) => (
  <button
    onClick={onClick}
    className={`flex items-center justify-between p-4 border ${
      active ? 'border-green-500 bg-green-500/10' : 'border-green-900/30 bg-black'
    } transition-all group hover:border-green-500/50 w-full text-left`}
  >
    <div className="flex items-center gap-3">
      <div className={`${active ? 'text-green-500' : 'text-green-900'} transition-colors`}>
        <Icon size={18} />
      </div>
      <span className={`text-[10px] font-bold uppercase tracking-wider ${active ? 'text-green-400' : 'text-green-800'}`}>
        {label}
      </span>
    </div>
    <div className={`w-10 h-5 border ${active ? 'border-green-500' : 'border-green-900'} relative p-0.5`}>
      <div className={`h-full w-4 bg-green-500 transition-all ${active ? 'ml-[18px] shadow-[0_0_10px_#22c55e]' : 'ml-0 opacity-20'}`} />
    </div>
  </button>
);

export default function PermissionsTab({ channel, onUpdate }: PermissionsTabProps) {
  const [msgPerms, setMsgPerms] = useState<MessagePermission>(
    channel.defaultCaption?.messagePermission || defaultMessagePermissions
  );
  const [btnPerms, setBtnPerms] = useState<ButtonsPermission>(
    channel.defaultCaption?.buttonsPermission || defaultButtonsPermissions
  );

  const toggleMsgPerm = async (key: keyof MessagePermission) => {
    const newPerms = { ...msgPerms, [key]: !msgPerms[key] };
    setMsgPerms(newPerms);
    try {
      await updateMessagePermissions(channel.id, newPerms);
      onUpdate({
        ...channel,
        defaultCaption: {
          ...(channel.defaultCaption || { captionId: '', caption: '', buttonsPermission: btnPerms }),
          messagePermission: newPerms
        }
      });
    } catch (err) {
      console.error("Failed to update message permissions", err);
      showToast.error("Erro ao atualizar permissão");
      setMsgPerms(msgPerms); // rollback
    }
  };

  const toggleBtnPerm = async (key: keyof ButtonsPermission) => {
    const newPerms = { ...btnPerms, [key]: !btnPerms[key] };
    setBtnPerms(newPerms);
    try {
      await updateButtonsPermissions(channel.id, newPerms);
      onUpdate({
        ...channel,
        defaultCaption: {
          ...(channel.defaultCaption || { captionId: '', caption: '', messagePermission: msgPerms }),
          buttonsPermission: newPerms
        }
      });
    } catch (err) {
      console.error("Failed to update button permissions", err);
      showToast.error("Erro ao atualizar permissão");
      setBtnPerms(btnPerms); // rollback
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      {/* Message Permissions */}
      <section className="space-y-6">
        <div className="flex items-center gap-2 text-xs font-black uppercase tracking-[0.3em] text-green-500/50">
          <MessageSquare size={14} />
          PERMISSOES_DE_LEGENDA
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <Switch 
            active={msgPerms.message} 
            label="Mensagens de Texto" 
            icon={Type} 
            onClick={() => toggleMsgPerm('message')} 
          />
          <Switch 
            active={msgPerms.linkPreview} 
            label="Preview de Link" 
            icon={Link} 
            onClick={() => toggleMsgPerm('linkPreview')} 
          />
          <Switch 
            active={msgPerms.photo} 
            label="Fotos" 
            icon={ImageIcon} 
            onClick={() => toggleMsgPerm('photo')} 
          />
          <Switch 
            active={msgPerms.video} 
            label="Vídeos" 
            icon={Video} 
            onClick={() => toggleMsgPerm('video')} 
          />
          <Switch 
            active={msgPerms.audio} 
            label="Áudio / Voz" 
            icon={Music} 
            onClick={() => toggleMsgPerm('audio')} 
          />
          <Switch 
            active={msgPerms.gif} 
            label="GIFs" 
            icon={PlayCircle} 
            onClick={() => toggleMsgPerm('gif')} 
          />
          <Switch 
            active={msgPerms.sticker} 
            label="Stickers" 
            icon={StickyNote} 
            onClick={() => toggleMsgPerm('sticker')} 
          />
        </div>
        
        <div className="p-4 border border-green-900/20 bg-green-900/5 flex gap-3">
          <Info size={16} className="text-green-700 shrink-0 mt-0.5" />
          <p className="text-[10px] text-green-800 leading-relaxed uppercase">
            Define quais tipos de mídia receberão a legenda padrão automaticamente ao serem postados.
          </p>
        </div>
      </section>

      {/* Button Permissions */}
      <section className="space-y-6">
        <div className="flex items-center gap-2 text-xs font-black uppercase tracking-[0.3em] text-green-500/50">
          <MousePointer2 size={14} />
          PERMISSOES_DE_BOTOES
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <Switch 
            active={btnPerms.message} 
            label="Mensagens de Texto" 
            icon={Type} 
            onClick={() => toggleBtnPerm('message')} 
          />
          <div className="hidden md:block" /> {/* Spacer */}
          <Switch 
            active={btnPerms.photo} 
            label="Fotos" 
            icon={ImageIcon} 
            onClick={() => toggleBtnPerm('photo')} 
          />
          <Switch 
            active={btnPerms.video} 
            label="Vídeos" 
            icon={Video} 
            onClick={() => toggleBtnPerm('video')} 
          />
          <Switch 
            active={btnPerms.audio} 
            label="Áudio / Voz" 
            icon={Music} 
            onClick={() => toggleBtnPerm('audio')} 
          />
          <Switch 
            active={btnPerms.gif} 
            label="GIFs" 
            icon={PlayCircle} 
            onClick={() => toggleBtnPerm('gif')} 
          />
          <Switch 
            active={btnPerms.sticker} 
            label="Stickers" 
            icon={StickyNote} 
            onClick={() => toggleBtnPerm('sticker')} 
          />
        </div>

        <div className="p-4 border border-green-900/20 bg-green-900/5 flex gap-3">
          <Shield size={16} className="text-green-700 shrink-0 mt-0.5" />
          <p className="text-[10px] text-green-800 leading-relaxed uppercase">
            Define quais tipos de mídia receberão os botões configurados automaticamente.
          </p>
        </div>
      </section>
    </div>
  );
}
