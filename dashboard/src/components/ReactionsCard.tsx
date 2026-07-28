import { useState, useEffect, memo } from 'react';
import { SmilePlus, X } from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';

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

        if (isEmoji(trimmed)) {
            const newSlots = [...slots];
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
            const reactionsString = slots.filter(s => s.trim() !== '').join(',');
            await onUpdate(reactionsString);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="content-card">
            <div className="content-card-header">
                <div className="content-card-icon"><SmilePlus size={18} /></div>
                <div className="min-w-0 flex-1">
                    <div className="content-card-title">Reações / Votos (Grid)</div>
                    <div className="content-card-desc">Adicione até 5 emojis para votação rápida.</div>
                </div>
            </div>

                <div>
                    <div className="grid grid-cols-5 gap-2 mb-4">
                        {slots.map((slot, index) => (
                            <div key={index} className="relative group">
                                <Input
                                    type="text"
                                    value={slot}
                                    onChange={(e) => handleSlotChange(index, e.target.value)}
                                    placeholder="+"
                                    className="aspect-square text-center text-xl p-0"
                                />
                                {slot && (
                                    <Button 
                                        variant="ghost"
                                        size="icon-xs"
                                        className="absolute -top-1.5 -right-1.5 bg-destructive text-destructive-foreground rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-destructive/90"
                                        onClick={() => handleClearSlot(index)}
                                    >
                                        <X size={10} />
                                    </Button>
                                )}
                            </div>
                        ))}
                    </div>
                    
                    <Button 
                        variant="default" 
                        className="w-full"
                        onClick={handleSave}
                        disabled={loading}
                    >
                        {loading ? 'Salvando...' : 'Salvar Reações'}
                    </Button>
                </div>
        </div>
    );
});
