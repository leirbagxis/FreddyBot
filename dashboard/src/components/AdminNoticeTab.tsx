import { Dispatch, SetStateAction, useState } from 'react';
import {
    Users, Hash, Globe, MousePointerClick,
    Trash2, Link2, MessageSquare, Plus, Image as ImageIcon
} from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';
import { NoticeButton } from '../api';
import { ConfirmModal } from './ConfirmModal';

interface AdminNoticeTabProps {
    noticeMessage: string;
    setNoticeMessage: Dispatch<SetStateAction<string>>;
    noticeImageUrl: string;
    setNoticeImageUrl: Dispatch<SetStateAction<string>>;
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
    noticeImageUrl, setNoticeImageUrl,
    noticeTarget, setNoticeTarget,
    noticeButtons, handleAddNoticeButton,
    updateNoticeButton, removeNoticeButton,
    handleSendNotice, isSendingNotice
}: AdminNoticeTabProps) {
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);

    const maxChars = noticeImageUrl.trim() ? 1024 : 4096;
    const isOverLimit = noticeMessage.length > maxChars;
    const hasEmptyButtons = noticeButtons.some(b => !b.text.trim() || !b.value.trim());
    const isReady = noticeMessage.trim().length > 0 && !isOverLimit && !hasEmptyButtons;

    const renderPreview = () => {
        // Convert some basic markdown to HTML for preview
        let htmlContent = noticeMessage
            .replace(/\*\*(.*?)\*\*/g, '<b>$1</b>')
            .replace(/__(.*?)__/g, '<i>$1</i>')
            .replace(/~~(.*?)~~/g, '<s>$1</s>')
            .replace(/\|\|(.*?)\|\|/g, '<span class="spoiler bg-[var(--surface)] text-transparent hover:text-[var(--text)] transition-colors px-1 rounded cursor-pointer">$1</span>')
            .replace(/`([^`]+)`/g, '<code class="bg-[var(--surface)] px-1 py-0.5 rounded text-[12px] font-mono text-[var(--accent)]">$1</code>')
            .replace(/\n/g, '<br/>');

        return (
            <div className="bg-[var(--surface)] shadow-sm p-3 rounded-2xl rounded-bl-sm max-w-[320px] w-full mx-auto text-[14px] text-[var(--text)] leading-relaxed">
                {noticeImageUrl && (
                    <img src={noticeImageUrl} alt="Preview" className="w-full rounded-xl mb-2 object-cover max-h-[200px] bg-[var(--background)]" onError={(e) => (e.currentTarget.style.display = 'none')} />
                )}
                <div dangerouslySetInnerHTML={{ __html: htmlContent || '<span class="text-[var(--hint)] opacity-50 font-medium">Sua mensagem aparecerá aqui...</span>' }} className="mb-2 break-words" />
                {noticeButtons.length > 0 && (
                    <div className="flex flex-col gap-1.5 mt-3 pt-2 border-t border-[var(--border)]">
                        {noticeButtons.map((btn, i) => (
                            <div key={i} className="bg-[var(--background)] hover:bg-[var(--border)] transition-colors border border-[var(--border)] rounded-xl py-2 px-3 text-center text-[var(--accent)] font-semibold text-[13px] cursor-pointer">
                                {btn.text || 'Botão'}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        );
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 animate-fade-in">
            {/* Form */}
            <div className="card space-y-4" style={{ padding: '16px' }}>
                <h3 className="text-[16px] font-bold">Configurar Disparo</h3>

                <div className="space-y-1.5">
                    <label className="text-[13px] font-semibold text-[var(--hint)] flex items-center gap-1.5">
                        <ImageIcon size={14} /> URL da Imagem / GIF (Opcional)
                    </label>
                    <input
                        type="text"
                        placeholder="https://exemplo.com/imagem.jpg"
                        value={noticeImageUrl}
                        onChange={(e) => setNoticeImageUrl(e.target.value)}
                        className="w-full bg-[var(--background)] text-[var(--text)] border border-[var(--border)] rounded-lg p-2.5 text-[13px] focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--hint)]"
                    />
                </div>

                <div className="space-y-1.5">
                    <div className="flex items-center justify-between">
                        <label className="text-[13px] font-semibold text-[var(--hint)]">Mensagem (Suporta Markdown)</label>
                        <span className={`text-[12px] font-medium ${isOverLimit ? 'text-[var(--danger)]' : 'text-[var(--hint)]'}`}>
                            {noticeMessage.length} / {maxChars}
                        </span>
                    </div>
                    <RichTextEditor
                        value={noticeMessage}
                        onChange={setNoticeMessage}
                        placeholder="Digite o conteúdo da mensagem..."
                        rows={6}
                    />
                </div>

                <div className="space-y-1.5">
                    <label className="text-[13px] font-semibold text-[var(--hint)]">Público-Alvo</label>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-1.5">
                        <button
                            onClick={() => setNoticeTarget('all')}
                            className={`flex items-center justify-center gap-2 p-2.5 rounded-xl border ${noticeTarget === 'all'
                                ? 'bg-[var(--accent-soft)] border-[var(--accent)] text-[var(--accent)]'
                                : 'bg-[var(--surface)] border-[var(--border)] text-[var(--hint)] hover:bg-[var(--border)]'
                                } transition-all font-semibold text-sm`}
                        >
                            <Globe size={16} /> Todos
                        </button>
                        <button
                            onClick={() => setNoticeTarget('channels')}
                            className={`flex items-center justify-center gap-2 p-2.5 rounded-xl border ${noticeTarget === 'channels'
                                ? 'bg-[var(--accent-soft)] border-[var(--accent)] text-[var(--accent)]'
                                : 'bg-[var(--surface)] border-[var(--border)] text-[var(--hint)] hover:bg-[var(--border)]'
                                } transition-all font-semibold text-sm`}
                        >
                            <Hash size={16} /> Canais
                        </button>
                        <button
                            onClick={() => setNoticeTarget('users')}
                            className={`flex items-center justify-center gap-2 p-2.5 rounded-xl border ${noticeTarget === 'users'
                                ? 'bg-[var(--accent-soft)] border-[var(--accent)] text-[var(--accent)]'
                                : 'bg-[var(--surface)] border-[var(--border)] text-[var(--hint)] hover:bg-[var(--border)]'
                                } transition-all font-semibold text-sm`}
                        >
                            <Users size={16} /> Usuários
                        </button>
                    </div>
                </div>

                <div className="space-y-3 pt-3 border-t border-[var(--border)] mt-2">
                    <div className="flex items-center justify-between">
                        <div>
                            <label className="text-[14px] font-bold text-[var(--text)] flex items-center gap-2">
                                <MousePointerClick size={16} className="text-[var(--accent)]" />
                                Botões Inline
                            </label>
                            <p className="text-[12px] text-[var(--hint)] mt-0.5">Adicione botões interativos ({noticeButtons.length}/5)</p>
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
                                        className="text-[var(--danger)]/70 hover:text-[var(--danger)] hover:bg-[var(--danger-soft)] p-1 rounded-lg transition-colors"
                                    >
                                        <Trash2 size={14} />
                                    </button>
                                </div>
                                <div className="p-3 space-y-2 flex flex-col sm:flex-row sm:space-y-0 sm:gap-2">
                                    <div className="flex-1">
                                        <input
                                            placeholder="Nome (Ex: Entrar)"
                                            className="w-full bg-[var(--background)] text-[var(--text)] border border-[var(--border)] rounded-lg p-2 text-[13px] focus:outline-none focus:border-[var(--accent)] placeholder:text-[var(--hint)]"
                                            value={btn.text}
                                            onChange={(e) => updateNoticeButton(idx, 'text', e.target.value)}
                                            maxLength={30}
                                        />
                                    </div>
                                    <div className="flex-[1.5]">
                                        <input
                                            placeholder={btn.type === 'url' ? "https://..." : "Comando"}
                                            className="w-full bg-[var(--background)] text-[var(--text)] border border-[var(--border)] rounded-lg p-2 text-[13px] focus:outline-none focus:border-[var(--accent)] placeholder:text-[var(--hint)]"
                                            value={btn.value}
                                            onChange={(e) => updateNoticeButton(idx, 'value', e.target.value)}
                                            maxLength={100}
                                        />
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                <button
                    className="btn w-full mt-4 bg-[var(--accent)] text-white hover:opacity-90 disabled:opacity-50 font-bold py-3.5 rounded-xl transition-all shadow-md shadow-[var(--accent)]/20"
                    onClick={() => setIsConfirmOpen(true)}
                    disabled={isSendingNotice || !isReady}
                >
                    {isSendingNotice ? 'Enviando...' : 'Revisar & Disparar Mensagem'}
                </button>
            </div>

            {/* Preview Panel */}
            <div className="card space-y-3 flex flex-col h-full" style={{ padding: '16px' }}>
                <h3 className="text-[16px] font-bold flex items-center justify-between">
                    Pré-visualização
                    <span className="text-[11px] font-medium bg-[var(--accent-soft)] text-[var(--accent)] px-2.5 py-1 rounded-full tracking-wide uppercase">Telegram View</span>
                </h3>
                <div className="flex-1 bg-gradient-to-br from-[var(--background)] to-[var(--surface)] border border-[var(--border)] rounded-2xl p-4 flex flex-col justify-center relative overflow-hidden min-h-[300px]">
                    <div className="absolute inset-0 opacity-10 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')]" />
                    <div className="relative z-10 w-full flex justify-center">
                        {renderPreview()}
                    </div>
                </div>
            </div>

            <ConfirmModal
                open={isConfirmOpen}
                onClose={() => setIsConfirmOpen(false)}
                onConfirm={handleSendNotice}
                title="Confirmar Disparo em Massa"
                message={`Você está prestes a enviar uma mensagem para ${noticeTarget === 'all' ? 'todos os usuários e canais cadastrados' : noticeTarget === 'channels' ? 'todos os canais cadastrados' : 'todos os usuários do bot'}. Tem certeza que deseja prosseguir?`}
                confirmText="Sim, Disparar Agora"
                danger={true}
            />
        </div>
    );
}
