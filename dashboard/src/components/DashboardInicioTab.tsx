import { useState, memo } from 'react';
import {
    Users, LogOut, ShieldCheck, Send
} from 'lucide-react';
import { Channel } from '../types';
import { ConfirmModal } from './ConfirmModal';
import { fetchUserInfo, transferChannel } from '../api';
import { useToast } from './Toast';
import { Card, CardContent } from './ui/card';
import { Button } from './ui/button';
import { Input } from './ui/input';

interface DashboardInicioTabProps {
    channel: Channel;
    displayName: string;
    getGreeting: () => string;
    getGreetingIcon: () => React.ReactNode;
    handleDisconnect: () => void;
    showDisconnect: boolean;
    setShowDisconnect: (open: boolean) => void;
    isDisconnecting: boolean;
    confirmDisconnect: () => void;
    showDisconnectSuccess: boolean;
    setShowDisconnectSuccess: (open: boolean) => void;
}

export const DashboardInicioTab = memo(({
    channel, displayName, getGreeting, getGreetingIcon,
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
            <Card className="animate-stagger-in">
                <CardContent className="pt-4">
                    {/* Header: Greeting & Emoji */}
                    <div className="flex items-center justify-between mb-5">
                        <div className="flex items-center gap-2">
                            <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                                {getGreetingIcon()}
                            </div>
                            <div>
                                <h2 className="text-[15px] font-bold leading-none">{getGreeting()}</h2>
                                <p className="text-[10px] text-muted-foreground mt-1 uppercase tracking-wider font-semibold">Painel de Controle</p>
                            </div>
                        </div>
                        <div className="flex items-center gap-1.5 bg-accent/10 px-2.5 py-1 rounded-lg">
                            <ShieldCheck size={12} className="text-accent" />
                            <span className="text-[11px] font-mono font-bold text-accent">{channel.ownerId}</span>
                        </div>
                    </div>

                    {/* User Info */}
                    <div className="flex items-center gap-3 p-3 bg-muted/50 border border-border rounded-2xl mb-4">
                        <div className="w-11 h-11 rounded-full flex items-center justify-center bg-primary text-primary-foreground font-bold text-lg flex-shrink-0 shadow-sm">
                            {displayName.charAt(0).toUpperCase()}
                        </div>
                        <div className="min-w-0 flex-1">
                            <h3 className="text-[16px] font-bold text-foreground truncate">{displayName}</h3>
                            <p className="text-[11px] text-muted-foreground truncate">Administrador do Canal</p>
                        </div>
                    </div>
                    
                    {/* Integrated Disconnect Action */}
                    <Button 
                        variant="destructive"
                        className="w-full"
                        onClick={handleDisconnect}
                    >
                        <LogOut size={16} />
                        Desconectar Bot
                    </Button>
                </CardContent>
            </Card>

            {/* Transferir Posse */}
            <div className="rounded-xl border border-border p-4 space-y-3">
                <div className="flex items-center gap-3">
                    <div className="section-icon" style={{ background: 'var(--warning-soft)', color: 'var(--warning)' }}>
                        <Users size={18} />
                    </div>
                    <div className="min-w-0 flex-1">
                        <h3 className="text-[15px] font-semibold truncate">Transferir Posse</h3>
                        <p className="text-xs truncate text-muted-foreground">Passe a administração para outro usuário</p>
                    </div>
                </div>
                <div className="flex gap-2 items-center">
                    <div className="relative flex-1">
                        <Users size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none" />
                        <Input
                            className="h-10 pl-9 rounded-xl"
                            placeholder="ID ou Username do novo dono"
                            value={transferInput}
                            onChange={(e) => setTransferInput(e.target.value)}
                        />
                    </div>
                    <Button
                        variant="default"
                        className="h-10 shrink-0"
                        onClick={handleTransferClick}
                        disabled={!transferInput.trim() || isTransferring}
                    >
                        {isTransferring ? 'Aguarde...' : <><Send size={16} /> Transferir</>}
                    </Button>
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
