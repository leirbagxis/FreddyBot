import { useState, useEffect } from 'react';
import { SmilePlus, X } from 'lucide-react';

interface ReactionsCardProps {
    reactions: string;
    onUpdate: (reactions: string) => Promise<void>;
}

export function ReactionsCard({ reactions, onUpdate }: ReactionsCardProps) {
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

    const handleSlotChange = (index: number, value: string) => {
        const newSlots = [...slots];
        // Take only the first character if it's an emoji/text, or let it be for now
        // Usually people might paste an emoji or type it.
        newSlots[index] = value;
        setSlots(newSlots);
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
}
