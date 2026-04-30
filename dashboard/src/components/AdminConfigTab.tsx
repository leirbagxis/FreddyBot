import { useState, useEffect } from 'react';
import { Settings, ShieldCheck, Construction } from 'lucide-react';
import { fetchServerConfig, updateServerConfig } from '../api';
import { ServerConfig } from '../types';
import { useToast } from './Toast';

export function AdminConfigTab() {
    const [config, setConfig] = useState<ServerConfig | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const toast = useToast();

    useEffect(() => {
        const loadConfig = async () => {
            try {
                const res = await fetchServerConfig();
                if (res.success) {
                    setConfig(res.config);
                }
            } catch (err) {
                toast('Erro ao carregar configurações', 'error');
            } finally {
                setLoading(false);
            }
        };
        loadConfig();
    }, [toast]);

    const handleToggle = async (field: 'maintence' | 'forceJoin') => {
        if (!config) return;
        
        const newConfig = { ...config, [field]: !config[field] };
        
        setSaving(true);
        try {
            const res = await updateServerConfig(newConfig.maintence, newConfig.forceJoin);
            if (res.success) {
                setConfig(res.config);
                toast(`Configuração atualizada: ${field === 'maintence' ? 'Manutenção' : 'Force Join'} ${newConfig[field] ? 'Ativado' : 'Desativado'}`, 'success');
            }
        } catch (err) {
            toast('Erro ao atualizar configuração', 'error');
        } finally {
            setSaving(false);
        }
    };

    if (loading) return <div className="p-8 text-center opacity-50">Carregando configurações...</div>;

    return (
        <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <div className="admin-welcome-card">
                <div className="flex items-center gap-4">
                    <div className="section-icon purple">
                        <Settings size={22} />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold">Configurações Globais</h2>
                        <p className="text-sm opacity-60">Gerencie o estado do bot e restrições de acesso.</p>
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Manutenção */}
                <div className="card">
                    <div className="section-header">
                        <div className="section-icon amber">
                            <Construction size={18} />
                        </div>
                        <div className="flex-1 min-w-0">
                            <h3 className="text-[15px] font-semibold truncate">Modo Manutenção</h3>
                            <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
                                {config?.maintence ? 'O bot está offline para usuários' : 'O bot está operando normalmente'}
                            </p>
                        </div>
                    </div>
                    <div className={`perm-row ${config?.maintence ? 'on' : ''}`} onClick={() => !saving && handleToggle('maintence')}>
                        <div className="flex items-center gap-3">
                            <span className="text-[13px] font-medium">Status da Manutenção</span>
                        </div>
                        <div className={`toggle ${config?.maintence ? 'on' : ''}`} />
                    </div>
                </div>

                {/* Force Join */}
                <div className="card">
                    <div className="section-header">
                        <div className="section-icon purple">
                            <ShieldCheck size={18} />
                        </div>
                        <div className="flex-1 min-w-0">
                            <h3 className="text-[15px] font-semibold truncate">Force Join (Inscrição Obrigatória)</h3>
                            <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
                                {config?.forceJoin ? 'Usuários devem entrar no canal oficial' : 'Acesso livre para todos'}
                            </p>
                        </div>
                    </div>
                    <div className={`perm-row ${config?.forceJoin ? 'on' : ''}`} onClick={() => !saving && handleToggle('forceJoin')}>
                        <div className="flex items-center gap-3">
                            <span className="text-[13px] font-medium">Status do Force Join</span>
                        </div>
                        <div className={`toggle ${config?.forceJoin ? 'on' : ''}`} />
                    </div>
                    <div className="mt-3 p-3 bg-[var(--surface)] rounded-xl border border-[var(--border)]">
                        <div className="text-[11px] font-bold opacity-40 uppercase mb-1">Canal de Verificação</div>
                        <div className="text-[12px] font-mono opacity-80">-1003767126116</div>
                    </div>
                </div>
            </div>
        </div>
    );
}
