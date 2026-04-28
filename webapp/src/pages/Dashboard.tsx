import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { fetchMe, getBotId, consumeItem, sellItemToShop, fetchMyLogs } from '../api';
import { useBotConfig } from '../hooks/useBotConfig';
import { showToast } from '../components/Toast';
import ConfirmationModal from '../components/ConfirmationModal';
import InputModal from '../components/InputModal';
import { Activity, Hash, UserMinus } from 'lucide-react';

// Sub-components
import UserOverview from '../components/dashboard/UserOverview';
import Inventory from '../components/dashboard/Inventory';
import FinanceHistory from '../components/dashboard/FinanceHistory';

interface InventoryItem {
  item_id: string;
  name: string;
  quantity: number;
  item_type: string;
}

interface UserProfile {
  telegram_user_id: number;
  username: string;
  first_name: string;
  last_name: string;
  balance: number;
  inventory: InventoryItem[];
}

export default function Dashboard() {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [logs, setLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [selectedItem, setSelectedItem] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  
  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const activeTab = query.get('tab') || 'overview';

  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    title: string;
    message: string;
    onConfirm: () => void;
    type: 'danger' | 'primary' | 'success';
  }>({
    isOpen: false,
    title: '',
    message: '',
    onConfirm: () => {},
    type: 'primary'
  });
  
  const [inputModal, setInputModal] = useState<{
    isOpen: boolean;
    itemId: string;
    itemName: string;
    title: string;
    label: string;
    placeholder: string;
    icon: any;
    successMessage: string;
  }>({
    isOpen: false,
    itemId: '',
    itemName: '',
    title: '',
    label: '',
    placeholder: '',
    icon: Hash,
    successMessage: '',
  });

  const botId = getBotId();
  const botConfig = useBotConfig();

  const getBotDisplayName = () => {
    const { bot_name, bot_username } = botConfig;
    if (bot_name && bot_name !== bot_username && bot_name !== bot_username.replace('@', '')) {
      return bot_name;
    }
    return (bot_username || 'Meu Bot').replace('@', '').toUpperCase();
  };

  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Bom dia';
    if (hour < 18) return 'Boa tarde';
    return 'Boa noite';
  };

  const loadData = async () => {
    try {
      const [profileData, logsData] = await Promise.all([
        fetchMe(botId),
        fetchMyLogs(botId, 20, 0)
      ]);
      setProfile(profileData);
      setLogs(logsData || []);
      if ((logsData || []).length < 20) setHasMore(false);
    } catch (err) {
      console.error("Failed to load dashboard data", err);
    } finally {
      setLoading(false);
    }
  };

  const handleLoadMore = async () => {
    if (loadingMore || !hasMore) return;
    setLoadingMore(true);
    try {
      const newLogs = await fetchMyLogs(botId, 20, logs.length);
      if (newLogs.length < 20) setHasMore(false);
      setLogs(prev => [...prev, ...newLogs]);
    } catch (err) {
      showToast.error("Erro ao carregar mais transações");
    } finally {
      setLoadingMore(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [botId]);

  const handleUseItem = async (itemId: string, itemName: string, itemType?: string) => {
    if (itemType === 'tag') {
      setInputModal({ 
        isOpen: true, 
        itemId, 
        itemName,
        title: 'Configurar Tag',
        label: 'Sua Nova Tag',
        placeholder: 'Ex: VIP, STAFF...',
        icon: Hash,
        successMessage: 'Pedido para a TAG "{value}" enviado!'
      });
      return;
    }

    if (itemType === 'ban') {
      setInputModal({ 
        isOpen: true, 
        itemId, 
        itemName,
        title: 'Banir Usuário',
        label: 'ID ou Username do Alvo',
        placeholder: 'Ex: 123456 ou @usuario',
        icon: UserMinus,
        successMessage: 'Solicitação de banimento para "{value}" enviada!'
      });
      return;
    }

    setConfirmModal({
      isOpen: true,
      title: 'Usar Item',
      message: `Tem certeza que deseja usar o item "${itemName}" agora?`,
      type: 'primary',
      onConfirm: async () => {
        setIsProcessing(true);
        setConfirmModal(prev => ({ ...prev, isOpen: false }));
        try {
          await consumeItem(botId, itemId);
          showToast.success(`Item "${itemName}" usado com sucesso!`);
          loadData();
        } catch (err) {
          showToast.error("Erro ao usar item");
        } finally {
          setIsProcessing(false);
          setSelectedItem(null);
        }
      }
    });
  };

  const handleInputConfirm = async (value: string) => {
    const { itemId, successMessage } = inputModal;
    setIsProcessing(true);
    setInputModal(prev => ({ ...prev, isOpen: false }));
    try {
      await consumeItem(botId, itemId, value);
      showToast.success(successMessage.replace('{value}', value));
      loadData();
    } catch (err) {
      showToast.error("Erro ao processar item");
    } finally {
      setIsProcessing(false);
      setSelectedItem(null);
    }
  };

  const handleSellToShop = async (itemId: string, itemName: string) => {
    setConfirmModal({
      isOpen: true,
      title: 'Vender Item',
      message: `Deseja vender o item "${itemName}" para a loja? Esta ação não pode ser desfeita.`,
      type: 'danger',
      onConfirm: async () => {
        setIsProcessing(true);
        setConfirmModal(prev => ({ ...prev, isOpen: false }));
        try {
          await sellItemToShop(botId, itemId);
          showToast.success(`Item "${itemName}" vendido!`);
          loadData();
        } catch (err) {
          showToast.error("Erro ao vender item");
        } finally {
          setIsProcessing(false);
          setSelectedItem(null);
        }
      }
    });
  };

  const getOperationLabel = (type: string) => {
    switch (type) {
      case 'purchase_web':
      case 'purchase': return 'Compra na Loja';
      case 'admin_add': return 'Crédito Manual';
      case 'admin_reduce': return 'Débito Manual';
      case 'reward': return 'Recompensa';
      case 'sell_shop_web':
      case 'sell_shop': return 'Venda de Item';
      case 'consume': return 'Uso de Item';
      case 'sell_player': return 'Transferência (Venda)';
      case 'buy_player': return 'Transferência (Compra)';
      default: return type.replace(/_/g, ' ').toUpperCase();
    }
  };

  const getLogDescription = (log: any) => {
    try {
      const meta = typeof log.metadata === 'string' ? JSON.parse(log.metadata) : log.metadata;
      if (meta?.item_name) return meta.item_name;
      if (meta?.reason) return meta.reason;
    } catch (e) {
      // ignore
    }
    return getOperationLabel(log.operation_type);
  };

  if (loading) return (
    <div className="flex flex-col items-center justify-center py-xl">
      <Activity className="animate-spin text-primary mb-md" size={32} />
      <div className="status-label tracking-[0.3em] text-primary">SYNC_IN_PROGRESS...</div>
    </div>
  );

  if (!profile) return (
    <div className="refined-card bg-red-950/20 p-xl text-center border border-red-500/20">
      <div className="status-label text-red-500">CONNECTION_ERROR</div>
      <div className="text-sm font-bold opacity-60 text-red-200">Não foi possível carregar seu perfil neural.</div>
    </div>
  );

  return (
    <div className="space-y-lg animate-fade-in pb-20 font-sans">
      <header className="mb-8 border-b border-white/5 pb-6 relative overflow-hidden">
        <div className="absolute top-0 right-0 w-32 h-32 bg-primary/5 blur-3xl -z-10 rounded-full" />
        <h2 className="text-4xl font-black text-white tracking-tighter truncate pr-4 uppercase glitch-text" data-text={`${getGreeting()}, ${profile.first_name}!`}>
          {getGreeting()}, {profile.first_name}!
        </h2>
        <p className="text-primary font-bold text-[9px] uppercase tracking-[0.3em] mt-2 flex items-center gap-2">
          <span className="w-1.5 h-1.5 rounded-full bg-primary animate-pulse" />
          STATUS_CENTRAL // <span className="text-white font-black">{getBotDisplayName()}</span>
        </p>
      </header>

      {activeTab === 'history' ? (
        <FinanceHistory 
          logs={logs}
          hasMore={hasMore}
          loadingMore={loadingMore}
          handleLoadMore={handleLoadMore}
          getLogDescription={getLogDescription}
          getOperationLabel={getOperationLabel}
        />
      ) : (
        <div className="space-y-lg">
          <UserOverview botId={botId} profile={profile} botConfig={botConfig} />
          <Inventory 
            inventory={profile.inventory || []}
            selectedItem={selectedItem}
            setSelectedItem={setSelectedItem}
            handleUseItem={handleUseItem}
            handleSellToShop={handleSellToShop}
          />
        </div>
      )}

      <ConfirmationModal
        isOpen={confirmModal.isOpen}
        title={confirmModal.title}
        message={confirmModal.message}
        type={confirmModal.type}
        onConfirm={confirmModal.onConfirm}
        onClose={() => setConfirmModal(prev => ({ ...prev, isOpen: false }))}
        isLoading={isProcessing}
      />

      {inputModal.isOpen && (
        <InputModal 
          title={inputModal.title}
          label={inputModal.label}
          placeholder={inputModal.placeholder}
          icon={inputModal.icon}
          itemName={inputModal.itemName}
          onClose={() => setInputModal(prev => ({ ...prev, isOpen: false }))}
          onConfirm={handleInputConfirm}
          isLoading={isProcessing}
        />
      )}

      <style>{`
        .stagger-item {
          animation: slide-up 0.5s cubic-bezier(0.165, 0.84, 0.44, 1) forwards;
          opacity: 0;
        }

        @keyframes slide-up {
          from { transform: translateY(20px); opacity: 0; }
          to { transform: translateY(0); opacity: 1; }
        }
      `}</style>
    </div>
  );
}
