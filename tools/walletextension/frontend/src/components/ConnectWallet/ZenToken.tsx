import { useAccount, useBalance } from 'wagmi';
import { Address } from 'viem';
import { Button } from '@/components/ui/button';
import { useLocalStorage } from 'usehooks-ts';
import { AlertCircle, Copy } from 'lucide-react';
import { useState } from 'react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import Link from 'next/link';
import { RiExternalLinkLine } from 'react-icons/ri';
import EthereumProvider from '@walletconnect/ethereum-provider';

type Props = {
    insufficientBalance?: boolean;
};

export default function ZenToken({ insufficientBalance }: Props) {
    const { address, connector } = useAccount();
    const [hasAddedToken, setHasAddedToken] = useLocalStorage('ADDED_ZEN_TOKEN_' + address, false);
    const [copied, setCopied] = useState(false);

    const { data: zenBalance, refetch: zenRefetch } = useBalance({
        address,
        token: process.env.NEXT_PUBLIC_ZEN_CONTRACT_ADDRESS as Address,
    });

    const handleAddToken = async () => {
        if (!connector || !address) return;

        try {
            const provider = (await connector.getProvider()) as EthereumProvider;

            const wasAdded = await provider.request({
                method: 'wallet_watchAsset',
                params: {
                    type: 'ERC20',
                    options: {
                        address: process.env.NEXT_PUBLIC_ZEN_CONTRACT_ADDRESS as string,
                        symbol: 'ZEN',
                        decimals: 18,
                    },
                },
            });

            if (wasAdded) {
                setHasAddedToken(true);
                zenRefetch();
            }
        } catch (error) {
            console.error('Failed to add token to wallet:', error);
        }
    };

    const showAddTokenCTA =
        address && zenBalance && Number(zenBalance.formatted) === 0 && !hasAddedToken;

    const shortenAddress = (address: string): string => {
        if (!address) return '';
        return `${address.slice(0, 6)}...${address.slice(-4)}`;
    };

    const copyAddressToClipboard = () => {
        if (process.env.NEXT_PUBLIC_ZEN_CONTRACT_ADDRESS) {
            navigator.clipboard.writeText(process.env.NEXT_PUBLIC_ZEN_CONTRACT_ADDRESS);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        }
    };

    if (!showAddTokenCTA && !insufficientBalance) return null;

    if (showAddTokenCTA) {
        return (
            <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Add ZEN</AlertTitle>
                <AlertDescription>
                    <p className="text-sm">
                        To see the token balance in your wallet you&#39;ll need to import the token.
                    </p>
                </AlertDescription>
                <div className="mt-2 flex flex-col items-center gap-2">
                    <div className="flex items-center justify-center gap-2">
                        <Button
                            onClick={handleAddToken}
                            size="sm"
                            variant="outline"
                            className="mt-2"
                        >
                            {process.env.NEXT_PUBLIC_ZEN_CONTRACT_ADDRESS &&
                                shortenAddress(process.env.NEXT_PUBLIC_ZEN_CONTRACT_ADDRESS)}
                        </Button>
                        <button
                            onClick={copyAddressToClipboard}
                            className="text-xs flex items-center p-1 hover:bg-muted rounded"
                            title="Copy full address to clipboard"
                            aria-label="Copy contract address"
                        >
                            <Copy size={14} />
                            {copied && <span className="ml-1 text-white/60">Copied!</span>}
                        </button>
                    </div>
                </div>
            </Alert>
        );
    }

    return (
        <Alert>
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Insufficient funds to play</AlertTitle>
            <AlertDescription>
                Visit the{' '}
                <Link
                    href="https://faucet.ten.xyz/"
                    target="_blank"
                    className="underline font-bold inline-flex gap-x-1 items-center text-white/80  hover:text-white"
                >
                    TEN Faucet
                    <RiExternalLinkLine />
                </Link>{' '}
                to top-up your wallet balance.
            </AlertDescription>
        </Alert>
    );
}
