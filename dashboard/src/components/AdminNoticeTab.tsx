import { Dispatch, SetStateAction, useState } from 'react';
import {
    Users, Hash, Globe, MousePointerClick,
    Trash2, Link2, MessageSquare, Plus, Image as ImageIcon,
    Send, Eye
} from 'lucide-react';
import { RichTextEditor } from './RichTextEditor';
import { NoticeButton, NoticeTarget } from '../api';
import { ConfirmModal } from './ConfirmModal';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Input } from './ui/input';
import { Textarea } from './ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';

interface AdminNoticeTabProps {
    noticeMessage: string;
    setNoticeMessage: Dispatch<SetStateAction<string>>;
    noticeImageUrl: string;
    setNoticeImageUrl: Dispatch<SetStateAction<string>>;
    noticeTarget: NoticeTarget;
    setNoticeTarget: Dispatch<SetStateAction<NoticeTarget>>;
    noticeTargetId: string;
    setNoticeTargetId: Dispatch<SetStateAction<string>>;
    noticeButtons: NoticeButton[];
    handleAddNoticeButton: () => void;
    updateNoticeButton: (index: number, field: keyof NoticeButton, value: string) => void;
    removeNoticeButton: (index: number) => void;
    handleSendNotice: () => void;
    isSendingNotice: boolean;
}

const targets: { id: NoticeTarget; label: string; icon: React.ReactNode; desc: string }[] = [
    { id: 'all', label: 'Todos', icon: <Globe size={16} />, desc: 'Todos os usuários e canais' },
    { id: 'channels', label: 'Canais', icon: <Hash size={16} />, desc: 'Apenas canais cadastrados' },
    { id: 'users', label: 'Usuários', icon: <Users size={16} />, desc: 'Apenas usuários do bot' },
    { id: 'single', label: 'Suporte', icon: <MousePointerClick size={16} />, desc: 'Mensagem para 1 usuário' },
    { id: 'user_ids', label: 'IDs Usuários', icon: <Users size={16} />, desc: 'Lista personalizada' },
    { id: 'channel_ids', label: 'IDs Canais', icon: <Hash size={16} />, desc: 'Lista personalizada' },
];

export function AdminNoticeTab({
    noticeMessage, setNoticeMessage,
    noticeImageUrl, setNoticeImageUrl,
    noticeTarget, setNoticeTarget,
    noticeTargetId, setNoticeTargetId,
    noticeButtons, handleAddNoticeButton,
    updateNoticeButton, removeNoticeButton,
    handleSendNotice, isSendingNotice
}: AdminNoticeTabProps) {
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);

    const maxChars = noticeImageUrl.trim() ? 1024 : 4096;
    const isOverLimit = noticeMessage.length > maxChars;
    const hasEmptyButtons = noticeButtons.some(b => !b.text.trim() || !b.value.trim());
    const specificTarget = noticeTarget === 'single' || noticeTarget === 'user_ids' || noticeTarget === 'channel_ids';
    const isReady = noticeMessage.trim().length > 0 && !isOverLimit && !hasEmptyButtons && (!specificTarget || noticeTargetId.trim().length > 5);

    const renderPreview = () => {
        let previewUrl = noticeImageUrl;
        if (noticeImageUrl && !noticeImageUrl.startsWith('http') && noticeImageUrl.length > 20) {
            previewUrl = `/api/admin/media-proxy/${noticeImageUrl}`;
        }

        const headerHtml = noticeTarget === 'single' || noticeTarget === 'user_ids' ? '<b>MENSAGEM DO SUPORTE</b><br/><br/>' : '';

        let htmlContent = noticeMessage
            .replace(/\*\*(.*?)\*\*/g, '<b>$1</b>')
            .replace(/__(.*?)__/g, '<i>$1</i>')
            .replace(/~~(.*?)~~/g, '<s>$1</s>')
            .replace(/\|\|(.*?)\|\|/g, '<span class="spoiler">$1</span>')
            .replace(/`([^`]+)`/g, '<code>$1</code>')
            .replace(/\n/g, '<br/>');

        return (
            <div className="bg-muted/30 shadow-sm p-3 rounded-2xl rounded-bl-sm max-w-[320px] w-full mx-auto text-[14px] text-foreground leading-relaxed">
                {noticeImageUrl && (
                    <img src={previewUrl} alt="Preview" className="w-full rounded-xl mb-2 object-contain max-h-[350px] bg-background" onError={(e) => (e.currentTarget.style.display = 'none')} />
                )}
                <div dangerouslySetInnerHTML={{ __html: headerHtml + (htmlContent || '<span class="text-muted-foreground/50 font-medium">Sua mensagem aparecerá aqui...</span>') }} className="mb-2 break-words" />
                {noticeButtons.length > 0 && (
                    <div className="flex flex-col gap-1.5 mt-3 pt-2 border-t border-border">
                        {noticeButtons.map((btn, i) => (
                            <div key={i} className="bg-background hover:bg-muted/50 transition-colors border border-border rounded-xl py-2 px-3 text-center text-accent font-semibold text-[13px] cursor-pointer">
                                {btn.text || 'Botão'}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        );
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* ─── Form ─── */}
            <div className="space-y-4">
                {/* Image URL */}
                <div className="rounded-xl border border-border p-4 space-y-2">
                    <label className="text-[12px] font-semibold text-muted-foreground flex items-center gap-1.5">
                        <ImageIcon size={14} /> URL da Imagem / GIF <Badge variant="secondary" className="text-[9px]">Opcional</Badge>
                    </label>
                    <Input
                        placeholder="https://exemplo.com/imagem.jpg"
                        value={noticeImageUrl}
                        onChange={(e) => setNoticeImageUrl(e.target.value)}
                        className="h-10 rounded-xl"
                    />
                </div>

                {/* Message */}
                <div className="rounded-xl border border-border p-4 space-y-2">
                    <div className="flex items-center justify-between">
                        <label className="text-[12px] font-semibold text-muted-foreground">Mensagem <span className="text-muted-foreground/50">(Markdown)</span></label>
                        <span className={`text-[11px] font-medium ${isOverLimit ? 'text-destructive' : 'text-muted-foreground'}`}>
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

                {/* Target */}
                <div className="rounded-xl border border-border p-4 space-y-3">
                    <label className="text-[12px] font-semibold text-muted-foreground">Público-Alvo</label>
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-1.5">
                        {targets.map((item) => (
                            <button
                                key={item.id}
                                onClick={() => setNoticeTarget(item.id)}
                                className={`flex flex-col items-center gap-1 p-2.5 rounded-xl border text-center transition-all ${
                                    noticeTarget === item.id
                                        ? 'border-accent bg-accent/10 text-accent'
                                        : 'border-border text-muted-foreground hover:border-muted-foreground/30'
                                }`}
                            >
                                <span className={noticeTarget === item.id ? 'text-accent' : 'text-muted-foreground/60'}>
                                    {item.icon}
                                </span>
                                <span className="text-[11px] font-semibold leading-tight">{item.label}</span>
                            </button>
                        ))}
                    </div>

                    {specificTarget && (
                        <div className="space-y-1.5 pt-1">
                            <label className="text-[11px] font-semibold text-muted-foreground">
                                {noticeTarget === 'channel_ids' ? 'IDs dos Canais' : noticeTarget === 'user_ids' ? 'IDs dos Usuários' : 'ID do Usuário'}
                            </label>
                            <Textarea
                                placeholder={noticeTarget === 'channel_ids' ? 'Ex: -1001234567890, -1009876543210' : 'Ex: 12345678, 987654321'}
                                value={noticeTargetId}
                                onChange={(e) => setNoticeTargetId(e.target.value)}
                                rows={noticeTarget === 'single' ? 1 : 3}
                                className="rounded-xl resize-none"
                            />
                            <p className="text-[11px] text-muted-foreground">
                                {noticeTarget === 'channel_ids' ? 'Canais específicos não recebem título de suporte.' : 'Separe IDs por vírgula, espaço ou quebra de linha.'}
                            </p>
                        </div>
                    )}
                </div>

                {/* Buttons */}
                <div className="rounded-xl border border-border p-4 space-y-3">
                    <div className="flex items-center justify-between">
                        <div>
                            <label className="text-[12px] font-semibold text-muted-foreground flex items-center gap-1.5">
                                <MousePointerClick size={14} /> Botões Inline
                            </label>
                            <p className="text-[11px] text-muted-foreground/60 mt-0.5">{noticeButtons.length}/5 adicionados</p>
                        </div>
                        <button
                            onClick={handleAddNoticeButton}
                            disabled={noticeButtons.length >= 5}
                            className="flex items-center justify-center w-8 h-8 rounded-full bg-accent/10 text-accent hover:bg-accent/20 transition-all disabled:opacity-30"
                            title="Adicionar Botão"
                        >
                            <Plus size={18} />
                        </button>
                    </div>

                    <div className="space-y-2">
                        {noticeButtons.map((btn, idx) => (
                            <div key={idx} className="rounded-xl border border-border overflow-hidden">
                                <div className="flex items-center justify-between px-3 py-2 border-b border-border bg-muted/20 gap-2">
                                    <div className="flex items-center gap-2 flex-1">
                                        {btn.type === 'url' ? <Link2 size={14} className="text-muted-foreground" /> : <MessageSquare size={14} className="text-muted-foreground" />}
                                        <Select value={btn.type} onValueChange={(v) => updateNoticeButton(idx, 'type', v)}>
                                            <SelectTrigger className="bg-transparent text-[12px] font-medium h-auto p-0 border-0 shadow-none gap-1">
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="url">Link Externo</SelectItem>
                                                <SelectItem value="callback">Callback</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                    <button
                                        onClick={() => removeNoticeButton(idx)}
                                        className="text-destructive/50 hover:text-destructive hover:bg-destructive/10 p-1 rounded-lg transition-colors"
                                    >
                                        <Trash2 size={14} />
                                    </button>
                                </div>
                                <div className="p-3 flex flex-col sm:flex-row gap-2">
                                    <Input
                                        placeholder="Nome (Ex: Entrar)"
                                        className="flex-1 h-9 rounded-lg"
                                        value={btn.text}
                                        onChange={(e) => updateNoticeButton(idx, 'text', e.target.value)}
                                        maxLength={30}
                                    />
                                    <Input
                                        placeholder={btn.type === 'url' ? 'https://...' : 'Comando'}
                                        className="flex-[1.5] h-9 rounded-lg"
                                        value={btn.value}
                                        onChange={(e) => updateNoticeButton(idx, 'value', e.target.value)}
                                        maxLength={100}
                                    />
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Send button */}
                <Button
                    variant="default"
                    className="w-full h-12 font-bold shadow-lg shadow-accent/20"
                    onClick={() => setIsConfirmOpen(true)}
                    disabled={isSendingNotice || !isReady}
                >
                    {isSendingNotice ? 'Enviando...' : <><Send size={18} /> Revisar &amp; Disparar</>}
                </Button>
            </div>

            {/* ─── Preview ─── */}
            <div className="rounded-xl border border-border p-4 flex flex-col min-h-[300px]">
                <div className="flex items-center justify-between mb-4">
                    <h3 className="text-[13px] font-bold flex items-center gap-2">
                        <Eye size={16} className="text-accent" /> Pré-visualização
                    </h3>
                    <Badge variant="secondary" className="text-[9px] tracking-wide uppercase">Telegram View</Badge>
                </div>
                <div className="flex-1 bg-muted/20 border border-border rounded-2xl p-4 flex items-center justify-center min-h-[280px]">
                    {renderPreview()}
                </div>
            </div>

            <ConfirmModal
                open={isConfirmOpen}
                onClose={() => setIsConfirmOpen(false)}
                onConfirm={handleSendNotice}
                title="Confirmar Disparo em Massa"
                message={`Você está prestes a enviar uma mensagem para ${
                    noticeTarget === 'all' ? 'todos os usuários e canais cadastrados' :
                    noticeTarget === 'channels' ? 'todos os canais cadastrados' :
                    noticeTarget === 'users' ? 'todos os usuários do bot' :
                    noticeTarget === 'channel_ids' ? 'os canais informados' :
                    'os usuários informados'
                }. Tem certeza?`}
                confirmText="Sim, Disparar Agora"
                danger={true}
            />
        </div>
    );
}
