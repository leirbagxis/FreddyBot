import { useState, useEffect } from 'react';
import {
  connectStart, connectVerify, connect2FA,
  connectStatus, connectDisconnect,
} from '../api';

type ConnectStep = 'idle' | 'phone' | 'code' | '2fa' | 'success' | 'error';

export default function TelegramConnect() {
  const [step, setStep] = useState<ConnectStep>('idle');
  const [phone, setPhone] = useState('');
  const [code, setCode] = useState('');
  const [password, setPassword] = useState('');
  const [errorMsg, setErrorMsg] = useState('');
  const [loading, setLoading] = useState(false);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const res = await connectStatus();
        if (res.connected) {
          setConnected(true);
          setStep('success');
        } else {
          setStep('phone');
        }
      } catch {
        setStep('phone');
      }
    })();
  }, []);

  const handlePhoneSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setErrorMsg('');

    const formatted = phone.startsWith('+') ? phone : `+${phone.replace(/\s/g, '')}`;
    if (formatted.length < 10) {
      setErrorMsg('Número inválido. Use o formato internacional (ex: +5511999999999)');
      setLoading(false);
      return;
    }

    try {
      await connectStart(formatted);
      setStep('code');
    } catch (err: any) {
      const msg = err?.message || 'Erro ao iniciar conexão';
      try {
        const parsed = JSON.parse(msg);
        setErrorMsg(parsed.message || msg);
      } catch {
        setErrorMsg(msg);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCodeSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setErrorMsg('');

    try {
      const res = await connectVerify(code);
      if (res?.data?.needs2FA) {
        setStep('2fa');
      } else {
        setConnected(true);
        setStep('success');
      }
    } catch (err: any) {
      const msg = err?.message || 'Código inválido';
      try {
        const parsed = JSON.parse(msg);
        setErrorMsg(parsed.message || msg);
      } catch {
        setErrorMsg(msg);
      }
    } finally {
      setLoading(false);
    }
  };

  const handle2FASubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setErrorMsg('');

    try {
      await connect2FA(password);
      setConnected(true);
      setStep('success');
    } catch (err: any) {
      const msg = err?.message || 'Senha inválida';
      try {
        const parsed = JSON.parse(msg);
        setErrorMsg(parsed.message || msg);
      } catch {
        setErrorMsg(msg);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleDisconnect = async () => {
    setLoading(true);
    try {
      await connectDisconnect();
      setConnected(false);
      setStep('phone');
      setPhone('');
      setCode('');
      setPassword('');
    } catch (err: any) {
      setErrorMsg('Erro ao desconectar');
    } finally {
      setLoading(false);
    }
  };

  const formatPhone = (val: string) => {
    const digits = val.replace(/\D/g, '');
    if (digits.startsWith('55') && digits.length > 2) {
      const country = digits.slice(0, 2);
      const rest = digits.slice(2);
      if (rest.length <= 2) return `+${country} (${rest}`;
      if (rest.length <= 7) return `+${country} (${rest.slice(0, 2)}) ${rest.slice(2)}`;
      return `+${country} (${rest.slice(0, 2)}) ${rest.slice(2, 7)}-${rest.slice(7, 11)}`;
    }
    return `+${digits}`;
  };

  if (step === 'idle') {
    return (
      <div className="app-layout">
        <div className="main-content" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
          <div className="auth-spinner" />
        </div>
      </div>
    );
  }

  if (step === 'success') {
    return (
      <div className="app-layout">
        <div className="main-content" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh', gap: 16, textAlign: 'center', padding: 24 }}>
          <div style={{ width: 72, height: 72, borderRadius: '50%', background: 'var(--accent-soft)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="var(--accent)" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="20 6 9 17 4 12" />
            </svg>
          </div>
          <h2 style={{ fontSize: 22, fontWeight: 800 }}>Conta Conectada</h2>
          <p style={{ fontSize: 15, color: 'var(--hint)', maxWidth: 300, lineHeight: 1.6 }}>
            Sua conta Telegram está conectada ao bot. Agora você pode usar recursos avançados.
          </p>
          <button
            className="btn btn-secondary"
            onClick={handleDisconnect}
            disabled={loading}
            style={{ marginTop: 8, minWidth: 200 }}
          >
            {loading ? 'Desconectando...' : '🔌 Desconectar Conta'}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="app-layout">
      <div className="main-content" style={{ padding: 24, maxWidth: 400, margin: '0 auto' }}>
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <div style={{ width: 72, height: 72, borderRadius: '50%', background: 'var(--accent-soft)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="var(--accent)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z" />
            </svg>
          </div>
          <h2 style={{ fontSize: 22, fontWeight: 800 }}>Conectar Conta</h2>
          <p style={{ fontSize: 14, color: 'var(--hint)', marginTop: 4 }}>
            {step === 'phone' && 'Digite seu número de telefone com código do país'}
            {step === 'code' && 'Digite o código enviado pelo Telegram'}
            {step === '2fa' && 'Digite sua senha de verificação em duas etapas'}
          </p>
        </div>

        {errorMsg && (
          <div style={{ padding: '12px 16px', borderRadius: 12, background: 'var(--danger-soft)', color: 'var(--danger)', fontSize: 13, marginBottom: 16, textAlign: 'center' }}>
            {errorMsg}
          </div>
        )}

        {step === 'phone' && (
          <form onSubmit={handlePhoneSubmit}>
            <div style={{ marginBottom: 16 }}>
              <input
                type="tel"
                inputMode="tel"
                placeholder="+5511999999999"
                value={phone}
                onChange={e => setPhone(formatPhone(e.target.value))}
                style={{
                  width: '100%',
                  padding: '14px 16px',
                  borderRadius: 12,
                  border: '1px solid var(--border)',
                  background: 'var(--surface)',
                  color: 'var(--text)',
                  fontSize: 16,
                  outline: 'none',
                  boxSizing: 'border-box',
                }}
                autoFocus
                disabled={loading}
              />
            </div>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading || phone.length < 5}
              style={{ width: '100%', justifyContent: 'center' }}
            >
              {loading ? 'Enviando...' : 'Enviar Código'}
            </button>
          </form>
        )}

        {step === 'code' && (
          <form onSubmit={handleCodeSubmit}>
            <div style={{ marginBottom: 16 }}>
              <input
                type="text"
                inputMode="numeric"
                placeholder="00000"
                value={code}
                onChange={e => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                style={{
                  width: '100%',
                  padding: '14px 16px',
                  borderRadius: 12,
                  border: '1px solid var(--border)',
                  background: 'var(--surface)',
                  color: 'var(--text)',
                  fontSize: 24,
                  textAlign: 'center',
                  letterSpacing: 8,
                  outline: 'none',
                  boxSizing: 'border-box',
                }}
                autoFocus
                disabled={loading}
              />
            </div>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading || code.length < 4}
              style={{ width: '100%', justifyContent: 'center' }}
            >
              {loading ? 'Verificando...' : 'Verificar Código'}
            </button>
          </form>
        )}

        {step === '2fa' && (
          <form onSubmit={handle2FASubmit}>
            <div style={{ marginBottom: 16 }}>
              <input
                type="password"
                placeholder="Sua senha 2FA"
                value={password}
                onChange={e => setPassword(e.target.value)}
                style={{
                  width: '100%',
                  padding: '14px 16px',
                  borderRadius: 12,
                  border: '1px solid var(--border)',
                  background: 'var(--surface)',
                  color: 'var(--text)',
                  fontSize: 16,
                  outline: 'none',
                  boxSizing: 'border-box',
                }}
                autoFocus
                disabled={loading}
              />
            </div>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading || !password}
              style={{ width: '100%', justifyContent: 'center' }}
            >
              {loading ? 'Autenticando...' : 'Confirmar'}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
