import { useState, useEffect, useCallback } from 'react';
import { UserCaptionTemplate, Button as TemplateButton } from '../types';
import { listUserCaptionTemplates, createUserCaptionTemplate, updateUserCaptionTemplate, deleteUserCaptionTemplate, createUserCaptionTemplateButton, updateUserCaptionTemplateButton, deleteUserCaptionTemplateButton, updateUserCaptionTemplateLayout } from '../api';
import { RichTextEditor } from './RichTextEditor';
import { CaptionPreview } from './CaptionPreview';
import { ButtonGrid } from './ButtonGrid';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Input } from './ui/input';
import { Plus, Trash2, Loader2, Check, ChevronDown, ChevronRight, Hash, Layers } from 'lucide-react';

interface Props {
  toast: (msg: string, type: 'success' | 'error' | 'info') => void;
}

export function UserTemplatesManager({ toast }: Props) {
  const [templates, setTemplates] = useState<UserCaptionTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [newName, setNewName] = useState('');
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const captions = await listUserCaptionTemplates();
      setTemplates(captions || []);
    } catch {
      toast('Erro ao carregar templates', 'error');
    } finally {
      setLoading(false);
    }
  }, [toast]);

  useEffect(() => { load(); }, [load]);

  const handleCreate = async () => {
    if (!newName.trim()) { toast('Digite um nome curto', 'error'); return; }
    setSaving(true);
    try {
      const tpl = await createUserCaptionTemplate(newName.trim(), '');
      if (tpl) {
        tpl.buttons = tpl.buttons || [];
        setTemplates(prev => [tpl, ...prev]);
      }
      setNewName('');
      toast('Template criado', 'success');
    } catch {
      toast('Erro ao criar template', 'error');
    } finally {
      setSaving(false);
    }
  };

  const handleUpdateCaption = async (id: string, code: string, caption: string) => {
    try {
      await updateUserCaptionTemplate(id, code, caption);
      setTemplates(prev => prev.map(t =>
        t.id === id ? { ...t, code, caption } : t
      ));
      toast('Template salvo', 'success');
    } catch {
      toast('Erro ao atualizar template', 'error');
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.Telegram?.WebApp) {
      if (!window.confirm('Excluir este template?')) return;
      executeDelete(id);
      return;
    }
    
    window.Telegram.WebApp.showConfirm('Excluir este template?', (confirmed) => {
      if (confirmed) executeDelete(id);
    });
  };

  const executeDelete = async (id: string) => {
    setDeleting(id);
    try {
      await deleteUserCaptionTemplate(id);
      setTemplates(prev => prev.filter(t => t.id !== id));
      toast('Template excluído', 'success');
    } catch {
      toast('Erro ao excluir template', 'error');
    } finally {
      setDeleting(null);
    }
  };

  const handleAddButton = async (templateId: string, btn: TemplateButton) => {
    try {
      await createUserCaptionTemplateButton(templateId, btn.nameButton, btn.buttonUrl);
      load(); // Reload to get the real DB ID
      toast('Botão adicionado', 'success');
    } catch {
      toast('Erro ao adicionar botão', 'error');
    }
  };

  const handleEditButton = async (templateId: string, buttonId: string, updates: any) => {
    try {
      await updateUserCaptionTemplateButton(templateId, buttonId, updates.nameButton, updates.buttonUrl || '');
      setTemplates(prev => prev.map(t => {
        if (t.id === templateId) {
          return { ...t, buttons: t.buttons.map(b => b.buttonId === buttonId ? { ...b, ...updates } : b) };
        }
        return t;
      }));
    } catch {
      toast('Erro ao atualizar botão', 'error');
    }
  };

  const handleDeleteButton = async (templateId: string, buttonId: string) => {
    try {
      await deleteUserCaptionTemplateButton(templateId, buttonId);
      setTemplates(prev => prev.map(t => {
        if (t.id === templateId) {
          return { ...t, buttons: t.buttons.filter(b => b.buttonId !== buttonId) };
        }
        return t;
      }));
    } catch {
      toast('Erro ao remover botão', 'error');
    }
  };

  const handleMoveButton = async (templateId: string, layout: { buttonId: string }[][]) => {
    try {
      await updateUserCaptionTemplateLayout(templateId, layout);
      
      setTemplates(prev => prev.map(t => {
        if (t.id !== templateId) return t;
        
        const positionMap = new Map<string, {x: number, y: number}>();
        layout.forEach((row, y) => {
          row.forEach((col, x) => {
            if (col && col.buttonId) {
              positionMap.set(col.buttonId, { x, y });
            }
          });
        });
        
        const newButtons = t.buttons.map(b => {
          const pos = positionMap.get(b.buttonId);
          if (pos) {
            return { ...b, positionX: pos.x, positionY: pos.y };
          }
          return b;
        });
        
        return { ...t, buttons: newButtons };
      }));
      
      toast('Posição salva', 'success');
    } catch {
      toast('Erro ao atualizar layout', 'error');
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between py-3 px-4">
        <CardTitle className="text-sm font-semibold flex items-center gap-2">
          <Layers className="h-4 w-4" />
          Meus Templates
        </CardTitle>
        <Badge variant="secondary" className="text-[10px]">{templates.length}</Badge>
      </CardHeader>
      <CardContent className="px-4 pb-4 space-y-3">
        <div className="flex flex-col gap-3 p-4 bg-muted/30 border border-border rounded-lg">
          <div className="space-y-1">
            <h4 className="text-sm font-medium">Criar novo template</h4>
            <p className="text-xs text-muted-foreground">Escolha um nome curto para identificar a legenda e seus botões.</p>
          </div>
          <div className="flex items-center gap-3">
            <Input
              type="text"
              placeholder="Ex: promo_blackfriday"
              value={newName}
              onChange={e => setNewName(e.target.value)}
              className="flex-1 h-9"
              onKeyDown={e => {
                if (e.key === 'Enter' && newName.trim() && !saving) {
                  e.preventDefault();
                  handleCreate();
                }
              }}
            />
            <Button variant="default" size="sm" onClick={handleCreate} disabled={saving || !newName.trim()} className="h-9 px-4">
              {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4 mr-2" />}
              {saving ? <span className="ml-2">Criando...</span> : <span>Adicionar</span>}
            </Button>
          </div>
        </div>

        {loading && (
          <div className="flex justify-center py-6"><Loader2 className="h-5 w-5 animate-spin text-muted-foreground" /></div>
        )}

        {!loading && templates.length === 0 && (
          <p className="text-xs text-muted-foreground text-center py-4">
            Nenhum template. Crie um com o campo acima ou crie pelo Post Builder no bot.
          </p>
        )}

        {templates.map(tpl => {
          const title = tpl.code;
          const subTitle = tpl.caption;
          const buttonsCount = tpl.buttons.length;
          
          return (
            <div key={tpl.id} className="rounded-lg border border-border overflow-hidden">
              <div
                className="flex items-center gap-2 p-3 cursor-pointer hover:bg-muted/30"
                onClick={() => setExpandedId(expandedId === tpl.id ? null : tpl.id)}
              >
                <div className="flex items-center justify-center size-7 rounded-md shrink-0" style={{ background: 'var(--accent-soft)' }}>
                  <Hash className="h-3.5 w-3.5" style={{ color: 'var(--accent)' }} />
                </div>
                <div className="flex-1 min-w-0">
                  <span className="text-sm font-semibold">{title}</span>
                  <p className="text-xs text-muted-foreground truncate mt-0.5">{subTitle}</p>
                </div>
                <div className="flex items-center gap-1 shrink-0">
                  <Badge variant="secondary" className="text-[10px]">{buttonsCount}</Badge>
                  <Button variant="ghost" size="sm" className="h-7 w-7 p-0" onClick={e => { e.stopPropagation(); handleDelete(tpl.id); }} disabled={deleting === tpl.id} title="Excluir">
                    {deleting === tpl.id ? <Loader2 className="h-3 w-3 animate-spin text-red-500" /> : <Trash2 className="h-3 w-3 text-red-500" />}
                  </Button>
                  {expandedId === tpl.id ? <ChevronDown className="h-4 w-4 text-muted-foreground" /> : <ChevronRight className="h-4 w-4 text-muted-foreground" />}
                </div>
              </div>

              {expandedId === tpl.id && (
                <TemplateEditor
                  template={tpl}
                  onUpdateCaption={handleUpdateCaption}
                  onAddButton={(btn) => handleAddButton(tpl.id, btn)}
                  onEditButton={(buttonId, updates) => handleEditButton(tpl.id, buttonId, updates)}
                  onDeleteButton={(buttonId) => handleDeleteButton(tpl.id, buttonId)}
                  onMoveButton={(layout) => handleMoveButton(tpl.id, layout)}
                  toast={toast}
                />
              )}
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}

function TemplateEditor({
  template, onUpdateCaption, onAddButton, onEditButton, onDeleteButton, onMoveButton, toast,
}: {
  template: UserCaptionTemplate;
  onUpdateCaption: (id: string, code: string, caption: string) => void;
  onAddButton: (btn: TemplateButton) => void;
  onEditButton: (buttonId: string, updates: any) => void;
  onDeleteButton: (buttonId: string) => void;
  onMoveButton: (layout: { buttonId: string }[][]) => void;
  toast: (msg: string, type: 'success' | 'error' | 'info') => void;
}) {
  const [code, setCode] = useState(template.code);
  const [caption, setCaption] = useState(template.caption);
  const [editing, setEditing] = useState(false);
  const [dirty, setDirty] = useState(false);

  useEffect(() => { setCode(template.code); setCaption(template.caption); setDirty(false); setEditing(false); }, [template]);

  const handleSaveCaption = () => {
    if (!code.trim()) { toast('Nome curto é obrigatório', 'error'); return; }
    onUpdateCaption(template.id, code.trim(), caption);
    setDirty(false);
    setEditing(false);
  };

  const handleCancelEdit = () => {
    setCaption(template.caption);
    setDirty(false);
    setEditing(false);
  };

  return (
    <div className="border-t border-border p-3 space-y-3 bg-muted/10">
      <div>
        <label className="text-xs font-medium text-muted-foreground mb-1 block">Nome curto</label>
        <input
          type="text"
          value={code}
          onChange={e => { setCode(e.target.value); setDirty(true); }}
          className="w-full px-3 py-1.5 text-sm border rounded-md bg-background"
        />
      </div>

      <div>
        <label className="text-xs font-medium text-muted-foreground mb-1 flex items-center justify-between">
          <span>Legenda (Opcional)</span>
          <span className="text-[10px] text-muted-foreground/60 font-normal">Use variáveis como {`{nome}`}</span>
        </label>
        {editing ? (
          <div className="space-y-3">
            <div className="border rounded-md bg-background overflow-hidden">
              <RichTextEditor
                value={caption || ''}
                onChange={(val) => { setCaption(val); setDirty(true); }}
                placeholder="Digite a legenda..."
              />
            </div>
            <div className="flex items-center justify-end gap-2">
              <Button variant="secondary" size="sm" onClick={handleCancelEdit}>
                Cancelar
              </Button>
              <Button size="sm" onClick={handleSaveCaption} disabled={!dirty} className="h-8">
                <Check className="h-3.5 w-3.5 mr-1" />
                Salvar
              </Button>
            </div>
          </div>
        ) : (
          <div onClick={() => setEditing(true)}>
            <CaptionPreview text={caption} />
          </div>
        )}
      </div>

      <div className="mt-4 border-t border-border pt-4">
        <ButtonGrid
          buttons={template.buttons}
          reactions=""
          reactionPosition={-1}
          channelId={0}
          onAdd={onAddButton}
          onEdit={onEditButton}
          onDelete={onDeleteButton}
          onMove={(buttonId, x, y) => {
            const moved = template.buttons.map(button => button.buttonId === buttonId ? { ...button, positionX: x, positionY: y } : button);
            const rows = Array.from({ length: Math.max(...moved.map(button => button.positionY), 0) + 1 }, (_, row) =>
              moved
                .filter(button => button.positionY === row)
                .sort((a, b) => a.positionX - b.positionX)
                .map(button => ({ buttonId: button.buttonId }))
            );
            onMoveButton(rows);
          }}
          onMoveReactions={() => {}}
          hideReactions={true}
        />
      </div>
    </div>
  );
}
