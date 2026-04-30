import { useState, useEffect } from 'react';
import { Package, Pencil, X, Check, Info } from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';

interface Props {
  caption: string;
  onUpdate?: (text: string) => void;
}

export function NewPackCaptionCard({ caption, onUpdate }: Props) {
  const [editing, setEditing] = useState(false);
  const [text, setText] = useState(caption);

  useEffect(() => { setText(caption); }, [caption]);

  const save = () => { if (text.trim()) { onUpdate?.(text); setEditing(false); } };
  const cancel = () => { setText(caption); setEditing(false); };

  return (
    <div className="card">
      <div className="section-header">
        <div className="section-icon amber"><Package size={18} /></div>
        <div className="flex-1 min-w-0">
          <h3 className="text-[15px] font-semibold">New Pack Caption</h3>
          <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>Template para novo pack</p>
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
            rows={8}
            placeholder="Template..."
          />
          <div className="flex items-center gap-2 py-1" style={{ color: 'var(--hint)' }}>
            <Info size={13} className="flex-shrink-0" />
            <span className="text-xs">Use <strong>$title</strong> e <strong>$link</strong> como variáveis</span>
          </div>
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
          {caption}
        </div>
      )}
    </div>
  );
}
