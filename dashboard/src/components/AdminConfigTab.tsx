import { useState, useEffect } from 'react';
import { Settings, ShieldCheck, Construction, FileText, Save, KeyRound, Code2 } from 'lucide-react';
import { fetchServerConfig, updateServerConfig } from '../api';
import { ServerConfig } from '../types';
import { useToast } from './Toast';
import { RichTextEditor } from './RichTextEditor';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Switch } from './ui/switch';
import { Textarea } from './ui/textarea';

export function AdminConfigTab() {
    const [config, setConfig] = useState<ServerConfig | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);

    const [globalDefault, setGlobalDefault] = useState('');
    const [globalNewPack, setGlobalNewPack] = useState('');
    const [fixedPostEnabled, setFixedPostEnabled] = useState(true);
    const [fixedPostKey, setFixedPostKey] = useState('legendasbot');
    const [fixedPostPayload, setFixedPostPayload] = useState('');

    const toast = useToast();

    useEffect(() => {
        const loadConfig = async () => {
            try {
                const res = await fetchServerConfig();
                if (res.success) {
                    const serverData = res.data || res.config;
                    if (serverData) {
                        setConfig(serverData);
                        setGlobalDefault(serverData.globalDefaultCaption || '');
                        setGlobalNewPack(serverData.globalNewPackCaption || '');
                        setFixedPostEnabled(Boolean(serverData.fixedPostBuilderEnabled));
                        setFixedPostKey(serverData.fixedPostBuilderKey || 'legendasbot');
                        setFixedPostPayload(serverData.fixedPostBuilderPayload || '');
                    }
                }
            } catch (err) {
                console.error("Erro ao carregar configurações:", err);
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
            globalNewPackCaption: overrides.globalNewPackCaption ?? globalNewPack,
            fixedPostBuilderEnabled: overrides.fixedPostBuilderEnabled ?? fixedPostEnabled,
            fixedPostBuilderKey: overrides.fixedPostBuilderKey ?? fixedPostKey,
            fixedPostBuilderPayload: overrides.fixedPostBuilderPayload ?? fixedPostPayload
        };

        setSaving(true);
        try {
            const res = await updateServerConfig(payload);
            if (res.success) {
                const serverData = res.data || res.config;
                if (serverData) {
                    setConfig(serverData);
                    setGlobalDefault(serverData.globalDefaultCaption || '');
                    setGlobalNewPack(serverData.globalNewPackCaption || '');
                    setFixedPostEnabled(Boolean(serverData.fixedPostBuilderEnabled));
                    setFixedPostKey(serverData.fixedPostBuilderKey || 'legendasbot');
                    setFixedPostPayload(serverData.fixedPostBuilderPayload || '');
                }
                toast('Configurações atualizadas com sucesso', 'success');
            }
        } catch (err) {
            toast('Erro ao atualizar configurações', 'error');
        } finally {
            setSaving(false);
        }
    };

    const handleToggle = (field: 'maintence' | 'forceJoin' | 'fixedPostBuilderEnabled') => {
        if (!config) return;
        if (field === 'fixedPostBuilderEnabled') {
            const next = !fixedPostEnabled;
            setFixedPostEnabled(next);
            handleSave({ fixedPostBuilderEnabled: next });
            return;
        }
        handleSave({ [field]: !config[field] });
    };

    if (loading) return (
        <div className="flex flex-col items-center py-12 gap-3">
            <div className="auth-spinner" />
            <p className="text-[13px] text-muted-foreground">Carregando configurações...</p>
        </div>
    );

    return (
        <div className="space-y-4 pb-12">
            {/* Header */}
            <div className="flex items-center gap-3">
                <div className="flex items-center justify-center size-10 rounded-xl shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                    <Settings size={20} />
                </div>
                <div>
                    <h2 className="text-base font-bold">Configurações Globais</h2>
                    <p className="text-xs text-muted-foreground">Gerencie o estado do bot e legendas iniciais</p>
                </div>
            </div>

            {/* Toggles Row */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {/* Manutenção */}
                <div className="rounded-xl border border-border p-4">
                    <div className="flex items-center gap-3 mb-3">
                        <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--warning-soft)', color: 'var(--warning)' }}>
                            <Construction size={16} />
                        </div>
                        <div className="flex-1 min-w-0">
                            <h3 className="text-[13px] font-semibold">Modo Manutenção</h3>
                            <p className="text-[11px] text-muted-foreground mt-0.5">
                                {config?.maintence ? 'Bot offline para usuários' : 'Operando normalmente'}
                            </p>
                        </div>
                    </div>
                    <div
                        className={`flex items-center justify-between px-4 py-3 rounded-xl cursor-pointer transition-all ${config?.maintence ? 'bg-warning/10' : 'bg-muted/30'}`}
                        onClick={() => !saving && handleToggle('maintence')}
                    >
                        <span className="text-[12px] font-medium">Ativar Manutenção</span>
                        <Switch checked={!!config?.maintence} onCheckedChange={() => !saving && handleToggle('maintence')} onClick={(e: React.MouseEvent) => e.stopPropagation()} />
                    </div>
                </div>

                {/* Force Join */}
                <div className="rounded-xl border border-border p-4">
                    <div className="flex items-center gap-3 mb-3">
                        <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                            <ShieldCheck size={16} />
                        </div>
                        <div className="flex-1 min-w-0">
                            <h3 className="text-[13px] font-semibold">Force Join</h3>
                            <p className="text-[11px] text-muted-foreground mt-0.5">
                                {config?.forceJoin ? 'Inscrição obrigatória' : 'Acesso livre'}
                            </p>
                        </div>
                    </div>
                    <div
                        className={`flex items-center justify-between px-4 py-3 rounded-xl cursor-pointer transition-all ${config?.forceJoin ? 'bg-accent/10' : 'bg-muted/30'}`}
                        onClick={() => !saving && handleToggle('forceJoin')}
                    >
                        <span className="text-[12px] font-medium">Exigir inscrição no canal</span>
                        <Switch checked={!!config?.forceJoin} onCheckedChange={() => !saving && handleToggle('forceJoin')} onClick={(e: React.MouseEvent) => e.stopPropagation()} />
                    </div>
                </div>
            </div>

            {/* Global Caption */}
            <div className="rounded-xl border border-border p-4">
                <div className="flex items-center gap-3 mb-3">
                    <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                        <FileText size={16} />
                    </div>
                    <div className="flex-1 min-w-0">
                        <h3 className="text-[13px] font-semibold">Legenda Padrão Global</h3>
                        <p className="text-[11px] text-muted-foreground mt-0.5">Preenche novos canais vinculados ao bot</p>
                    </div>
                </div>
                <RichTextEditor
                    value={globalDefault}
                    onChange={setGlobalDefault}
                    placeholder="Ex: @legendasbot [t.me/legendasbot](https://t.me/botusername)  ‹"
                />
            </div>

            {/* New Pack Caption */}
            <div className="rounded-xl border border-border p-4">
                <div className="flex items-center gap-3 mb-3">
                    <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--warning-soft)', color: 'var(--warning)' }}>
                        <FileText size={16} />
                    </div>
                    <div className="flex-1 min-w-0">
                        <h3 className="text-[13px] font-semibold">Legenda de Novo Pack (Global)</h3>
                        <p className="text-[11px] text-muted-foreground mt-0.5">Valor inicial para mensagem de pack padrão</p>
                    </div>
                </div>
                <RichTextEditor
                    value={globalNewPack}
                    onChange={setGlobalNewPack}
                    placeholder="Texto inicial para novos packs..."
                />
            </div>

            {/* PostBuilder Fixo */}
            <div className="rounded-xl border border-border p-4">
                <div className="flex items-center gap-3 mb-3">
                    <div className="flex items-center justify-center size-9 rounded-lg shrink-0" style={{ background: 'var(--accent-soft)', color: 'var(--accent)' }}>
                        <Code2 size={16} />
                    </div>
                    <div className="flex-1 min-w-0">
                        <h3 className="text-[13px] font-semibold">PostBuilder Fixo</h3>
                        <p className="text-[11px] text-muted-foreground mt-0.5">Post permanente usado no inline com chave fixa</p>
                    </div>
                </div>

                <div
                    className={`flex items-center justify-between px-4 py-3 rounded-xl cursor-pointer transition-all mb-4 ${fixedPostEnabled ? 'bg-accent/10' : 'bg-muted/30'}`}
                    onClick={() => !saving && handleToggle('fixedPostBuilderEnabled')}
                >
                    <span className="text-[12px] font-medium">Postagem fixa ativa</span>
                    <Switch checked={fixedPostEnabled} onCheckedChange={() => !saving && handleToggle('fixedPostBuilderEnabled')} onClick={(e: React.MouseEvent) => e.stopPropagation()} />
                </div>

                <div className="space-y-3">
                    <div>
                        <label className="text-[11px] font-semibold text-muted-foreground flex items-center gap-1.5 mb-1.5">
                            <KeyRound size={13} /> Key fixa
                        </label>
                        <Input
                            value={fixedPostKey}
                            onChange={(e) => setFixedPostKey(e.target.value)}
                            placeholder="legendasbot"
                            disabled={saving}
                            className="h-9"
                        />
                    </div>

                    <div>
                        <label className="text-[11px] font-semibold text-muted-foreground flex items-center gap-1.5 mb-1.5">
                            <Code2 size={13} /> Payload JSON
                        </label>
                        <Textarea
                            value={fixedPostPayload}
                            onChange={(e) => setFixedPostPayload(e.target.value)}
                            className="min-h-[220px] rounded-xl font-mono text-xs resize-y"
                            placeholder='{ "media_type": "photo", "media_file_id": "..." }'
                            disabled={saving}
                        />
                    </div>

                    <p className="text-[11px] text-muted-foreground leading-relaxed">
                        Uso inline: <code className="bg-muted/30 px-1.5 py-0.5 rounded text-accent text-[10px]">@FreddyCaptionBot pb {fixedPostKey || 'legendasbot'}</code>. Quando desativado, a chave é removida do Redis.
                    </p>
                </div>
            </div>

            {/* Save All */}
            <Button
                variant="default"
                className="w-full h-12 font-bold shadow-lg shadow-accent/20"
                onClick={() => !saving && handleSave()}
                disabled={saving}
            >
                <Save size={18} />
                {saving ? 'Salvando...' : 'Salvar Legendas Globais'}
            </Button>
        </div>
    );
}
