import { AlertTriangle } from 'lucide-react';

interface ConfirmModalProps {
    open: boolean;
    onClose: () => void;
    onConfirm: () => void;
    title: string;
    message: string;
    confirmText?: string;
    danger?: boolean;
    alertOnly?: boolean;
}

export function ConfirmModal({
    open, onClose, onConfirm, title, message, confirmText, danger, alertOnly
}: ConfirmModalProps) {
    if (!open) return null;
    return (
        <div className="overlay" onClick={onClose}>
            <div className="dialog confirm-dialog" onClick={e => e.stopPropagation()}>
                <div className="dialog-handle" />
                <div className="confirm-icon-wrap" style={{ background: danger ? 'var(--danger-soft)' : 'var(--accent-soft)' }}>
                    <AlertTriangle size={28} style={{ color: danger ? 'var(--danger)' : 'var(--accent)' }} />
                </div>
                <h3 className="confirm-title">{title}</h3>
                <p className="confirm-message">{message}</p>
                <div className="confirm-actions">
                    {!alertOnly && (
                        <button className="btn btn-secondary flex-1" onClick={onClose}>
                            Cancelar
                        </button>
                    )}
                    <button
                        className={`btn flex-1 ${danger ? 'btn-danger-solid' : 'btn-primary'}`}
                        onClick={() => { onConfirm(); onClose(); }}
                    >
                        {confirmText || 'Confirmar'}
                    </button>
                </div>
            </div>
        </div>
    );
}
