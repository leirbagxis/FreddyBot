import { useState, memo } from 'react';
import {
    Users, LogOut, ShieldCheck
} from 'lucide-react';
import { Channel } from '../types';
import { ConfirmModal } from './ConfirmModal';
import { fetchUserInfo, transferChannel } from '../api';
import { useToast } from './Toast';

interface DashboardInicioTabProps {
    channel: Channel;
    displayName: string;
    getGreeting: () => string;
    getGreetingEmoji: () => string;
    handleDisconnect: () => void;
    showDisconnect: boolean;
    setShowDisconnect: (open: boolean) => void;
    isDisconnecting: boolean;
    confirmDisconnect: () => void;
    showDisconnectSuccess: boolean;
    setShowDisconnectSuccess: (open: boolean) => void;
}

export const DashboardInicioTab = memo(({
    channel, displayName, getGreeting, getGreetingEmoji,
    handleDisconnect, showDisconnect, setShowDisconnect, isDisconnecting, confirmDisconnect,
    showDisconnectSuccess, setShowDisconnectSuccess,
}: DashboardInicioTabProps) => {
    const [transferInput, setTransferInput] = useState('');
    const [isTransferring, setIsTransferring] = useState(false);
    const [showTransferConfirm, setShowTransferConfirm] = useState(false);
    const [transferNewOwnerName, setTransferNewOwnerName] = useState('');
    const [transferNewOwnerId, setTransferNewOwnerId] = useState<number | null>(null);
    const [showTransferError, setShowTransferError] = useState(false);
    const [transferErrorMessage, setTransferErrorMessage] = useState('');
    const [showTransferSuccess, setShowTransferSuccess] = useState(false);
    const toast = useToast();

    const handleTransferClick = async () => {
        const newOwner = transferInput.trim();
        if (!newOwner) {
            toast('Digite o ID ou Username do novo dono', 'error');
            return;
        }

        setIsTransferring(true);
        setTransferErrorMessage('');
        setShowTransferError(false);

        try {
            const resp = await fetchUserInfo(newOwner);
            const isSuccess = resp && (resp.success || resp.succes) && resp.user;

            if (isSuccess) {
                setTransferNewOwnerName(resp.user.first_name);
                setTransferNewOwnerId(resp.user.id);
                setShowTransferConfirm(true);
            } else {
                setTransferErrorMessage(`Não foi possível encontrar nenhum usuário com o ID ou Username "${newOwner}". Por favor, verifique e tente novamente.`);
                setShowTransferError(true);
            }
        } catch {
            setTransferErrorMessage(`Ocorreu um erro ao buscar as informações do usuário. Tente novamente.`);
            setShowTransferError(true);
        } finally {
            setIsTransferring(false);
        }
    };

    const confirmTransfer = async () => {
        try {
            if (!channel?.ownerId) throw new Error("Owner ID not found");
            if (!transferNewOwnerId) throw new Error("New owner ID not found");

            await transferChannel(channel.ownerId, transferNewOwnerId, channel.id);
            setShowTransferSuccess(true);
            setTransferInput('');
            setShowTransferConfirm(false);
            setTransferNewOwnerName('');
            setTransferNewOwnerId(null);
        } catch (err: any) {
            if (err instanceof Error) {
                try {
                    const parsedErr = JSON.parse(err.message);
                    setTransferErrorMessage(parsedErr.message || 'Erro ao passar a posse para o novo usuário.');
                } catch {
                    setTransferErrorMessage(err.message || 'Erro ao passar a posse para o novo usuário.');
                }
            } else {
                setTransferErrorMessage('Erro desconhecido ao transferir o canal');
            }
            setShowTransferConfirm(false);
            setShowTransferError(true);
        }
    };

    return (
        <div className="space-y-3 tab-content-wrapper">
            
            {/* Unified Identity Card */}
            <div className="card animate-stagger-in">
                {/* Header: Greeting & Emoji */}
                <div className="flex items-center justify-between mb-5">
                    <div className="flex items-center gap-2">
                        <span className="text-xl">{getGreetingEmoji()}</span>
                        <div>
                            <h2 className="text-[15px] font-bold leading-none">{getGreeting()}</h2>
                            <p className="text-[10px] text-[var(--hint)] mt-1 uppercase tracking-wider font-semibold">Painel de Controle</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-1.5 bg-[var(--accent-soft)] px-2.5 py-1 rounded-lg">
                        <ShieldCheck size={12} className="text-[var(--accent)]" />
                        <span className="text-[11px] font-mono font-bold text-[var(--accent)]">{channel.ownerId}</span>
                    </div>
                </div>

                {/* User Info */}
                <div className="flex items-center gap-3 p-3 bg-[var(--surface)] border border-[var(--border)] rounded-2xl mb-4">
                    <div className="w-11 h-11 rounded-full flex items-center justify-center bg-[var(--accent)] text-white font-bold text-lg flex-shrink-0 shadow-sm">
                        {displayName.charAt(0).toUpperCase()}
                    </div>
                    <div className="min-w-0 flex-1">
                        <h3 className="text-[16px] font-bold text-[var(--text)] truncate">{displayName}</h3>
                        <p className="text-[11px] text-[var(--hint)] truncate">Administrador do Canal</p>
                    </div>
                </div>
                
                {/* Integrated Disconnect Action */}
                <button 
                    className="w-full flex items-center justify-center gap-2 py-3 px-4 rounded-xl bg-[var(--danger-soft)] text-[var(--danger)] text-[13px] font-bold hover:opacity-80 transition-all active:scale-[0.98]" 
                    onClick={handleDisconnect}
                >
                    <LogOut size={16} />
                    <span>Desconectar Bot</span>
                </button>
            </div>

            {/* Transferir Posse */}
            <div className="card">
                <div className="section-header">
                    <div className="section-icon purple" style={{ background: 'var(--warning-soft)', color: 'var(--warning)' }}>
                        <Users size={18} />
                    </div>
                    <div className="min-w-0 flex-1">
                        <h3 className="text-[15px] font-semibold truncate">Transferir Posse</h3>
                        <p className="text-xs truncate" style={{ color: 'var(--hint)' }}>Passe a administração para outro usuário</p>
                    </div>
                </div>
                <div className="flex flex-col gap-3 mt-3">
                    <input
                        type="text"
                        placeholder="ID ou Username do novo dono"
                        className="w-full bg-[var(--background)] text-[var(--text)] border border-[var(--border)] rounded-xl py-2 px-3 focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
                        value={transferInput}
                        onChange={(e) => setTransferInput(e.target.value)}
                    />
                    <button
                        className="btn btn-primary"
                        onClick={handleTransferClick}
                        disabled={!transferInput.trim() || isTransferring}
                    >
                        {isTransferring ? 'Aguarde...' : 'Enviar'}
                    </button>
                </div>
            </div>

            {/* Disconnect Confirm Modal */}
            <ConfirmModal
                open={showDisconnect}
                onClose={() => !isDisconnecting && setShowDisconnect(false)}
                onConfirm={confirmDisconnect}
                title="Desconectar Bot"
                message="Tem certeza que deseja desconectar o bot deste canal? Todas as configurações serão perdidas."
                confirmText={isDisconnecting ? "Desconectando..." : "Desconectar"}
                danger
            />

            {/* Disconnect Success Modal */}
            <ConfirmModal
                open={showDisconnectSuccess}
                onClose={() => { }}
                onConfirm={() => {
                    const tg = window.Telegram?.WebApp;
                    if (tg) {
                        tg.close();
                    }
                }}
                title="Desconectado"
                message="O bot foi desconectado com sucesso. Esta janela será fechada."
                confirmText="Fechar"
            />

            {/* Transfer Confirm Modal */}
            <ConfirmModal
                open={showTransferConfirm}
                onClose={() => !isTransferring && setShowTransferConfirm(false)}
                onConfirm={confirmTransfer}
                title="Confirmar Transferência"
                message={`Você tem certeza que deseja transferir a posse para ${transferNewOwnerName}? Você perderá o acesso de dono.`}
                confirmText={isTransferring ? "Transferindo..." : "Confirmar"}
                danger
            />

            {/* Transfer Success Modal */}
            <ConfirmModal
                open={showTransferSuccess}
                onClose={() => { }}
                onConfirm={() => {
                    const tg = window.Telegram?.WebApp;
                    if (tg) {
                        tg.close();
                    }
                }}
                title="Sucesso"
                message="Posse transferida com sucesso. O bot foi reiniciado e esta janela será fechada."
                confirmText="Fechar"
            />

            {/* Transfer Error Modal */}
            <ConfirmModal
                open={showTransferError}
                onClose={() => setShowTransferError(false)}
                onConfirm={() => setShowTransferError(false)}
                title="Erro na Transferência"
                message={transferErrorMessage}
                confirmText="Ok"
                danger
            />
        </div>
    );
});
