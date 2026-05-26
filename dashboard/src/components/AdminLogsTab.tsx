import { useEffect, useMemo, useState } from 'react';
import { Calendar, ChevronDown, ChevronRight, Hash, RefreshCcw, Search } from 'lucide-react';
import { fetchAdminLogs } from '../api';
import { AdminLogsFilters, ChannelEvent } from '../types';
import { useToast } from './Toast';

interface AdminLogsTabProps {
  navigateToChannel: (id: number) => void;
  initialChannelId?: string;
}

const sourceLabels: Record<string, string> = {
  channel_post: 'Postagem',
  post_builder: 'PostBuilder',
};

const statusLabels: Record<string, string> = {
  success: 'Sucesso',
  error: 'Erro',
  skipped: 'Ignorado',
  info: 'Info',
};

function eventLabel(value: string): string {
  return value.replaceAll('_', ' ');
}

function formatDate(value: string): string {
  if (!value) return '-';
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(value));
}

function parseMetadata(event: ChannelEvent): string {
  if (!event.metadata) return '{}';
  try {
    return JSON.stringify(JSON.parse(event.metadata), null, 2);
  } catch {
    return event.metadata;
  }
}

function statusBadgeClass(status: string): string {
  switch (status) {
    case 'success':
      return 'badge-success';
    case 'info':
      return 'badge-info';
    case 'error':
      return 'badge-danger';
    case 'skipped':
      return 'badge-warning';
    default:
      return 'badge-muted';
  }
}

export function AdminLogsTab({ navigateToChannel, initialChannelId = '' }: AdminLogsTabProps) {
  const toast = useToast();
  const [filters, setFilters] = useState<AdminLogsFilters>({ limit: 50, offset: 0, channelId: initialChannelId });
  const [events, setEvents] = useState<ChannelEvent[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const page = useMemo(() => Math.floor((filters.offset || 0) / (filters.limit || 50)) + 1, [filters.offset, filters.limit]);
  const pageCount = useMemo(() => Math.max(1, Math.ceil(total / (filters.limit || 50))), [total, filters.limit]);

  const loadLogs = async (nextFilters = filters) => {
    setLoading(true);
    try {
      const data = await fetchAdminLogs(nextFilters);
      setEvents(data.events || []);
      setTotal(data.total || 0);
    } catch (err: any) {
      toast(err.message || 'Erro ao carregar logs', 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadLogs(filters);
  }, []);

  const updateFilter = (key: keyof AdminLogsFilters, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value, offset: 0 }));
  };

  const applyFilters = () => {
    const next = { ...filters, offset: 0 };
    setFilters(next);
    loadLogs(next);
  };

  const goToPage = (direction: 'prev' | 'next') => {
    const limit = filters.limit || 50;
    const offset = Math.max(0, (filters.offset || 0) + (direction === 'next' ? limit : -limit));
    const next = { ...filters, offset };
    setFilters(next);
    loadLogs(next);
  };

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        <div className="search-bar-container relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2" size={18} style={{ color: 'var(--hint)' }} />
          <input className="admin-search-input input" placeholder="Buscar título, erro ou metadata" value={filters.q || ''} onChange={e => updateFilter('q', e.target.value)} />
        </div>
        <div className="search-bar-container relative">
          <Hash className="absolute left-4 top-1/2 -translate-y-1/2" size={18} style={{ color: 'var(--hint)' }} />
          <input className="admin-search-input input" placeholder="ID do canal" value={filters.channelId || ''} onChange={e => updateFilter('channelId', e.target.value)} />
        </div>
        <select className="input" value={filters.source || ''} onChange={e => updateFilter('source', e.target.value)}>
          <option value="">Todas as origens</option>
          <option value="channel_post">Postagens</option>
          <option value="post_builder">PostBuilder</option>
        </select>
        <select className="input" value={filters.status || ''} onChange={e => updateFilter('status', e.target.value)}>
          <option value="">Todos os status</option>
          <option value="success">Sucesso</option>
          <option value="error">Erro</option>
          <option value="skipped">Ignorado</option>
          <option value="info">Info</option>
        </select>
        <div className="search-bar-container relative">
          <Calendar className="absolute left-4 top-1/2 -translate-y-1/2" size={18} style={{ color: 'var(--hint)' }} />
          <input className="admin-search-input input" type="date" value={filters.dateFrom || ''} onChange={e => updateFilter('dateFrom', e.target.value)} />
        </div>
        <div className="search-bar-container relative">
          <Calendar className="absolute left-4 top-1/2 -translate-y-1/2" size={18} style={{ color: 'var(--hint)' }} />
          <input className="admin-search-input input" type="date" value={filters.dateTo || ''} onChange={e => updateFilter('dateTo', e.target.value)} />
        </div>
      </div>

      <div className="flex gap-2">
        <button className="btn btn-primary flex-1" onClick={applyFilters} disabled={loading}>
          <Search size={18} /> Buscar
        </button>
        <button className="btn btn-secondary" onClick={() => loadLogs(filters)} disabled={loading}>
          <RefreshCcw size={18} />
        </button>
      </div>

      <div className="flex items-center justify-between text-xs" style={{ color: 'var(--hint)' }}>
        <span>{total} eventos</span>
        <span>Página {page} de {pageCount}</span>
      </div>

      <div className="space-y-3">
        {events.length === 0 && !loading ? (
          <div className="card text-center py-8" style={{ color: 'var(--hint)' }}>Nenhum log encontrado</div>
        ) : events.map(event => {
          const expanded = expandedId === event.id;
          return (
            <div key={event.id} className="admin-list-item p-4">
              <button className="flex items-start w-full text-left gap-3" onClick={() => setExpandedId(expanded ? null : event.id)}>
                <div className="section-icon purple mt-0.5">{expanded ? <ChevronDown size={18} /> : <ChevronRight size={18} />}</div>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap gap-2 mb-1">
                    <span className={`badge ${statusBadgeClass(event.status)}`}>{statusLabels[event.status] || event.status}</span>
                    <span className="badge badge-muted">{sourceLabels[event.source] || event.source}</span>
                  </div>
                  <h3 className="text-[15px] font-semibold truncate">{eventLabel(event.eventType)}</h3>
                  <p className="text-xs truncate mt-0.5" style={{ color: 'var(--hint)' }}>
                    {event.channelTitle || 'Sem canal'} {event.channelId ? `(${event.channelId})` : ''} • {formatDate(event.created_at)}
                  </p>
                  {event.errorMessage && <p className="text-xs mt-1 text-[var(--danger)] truncate">{event.errorMessage}</p>}
                </div>
              </button>

              {expanded && (
                <div className="mt-3 space-y-3">
                  <div className="grid grid-cols-2 gap-2 text-xs" style={{ color: 'var(--hint)' }}>
                    <span>Owner: {event.ownerId || '-'}</span>
                    <span>Actor: {event.actorId || '-'}</span>
                    <span>Mensagem: {event.telegramMessageId || '-'}</span>
                    <span>Sessão: {event.sessionId || '-'}</span>
                  </div>
                  {event.channelId !== 0 && (
                    <button className="btn btn-secondary w-full" onClick={() => navigateToChannel(event.channelId)}>
                      <Hash size={18} /> Abrir canal
                    </button>
                  )}
                  <pre className="text-xs overflow-auto rounded-md p-3" style={{ background: 'var(--bg-secondary)', color: 'var(--text)', maxHeight: 220 }}>{parseMetadata(event)}</pre>
                </div>
              )}
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-2 gap-2">
        <button className="btn btn-secondary" disabled={loading || (filters.offset || 0) === 0} onClick={() => goToPage('prev')}>Anterior</button>
        <button className="btn btn-secondary" disabled={loading || page >= pageCount} onClick={() => goToPage('next')}>Próxima</button>
      </div>
    </div>
  );
}
