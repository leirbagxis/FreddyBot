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
            <DialogContent className="sm:max-w-md p-6 text-center" showCloseButton={false}>
                <DialogHeader className="items-center gap-3">
                    <div
                        className="flex size-[64px] items-center justify-center rounded-2xl transition-transform hover:scale-105"
                        style={{
                            background: danger ? 'var(--danger-soft)' : 'var(--accent-soft)',
                            border: `1px solid ${danger ? 'rgba(232, 62, 62, 0.25)' : 'rgba(167, 139, 250, 0.25)'}`
                        }}
                    >
                        <AlertTriangle size={30} style={{ color: danger ? 'var(--danger)' : 'var(--accent)' }} />
                    </div>
                    <DialogTitle className="text-[20px] font-bold text-foreground tracking-tight mt-1">{title}</DialogTitle>
                    <DialogDescription className="text-[14px] leading-relaxed px-2 text-muted-foreground">
                        {message}
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter className="sm:justify-center gap-3 mt-4 pt-1 -mx-0 -mb-0 bg-transparent border-t-0 p-0">
                    {!alertOnly && (
                        <DialogClose render={
                            <Button variant="secondary" className="flex-1 rounded-xl h-12 px-5 py-3 text-[15px] font-semibold transition-all shadow-sm">
                                Cancelar
                            </Button>
                        } />
                    )}
                    <Button
                        variant={danger ? "destructive" : "default"}
                        className="flex-1 rounded-xl h-12 px-5 py-3 text-[15px] font-semibold transition-all shadow-md hover:shadow-lg"
                        onClick={() => { onConfirm(); onClose(); }}
                    >
                        {confirmText || 'Confirmar'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
