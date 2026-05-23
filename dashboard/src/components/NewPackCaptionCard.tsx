import { useState, useEffect, memo } from 'react';
import { Package, Pencil, X, Check, Info } from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';

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
          <div className="space-y-2">
            <div
              className={`perm-row ${messageBtn ? 'on' : ''}`}
              onClick={() => setMessageBtn(v => !v)}
            >
              <div>
                <span className="text-[13px] font-medium">Botão na mensagem do bot</span>
                <p className="text-[11px] mt-0.5" style={{ color: 'var(--hint)' }}>Mostra o botão do pack na mensagem editada.</p>
              </div>
              <div className={`toggle ${messageBtn ? 'on' : ''}`} />
            </div>

            <div
              className={`perm-row ${stickerBtn ? 'on' : ''}`}
              onClick={() => setStickerBtn(v => !v)}
            >
              <div>
                <span className="text-[13px] font-medium">Botão no sticker do pack</span>
                <p className="text-[11px] mt-0.5" style={{ color: 'var(--hint)' }}>Mostra o botão abaixo do sticker enviado.</p>
              </div>
              <div className={`toggle ${stickerBtn ? 'on' : ''}`} />
            </div>

            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                className={`btn btn-sm ${position === 'above' ? 'btn-primary' : 'btn-secondary'}`}
                onClick={() => setPosition('above')}
              >
                Mensagem acima
              </button>
              <button
                type="button"
                className={`btn btn-sm ${position === 'below' ? 'btn-primary' : 'btn-secondary'}`}
                onClick={() => setPosition('below')}
              >
                Mensagem abaixo
              </button>
            </div>

            {position === 'below' && (
              <div
                className={`perm-row ${replySticker ? 'on' : ''}`}
                onClick={() => setReplySticker(v => !v)}
              >
                <div>
                  <span className="text-[13px] font-medium">Marcar Sticker</span>
                  <p className="text-[11px] mt-0.5" style={{ color: 'var(--hint)' }}>Envia a mensagem respondendo ao sticker do pack.</p>
                </div>
                <div className={`toggle ${replySticker ? 'on' : ''}`} />
              </div>
            )}
          </div>

          <div className="flex items-start gap-2 py-1" style={{ color: 'var(--hint)' }}>
            <button
              type="button"
              className={`icon-btn ${showHelp ? 'accent' : ''}`}
              onClick={() => setShowHelp(v => !v)}
              title="Variáveis disponíveis"
            >
              <Info size={13} />
            </button>
            {showHelp && (
              <span className="text-xs pt-1">Use <strong>$name</strong>, <strong>$title</strong>, <strong>$link</strong> e <strong>$count</strong>. Ex: [abrir pack]($link)</span>
            )}
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
});
