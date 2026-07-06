import { useEffect, useMemo, useState } from 'react';
import { Calendar, ChevronDown, ChevronRight, Hash, RefreshCcw, Search, ChevronLeft, ChevronRight as ChevronRightIcon } from 'lucide-react';
import { fetchAdminLogs } from '../api';
import { AdminLogsFilters, ChannelEvent } from '../types';
import { useToast } from './Toast';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Input } from './ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';

interface AdminLogsTabProps {
  navigateToChannel: (id: number) => void;
  initialChannelId?: string;
}

const sourceLabels: Record<string, string> = {
  channel_post: 'Postagem',
  post_builder: 'PostBuilder',
};

const statusConfig: Record<string, { label: string; color: string; bg: string }> = {
  success: { label: 'Sucesso', color: 'var(--success)', bg: 'var(--success-soft)' },
  error: { label: 'Erro', color: 'var(--danger)', bg: 'var(--danger-soft)' },
  info: { label: 'Info', color: 'var(--accent)', bg: 'var(--accent-soft)' },
  skipped: { label: 'Ignorado', color: 'var(--warning)', bg: 'var(--warning-soft)' },
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
      {/* Filter bar */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-2">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
          <Input
            className="h-10 pl-9 rounded-xl"
            placeholder="Buscar..."
            value={filters.q || ''}
            onChange={e => updateFilter('q', e.target.value)}
          />
        </div>
        <div className="relative">
          <Hash className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
          <Input
            className="h-10 pl-9 rounded-xl"
            placeholder="ID do canal"
            value={filters.channelId || ''}
            onChange={e => updateFilter('channelId', e.target.value)}
          />
        </div>
        <Select value={filters.source || ''} onValueChange={v => updateFilter('source', v)}>
          <SelectTrigger className="w-full h-10 rounded-xl">
            <SelectValue placeholder="Todas origens" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">Todas origens</SelectItem>
            <SelectItem value="channel_post">Postagens</SelectItem>
            <SelectItem value="post_builder">PostBuilder</SelectItem>
          </SelectContent>
        </Select>
        <Select value={filters.status || ''} onValueChange={v => updateFilter('status', v)}>
          <SelectTrigger className="w-full h-10 rounded-xl">
            <SelectValue placeholder="Todos status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">Todos status</SelectItem>
            <SelectItem value="success">Sucesso</SelectItem>
            <SelectItem value="error">Erro</SelectItem>
            <SelectItem value="skipped">Ignorado</SelectItem>
            <SelectItem value="info">Info</SelectItem>
          </SelectContent>
        </Select>
        <div className="relative">
          <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
          <Input
            className="h-10 pl-9 rounded-xl"
            type="date"
            value={filters.dateFrom || ''}
            onChange={e => updateFilter('dateFrom', e.target.value)}
          />
        </div>
        <div className="relative">
          <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
          <Input
            className="h-10 pl-9 rounded-xl"
            type="date"
            value={filters.dateTo || ''}
            onChange={e => updateFilter('dateTo', e.target.value)}
          />
        </div>
        <Button variant="default" className="h-10" onClick={applyFilters} disabled={loading}>
          <Search size={16} /> Buscar
        </Button>
        <Button variant="secondary" className="h-10" onClick={() => loadLogs(filters)} disabled={loading}>
          <RefreshCcw size={16} /> Atualizar
        </Button>
      </div>

      {/* Status bar */}
      <div className="flex items-center justify-between text-xs text-muted-foreground">
        <span className="font-medium">{total} eventos encontrados</span>
        <span>Página {page} de {pageCount}</span>
      </div>

      {/* Events list */}
      <div className="space-y-1.5">
        {events.length === 0 && !loading ? (
          <div className="flex flex-col items-center py-10 text-muted-foreground rounded-xl border border-border">
            <Search size={28} className="opacity-30 mb-2" />
            <p className="text-[13px] font-medium">Nenhum log encontrado</p>
            <p className="text-[11px] text-muted-foreground/60 mt-1">Tente ajustar os filtros ou carregar mais dados</p>
          </div>
        ) : events.map(event => {
          const expanded = expandedId === event.id;
          const sc = statusConfig[event.status] || { label: event.status, color: 'var(--hint)', bg: 'transparent' };
          return (
            <div key={event.id} className="rounded-xl border border-border overflow-hidden transition-all">
              {/* Main row */}
              <button
                className="flex items-start w-full text-left gap-3 p-3 hover:bg-muted/20 transition-colors"
                onClick={() => setExpandedId(expanded ? null : event.id)}
              >
                <div className="shrink-0 mt-0.5 text-muted-foreground">
                  {expanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                </div>
                <div className="shrink-0 mt-1" style={{ width: 8, height: 8, borderRadius: '50%', background: sc.color }} />
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="text-[13px] font-semibold">{eventLabel(event.eventType)}</span>
                    <Badge
                      variant="secondary"
                      className="text-[9px] h-[18px]"
                      style={{ background: sc.bg, color: sc.color, border: 'none' }}
                    >
                      {sc.label}
                    </Badge>
                    <Badge variant="secondary" className="text-[9px] h-[18px]">
                      {sourceLabels[event.source] || event.source}
                    </Badge>
                    <span className="text-[10px] text-muted-foreground ml-auto">{formatDate(event.created_at)}</span>
                  </div>
                  <p className="text-[11px] text-muted-foreground mt-0.5 truncate">
                    {event.channelTitle || 'Sem canal'}{event.channelId ? ` (${event.channelId})` : ''}
                  </p>
                  {event.errorMessage && (
                    <p className="text-[11px] mt-1 text-destructive truncate">{event.errorMessage}</p>
                  )}
                </div>
              </button>

              {/* Expanded details */}
              {expanded && (
                <div className="border-t border-border px-3 py-3 space-y-3 bg-muted/10">
                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs text-muted-foreground">
                    <div className="rounded-lg bg-muted/20 px-3 py-2">
                      <span className="block text-[10px] font-semibold text-muted-foreground/60 uppercase tracking-wide">Owner</span>
                      <span className="font-medium">{event.ownerId || '-'}</span>
                    </div>
                    <div className="rounded-lg bg-muted/20 px-3 py-2">
                      <span className="block text-[10px] font-semibold text-muted-foreground/60 uppercase tracking-wide">Actor</span>
                      <span className="font-medium">{event.actorId || '-'}</span>
                    </div>
                    <div className="rounded-lg bg-muted/20 px-3 py-2">
                      <span className="block text-[10px] font-semibold text-muted-foreground/60 uppercase tracking-wide">Mensagem ID</span>
                      <span className="font-medium">{event.telegramMessageId || '-'}</span>
                    </div>
                    <div className="rounded-lg bg-muted/20 px-3 py-2">
                      <span className="block text-[10px] font-semibold text-muted-foreground/60 uppercase tracking-wide">Sessão</span>
                      <span className="font-medium truncate block">{event.sessionId || '-'}</span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    {event.channelId !== 0 && (
                      <Button variant="secondary" size="sm" onClick={() => navigateToChannel(event.channelId)}>
                        <Hash size={14} /> Abrir canal
                      </Button>
                    )}
                  </div>
                  <pre
                    className="text-[11px] overflow-auto rounded-xl p-3 leading-relaxed max-h-[200px] border border-border"
                    style={{ background: 'var(--muted)' }}
                  >
                    {parseMetadata(event)}
                  </pre>
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* Pagination */}
      {pageCount > 1 && (
        <div className="flex items-center justify-center gap-2 pt-2">
          <Button
            variant="secondary"
            size="sm"
            disabled={loading || (filters.offset || 0) === 0}
            onClick={() => goToPage('prev')}
          >
            <ChevronLeft size={16} /> Anterior
          </Button>
          <div className="flex items-center gap-1 text-xs text-muted-foreground px-3">
            {Array.from({ length: Math.min(pageCount, 5) }, (_, i) => {
              let pageNum: number;
              if (pageCount <= 5) {
                pageNum = i + 1;
              } else if (page <= 3) {
                pageNum = i + 1;
              } else if (page >= pageCount - 2) {
                pageNum = pageCount - 4 + i;
              } else {
                pageNum = page - 2 + i;
              }
              return (
                <button
                  key={pageNum}
                  onClick={() => {
                    const offset = (pageNum - 1) * (filters.limit || 50);
                    const next = { ...filters, offset };
                    setFilters(next);
                    loadLogs(next);
                  }}
                  className={`w-7 h-7 rounded-lg text-[11px] font-semibold transition-colors ${
                    pageNum === page
                      ? 'bg-accent text-accent-foreground'
                      : 'hover:bg-muted/30 text-muted-foreground'
                  }`}
                >
                  {pageNum}
                </button>
              );
            })}
          </div>
          <Button
            variant="secondary"
            size="sm"
            disabled={loading || page >= pageCount}
            onClick={() => goToPage('next')}
          >
            Próxima <ChevronRightIcon size={16} />
          </Button>
        </div>
      )}
    </div>
  );
}
