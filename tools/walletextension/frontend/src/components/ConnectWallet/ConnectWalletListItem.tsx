import Image from 'next/image';
import { Wallet } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Connector } from 'wagmi';
import { supportedWallets } from '@/lib/supportedWallets';

type Props = {
    connector: Connector;
    supported: boolean;
    onClick: () => void;
};

export default function ConnectWalletListItem({ connector, supported, onClick }: Props) {
    const handleClick = () => {
        onClick();
    };
    const matchedWallet = supportedWallets.find((wallet) => wallet.name === connector.name);
    const icon = connector.icon?.trimStart().trimEnd() || matchedWallet?.logo;

    return (
        <Button
            variant="outline"
            className="w-full justify-start gap-4 h-14 relative"
            onClick={handleClick}
            disabled={!supported}
        >
            <div>
                {icon ? (
                    <Image
                        src={icon}
                        height={48}
                        width={48}
                        alt={connector.name}
                        className="w-[32px]"
                        unoptimized
                    />
                ) : (
                    <Wallet className="h-6 w-6" />
                )}
            </div>
            <div className="flex flex-col items-start flex-grow">
                <span className="font-medium">
                    {connector.name === 'Injected' ? 'Browser Wallet' : connector.name}
                </span>
                <span className="text-sm text-muted-foreground">
                    {connector.type === 'injected' ? 'Browser Extension' : connector.type}
                </span>
            </div>
        </Button>
    );
}
