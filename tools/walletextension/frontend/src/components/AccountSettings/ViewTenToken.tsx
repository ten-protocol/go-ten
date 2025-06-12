import { useState } from 'react';
import { Copy, Check } from 'lucide-react';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { useLocalStorage } from 'usehooks-ts';

type Props = {
    isOpen: boolean;
    onChange: (open: boolean) => void;
};

export default function ViewTenToken({ isOpen, onChange }: Props) {
    const [tenToken] = useLocalStorage('ten_token', '');
    const [isCopied, setIsCopied] = useState(false);
    const rpcUrl = `${process.env.NEXT_PUBLIC_GATEWAY_URL}/v1/?token=${tenToken}`;

    const copyAddressToClipboard = async () => {
        if (tenToken) {
            try {
                await navigator.clipboard.writeText(rpcUrl);
                setIsCopied(true);
                setTimeout(() => setIsCopied(false), 1000);
            } catch (error) {
                console.error('Failed to copy address:', error);
            }
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={onChange}>
            <DialogContent className="bg-neutral-900">
                <DialogHeader>
                    <DialogTitle>TEN RPC TOKEN</DialogTitle>
                    <DialogDescription>
                        Your TEN RPC token is shown below. Click the copy button to copy it to your
                        clipboard.
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div className="flex items-center gap-2 p-3 bg-neutral-800 rounded-md">
                        <code className="flex-1 text-sm break-all text-green-400 font-mono">
                            {rpcUrl}
                        </code>
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={copyAddressToClipboard}
                            className="shrink-0"
                            title="Copy address to clipboard"
                        >
                            {isCopied ? (
                                <Check className="h-4 w-4 text-green-500" />
                            ) : (
                                <Copy className="h-4 w-4" />
                            )}
                        </Button>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}
