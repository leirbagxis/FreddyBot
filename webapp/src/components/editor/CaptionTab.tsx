import { useState } from 'react';
import { updateDefaultCaption, updateNewPackCaption } from '../../api';
import { Channel } from '../../types';
import { Save, AlertCircle, Info } from 'lucide-react';
import { showToast } from '../Toast';

interface CaptionTabProps {
  channel: Channel;
  onUpdate: (channel: Channel) => void;
}

export default function CaptionTab({ channel, onUpdate }: CaptionTabProps) {
  const [defaultCaption, setDefaultCaption] = useState(channel.defaultCaption?.caption || '');
  const [newPackCaption, setNewPackCaption] = useState(channel.newPackCaption || '');
  const [savingDefault, setSavingDefault] = useState(false);
  const [savingNewPack, setSavingNewPack] = useState(false);

  const handleSaveDefault = async () => {
    setSavingDefault(true);
    try {
      await updateDefaultCaption(channel.id, defaultCaption);
      showToast.success('Legenda padrão atualizada!');
      onUpdate({
        ...channel,
        defaultCaption: {
          ...(channel.defaultCaption || { captionId: '', buttonsPermission: {}, messagePermission: {} }),
          caption: defaultCaption
        }
      });
    } catch (err) {
      console.error("Failed to save default caption", err);
      showToast.error('Erro ao salvar legenda padrão');
    } finally {
      setSavingDefault(false);
    }
  };

  const handleSaveNewPack = async () => {
    setSavingNewPack(true);
    try {
      await updateNewPackCaption(channel.id, newPackCaption);
      showToast.success('Legenda de novo pack atualizada!');
      onUpdate({
        ...channel,
        newPackCaption: newPackCaption
      });
    } catch (err) {
      console.error("Failed to save new pack caption", err);
      showToast.error('Erro ao salvar legenda de novo pack');
    } finally {
      setSavingNewPack(false);
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
      {/* Default Caption Section */}
      <section className="border border-green-900 p-6 bg-green-950/10">
        <div className="flex items-center gap-2 mb-4 text-xs font-bold uppercase tracking-widest">
          <Info size={14} className="text-green-500" />
          LEGENDA_PADRAO
        </div>
        <p className="text-[10px] opacity-60 mb-4 leading-relaxed">
          Esta legenda será aplicada automaticamente a todas as postagens enviadas para o canal.
          Você pode usar HTML do Telegram (ex: &lt;b&gt;bold&lt;/b&gt;).
        </p>
        <textarea
          value={defaultCaption}
          onChange={(e) => setDefaultCaption(e.target.value)}
          placeholder="Digite a legenda padrão..."
          className="w-full h-48 bg-black border border-green-900 focus:border-green-500 outline-none p-4 text-sm resize-none transition-colors"
        />
        <button
          onClick={handleSaveDefault}
          disabled={savingDefault}
          className="mt-4 w-full flex items-center justify-center gap-2 py-3 border border-green-500 hover:bg-green-500 hover:text-black transition-all text-xs font-black disabled:opacity-50"
        >
          {savingDefault ? 'EXECUTING_SAVE...' : '[ SAVE_DEFAULT_CAPTION ]'}
          {!savingDefault && <Save size={14} />}
        </button>
      </section>

      {/* New Pack Caption Section */}
      <section className="border border-green-900 p-6 bg-green-950/10">
        <div className="flex items-center gap-2 mb-4 text-xs font-bold uppercase tracking-widest">
          <AlertCircle size={14} className="text-amber-500" />
          LEGENDA_NOVO_PACK
        </div>
        <p className="text-[10px] opacity-60 mb-4 leading-relaxed">
          Legenda especial para quando um NOVO PACK é detectado.
          Geralmente contém informações sobre o conteúdo e links.
        </p>
        <textarea
          value={newPackCaption}
          onChange={(e) => setNewPackCaption(e.target.value)}
          placeholder="Digite a legenda de novo pack..."
          className="w-full h-48 bg-black border border-green-900 focus:border-amber-500/50 outline-none p-4 text-sm resize-none transition-colors"
        />
        <button
          onClick={handleSaveNewPack}
          disabled={savingNewPack}
          className="mt-4 w-full flex items-center justify-center gap-2 py-3 border border-amber-500 text-amber-500 hover:bg-amber-500 hover:text-black transition-all text-xs font-black disabled:opacity-50"
        >
          {savingNewPack ? 'EXECUTING_SAVE...' : '[ SAVE_NEW_PACK_CAPTION ]'}
          {!savingNewPack && <Save size={14} />}
        </button>
      </section>

      {/* Formatting Tips */}
      <div className="lg:col-span-2 border-t border-green-900/30 pt-6 mt-4">
        <div className="text-[10px] opacity-40 uppercase mb-2">Formatting_Tips:</div>
        <div className="flex flex-wrap gap-x-8 gap-y-2 text-[10px] font-mono opacity-60">
          <div>&lt;b&gt;Negrito&lt;/b&gt;</div>
          <div>&lt;i&gt;Itálico&lt;/i&gt;</div>
          <div>&lt;u&gt;Sublinhado&lt;/u&gt;</div>
          <div>&lt;code&gt;Monospace&lt;/code&gt;</div>
          <div>&lt;a href="..."&gt;Link&lt;/a&gt;</div>
        </div>
      </div>
    </div>
  );
}
