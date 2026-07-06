import { useState, useEffect, memo } from 'react';
import { Caption } from '../types';
import { FileText, Pencil, X, Check } from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';
import { Card, CardContent } from './ui/card';
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
    <Card>
      <CardContent className="pt-4">
        <div className="flex items-center gap-3 mb-4">
          <div className="section-icon purple"><FileText size={18} /></div>
          <div className="flex-1 min-w-0">
            <h3 className="text-[15px] font-semibold">Caption Padrão</h3>
            <p className="text-xs mt-0.5 text-muted-foreground">Aplicada em todas as mensagens</p>
          </div>
          {!editing && (
            <Button variant="ghost" size="icon" className="text-accent" onClick={() => setEditing(true)}>
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
          <div className="caption-preview" onClick={() => setEditing(true)}>
            {caption.caption || <span style={{ opacity: 0.3, fontStyle: 'italic' }}>Sem caption definida</span>}
          </div>
        )}
      </CardContent>
    </Card>
  );
});
