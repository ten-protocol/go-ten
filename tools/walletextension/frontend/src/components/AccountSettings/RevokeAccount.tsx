import { useAccount } from 'wagmi';
import { useTenChainAuth } from '@/hooks/useTenChainAuth';
import { useTokenFromCookie } from '@/hooks/useTokenFromCookie';
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
    const [tenToken, setTenTokenToCookie] = useTokenFromCookie();

    const handleRevokeAccount = async () => {
        if (tenToken) {
            revokeAccount();
            await setTenTokenToCookie(''); // Clear token from cookie
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
