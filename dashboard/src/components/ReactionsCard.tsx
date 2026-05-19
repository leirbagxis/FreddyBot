import { useState, useEffect, memo } from 'react';
import { SmilePlus, X } from 'lucide-react';

interface ReactionsCardProps {
    reactions: string;
    onUpdate: (reactions: string) => Promise<void>;
}

export const ReactionsCard = memo(({ reactions, onUpdate }: ReactionsCardProps) => {
    const [slots, setSlots] = useState<string[]>(['', '', '', '', '']);
    const [loading, setLoading] = useState(false);

    // Initialize slots from comma-separated string
    useEffect(() => {
        if (reactions) {
            const split = reactions.split(',').map(s => s.trim());
            const newSlots = ['', '', '', '', ''];
            for (let i = 0; i < 5; i++) {
                if (split[i]) newSlots[i] = split[i];
            }
            setSlots(newSlots);
        }
    }, [reactions]);

    const isEmoji = (str: string) => {
        // Regex para detectar se a string contém APENAS emojis (incluindo variações de colos, etc)
        const emojiRegex = /^(\u00a9|\u00ae|[\u2000-\u3300]|\ud83c[\ud000-\udfff]|\ud83d[\ud000-\udfff]|\ud83e[\ud000-\udfff])+$/;
        return emojiRegex.test(str);
    };

    const handleSlotChange = (index: number, value: string) => {
        const trimmed = value.trim();
        if (trimmed === '') {
            const newSlots = [...slots];
            newSlots[index] = '';
            setSlots(newSlots);
            return;
        }

        // Se for um emoji válido, aceita. Caso contrário, ignora (ou pega só o emoji se colarem texto+emoji)
        // Para simplificar, vamos validar se o que foi digitado/colado contém emoji
        if (isEmoji(trimmed)) {
            const newSlots = [...slots];
            // Se colarem vários emojis, pegamos apenas o primeiro símbolo (que pode ser composto)
            // Usando Array.from para lidar corretamente com surrogate pairs de emojis
            const emojis = Array.from(trimmed);
            newSlots[index] = emojis[0];
            setSlots(newSlots);
        }
    };

    const handleClearSlot = (index: number) => {
        const newSlots = [...slots];
        newSlots[index] = '';
        setSlots(newSlots);
    };

    const handleSave = async () => {
        setLoading(true);
        try {
            // Filter out empty slots and join by comma
            const reactionsString = slots.filter(s => s.trim() !== '').join(',');
            await onUpdate(reactionsString);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="card">
            <div className="section-header">
                <div className="section-icon blue">
                    <SmilePlus size={18} />
                </div>
                <div className="min-w-0 flex-1">
                    <h3 className="text-[15px] font-semibold truncate">Reações / Votos (Grid)</h3>
                    <p className="text-xs truncate" style={{ color: 'var(--hint)' }}>Adicione até 5 emojis para votação rápida.</p>
                </div>
            </div>

            <div className="mt-4">
                <div className="grid grid-cols-5 gap-2 mb-4">
                    {slots.map((slot, index) => (
                        <div key={index} className="relative group">
                            <input
                                type="text"
                                value={slot}
                                onChange={(e) => handleSlotChange(index, e.target.value)}
                                placeholder="+"
                                className="w-full aspect-square text-center text-xl bg-[var(--background)] text-[var(--text)] border border-[var(--border)] rounded-xl focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
                            />
                            {slot && (
                                <button 
                                    onClick={() => handleClearSlot(index)}
                                    className="absolute -top-1 -right-1 bg-[var(--danger)] text-white rounded-full p-0.5 opacity-0 group-hover:opacity-100 transition-opacity"
                                >
                                    <X size={10} />
                                </button>
                            )}
                        </div>
                    ))}
                </div>
                
                <button 
                    className="btn btn-primary w-full" 
                    onClick={handleSave}
                    disabled={loading}
                >
                    {loading ? 'Salvando...' : 'Salvar Reações'}
                </button>
            </div>
        </div>
    );
});
