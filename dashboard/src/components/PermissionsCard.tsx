import { Permission } from '../types';
import {
  MessageSquare, Headphones, Video, Image, Smile, Film, Link2, Zap, FileText
} from 'lucide-react';
import { memo, type ReactNode } from 'react';

interface Props {
  title: string;
  icon: ReactNode;
  permission: Permission;
  onToggle?: (field: string, value: boolean) => void;
}

const fields: { key: string; label: string; icon: ReactNode }[] = [
  { key: 'message', label: 'Mensagem', icon: <MessageSquare size={16} /> },
  { key: 'audio', label: 'Áudio', icon: <Headphones size={16} /> },
  { key: 'video', label: 'Vídeo', icon: <Video size={16} /> },
  { key: 'photo', label: 'Foto', icon: <Image size={16} /> },
  { key: 'document', label: 'Arquivo', icon: <FileText size={16} /> },
  { key: 'sticker', label: 'Sticker', icon: <Smile size={16} /> },
  { key: 'gif', label: 'GIF', icon: <Film size={16} /> },
  { key: 'linkPreview', label: 'Link Preview', icon: <Link2 size={16} /> },
];

export function PermissionsCard({ title, icon, permission, onToggle }: Props) {
  const isMessagePerm = title.toLowerCase().includes('mensagem');
  const available = fields.filter(f => {
    if (f.key === 'linkPreview') return isMessagePerm;
    if (!permission) return false;
    return (f.key in permission);
  });
  const perm = (permission || {}) as unknown as Record<string, unknown>;
  const active = available.filter(f => perm[f.key] === true).length;

  return (
    <div className="card">
      <div className="section-header">
        <div className="section-icon purple">{icon}</div>
        <div className="flex-1 min-w-0">
          <h3 className="text-[15px] font-semibold truncate">{title}</h3>
          <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
            {active} de {available.length} ativas
          </p>
        </div>
        <span className="badge badge-accent">{active}/{available.length}</span>
      </div>

      <div className="space-y-2">
        {available.map((f, index) => {
          const isOn = perm[f.key] === true;
          return (
            <div
              key={f.key}
              className={`perm-row animate-stagger-in ${isOn ? 'on' : ''}`}
              style={{ animationDelay: `${index * 0.05}s` }}
              onClick={() => onToggle?.(f.key, !isOn)}
            >
              <div className="flex items-center gap-3 min-w-0">
                <span
                  className="flex-shrink-0"
                  style={{ color: isOn ? 'var(--accent)' : 'var(--hint)', opacity: isOn ? 1 : 0.4 }}
                >
                  {f.icon}
                </span>
                <span className="text-[13px] font-medium">{f.label}</span>
              </div>
              <div className={`toggle ${isOn ? 'on' : ''}`} />
            </div>
          );
        })}
      </div>
    </div>
  );
}
