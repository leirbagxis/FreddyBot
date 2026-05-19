import { useState, useEffect } from 'react';
import { Settings, ShieldCheck, Construction, FileText, PackagePlus, Save } from 'lucide-react';
import { fetchServerConfig, updateServerConfig } from '../api';
import { ServerConfig } from '../types';
import { useToast } from './Toast';
import { RichTextEditor } from './RichTextEditor';

export function AdminConfigTab() {
    const [config, setConfig] = useState<ServerConfig | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    
    // Estados locais para os editores para evitar re-render pesado da tab inteira a cada caractere
    const [globalDefault, setGlobalDefault] = useState('');
    const [globalNewPack, setGlobalNewPack] = useState('');

    const toast = useToast();

    useEffect(() => {
        const loadConfig = async () => {
            try {
                const res = await fetchServerConfig();
                if (res.success) {
                    // O backend retorna os dados dentro da propriedade 'data' (NewSuccessResponse)
                    const serverData = res.data || res.config; 
                    if (serverData) {
                        setConfig(serverData);
                        setGlobalDefault(serverData.globalDefaultCaption || '');
                        setGlobalNewPack(serverData.globalNewPackCaption || '');
                    }
                }
            } catch (err) {
                console.error("Erro ao carregar configurações Admin:", err);
                toast('Erro ao carregar configurações', 'error');
            } finally {
                setLoading(false);
            }
        };
        loadConfig();
    }, [toast]);

    const handleSave = async (overrides: Partial<ServerConfig> = {}) => {
        if (!config) return;
        
        const payload = {
            maintence: overrides.maintence ?? config.maintence,
            forceJoin: overrides.forceJoin ?? config.forceJoin,
            globalDefaultCaption: overrides.globalDefaultCaption ?? globalDefault,
            globalNewPackCaption: overrides.globalNewPackCaption ?? globalNewPack
        };
        
        setSaving(true);
        try {
            const res = await updateServerConfig(
                payload.maintence, 
                payload.forceJoin, 
                payload.globalDefaultCaption, 
                payload.globalNewPackCaption
            );
            if (res.success) {
                const serverData = res.data || res.config;
                if (serverData) {
                    setConfig(serverData);
                }
                toast('Configurações atualizadas com sucesso', 'success');
            }
        } catch (err) {
            toast('Erro ao atualizar configurações', 'error');
        } finally {
            setSaving(false);
        }
    };

    const handleToggle = (field: 'maintence' | 'forceJoin') => {
        if (!config) return;
        handleSave({ [field]: !config[field] });
    };

    if (loading) return <div className="p-8 text-center opacity-50">Carregando configurações...</div>;

    return (
        <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500 pb-20">
            <div className="admin-welcome-card">
                <div className="flex items-center gap-4">
                    <div className="section-icon purple">
                        <Settings size={22} />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold">Configurações Globais</h2>
                        <p className="text-sm opacity-60">Gerencie o estado do bot e as legendas iniciais de novos canais.</p>
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
                </div>
            </div>

            {/* Legenda Padrão Global */}
            <div className="card">
                <div className="section-header">
                    <div className="section-icon purple">
                        <FileText size={18} />
                    </div>
                    <div className="flex-1 min-w-0">
                        <h3 className="text-[15px] font-semibold truncate">Legenda Padrão (Global)</h3>
                        <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
                            Usada para preencher novos canais vinculados ao bot.
                        </p>
                    </div>
                </div>
                <div className="p-4 bg-[var(--surface)] rounded-2xl border border-[var(--border)] mt-2">
                    <RichTextEditor 
                        value={globalDefault}
                        onChange={setGlobalDefault}
                        placeholder="Ex: 🐈‍⠀៹ [t.me/legendasbot](https://t.me/botusername)  ‹"
                    />
                </div>
            </div>

            {/* Legenda Novo Pack Global */}
            <div className="card">
                <div className="section-header">
                    <div className="section-icon amber">
                        <PackagePlus size={18} />
                    </div>
                    <div className="flex-1 min-w-0">
                        <h3 className="text-[15px] font-semibold truncate">Legenda de Novo Pack (Global)</h3>
                        <p className="text-xs mt-0.5" style={{ color: 'var(--hint)' }}>
                            Usada como valor inicial para a mensagem de pack padrão.
                        </p>
                    </div>
                </div>
                <div className="p-4 bg-[var(--surface)] rounded-2xl border border-[var(--border)] mt-2">
                    <RichTextEditor 
                        value={globalNewPack}
                        onChange={setGlobalNewPack}
                        placeholder="Texto inicial para novos packs..."
                    />
                </div>
            </div>

            {/* Botão Salvar Geral */}
            <div className="pt-4 pb-12">
                <button 
                    className={`btn-primary w-full shadow-2xl flex items-center justify-center gap-2 h-12 rounded-2xl transition-all active:scale-95 ${saving ? 'opacity-70 grayscale' : ''}`}
                    onClick={() => !saving && handleSave()}
                    disabled={saving}
                >
                    <Save size={20} />
                    <span className="font-bold">{saving ? 'Salvando...' : 'Salvar Legendas Globais'}</span>
                </button>
            </div>
        </div>
    );
}
