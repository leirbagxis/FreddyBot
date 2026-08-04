import { useState, useEffect } from 'react';
import { ScheduledPost } from '../types';
import { fetchMySchedules, updateScheduleStatus, deleteSchedule, updateScheduleTime } from '../api';
import { useToast } from './Toast';
import { Card, CardContent } from './ui/card';
import { Badge } from './ui/badge';
import { Button } from './ui/button';
import { Clock, Trash2, Pause, Play, Calendar, Edit, Pin, PinOff } from 'lucide-react';

interface ScheduleTabProps {
  channelId: number;
}

const scheduleTypeLabels: Record<string, string> = {
  once: 'Único',
  daily: 'Diário',
  weekly: 'Semanal',
  queue: 'Fila',
};

const statusColors: Record<string, string> = {
  pending: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  scheduled: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  paused: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  cancelled: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  completed: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  sent: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200',
  error: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
};

export function ScheduleTab({ channelId }: ScheduleTabProps) {
  const [schedules, setSchedules] = useState<ScheduledPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingSchedule, setEditingSchedule] = useState<ScheduledPost | null>(null);
  const [editDate, setEditDate] = useState('');
  const [editTime, setEditTime] = useState('');
  const [saving, setSaving] = useState(false);
  const toast = useToast();

  useEffect(() => {
    loadSchedules();
  }, [channelId]);

  const loadSchedules = async () => {
    setLoading(true);
    try {
      const all = await fetchMySchedules();
      setSchedules(all.filter((s: ScheduledPost) => s.channelId === channelId));
    } catch {
      toast('Erro ao carregar agendamentos', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleTogglePause = async (schedule: ScheduledPost) => {
    const newStatus = schedule.status === 'paused' ? 'pending' : 'paused';
    try {
      await updateScheduleStatus(schedule.id, newStatus);
      setSchedules(prev =>
        prev.map(s => (s.id === schedule.id ? { ...s, status: newStatus } : s))
      );
      toast(newStatus === 'paused' ? 'Agendamento pausado' : 'Agendamento retomado', 'success');
    } catch {
      toast('Erro ao atualizar agendamento', 'error');
    }
  };

  const handleDelete = async (schedule: ScheduledPost) => {
    try {
      await deleteSchedule(schedule.id);
      setSchedules(prev => prev.filter(s => s.id !== schedule.id));
      toast('Agendamento removido', 'success');
    } catch {
      toast('Erro ao remover agendamento', 'error');
    }
  };

  const handleTogglePin = async (schedule: ScheduledPost) => {
    const newPin = !schedule.pinMessage;
    try {
      await updateScheduleTime(schedule.id, { pinMessage: newPin });
      setSchedules(prev =>
        prev.map(s => (s.id === schedule.id ? { ...s, pinMessage: newPin } : s))
      );
      toast(newPin ? 'Mensagem será fixada no canal' : 'Mensagem não será mais fixada', 'success');
    } catch {
      toast('Erro ao alterar fixação', 'error');
    }
  };

  const openEditModal = (schedule: ScheduledPost) => {
    setEditingSchedule(schedule);
    if (schedule.nextRunAt) {
      const d = new Date(schedule.nextRunAt);
      setEditDate(d.toISOString().split('T')[0]);
      setEditTime(d.toTimeString().slice(0, 5));
    }
  };

  const closeEditModal = () => {
    setEditingSchedule(null);
    setEditDate('');
    setEditTime('');
  };

  const handleSaveEdit = async () => {
    if (!editingSchedule) return;
    if (!editDate || !editTime) {
      toast('Preencha data e horário', 'error');
      return;
    }

    setSaving(true);
    try {
      const [year, month, day] = editDate.split('-').map(Number);
      const [hours, minutes] = editTime.split(':').map(Number);
      
      const dateStr = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}T${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:00-03:00`;
      
      await updateScheduleTime(editingSchedule.id, { nextRunAt: dateStr });
      
      setSchedules(prev =>
        prev.map(s => {
          if (s.id === editingSchedule.id) {
            return {
              ...s,
              nextRunAt: new Date(dateStr).toISOString(),
              scheduleTime: editingSchedule.scheduleType !== 'once' ? editTime : s.scheduleTime,
            };
          }
          return s;
        })
      );
      
      toast('Agendamento atualizado com sucesso', 'success');
      closeEditModal();
    } catch {
      toast('Erro ao atualizar agendamento', 'error');
    } finally {
      setSaving(false);
    }
  };

  const formatNextRun = (dateStr: string) => {
    try {
      const d = new Date(dateStr);
      return d.toLocaleString('pt-BR', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' });
    } catch {
      return dateStr;
    }
  };

  if (loading) {
    return (
      <div className="tab-content-wrapper flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--accent)]" />
      </div>
    );
  }

  if (schedules.length === 0) {
    return (
      <div className="tab-content-wrapper">
        <Card>
          <CardContent className="p-6 text-center text-muted-foreground">
            <Calendar className="mx-auto mb-3 h-10 w-10 opacity-40" />
            <p className="text-sm">Nenhum agendamento encontrado para este canal.</p>
            <p className="text-xs mt-1">Use o Post Builder no Telegram para criar agendamentos.</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="tab-content-wrapper space-y-3">
      <div className="flex items-center justify-between mb-2 px-1">
        <h3 className="text-sm font-semibold text-foreground">
          Agendamentos Ativos
        </h3>
        <Badge variant="secondary">{schedules.length}</Badge>
      </div>

      {schedules.map(schedule => (
        <Card key={schedule.id} className="transition-all hover:shadow-sm">
          <CardContent className="p-4">
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className="font-medium text-sm truncate">{schedule.channelTitle}</span>
                  <Badge variant="outline" className="text-[10px] shrink-0">
                    {scheduleTypeLabels[schedule.scheduleType] || schedule.scheduleType}
                  </Badge>
                  <Badge className={`text-[10px] shrink-0 ${statusColors[schedule.status] || ''}`}>
                    {schedule.status}
                  </Badge>
                </div>

                <div className="flex items-center gap-2 text-xs text-muted-foreground mt-1">
                  <Clock className="h-3 w-3" />
                  <span>Próximo: {formatNextRun(schedule.nextRunAt)}</span>
                </div>

                {schedule.scheduleTime && (
                  <p className="text-xs text-muted-foreground mt-1">
                    Horário: {schedule.scheduleTime}
                    {schedule.scheduleDays && ` | Dias: ${schedule.scheduleDays}`}
                  </p>
                )}

                {schedule.lastError && (
                  <p className="text-xs text-red-500 mt-1 truncate">
                    Erro: {schedule.lastError}
                  </p>
                )}
              </div>

              <div className="flex items-center gap-1 ml-2 shrink-0">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={() => openEditModal(schedule)}
                  title="Editar data/horário"
                >
                  <Edit className="h-4 w-4 text-blue-500" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={() => handleTogglePin(schedule)}
                  title={schedule.pinMessage ? 'Fixando mensagem' : 'Não fixar mensagem'}
                >
                  {schedule.pinMessage ? (
                    <PinOff className="h-4 w-4 text-purple-500" />
                  ) : (
                    <Pin className="h-4 w-4 text-gray-400" />
                  )}
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={() => handleTogglePause(schedule)}
                  title={schedule.status === 'paused' ? 'Retomar' : 'Pausar'}
                >
                  {schedule.status === 'paused' ? (
                    <Play className="h-4 w-4 text-green-500" />
                  ) : (
                    <Pause className="h-4 w-4 text-yellow-500" />
                  )}
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={() => handleDelete(schedule)}
                  title="Remover"
                >
                  <Trash2 className="h-4 w-4 text-red-500" />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}

      {editingSchedule && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="bg-background rounded-lg p-6 w-full max-w-md mx-4 shadow-lg">
            <h3 className="text-lg font-semibold mb-4">Editar Agendamento</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Data</label>
                <input
                  type="date"
                  value={editDate}
                  onChange={(e) => setEditDate(e.target.value)}
                  className="w-full px-3 py-2 border rounded-md bg-background text-foreground"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-1">Horário</label>
                <input
                  type="time"
                  value={editTime}
                  onChange={(e) => setEditTime(e.target.value)}
                  className="w-full px-3 py-2 border rounded-md bg-background text-foreground"
                />
              </div>
            </div>

            <div className="flex justify-end gap-2 mt-6">
              <Button
                variant="outline"
                onClick={closeEditModal}
                disabled={saving}
              >
                Cancelar
              </Button>
              <Button
                onClick={handleSaveEdit}
                disabled={saving}
              >
                {saving ? 'Salvando...' : 'Salvar'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
