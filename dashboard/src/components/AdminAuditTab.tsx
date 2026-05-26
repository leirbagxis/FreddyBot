import { Dispatch, SetStateAction, useState } from 'react';
import { Hash, ShieldAlert, ChevronRight, User as UserIcon, Zap, Trash2, ShieldCheck, CheckCircle2, Sparkles } from 'lucide-react';
import { bulkDeleteChannels } from '../api';
import { AuditResult, Channel } from '../types';
import { useToast } from './Toast';
import { ConfirmModal } from './ConfirmModal';

interface AdminAuditTabProps {
    navigateToChannel: (id: number) => void;
    onOpenUser: (id: number) => void;
    results: AuditResult[] | null;
    setResults: Dispatch<SetStateAction<AuditResult[] | null>>;
    loading: boolean;
    onRunAudit: () => void;
}

export function AdminAuditTab({ navigateToChannel, onOpenUser, results, setResults, loading, onRunAudit }: AdminAuditTabProps) {
    const [deletingId, setDeletingId] = useState<number | null>(null);
    const [confirmDelete, setConfirmDelete] = useState<{ userId: number, channels: number[], name: string } | null>(null);
    const toast = useToast();

    const handleRunAudit = () => {
        onRunAudit();
    };

    const handleBulkDelete = async () => {
        if (!confirmDelete) return;
        setDeletingId(confirmDelete.userId);
        const { userId, channels } = confirmDelete;
        setConfirmDelete(null);

        try {
            const res = await bulkDeleteChannels(userId, channels);
            if (res.success) {
                toast(`Remoção concluída: ${res.data.deletedCount} canais limpos`, "success");
                // Atualizar lista local removendo o usuário
                setResults(prev => prev ? prev.filter(r => r.userId !== userId) : null);
            } else {
                throw new Error(res.message || "Erro ao excluir canais");
            }
        } catch (err: any) {
            toast(err.message || "Erro na exclusão em massa", "error");
        } finally {
            setDeletingId(null);
        }
    };

    return (
        <div className="space-y-4">
            <div className="card" style={{ padding: '24px' }}>
                <div className="flex flex-col items-center text-center">
                    <div className="section-icon red mb-4" style={{ width: 64, height: 64 }}>
                        <ShieldAlert size={32} />
                    </div>
                    <h2 className="text-xl font-bold mb-2">Auditoria Ativa</h2>
                    <p className="text-sm opacity-70 mb-6 max-w-sm">
                        Esta ferramenta realiza uma varredura em tempo real em todos os canais do banco para identificar onde o bot legado <b>@XavolaBot</b> ainda possui permissões de administrador.
                    </p>
                    
                    <button 
                        className={`btn ${loading ? 'btn-ghost pointer-events-none' : 'btn-danger'} w-full max-w-xs`}
                        onClick={handleRunAudit}
                        disabled={loading}
                    >
                        {loading ? (
                            <>
                                <div className="auth-spinner mr-2" style={{ width: 16, height: 16 }} />
                                Varrendo canais...
                            </>
                        ) : (
                            <>
                                <Zap size={18} className="mr-2" />
                                Iniciar Varredura Agora
                            </>
                        )}
                    </button>
                </div>
            </div>

            {results && results.length > 0 && (
                <div className="space-y-12 mt-8">
                    <h3 className="text-[15px] font-bold px-1" style={{ color: 'var(--hint)' }}>
                        Usuários com XavolaBot detectado ({results.length})
                    </h3>
                    
                    {results.map((result) => (
                        <div key={result.userId} className="space-y-4 animate-in fade-in slide-in-from-bottom-2 duration-500">
                            <div className="flex items-center justify-between px-1 bg-[var(--surface)] p-3 rounded-2xl border border-[var(--border)] shadow-sm">
                                <button
                                    className="flex min-w-0 flex-1 items-center gap-3 text-left transition-opacity hover:opacity-80"
                                    onClick={() => onOpenUser(result.userId)}
                                    title="Abrir usuário"
                                >
                                    <div className="section-icon purple sm"><UserIcon size={14} /></div>
                                    <div className="flex min-w-0 flex-col">
                                        <span className="truncate text-sm font-bold">{result.firstName}</span>
                                        <span className="truncate text-[10px] opacity-40">ID: {result.userId} • {result.channels.length} canais</span>
                                    </div>
                                    <ChevronRight size={16} className="stat-arrow" />
                                </button>
                                <button 
                                    className={`btn btn-danger sm ${deletingId === result.userId ? 'loading' : ''}`}
                                    disabled={deletingId !== null}
                                    onClick={() => setConfirmDelete({ 
                                        userId: result.userId, 
                                        channels: result.channels.map(c => c.id),
                                        name: result.firstName
                                    })}
                                    title="Remover todos estes canais"
                                >
                                    {deletingId === result.userId ? (
                                        <div className="auth-spinner" style={{ width: 14, height: 14 }} />
                                    ) : (
                                        <>
                                            <Trash2 size={14} className="mr-1.5" />
                                            Limpar Tudo
                                        </>
                                    )}
                                </button>
                            </div>
                            
                            <div className="grid gap-2 pl-4 border-l-2 border-[var(--border)] ml-4">
                                {result.channels.map((c: Channel) => (
                                    <button 
                                        key={c.id} 
                                        className="admin-list-item flex items-center w-full text-left p-4 opacity-80 hover:opacity-100" 
                                        onClick={() => navigateToChannel(c.id)}
                                    >
                                        <div className="section-icon purple mr-3" style={{ transform: 'scale(0.8)' }}><Hash size={18} /></div>
                                        <div className="min-w-0 flex-1">
                                            <h3 className="text-[13px] font-semibold truncate">{c.title}</h3>
                                            <p className="text-[10px] truncate mt-0.5 opacity-50">ID: {c.id}</p>
                                        </div>
                                        <ChevronRight size={16} className="stat-arrow" />
                                    </button>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {results && results.length === 0 && !loading && (
                <div className="card animate-in fade-in slide-in-from-bottom-2 duration-500" style={{ padding: '28px 20px' }}>
                    <div className="flex flex-col items-center text-center">
                        <div className="section-icon green mb-4" style={{ width: 68, height: 68, borderRadius: 20 }}>
                            <ShieldCheck size={34} />
                        </div>

                        <div className="inline-flex items-center gap-1.5 mb-3 px-3 py-1 rounded-full border border-[var(--border)] bg-[var(--surface)] text-[11px] font-bold" style={{ color: 'var(--success)' }}>
                            <CheckCircle2 size={13} />
                            Auditoria concluida
                        </div>

                        <h3 className="text-lg font-bold mb-2">Nenhum XavolaBot encontrado</h3>
                        <p className="text-sm leading-relaxed max-w-sm mb-5" style={{ color: 'var(--text-secondary)' }}>
                            A varredura terminou e nenhum canal do banco possui o bot legado com permissões de administrador.
                        </p>

                        <div className="grid gap-2 w-full max-w-sm">
                            <div className="flex items-center gap-3 rounded-2xl border border-[var(--border)] bg-[var(--surface)] px-4 py-3 text-left">
                                <div className="section-icon green sm"><CheckCircle2 size={14} /></div>
                                <div className="min-w-0">
                                    <p className="text-[12px] font-bold">Canais verificados</p>
                                    <p className="text-[10px] opacity-50">Nenhuma permissão legada detectada.</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-3 rounded-2xl border border-[var(--border)] bg-[var(--surface)] px-4 py-3 text-left">
                                <div className="section-icon purple sm"><Sparkles size={14} /></div>
                                <div className="min-w-0">
                                    <p className="text-[12px] font-bold">Nenhuma ação necessária</p>
                                    <p className="text-[10px] opacity-50">A lista de limpeza permanece vazia.</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            <ConfirmModal
                open={confirmDelete !== null}
                onClose={() => setConfirmDelete(null)}
                onConfirm={handleBulkDelete}
                title="Confirmar Exclusão em Massa"
                message={`Você está prestes a remover permanentemente os ${confirmDelete?.channels.length} canais de "${confirmDelete?.name}". O bot sairá dos canais e todos os dados serão apagados. Deseja continuar?`}
                confirmText="Sim, Excluir Tudo"
                danger={true}
            />
        </div>
    );
}

