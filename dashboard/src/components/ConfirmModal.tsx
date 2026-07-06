import { AlertTriangle } from 'lucide-react';
import {
  Dialog, DialogContent, DialogHeader, DialogTitle,
  DialogDescription, DialogFooter, DialogClose,
} from './ui/dialog';
import { Button } from './ui/button';

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
    return (
        <Dialog open={open} onOpenChange={(open) => { if (!open) onClose(); }}>
            <DialogContent className="sm:max-w-sm text-center" showCloseButton={false}>
                <DialogHeader className="items-center gap-3">
                    <div
                        className="flex size-[60px] items-center justify-center rounded-full"
                        style={{ background: danger ? 'var(--danger-soft)' : 'var(--accent-soft)' }}
                    >
                        <AlertTriangle size={28} style={{ color: danger ? 'var(--danger)' : 'var(--accent)' }} />
                    </div>
                    <DialogTitle className="text-[18px]">{title}</DialogTitle>
                    <DialogDescription className="text-[14px] leading-relaxed px-2">
                        {message}
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter className="sm:justify-center gap-2">
                    {!alertOnly && (
                        <DialogClose render={<Button variant="secondary" className="flex-1">Cancelar</Button>} />
                    )}
                    <Button
                        variant={danger ? "destructive" : "default"}
                        className="flex-1"
                        onClick={() => { onConfirm(); onClose(); }}
                    >
                        {confirmText || 'Confirmar'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
