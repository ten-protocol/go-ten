import { useAccount } from 'wagmi';
import { useTenChainAuth } from '@/hooks/useTenChainAuth';
import { useTenToken } from '@/contexts/TenTokenContext';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from '@/components/ui/alert-dialog';

type Props = {
    isOpen: boolean;
    onChange: (isOpen: boolean) => void;
};

export default function RevokeAccount({ isOpen, onChange }: Props) {
    const { address } = useAccount();
    const { revokeAccount } = useTenChainAuth(address);
    const { token: tenToken, setToken: setTenToken } = useTenToken();

    const handleRevokeAccount = async () => {
        if (tenToken) {
            console.log('[RevokeAccount] Revoking account and clearing token');
            try {
                await revokeAccount();
                await setTenToken(''); // Clear the token from cookie
                console.log('[RevokeAccount] Account revoked and token cleared successfully');
            } catch (error) {
                console.log('[RevokeAccount] Error during revocation:', error);
            }
        }
        onChange(false);
    };

    return (
        <AlertDialog open={isOpen} onOpenChange={onChange}>
            <AlertDialogContent className="bg-neutral-900">
                <AlertDialogHeader>
                    <AlertDialogTitle>Revoke private key</AlertDialogTitle>
                    <AlertDialogDescription>
                        This action cannot be undone and will revoke your authentication key. Your
                        wallet will no longer be able to interact with dapps on TEN protocol.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction onClick={handleRevokeAccount}>Revoke</AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
}
