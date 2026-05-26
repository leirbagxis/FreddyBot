import { useState, useRef, useCallback, useEffect } from 'react';
import { Button } from '../types';
import {
  Plus, Minus, GripVertical, Pencil, Trash2, ExternalLink,
  X, Check, AlertTriangle, Grid3X3
} from 'lucide-react';

interface Props {
  buttons: Button[];
  reactions: string;
  reactionPosition: number;
  channelId: number;
  onAdd: (button: Button) => void;
  onDelete: (buttonId: string) => void;
  onEdit: (buttonId: string, updates: Partial<Button>) => void;
  onMove: (buttonId: string, x: number, y: number) => void;
  onMoveReactions: (y: number) => void;
}

export function ButtonGrid({ buttons, reactions, reactionPosition, channelId, onAdd, onDelete, onEdit, onMove, onMoveReactions }: Props) {
  const [cols, setCols] = useState(() => Math.max(4, buttons.reduce((m, b) => Math.max(m, b.positionX), 0) + 1));
  const [rows, setRows] = useState(() => {
    const maxBtnY = buttons.reduce((m, b) => Math.max(m, b.positionY), -1);
    return Math.max(3, Math.max(maxBtnY, reactionPosition) + 2);
  });

  // Sync grid dimensions when props change (e.g. after move or external update)
  useEffect(() => {
    const maxBRow = buttons.reduce((m, b) => Math.max(m, b.positionY), -1);
    const neededRows = Math.max(maxBRow, reactionPosition) + 2;
    setRows(prev => Math.max(prev, neededRows));

    const maxBCol = buttons.reduce((m, b) => Math.max(m, b.positionX), -1);
    const neededCols = maxBCol + 1;
    setCols(prev => Math.max(prev, neededCols));
  }, [buttons, reactionPosition]);

  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState('');
  const [editUrl, setEditUrl] = useState('');
  const [addingAt, setAddingAt] = useState<{ x: number; y: number } | null>(null);
  const [newName, setNewName] = useState('');
  const [newUrl, setNewUrl] = useState('');
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null);

  const [dragBtnId, setDragBtnId] = useState<string | null>(null);
  const [dragOverKey, setDragOverKey] = useState<string | null>(null);
  const gridRef = useRef<HTMLDivElement>(null);
  const cloneRef = useRef<HTMLDivElement | null>(null);
  const cellRectsRef = useRef<Map<string, DOMRect>>(new Map());
  const longPressRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isDragging = useRef(false);
  const touchStart = useRef({ x: 0, y: 0 });
  const raf = useRef(0);

  const btnAt = (x: number, y: number) => buttons.find(b => b.positionX === x && b.positionY === y);
  const cellKey = (x: number, y: number) => `${x},${y}`;

  const processUrl = (url: string) => {
    const u = url.trim();
    const lower = u.toLowerCase();
    if (u.startsWith('@')) return `https://t.me/${u.slice(1).replace(/^\/+/, '')}`;
    if (lower.startsWith('t.me/')) return `https://t.me/${u.slice(5).replace(/^\/+/, '')}`;
    if (lower.startsWith('telegram.me/')) return `https://t.me/${u.slice('telegram.me/'.length).replace(/^\/+/, '')}`;
    if (!u.includes('://') && u.includes('.')) return `https://${u}`;
    return u;
  };

  const validateUrl = (url: string) => {
    const u = processUrl(url);
    if (!u) return false;
    try {
      const parsed = new URL(u);
      if (parsed.protocol === 'tg:') return true;
      if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') return false;
      if (!parsed.hostname) return false;
      if ((parsed.hostname === 't.me' || parsed.hostname === 'telegram.me') && parsed.pathname.replace(/\//g, '') === '') return false;
      return true;
    } catch {
      return false;
    }
  };

  const cacheCells = useCallback(() => {
    cellRectsRef.current.clear();
    gridRef.current?.querySelectorAll('[data-cell]').forEach(el => {
      const k = el.getAttribute('data-cell');
      if (k) cellRectsRef.current.set(k, el.getBoundingClientRect());
    });
  }, []);

  const findCell = (px: number, py: number): string | null => {
    for (const [k, r] of cellRectsRef.current) {
      if (px >= r.left && px <= r.right && py >= r.top && py <= r.bottom) {
        return k;
      }
    }
    return null;
  };

  const onDragStart = useCallback((e: React.DragEvent, id: string) => {
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', id);
    setDragBtnId(id);
    
    // For Desktop: Ensure the drag ghost is the full element
    if (id === 'REACTIONS_ROW' && e.currentTarget) {
      // Some browsers need this to correctly render the ghost of a spanned element
      // const target = e.currentTarget as HTMLElement;
      // e.dataTransfer.setDragImage(target, target.offsetWidth / 2, target.offsetHeight / 2);
    }
  }, []);

  const onDragEnd = useCallback(() => { setDragBtnId(null); setDragOverKey(null); }, []);

  const onCellDragOver = useCallback((e: React.DragEvent, x: number, y: number) => {
    e.preventDefault();
    setDragOverKey(cellKey(x, y));
  }, []);

  const onCellDrop = useCallback((e: React.DragEvent, x: number, y: number) => {
    e.preventDefault();
    const id = e.dataTransfer.getData('text/plain');
    if (id === 'REACTIONS_ROW') {
      onMoveReactions(y);
    } else if (id) {
      if (!btnAt(x, y)) {
        onMove(id, x, y);
      }
    }
    setDragBtnId(null);
    setDragOverKey(null);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [buttons, onMove, onMoveReactions]);

  const cleanTouch = useCallback(() => {
    if (longPressRef.current) { clearTimeout(longPressRef.current); longPressRef.current = null; }
    if (cloneRef.current) { cloneRef.current.remove(); cloneRef.current = null; }
    cancelAnimationFrame(raf.current);
    isDragging.current = false;
    setDragBtnId(null);
    setDragOverKey(null);
  }, []);

  const onTouchStart = useCallback((e: React.TouchEvent, id: string) => {
    const t = e.touches[0];
    const target = (e.target as HTMLElement).closest('[draggable]') as HTMLElement;
    if (!target) return;

    touchStart.current = { x: t.clientX, y: t.clientY };
    isDragging.current = false;
    
    longPressRef.current = setTimeout(() => {
      isDragging.current = true;
      setDragBtnId(id);
      setSelectedId(null);
      cacheCells();
      
      const c = target.cloneNode(true) as HTMLDivElement;
      c.className = `${target.className} drag-clone`;
      c.style.position = 'fixed';
      c.style.margin = '0';
      c.style.zIndex = '25000';
      c.style.width = `${target.offsetWidth}px`;
      c.style.height = `${target.offsetHeight}px`;
      c.style.pointerEvents = 'none';
      c.style.opacity = '0.8';
      c.style.background = id === 'REACTIONS_ROW' ? 'var(--accent-soft)' : 'var(--bg-secondary)';
      c.style.boxShadow = '0 10px 30px rgba(0,0,0,0.4)';
      c.style.borderRadius = '8px';
      c.style.transform = `translate3d(${t.clientX - target.offsetWidth / 2}px, ${t.clientY - target.offsetHeight / 2}px, 0)`;
      document.body.appendChild(c);
      cloneRef.current = c;

      if (navigator.vibrate) navigator.vibrate(15);
    }, 250);
  }, [cacheCells]);

  const onTouchMove = useCallback((e: React.TouchEvent) => {
    const t = e.touches[0];
    if (!isDragging.current) {
      const dx = Math.abs(t.clientX - touchStart.current.x);
      const dy = Math.abs(t.clientY - touchStart.current.y);
      if (dx > 8 || dy > 8) {
        if (longPressRef.current) { clearTimeout(longPressRef.current); longPressRef.current = null; }
      }
      return;
    }
    e.preventDefault();
    cancelAnimationFrame(raf.current);
    raf.current = requestAnimationFrame(() => {
      if (cloneRef.current) {
        cloneRef.current.style.transform = `translate3d(${t.clientX - cloneRef.current.offsetWidth / 2}px, ${t.clientY - cloneRef.current.offsetHeight / 2}px, 0)`;
      }
      setDragOverKey(findCell(t.clientX, t.clientY));
    });
  }, []);

  const onTouchEnd = useCallback(() => {
    if (isDragging.current && dragBtnId && dragOverKey) {
      const [xs, ys] = dragOverKey.split(',');
      const x = +xs, y = +ys;
      if (dragBtnId === 'REACTIONS_ROW') {
        onMoveReactions(y);
      } else if (!btnAt(x, y)) {
        onMove(dragBtnId, x, y);
      }
    }
    cleanTouch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dragBtnId, dragOverKey, buttons, onMove, onMoveReactions, cleanTouch]);

  const onCellClick = (x: number, y: number) => {
    if (isDragging.current) return;
    const b = btnAt(x, y);
    if (b) {
      setSelectedId(prev => prev === b.buttonId ? null : b.buttonId);
      setAddingAt(null); setEditingId(null);
    } else {
      setAddingAt({ x, y });
      setNewName(''); setNewUrl('');
      setSelectedId(null); setEditingId(null);
    }
  };

  const doAdd = () => {
    const finalUrl = processUrl(newUrl);
    if (!addingAt || !newName.trim() || !validateUrl(finalUrl)) return;
    onAdd({
      buttonId: '', // Temporarily empty, will be assigned real ID from DB via App.tsx
      nameButton: newName.trim(), buttonUrl: finalUrl,
      positionX: addingAt.x, positionY: addingAt.y,
      ownerChannelId: channelId,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    });
    setAddingAt(null);
  };

  const startEdit = (b: Button) => {
    setEditingId(b.buttonId);
    setEditName(b.nameButton);
    setEditUrl(b.buttonUrl);
    setSelectedId(null); setAddingAt(null);
  };

  const doEdit = () => {
    const finalUrl = processUrl(editUrl);
    if (!editingId || !editName.trim() || !validateUrl(finalUrl)) return;
    onEdit(editingId, { nameButton: editName.trim(), buttonUrl: finalUrl });
    setEditingId(null);
  };

  const adjustGrid = (axis: 'col' | 'row', dir: 1 | -1) => {
    if (axis === 'col') {
      const mx = buttons.reduce((m, b) => Math.max(m, b.positionX), -1);
      setCols(c => Math.max(Math.min(c + dir, 8), mx + 1, 1));
    } else {
      const mx = buttons.reduce((m, b) => Math.max(m, b.positionY), -1);
      const mxUsed = Math.max(mx, reactionPosition);
      setRows(r => Math.max(Math.min(r + dir, 10), mxUsed + 1, 1));
    }
  };

  const selBtn = selectedId ? buttons.find(b => b.buttonId === selectedId) : null;

  return (
    <div className="button-grid-content">
      <div className="card animate-stagger-in" style={{ animationDelay: '0s' }}>
        <div className="section-header">
          <div className="section-icon purple"><Grid3X3 size={18} /></div>
          <div className="flex-1 min-w-0">
            <h3 className="text-[15px] font-semibold">Botões Inline</h3>
            <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
              {buttons.length} botão(ões) • Segure para arrastar
            </p>
          </div>
        </div>

        {/* Grid controls */}
        <div className="flex items-center justify-between mb-3 gap-3">
          <div className="flex items-center gap-2">
            <span className="text-xs font-medium" style={{ color: 'var(--hint)' }}>Col</span>
            <button onClick={() => adjustGrid('col', -1)} className="icon-btn" style={{ width: 32, height: 32, borderRadius: 8 }}><Minus size={13} /></button>
            <span className="text-sm font-bold" style={{ width: 18, textAlign: 'center' }}>{cols}</span>
            <button onClick={() => adjustGrid('col', 1)} className="icon-btn" style={{ width: 32, height: 32, borderRadius: 8 }}><Plus size={13} /></button>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs font-medium" style={{ color: 'var(--hint)' }}>Lin</span>
            <button onClick={() => adjustGrid('row', -1)} className="icon-btn" style={{ width: 32, height: 32, borderRadius: 8 }}><Minus size={13} /></button>
            <span className="text-sm font-bold" style={{ width: 18, textAlign: 'center' }}>{rows}</span>
            <button onClick={() => adjustGrid('row', 1)} className="icon-btn" style={{ width: 32, height: 32, borderRadius: 8 }}><Plus size={13} /></button>
          </div>
        </div>

        {/* Grid */}
        <div className="btn-grid-wrapper" ref={gridRef}>
          <div
            className="grid gap-2"
            style={{
              gridTemplateColumns: `repeat(${cols}, 1fr)`,
              minWidth: cols > 4 ? cols * 76 : undefined,
              position: 'relative'
            }}
          >
            {/* 1. Render empty cells for drop targets */}
            {Array.from({ length: rows }, (_, y) =>
              Array.from({ length: cols }, (_, x) => {
                const k = cellKey(x, y);
                const hasBtn = !!btnAt(x, y);
                const isReac = y === reactionPosition;

                // When dragging reactions, highlight the WHOLE row
                const isOver = dragOverKey && (
                  dragBtnId === 'REACTIONS_ROW' 
                    ? dragOverKey.split(',')[1] === String(y)
                    : dragOverKey === k
                );
                
                return (
                  <div
                    key={`cell-${k}`}
                    data-cell={k}
                    className={`grid-cell ${isOver ? 'drag-over' : ''}`}
                    style={{
                      gridColumn: x + 1,
                      gridRow: y + 1,
                      height: '64px',
                      zIndex: 1, // Base layer for drop targets
                      pointerEvents: 'auto' // Must be auto to receive drops!
                    }}
                    onClick={() => {
                      if (isDragging.current || y === reactionPosition) return;
                      onCellClick(x, y);
                    }}
                    onDragOver={e => onCellDragOver(e, x, y)}
                    onDrop={e => onCellDrop(e, x, y)}
                    onDragLeave={() => setDragOverKey(null)}
                  >
                    {!hasBtn && !isReac && <Plus size={14} style={{ color: 'var(--hint)', opacity: 0.15, pointerEvents: 'none' }} />}
                  </div>
                );
              })
            )}

            {/* 2. Render buttons */}
            {buttons.map(b => {
              const isSel = selectedId === b.buttonId;
              const isSource = dragBtnId === b.buttonId;
              // Visual position: we use the raw Y from DB. 
              // If we want to show it BELOW the reactions if it's at the same row, we offset it.
              // BUT the user says: "em hipotese alguma um botao pode ocupar a linha x=0" (if x=0 is reactions).
              // So we just render at their exact coordinates.
              
              return (
                <div
                  key={`btn-${b.buttonId}`}
                  draggable
                  onDragStart={e => onDragStart(e, b.buttonId)}
                  onDragEnd={onDragEnd}
                  onTouchStart={e => onTouchStart(e, b.buttonId)}
                  onTouchMove={e => onTouchMove(e)}
                  onTouchEnd={onTouchEnd}
                  onClick={(e) => { 
                    e.stopPropagation(); 
                    if (!isDragging.current) onCellClick(b.positionX, b.positionY); 
                  }}
                  className={`grid-cell occupied ${isSource ? 'drag-source' : ''} ${isSel ? 'selected' : ''}`}
                  style={{
                    gridColumn: b.positionX + 1,
                    gridRow: b.positionY + 1,
                    zIndex: 10,
                    cursor: 'grab'
                  }}
                >
                  <div className="flex flex-col items-center justify-center w-full h-full select-none min-w-0" style={{ pointerEvents: 'none' }}>
                    <GripVertical size={12} style={{ color: 'var(--hint)', opacity: 0.25, marginBottom: 2 }} />
                    <span className="text-[11px] font-semibold truncate max-w-full px-1 leading-tight text-center">{b.nameButton}</span>
                  </div>
                </div>
              );
            })}

            {/* 3. Render Reactions Plate */}
            {(() => {
              const reactionsList = (reactions || '').split(',').filter(r => r.trim() !== '');
              const isSource = dragBtnId === 'REACTIONS_ROW';
              return (
                <div
                  key="reactions-plate"
                  draggable
                  onDragStart={e => onDragStart(e, 'REACTIONS_ROW')}
                  onDragEnd={onDragEnd}
                  onTouchStart={e => onTouchStart(e, 'REACTIONS_ROW')}
                  onTouchMove={e => onTouchMove(e)}
                  onTouchEnd={onTouchEnd}
                  className={`grid-cell reactions-plate ${isSource ? 'drag-source' : ''}`}
                  style={{
                    gridColumn: `1 / span ${cols}`,
                    gridRow: reactionPosition + 1,
                    background: 'rgba(var(--accent-rgb), 0.08)', // Darker, more subtle
                    border: '1.5px dashed rgba(var(--accent-rgb), 0.3)', // Less bright, thinner
                    borderRadius: '10px',
                    zIndex: 20,
                    height: '64px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    cursor: 'grab',
                    touchAction: 'none',
                    boxShadow: 'inset 0 0 10px rgba(0,0,0,0.1)' // Soft inner shadow for depth
                  }}
                >
                  <div className="flex items-center justify-center gap-3 w-full h-full text-[var(--accent)] font-bold opacity-70">
                    <GripVertical size={16} />
                    <span className="text-[10px] uppercase tracking-[0.2em] font-black">Reações</span>
                    <div className="flex gap-2">
                      {reactionsList.length > 0 ? reactionsList.map((r, i) => (
                        <span key={i} className="text-xl bg-[var(--card)] px-2 py-1 rounded-lg border border-[var(--border)] shadow-sm">{r}</span>
                      )) : <span className="text-[11px] opacity-40 font-normal italic">Nenhuma configurada</span>}
                    </div>
                  </div>
                </div>
              );
            })()}
          </div>
        </div>

        {/* Selected detail */}
        {selBtn && !editingId && (
          <div className="form-area">
            <div className="flex items-center justify-between gap-3 mb-3">
              <div className="min-w-0 flex-1">
                <h4 className="text-sm font-semibold truncate">{selBtn.nameButton}</h4>
                <p className="text-xs mt-1 truncate" style={{ color: 'var(--hint)', fontFamily: 'monospace' }}>{selBtn.buttonUrl}</p>
              </div>
            </div>
            <div className="flex gap-2 flex-wrap">
              <button className="btn btn-secondary btn-sm flex-1" onClick={() => startEdit(selBtn)}>
                <Pencil size={13} /> Editar
              </button>
              <a
                href={selBtn.buttonUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="btn btn-secondary btn-sm flex-1"
                style={{ textDecoration: 'none' }}
              >
                <ExternalLink size={13} /> Abrir
              </a>
              <button className="btn btn-danger btn-sm flex-1" onClick={() => setConfirmDeleteId(selBtn.buttonId)}>
                <Trash2 size={13} /> Excluir
              </button>
            </div>
          </div>
        )}

        {/* Edit form */}
        {editingId && (() => {
          const b = buttons.find(bt => bt.buttonId === editingId);
          if (!b) return null;
          const isValidBtn = editName.trim().length > 0 && validateUrl(editUrl);
          return (
            <div className="form-area space-y-3">
              <div className="flex items-center gap-3">
                <Pencil size={13} style={{ color: 'var(--accent)' }} />
                <span className="text-sm font-semibold">Editando "{b.nameButton}"</span>
              </div>
              <input className="input" value={editName} onChange={e => setEditName(e.target.value)} placeholder="Nome" autoFocus />
              <input
                className="input"
                style={{ fontFamily: 'monospace', fontSize: 14 }}
                value={editUrl}
                onChange={e => setEditUrl(e.target.value)}
                onBlur={e => setEditUrl(processUrl(e.target.value))}
                placeholder="https://t.me/username..."
              />
              {!validateUrl(editUrl) && editUrl.trim().length > 0 && (
                <p className="text-xs mt-1" style={{ color: 'var(--danger)' }}>Username do Telegram deve ter no mínimo 5 caracteres.</p>
              )}
              <div className="flex gap-2 justify-end mt-2">
                <button className="btn btn-secondary btn-sm" onClick={() => setEditingId(null)}>
                  <X size={13} /> Cancelar
                </button>
                <button className="btn btn-primary btn-sm" onClick={doEdit} disabled={!isValidBtn}>
                  <Check size={13} /> Salvar
                </button>
              </div>
            </div>
          );
        })()}

        {/* Add form */}
        {addingAt && (() => {
          const isValidBtn = newName.trim().length > 0 && validateUrl(newUrl);
          return (
            <div className="form-area space-y-3">
              <div className="flex items-center gap-3">
                <Plus size={13} style={{ color: 'var(--success)' }} />
                <span className="text-sm font-semibold">Novo botão</span>
              </div>
              <input className="input" value={newName} onChange={e => setNewName(e.target.value)} placeholder="Nome do botão" autoFocus />
              <input
                className="input"
                style={{ fontFamily: 'monospace', fontSize: 14 }}
                value={newUrl}
                onChange={e => setNewUrl(e.target.value)}
                onBlur={e => setNewUrl(processUrl(e.target.value))}
                placeholder="https://t.me/username..."
              />
              {!validateUrl(newUrl) && newUrl.trim().length > 0 && (
                <p className="text-xs mt-1" style={{ color: 'var(--danger)' }}>Username do Telegram deve ter no mínimo 5 caracteres.</p>
              )}
              <div className="flex gap-2 justify-end mt-2">
                <button className="btn btn-secondary btn-sm" onClick={() => setAddingAt(null)}>
                  <X size={13} /> Cancelar
                </button>
                <button className="btn btn-primary btn-sm" onClick={doAdd} disabled={!isValidBtn}>
                  <Check size={13} /> Adicionar
                </button>
              </div>
            </div>
          );
        })()}
      </div>

      {/* Delete confirm */}
      {confirmDeleteId && (
        <div className="overlay" onClick={() => setConfirmDeleteId(null)}>
          <div className="dialog" onClick={e => e.stopPropagation()}>
            <div className="dialog-handle" />
            <div className="flex items-center gap-4 mb-4">
              <div className="section-icon rose">
                <AlertTriangle size={20} />
              </div>
              <div className="min-w-0">
                <p className="font-semibold text-[15px]">Excluir Botão</p>
                <p className="text-sm mt-1" style={{ color: 'var(--hint)' }}>
                  Excluir "{buttons.find(b => b.buttonId === confirmDeleteId)?.nameButton}"?
                </p>
              </div>
            </div>
            <div className="flex gap-3 mt-5">
              <button className="btn btn-secondary flex-1" onClick={() => setConfirmDeleteId(null)}>
                Cancelar
              </button>
              <button
                className="btn btn-danger flex-1"
                onClick={() => { onDelete(confirmDeleteId); setConfirmDeleteId(null); setSelectedId(null); }}
              >
                <Trash2 size={15} /> Excluir
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
