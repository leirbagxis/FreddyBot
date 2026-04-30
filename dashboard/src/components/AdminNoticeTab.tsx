import { Dispatch, SetStateAction } from 'react';
import {
    Users, Hash, Globe, MousePointerClick,
    Trash2, Link2, MessageSquare, Plus
} from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';
import { NoticeButton } from '../api';

interface AdminNoticeTabProps {
    noticeMessage: string;
    setNoticeMessage: Dispatch<SetStateAction<string>>;
    noticeTarget: 'channels' | 'users' | 'all';
    setNoticeTarget: Dispatch<SetStateAction<'channels' | 'users' | 'all'>>;
    noticeButtons: NoticeButton[];
    handleAddNoticeButton: () => void;
    updateNoticeButton: (index: number, field: keyof NoticeButton, value: string) => void;
    removeNoticeButton: (index: number) => void;
    handleSendNotice: () => void;
    isSendingNotice: boolean;
}

export function AdminNoticeTab({
    noticeMessage, setNoticeMessage,
    noticeTarget, setNoticeTarget,
    noticeButtons, handleAddNoticeButton,
    updateNoticeButton, removeNoticeButton,
    handleSendNotice, isSendingNotice
}: AdminNoticeTabProps) {
    return (
        <div className="card space-y-3" style={{ padding: '16px' }}>
            <h3 className="text-[15px] font-bold">Nova Notificação Central</h3>

            <div className="space-y-1.5">
                <label className="text-[13px] font-semibold text-[var(--hint)]">Mensagem (Modo Avançado)</label>
                <RichTextEditor
                    value={noticeMessage}
                    onChange={setNoticeMessage}
                    placeholder="Digite o conteúdo da mensagem..."
                    rows={5}
                />
            </div>

            <div className="space-y-1.5">
                <label className="text-[13px] font-semibold text-[var(--hint)]">Destinatários</label>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-1.5">
                    <button
                        onClick={() => setNoticeTarget('all')}
                        className={`flex items-center justify-center gap-2 p-3 rounded-xl border ${noticeTarget === 'all'
                            ? 'bg-[var(--accent-soft)] border-[var(--accent)] text-[var(--accent)]'
                            : 'bg-[var(--surface)] border-[var(--border)] text-[var(--hint)] hover:bg-[var(--border)]'
                            } transition-all font-semibold text-sm`}
                    >
                        <Globe size={18} /> Todos
                    </button>
                    <button
                        onClick={() => setNoticeTarget('channels')}
                        className={`flex items-center justify-center gap-2 p-3 rounded-xl border ${noticeTarget === 'channels'
                            ? 'bg-[var(--accent-soft)] border-[var(--accent)] text-[var(--accent)]'
                            : 'bg-[var(--surface)] border-[var(--border)] text-[var(--hint)] hover:bg-[var(--border)]'
                            } transition-all font-semibold text-sm`}
                    >
                        <Hash size={18} /> Canais
                    </button>
                    <button
                        onClick={() => setNoticeTarget('users')}
                        className={`flex items-center justify-center gap-2 p-3 rounded-xl border ${noticeTarget === 'users'
                            ? 'bg-[var(--accent-soft)] border-[var(--accent)] text-[var(--accent)]'
                            : 'bg-[var(--surface)] border-[var(--border)] text-[var(--hint)] hover:bg-[var(--border)]'
                            } transition-all font-semibold text-sm`}
                    >
                        <Users size={16} /> Usuários
                    </button>
                </div>
            </div>

            <div className="space-y-3 pt-2 border-t border-[var(--border)] mt-2">
                <div className="flex items-center justify-between">
                    <div>
                        <label className="text-[14px] font-bold text-[var(--text)] flex items-center gap-2">
                            <MousePointerClick size={16} className="text-[var(--accent)]" />
                            Botões Anexados
                        </label>
                        <p className="text-[12px] text-[var(--hint)] mt-0.5">Adicione botões interativos abaixo da mensagem principal ({noticeButtons.length}/5)</p>
                    </div>
                    <button
                        onClick={handleAddNoticeButton}
                        disabled={noticeButtons.length >= 5}
                        className="flex items-center justify-center w-8 h-8 rounded-full bg-[var(--accent-soft)] text-[var(--accent)] hover:opacity-80 transition-all disabled:opacity-50"
                        title="Adicionar Botão"
                    >
                        <Plus size={18} />
                    </button>
                </div>

                <div className="space-y-2">
                    {noticeButtons.map((btn, idx) => (
                        <div key={idx} className="group relative focus-within:ring-2 focus-within:ring-[var(--accent)] rounded-xl border border-[var(--border)] bg-[var(--surface)] overflow-hidden transition-all">
                            {/* Header / Type Selector & Delete */}
                            <div className="flex items-center justify-between bg-[var(--background)] px-3 py-2 border-b border-[var(--border)] gap-2">
                                <div className="flex items-center flex-1">
                                    {btn.type === 'url' ? <Link2 size={14} className="text-[var(--hint)] mr-2" /> : <MessageSquare size={14} className="text-[var(--hint)] mr-2" />}
                                    <select
                                        className="bg-transparent text-[13px] font-medium text-[var(--text)] focus:outline-none flex-1"
                                        value={btn.type}
                                        onChange={(e) => updateNoticeButton(idx, 'type', e.target.value)}
                                    >
                                        <option value="url">Link Externo</option>
                                        <option value="callback">Ação Interna (Callback)</option>
                                    </select>
                                </div>
                                <button
                                    onClick={() => removeNoticeButton(idx)}
                                    className="text-[var(--danger)]/70 hover:text-[var(--danger)] hover:bg-[var(--danger-soft)] p-1.5 rounded-lg transition-colors"
                                    title="Remover botão"
                                >
                                    <Trash2 size={14} />
                                </button>
                            </div>

                            {/* Inputs */}
                            <div className="p-3 space-y-2 flex flex-col sm:flex-row sm:space-y-0 sm:gap-2">
                                <div className="flex-1">
                                    <label className="text-[11px] uppercase font-bold text-[var(--hint)] mb-1 block">Nome do Botão</label>
                                    <input
                                        placeholder="Ex: Entrar no Grupo"
                                        className="w-full bg-[var(--background)] text-[var(--text)] border border-[var(--border)] rounded-lg p-2.5 text-[13px] focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--hint)]"
                                        value={btn.text}
                                        onChange={(e) => updateNoticeButton(idx, 'text', e.target.value)}
                                        maxLength={30}
                                    />
                                </div>
                                <div className="flex-[1.5]">
                                    <label className="text-[11px] uppercase font-bold text-[var(--hint)] mb-1 block">Destino (Valor)</label>
                                    <input
                                        placeholder={btn.type === 'url' ? "https://t.me/seu_link" : "Comando ou texto"}
                                        className="w-full bg-[var(--background)] text-[var(--text)] border border-[var(--border)] rounded-lg p-2.5 text-[13px] focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--hint)]"
                                        value={btn.value}
                                        onChange={(e) => updateNoticeButton(idx, 'value', e.target.value)}
                                        maxLength={100}
                                    />
                                </div>
                            </div>
                        </div>
                    ))}
                    {noticeButtons.length === 0 && (
                        <div className="flex flex-col items-center justify-center p-6 border-2 border-dashed border-[var(--border)] rounded-xl text-[var(--hint)]">
                            <MousePointerClick size={24} className="mb-2 opacity-50" />
                            <p className="text-[13px] font-medium">Nenhum botão anexado</p>
                            <p className="text-[12px] opacity-70">Clique no + para adicionar elementos interativos</p>
                        </div>
                    )}
                </div>
            </div>

            <button
                className="btn btn-primary w-full mt-4"
                onClick={handleSendNotice}
                disabled={isSendingNotice}
            >
                {isSendingNotice ? 'Enviando...' : 'Disparar Mensagem'}
            </button>
        </div>
    );
}
