import { useState, useEffect, memo } from 'react';
import { Package, Pencil, X, Check, Info } from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';
import { CaptionPreview } from './CaptionPreview';
import { Button } from './ui/button';
import { Switch } from './ui/switch';

interface Props {
  caption: string;
  messageButtons: boolean;
  stickerButtons: boolean;
  messagePosition: 'above' | 'below';
  replyToSticker: boolean;
  onUpdate?: (settings: { caption: string; messageButtons: boolean; stickerButtons: boolean; messagePosition: 'above' | 'below'; replyToSticker: boolean }) => void;
}

export const NewPackCaptionCard = memo(({ caption, messageButtons, stickerButtons, messagePosition, replyToSticker, onUpdate }: Props) => {
  const [editing, setEditing] = useState(false);
  const [text, setText] = useState(caption);
  const [messageBtn, setMessageBtn] = useState(messageButtons);
  const [stickerBtn, setStickerBtn] = useState(stickerButtons);
  const [position, setPosition] = useState<'above' | 'below'>(messagePosition);
  const [replySticker, setReplySticker] = useState(replyToSticker);
  const [showHelp, setShowHelp] = useState(false);

  useEffect(() => { setText(caption); }, [caption]);
  useEffect(() => { setMessageBtn(messageButtons); }, [messageButtons]);
  useEffect(() => { setStickerBtn(stickerButtons); }, [stickerButtons]);
  useEffect(() => { setPosition(messagePosition); }, [messagePosition]);
  useEffect(() => { setReplySticker(replyToSticker); }, [replyToSticker]);

  const save = () => {
    if (text.trim()) {
      onUpdate?.({ caption: text, messageButtons: messageBtn, stickerButtons: stickerBtn, messagePosition: position, replyToSticker: position === 'below' && replySticker });
      setShowHelp(false);
      setEditing(false);
    }
  };
  const cancel = () => {
    setText(caption);
    setMessageBtn(messageButtons);
    setStickerBtn(stickerButtons);
    setPosition(messagePosition);
    setReplySticker(replyToSticker);
    setShowHelp(false);
    setEditing(false);
  };

  return (
    <div className="content-card">
      <div className="content-card-header">
        <div className="content-card-icon" style={{ background: 'var(--warning-soft)', color: 'var(--warning)' }}>
          <Package size={18} />
        </div>
        <div className="flex-1 min-w-0">
          <div className="content-card-title">New Pack Caption</div>
          <div className="content-card-desc">Template para novo pack</div>
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
              rows={8}
              placeholder="Template..."
            />
            <div className="space-y-2">
              <div
                className={`flex items-center justify-between px-[18px] py-3 rounded-[20px] gap-3 min-h-[52px] cursor-pointer transition-all ${messageBtn ? 'bg-accent/10' : 'bg-muted/50'}`}
                onClick={() => setMessageBtn(v => !v)}
              >
                <div>
                  <span className="text-[13px] font-medium">Botão na mensagem do bot</span>
                  <p className="text-[11px] mt-0.5 text-muted-foreground">Mostra o botão do pack na mensagem editada.</p>
                </div>
                <Switch
                  checked={messageBtn}
                  onCheckedChange={(c) => setMessageBtn(c)}
                  onClick={(e: React.MouseEvent) => e.stopPropagation()}
                />
              </div>

              <div
                className={`flex items-center justify-between px-[18px] py-3 rounded-[20px] gap-3 min-h-[52px] cursor-pointer transition-all ${stickerBtn ? 'bg-accent/10' : 'bg-muted/50'}`}
                onClick={() => setStickerBtn(v => !v)}
              >
                <div>
                  <span className="text-[13px] font-medium">Botão no sticker do pack</span>
                  <p className="text-[11px] mt-0.5 text-muted-foreground">Mostra o botão abaixo do sticker enviado.</p>
                </div>
                <Switch
                  checked={stickerBtn}
                  onCheckedChange={(c) => setStickerBtn(c)}
                  onClick={(e: React.MouseEvent) => e.stopPropagation()}
                />
              </div>

              <div className="grid grid-cols-2 gap-2">
                <Button
                  size="sm"
                  variant={position === 'above' ? 'default' : 'secondary'}
                  onClick={() => setPosition('above')}
                >
                  Mensagem acima
                </Button>
                <Button
                  size="sm"
                  variant={position === 'below' ? 'default' : 'secondary'}
                  onClick={() => setPosition('below')}
                >
                  Mensagem abaixo
                </Button>
              </div>

              {position === 'below' && (
                <div
                  className={`flex items-center justify-between px-[18px] py-3 rounded-[20px] gap-3 min-h-[52px] cursor-pointer transition-all ${replySticker ? 'bg-accent/10' : 'bg-muted/50'}`}
                  onClick={() => setReplySticker(v => !v)}
                >
                  <div>
                    <span className="text-[13px] font-medium">Marcar Sticker</span>
                    <p className="text-[11px] mt-0.5 text-muted-foreground">Envia a mensagem respondendo ao sticker do pack.</p>
                  </div>
                  <Switch
                    checked={replySticker}
                    onCheckedChange={(c) => setReplySticker(c)}
                    onClick={(e: React.MouseEvent) => e.stopPropagation()}
                  />
                </div>
              )}
            </div>

            <div className="flex items-start gap-2 py-1 text-muted-foreground">
              <Button
                variant="ghost"
                size="icon-xs"
                className={showHelp ? 'text-accent' : ''}
                onClick={() => setShowHelp(v => !v)}
                title="Variáveis disponíveis"
              >
                <Info size={13} />
              </Button>
              {showHelp && (
                <span className="text-xs pt-1">Use <strong>$name</strong>, <strong>$title</strong>, <strong>$link</strong> e <strong>$count</strong>. Ex: [abrir pack]($link)</span>
              )}
            </div>
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
          <div className="cursor-pointer" onClick={() => setEditing(true)}>
            <CaptionPreview text={caption} />
          </div>
        )}
    </div>
  );
});
