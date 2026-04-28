import { useState, useEffect } from 'react';
import { Plus, Trash2, Edit2, ExternalLink, Save, Loader2, GripVertical, ChevronUp, ChevronDown, ChevronLeft, ChevronRight } from 'lucide-react';
import { createButton, updateButton, deleteButton, updateButtonsLayout } from '../../api';
import { Channel, Button } from '../../types';
import { showToast } from '../Toast';

interface ButtonsTabProps {
  channel: Channel;
  onUpdate: (channel: Channel) => void;
}

interface ButtonForm {
  name: string;
  url: string;
}

export default function ButtonsTab({ channel, onUpdate }: ButtonsTabProps) {
  const [buttons, setButtons] = useState<Button[]>(channel.buttons || []);
  const [isAdding, setIsAdding] = useState(false);
  const [editingButton, setEditingButton] = useState<Button | null>(null);
  const [form, setForm] = useState<ButtonForm>({ name: '', url: '' });
  const [loading, setLoading] = useState(false);
  const [layoutSaving, setLayoutSaving] = useState(false);

  useEffect(() => {
    setButtons(channel.buttons || []);
  }, [channel.buttons]);

  const handleAddButton = async () => {
    if (!form.name || !form.url) {
      showToast.error('Preencha o nome e a URL');
      return;
    }

    setLoading(true);
    try {
      const response = await createButton(channel.id, form.name, form.url);
      if (response.success) {
        showToast.success('Botão criado com sucesso!');
        // O backend retorna o novo botão em response.data
        const newButton = response.data;
        const updatedButtons = [...buttons, newButton];
        setButtons(updatedButtons);
        onUpdate({ ...channel, buttons: updatedButtons });
        setIsAdding(false);
        setForm({ name: '', url: '' });
      } else {
        showToast.error(response.message || 'Erro ao criar botão');
      }
    } catch (err) {
      console.error('Failed to add button', err);
      showToast.error('Erro ao conectar com o servidor');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateButton = async () => {
    if (!editingButton || !form.name || !form.url) return;

    setLoading(true);
    try {
      const response = await updateButton(channel.id, editingButton.buttonId, form.name, form.url);
      if (response.success) {
        showToast.success('Botão atualizado!');
        const updatedButtons = buttons.map(b => 
          b.buttonId === editingButton.buttonId ? { ...b, nameButton: form.name, buttonUrl: form.url } : b
        );
        setButtons(updatedButtons);
        onUpdate({ ...channel, buttons: updatedButtons });
        setEditingButton(null);
        setForm({ name: '', url: '' });
      } else {
        showToast.error(response.message || 'Erro ao atualizar botão');
      }
    } catch (err) {
      console.error('Failed to update button', err);
      showToast.error('Erro ao conectar com o servidor');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteButton = async (buttonId: string) => {
    if (!confirm('Tem certeza que deseja excluir este botão?')) return;

    setLoading(true);
    try {
      const response = await deleteButton(channel.id, buttonId);
      if (response.success) {
        showToast.success('Botão excluído!');
        const updatedButtons = buttons.filter(b => b.buttonId !== buttonId);
        setButtons(updatedButtons);
        onUpdate({ ...channel, buttons: updatedButtons });
      } else {
        showToast.error(response.message || 'Erro ao excluir botão');
      }
    } catch (err) {
      console.error('Failed to delete button', err);
      showToast.error('Erro ao conectar com o servidor');
    } finally {
      setLoading(false);
    }
  };

  const startEdit = (button: Button) => {
    setEditingButton(button);
    setForm({ name: button.nameButton, url: button.buttonUrl });
    setIsAdding(false);
  };

  const cancelForm = () => {
    setIsAdding(false);
    setEditingButton(null);
    setForm({ name: '', url: '' });
  };

  // Organizar botões em linhas e colunas para visualização
  const rows: Button[][] = [];
  buttons.forEach(btn => {
    const y = btn.positionY || 0;
    if (!rows[y]) rows[y] = [];
    rows[y][btn.positionX || 0] = btn;
  });

  // Limpar buracos em linhas/colunas
  const cleanRows = rows.map(row => row?.filter(btn => !!btn) || []).filter(row => row.length > 0);

  const handleSaveLayout = async () => {
    setLayoutSaving(true);
    try {
      // O layout deve ser [][]ButtonLayoutItem
      const layoutData = cleanRows.map(row => 
        row.map(btn => ({ id: btn.buttonId }))
      );

      const response = await updateButtonsLayout(channel.id, { layout: layoutData });
      if (response.success) {
        showToast.success('Layout salvo com sucesso!');
        // O backend deve retornar os botões com as novas posições
        if (response.data && response.data.buttons) {
           onUpdate({ ...channel, buttons: response.data.buttons });
        }
      } else {
        showToast.error(response.message || 'Erro ao salvar layout');
      }
    } catch (err) {
      console.error('Failed to save layout', err);
      showToast.error('Erro ao conectar com o servidor');
    } finally {
      setLayoutSaving(false);
    }
  };

  const moveButton = (btnId: string, direction: 'up' | 'down' | 'left' | 'right') => {
    const newButtons = [...buttons];
    const btnIndex = newButtons.findIndex(b => b.buttonId === btnId);
    if (btnIndex === -1) return;

    const btn = { ...newButtons[btnIndex] };
    
    if (direction === 'up' && btn.positionY > 0) btn.positionY--;
    if (direction === 'down') btn.positionY++;
    if (direction === 'left' && btn.positionX > 0) btn.positionX--;
    if (direction === 'right') btn.positionX++;

    newButtons[btnIndex] = btn;
    setButtons(newButtons);
    // Nota: Aqui apenas atualizamos o estado local, o usuário deve clicar em Salvar Layout
  };

  return (
    <div className="space-y-8">
      {/* Botões Management Section */}
      <section className="border border-green-900 p-6 bg-green-950/10">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2 text-xs font-bold uppercase tracking-widest">
            <GripVertical size={14} className="text-green-500" />
            GERENCIAR_BOTOES
          </div>
          {!isAdding && !editingButton && (
            <button
              onClick={() => setIsAdding(true)}
              className="flex items-center gap-2 px-4 py-2 border border-green-500 text-green-500 hover:bg-green-500 hover:text-black transition-all text-[10px] font-black"
            >
              <Plus size={14} />
              ADICIONAR_BOTAO
            </button>
          )}
        </div>

        {/* Form para Adicionar/Editar */}
        {(isAdding || editingButton) && (
          <div className="mb-8 p-4 border border-green-500/30 bg-black animate-in slide-in-from-top-4 duration-300">
            <div className="text-[10px] font-bold text-green-500 uppercase mb-4">
              {editingButton ? `EDITANDO_BOTAO: ${editingButton.buttonId}` : 'NOVO_BOTAO'}
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div className="space-y-1">
                <label className="text-[9px] uppercase opacity-40">Nome do Botão</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Ex: 📂 Canal VIP"
                  className="w-full bg-black border border-green-900 focus:border-green-500 outline-none p-3 text-sm transition-colors"
                />
              </div>
              <div className="space-y-1">
                <label className="text-[9px] uppercase opacity-40">URL do Botão</label>
                <input
                  type="text"
                  value={form.url}
                  onChange={(e) => setForm({ ...form, url: e.target.value })}
                  placeholder="https://t.me/..."
                  className="w-full bg-black border border-green-900 focus:border-green-500 outline-none p-3 text-sm transition-colors"
                />
              </div>
            </div>
            <div className="flex gap-2">
              <button
                onClick={editingButton ? handleUpdateButton : handleAddButton}
                disabled={loading}
                className="flex-1 flex items-center justify-center gap-2 py-3 bg-green-500 text-black hover:bg-green-400 transition-all text-xs font-black disabled:opacity-50"
              >
                {loading ? <Loader2 className="animate-spin" size={14} /> : editingButton ? '[ UPDATE ]' : '[ CREATE ]'}
              </button>
              <button
                onClick={cancelForm}
                className="px-6 py-3 border border-red-900 text-red-500 hover:bg-red-900 hover:text-white transition-all text-xs font-black"
              >
                CANCELAR
              </button>
            </div>
          </div>
        )}

        {/* Lista de Botões */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {buttons.length === 0 && !isAdding && (
            <div className="col-span-full py-12 text-center border border-dashed border-green-900 opacity-40 text-xs">
              NENHUM_BOTAO_ENCONTRADO_NA_DATABASE
            </div>
          )}
          {buttons.map((btn) => (
            <div key={btn.buttonId} className="group border border-green-900 p-4 bg-black hover:border-green-500 transition-colors">
              <div className="flex justify-between items-start mb-2">
                <div className="font-bold text-sm truncate pr-2">{btn.nameButton}</div>
                <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button onClick={() => startEdit(btn)} className="p-1 hover:text-green-400"><Edit2 size={12} /></button>
                  <button onClick={() => handleDeleteButton(btn.buttonId)} className="p-1 hover:text-red-400"><Trash2 size={12} /></button>
                </div>
              </div>
              <div className="text-[10px] opacity-40 truncate mb-4 flex items-center gap-1">
                <ExternalLink size={10} />
                {btn.buttonUrl}
              </div>
              <div className="flex items-center justify-between">
                <div className="text-[9px] font-mono border border-green-900/50 px-2 py-1 text-green-500/50">
                  POS_X:{btn.positionX} POS_Y:{btn.positionY}
                </div>
                <div className="flex gap-1">
                  <button onClick={() => moveButton(btn.buttonId, 'left')} className="p-1 border border-green-900/30 hover:border-green-500 text-[8px]"><ChevronLeft size={10}/></button>
                  <button onClick={() => moveButton(btn.buttonId, 'right')} className="p-1 border border-green-900/30 hover:border-green-500 text-[8px]"><ChevronRight size={10}/></button>
                  <button onClick={() => moveButton(btn.buttonId, 'up')} className="p-1 border border-green-900/30 hover:border-green-500 text-[8px]"><ChevronUp size={10}/></button>
                  <button onClick={() => moveButton(btn.buttonId, 'down')} className="p-1 border border-green-900/30 hover:border-green-500 text-[8px]"><ChevronDown size={10}/></button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Visual Layout Grid Section */}
      <section className="border border-green-900 p-6 bg-green-950/10">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2 text-xs font-bold uppercase tracking-widest">
            <Save size={14} className="text-green-500" />
            LAYOUT_VISUAL_GRID
          </div>
          <button
            onClick={handleSaveLayout}
            disabled={layoutSaving || buttons.length === 0}
            className="flex items-center gap-2 px-6 py-2 border border-green-500 text-green-500 hover:bg-green-500 hover:text-black transition-all text-[10px] font-black disabled:opacity-30"
          >
            {layoutSaving ? <Loader2 className="animate-spin" size={14} /> : '[ SALVAR_LAYOUT ]'}
          </button>
        </div>

        <div className="bg-black/50 border border-green-900/30 p-4 min-h-[200px] flex flex-col gap-2">
          {cleanRows.length === 0 ? (
             <div className="h-40 flex items-center justify-center text-[10px] uppercase opacity-20">
               Grid_Empty
             </div>
          ) : (
            cleanRows.map((row, y) => (
              <div key={y} className="flex gap-2">
                {row.map((btn) => (
                  <div 
                    key={btn.buttonId} 
                    className="flex-1 min-w-0 border border-green-500/50 bg-green-500/5 p-2 text-center text-[10px] font-bold truncate"
                  >
                    {btn.nameButton}
                  </div>
                ))}
              </div>
            ))
          )}
        </div>
        
        <p className="mt-4 text-[9px] opacity-40 leading-relaxed uppercase">
          &gt; Use as setas nos botões acima para ajustar a posição (X=Coluna, Y=Linha).<br/>
          &gt; O layout acima é uma prévia de como os botões aparecerão no Telegram.<br/>
          &gt; Não esqueça de clicar em "SALVAR_LAYOUT" após as alterações.
        </p>
      </section>
    </div>
  );
}
