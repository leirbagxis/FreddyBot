import { useState, useEffect } from 'react';
import { Save } from 'lucide-react';
import { fetchServerConfig, updateServerConfig } from '../api';
import { ServerConfig } from '../types';
import { useToast } from './Toast';
import { RichTextEditor } from './RichTextEditor';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Switch } from './ui/switch';

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
            } catch {
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
        <div className="admin-config pb-16">
            <div className="admin-config-card">

                {/* ── Header ── */}
                <div className="cfg-header">
                    <span className="cfg-header-title">Configurações Globais</span>
                    <span className="cfg-header-sub">Estado do bot e legendas padrão</span>
                </div>

                {/* ── Sistema ── */}
                <div className="cfg-section">
                    <span className="cfg-section-label">Sistema</span>

                    <div className="cfg-row" onClick={() => !saving && handleToggle('maintence')}>
                        <div className="cfg-row-text">
                            <span className="cfg-row-title">Manutenção</span>
                            <span className="cfg-row-desc">
                                {config?.maintence ? 'Bot offline para usuários' : 'Operando normalmente'}
                            </span>
                        </div>
                        <Switch
                            checked={!!config?.maintence}
                            onCheckedChange={() => !saving && handleToggle('maintence')}
                            onClick={(e: React.MouseEvent) => e.stopPropagation()}
                        />
                    </div>

                    <div className="cfg-divider" />

                    <div className="cfg-row" onClick={() => !saving && handleToggle('forceJoin')}>
                        <div className="cfg-row-text">
                            <span className="cfg-row-title">Force Join</span>
                            <span className="cfg-row-desc">
                                {config?.forceJoin ? 'Inscrição obrigatória' : 'Acesso livre'}
                            </span>
                        </div>
                        <Switch
                            checked={!!config?.forceJoin}
                            onCheckedChange={() => !saving && handleToggle('forceJoin')}
                            onClick={(e: React.MouseEvent) => e.stopPropagation()}
                        />
                    </div>
                </div>

                <div className="cfg-divider-full" />

                {/* ── Legendas ── */}
                <div className="cfg-section">
                    <span className="cfg-section-label">Legendas</span>

                    <div className="cfg-editor">
                        <div className="cfg-editor-header">
                            <span className="cfg-editor-title">Legenda Padrão Global</span>
                            <span className="cfg-editor-desc">Preenche novos canais vinculados ao bot</span>
                        </div>
                        <RichTextEditor
                            value={globalDefault}
                            onChange={setGlobalDefault}
                            placeholder="Ex: @legendasbot [t.me/legendasbot](https://t.me/botusername)"
                        />
                    </div>

                    <div className="cfg-editor">
                        <div className="cfg-editor-header">
                            <span className="cfg-editor-title">Legenda de Novo Pack</span>
                            <span className="cfg-editor-desc">Valor inicial para mensagem de pack padrão</span>
                        </div>
                        <RichTextEditor
                            value={globalNewPack}
                            onChange={setGlobalNewPack}
                            placeholder="Texto inicial para novos packs..."
                        />
                    </div>
                </div>

                <div className="cfg-divider-full" />

                {/* ── PostBuilder ── */}
                <div className="cfg-section">
                    <span className="cfg-section-label">PostBuilder</span>

                    <div className="cfg-row" onClick={() => !saving && handleToggle('fixedPostBuilderEnabled')}>
                        <div className="cfg-row-text">
                            <span className="cfg-row-title">Postagem fixa</span>
                            <span className="cfg-row-desc">Post permanente usado no inline com chave fixa</span>
                        </div>
                        <Switch
                            checked={fixedPostEnabled}
                            onCheckedChange={() => !saving && handleToggle('fixedPostBuilderEnabled')}
                            onClick={(e: React.MouseEvent) => e.stopPropagation()}
                        />
                    </div>

                    <div className="cfg-fields">
                        <div className="cfg-field">
                            <label className="cfg-field-label">Chave fixa</label>
                            <Input
                                value={fixedPostKey}
                                onChange={(e) => setFixedPostKey(e.target.value)}
                                placeholder="legendasbot"
                                disabled={saving}
                                className="h-9"
                            />
                        </div>

                        <div className="cfg-field">
                            <label className="cfg-field-label">Payload JSON</label>
                            <textarea
                                value={fixedPostPayload}
                                onChange={(e) => setFixedPostPayload(e.target.value)}
                                className="cfg-textarea"
                                placeholder='{ "media_type": "photo", "media_file_id": "..." }'
                                disabled={saving}
                            />
                        </div>

                        <p className="cfg-hint">
                            Uso inline:{' '}
                            <code className="cfg-code">
                                @FreddyCaptionBot pb {fixedPostKey || 'legendasbot'}
                            </code>
                            . Quando desativado, a chave é removida do Redis.
                        </p>
                    </div>
                </div>

                {/* ── Save ── */}
                <div className="cfg-footer">
                    <Button
                        variant="default"
                        className="cfg-save-btn"
                        onClick={() => !saving && handleSave()}
                        disabled={saving}
                    >
                        <Save size={15} />
                        {saving ? 'Salvando...' : 'Salvar'}
                    </Button>
                </div>
            </div>
        </div>
    );
}
