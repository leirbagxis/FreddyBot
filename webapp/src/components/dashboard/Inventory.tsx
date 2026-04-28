import { Box, Play, ArrowRight, Coins } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: any[]) {
  return twMerge(clsx(inputs));
}

interface InventoryProps {
  inventory: any[];
  selectedItem: string | null;
  setSelectedItem: (id: string | null) => void;
  handleUseItem: (itemId: string, itemName: string, itemType?: string) => void;
  handleSellToShop: (itemId: string, itemName: string) => void;
}

export default function Inventory({ 
  inventory, 
  selectedItem, 
  setSelectedItem, 
  handleUseItem, 
  handleSellToShop 
}: InventoryProps) {
  return (
    <div className="refined-card border-white/5">
      <header className="flex items-center justify-between mb-xl border-b border-white/5 pb-md">
        <div className="flex items-center gap-md">
          <Box className="text-primary animate-pulse" size={20} />
          <h4 className="text-lg font-black text-white uppercase tracking-tighter">Inventário_Neural</h4>
        </div>
        <span className="bg-primary/10 border border-primary/20 text-primary px-4 py-1 rounded-full text-[9px] font-black uppercase tracking-widest">
          {inventory.length} units
        </span>
      </header>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-md">
        {inventory.length === 0 ? (
          <div className="col-span-full py-xl text-center text-white/10 font-black text-xs uppercase tracking-[0.3em] italic">-- EMPTY_STORAGE --</div>
        ) : (
          inventory.map((item, index) => (
            <div 
              key={item.item_id} 
              className={cn(
                "flex flex-col p-0 bg-white/5 rounded-3xl border-2 transition-all group overflow-hidden stagger-item",
                selectedItem === item.item_id ? "border-primary shadow-neon scale-[1.02] z-10" : "border-white/5 hover:border-primary/20"
              )}
              style={{ animationDelay: `${index * 50}ms` }}
            >
              <div className="flex justify-between items-center p-lg cursor-pointer" onClick={() => setSelectedItem(selectedItem === item.item_id ? null : item.item_id)}>
                <div className="flex flex-col">
                  <span className="text-sm font-black uppercase text-white/90 group-hover:text-primary transition-colors">{item.name}</span>
                  <span className="text-[8px] font-bold text-white/30 uppercase tracking-widest mt-1">SIG: {item.item_id.substring(0, 8)}</span>
                </div>
                <div className="bg-primary text-black px-3 py-1 rounded-full text-[10px] font-black shadow-neon">x{item.quantity}</div>
              </div>
              {selectedItem === item.item_id && (
                <div className="p-4 bg-black/40 backdrop-blur-md border-t border-white/5 flex flex-col gap-2 animate-in slide-in-from-top-2 duration-300">
                  <button onClick={() => handleUseItem(item.item_id, item.name, item.item_type)} className="flex items-center justify-between p-3 rounded-2xl bg-primary text-black hover:bg-white transition-all group/btn shadow-neon">
                    <div className="flex items-center gap-3"><Play size={14} className="fill-current" /><span className="text-[10px] font-black uppercase tracking-widest">EXECUTAR_ITEM</span></div>
                    <ArrowRight size={14} />
                  </button>
                  <button onClick={() => handleSellToShop(item.item_id, item.name)} className="flex items-center justify-center gap-2 p-3 rounded-2xl bg-white/5 text-white/60 hover:bg-accent hover:text-white transition-all border border-white/5">
                    <Coins size={14} /><span className="text-[10px] font-black uppercase tracking-widest">LIQUIDAR_VALOR</span>
                  </button>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}
