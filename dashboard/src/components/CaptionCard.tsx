import { useState, useEffect, memo } from 'react';
import { Caption } from '../types';
import { FileText, Pencil, X, Check } from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';
import { CaptionPreview } from './CaptionPreview';
import { Button } from './ui/button';

interface Props {
  caption: Caption;
  onUpdate?: (text: string) => void;
}

export const CaptionCard = memo(({ caption, onUpdate }: Props) => {
  const [editing, setEditing] = useState(false);
  const [text, setText] = useState(caption.caption);

  useEffect(() => { setText(caption.caption); }, [caption.caption]);

  const save = () => { if (text.trim()) { onUpdate?.(text); setEditing(false); } };
  const cancel = () => { setText(caption.caption); setEditing(false); };

  return (
    <div className="content-card">
      <div className="content-card-header">
        <div className="content-card-icon"><FileText size={18} /></div>
        <div className="flex-1 min-w-0">
          <div className="content-card-title">Caption Padrão</div>
          <div className="content-card-desc">Aplicada em todas as mensagens</div>
        </div>
        {!editing && (
          <Button variant="ghost" size="icon" className="text-accent shrink-0" onClick={() => setEditing(true)}>
            <Pencil size={15} />
          </Button>
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
            <Button variant="secondary" size="sm" onClick={cancel}>
              <X size={14} /> Cancelar
            </Button>
            <Button variant="default" size="sm" onClick={save}>
              <Check size={14} /> Salvar
            </Button>
          </div>
        </div>
      ) : (
        <div onClick={() => setEditing(true)}>
          <CaptionPreview text={caption.caption} />
        </div>
      )}
    </div>
  );
});
