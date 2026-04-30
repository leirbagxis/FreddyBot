import { useState, useEffect } from 'react';
import { Caption } from '../types';
import { FileText, Pencil, X, Check } from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';

interface Props {
  caption: Caption;
  onUpdate?: (text: string) => void;
}

export function CaptionCard({ caption, onUpdate }: Props) {
  const [editing, setEditing] = useState(false);
  const [text, setText] = useState(caption.caption);

  useEffect(() => { setText(caption.caption); }, [caption.caption]);

  const save = () => { if (text.trim()) { onUpdate?.(text); setEditing(false); } };
  const cancel = () => { setText(caption.caption); setEditing(false); };

  return (
    <div className="card">
      <div className="section-header">
        <div className="section-icon purple"><FileText size={18} /></div>
        <div className="flex-1 min-w-0">
          <h3 className="text-[15px] font-semibold">Caption Padrão</h3>
          <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>Aplicada em todas as mensagens</p>
        </div>
        {!editing && (
          <button className="icon-btn accent" onClick={() => setEditing(true)}>
            <Pencil size={15} />
          </button>
        )}
      </div>

      {editing ? (
        <div className="space-y-3">
          <RichTextEditor
            value={text}
            onChange={setText}
            rows={5}
            placeholder="Caption padrão..."
          />
          <div className="flex items-center justify-end gap-2">
            <button className="btn btn-secondary btn-sm" onClick={cancel}>
              <X size={14} /> Cancelar
            </button>
            <button className="btn btn-primary btn-sm" onClick={save}>
              <Check size={14} /> Salvar
            </button>
          </div>
        </div>
      ) : (
        <div className="caption-preview" onClick={() => setEditing(true)}>
          {caption.caption || <span style={{ opacity: 0.3, fontStyle: 'italic' }}>Sem caption definida</span>}
        </div>
      )}
    </div>
  );
}
