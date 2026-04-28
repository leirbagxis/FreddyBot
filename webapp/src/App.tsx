import { useEffect, useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Admin from './pages/Admin';
import ChannelEditor from './pages/ChannelEditor';
import { login } from './api';

declare global {
  interface Window {
    Telegram: any;
  }
}

function App() {
  const [isAuthed, setIsAuthed] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const authenticate = async () => {
      try {
        const tg = window.Telegram?.WebApp;
        if (tg) {
          tg.ready();
          tg.expand();
        }
        
        // Extração robusta do channelId (antigo botId)
        const pathname = window.location.pathname;
        const match = pathname.match(/dashboard\/([^/]+)/);
        const channelIdStr = match ? match[1] : null;
        const channelId = channelIdStr ? parseInt(channelIdStr, 10) : NaN;

        if (tg?.initData && !isNaN(channelId)) {
          try {
            const user = tg.initDataUnsafe?.user;
            const response = await login(channelId, user);
            
            if (response.success) {
              setIsAuthed(true);
            } else {
              setError(response.message || "Falha na autenticação");
            }
          } catch (err: any) {
            console.error("Auth failed", err);
            if (err.response?.status === 401 || err.response?.status === 403) {
              setError("Acesso negado. Você não tem permissão para acessar este canal.");
            } else {
              setError(`Falha na autenticação: ${err.response?.data?.message || err.message}`);
            }
          }
        } else {
          // Dev mode
          setIsAuthed(true); 
        }
      } catch (e: any) {
        console.error("Init error", e);
        setError(`Erro de inicialização: ${e.message}`);
      } finally {
        setLoading(false);
      }
    };

    authenticate();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-surface flex flex-col items-center justify-center">
        <div className="w-12 h-12 border-4 border-black border-t-primary rounded-full animate-spin mb-4"></div>
        <div className="text-black font-black uppercase text-[10px] tracking-[0.3em]">
          Iniciando Sessão...
        </div>
      </div>
    );
  }

  if (error && !isAuthed) {
    return (
      <div className="min-h-screen bg-surface flex flex-col items-center justify-center p-10">
        <div className="bg-red-50 border-2 border-red-100 p-8 rounded-[2rem] text-center max-w-sm shadow-xl">
          <div className="text-red-500 font-black uppercase text-xs tracking-[0.2em] mb-4">Erro Crítico</div>
          <div className="text-red-900 text-sm font-bold uppercase tracking-tight leading-relaxed">{error}</div>
          <button onClick={() => window.location.reload()} className="mt-6 px-6 py-2 bg-red-600 text-white rounded-full text-[10px] font-black uppercase tracking-widest">Tentar Novamente</button>
        </div>
      </div>
    );
  }

  return (
    <BrowserRouter basename="/dashboard">
      <Routes>
        <Route path="/:botId" element={<Layout />}>
          <Route index element={<Dashboard />} />
          <Route path="admin" element={<Admin />} />
          <Route path="channel/:channelId" element={<ChannelEditor />} />
        </Route>
        
        {/* Redirecionamentos para quando não houver botId na URL */}
        <Route path="/" element={<Navigate to="/default" replace />} />
        <Route path="*" element={<Navigate to="/default" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
