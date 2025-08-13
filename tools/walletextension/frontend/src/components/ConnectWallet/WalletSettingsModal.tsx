'use client';

import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Loader2, AlertCircle, Wallet } from 'lucide-react';
import { toast } from 'sonner';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { useAccount, useDisconnect, useSwitchChain, useBalance } from 'wagmi';
import { Address, formatEther } from 'viem';
import { Skeleton } from '@/components/ui/skeleton';
import numeral from 'numeral';
import {
    nativeCurrency,
    tenChainIDDecimal,
    tenGatewayAddress,
    tenNetworkName,
} from '@/lib/constants';
import Image from 'next/image';
import { useTokenFromCookie } from '@/hooks/useTokenFromCookie';
import { joinTestnet } from '@/api/gateway';
import { useState } from 'react';

type Props = {
    isOpen: boolean;
    onOpenChange: (open: boolean) => void;
};

export default function WalletSettingsModal({ isOpen, onOpenChange }: Props) {
    const { address, chain, connector } = useAccount();
    const { disconnect } = useDisconnect();
    const {
        switchChain,
        isPending: isSwitchingChain,
        error: switchChainError,
        reset: switchChainReset,
    } = useSwitchChain();
    const [missingKeyError, setMissingKeyError] = useState(false);
    const [tenToken, setTenTokenToCookie] = useTokenFromCookie();

    const { data: ethBalance, isLoading: isLoadingEthBalance } = useBalance({
        address,
        chainId: tenChainIDDecimal,
    });

    const { data: zenBalance, isLoading: isLoadingZenBalance } = useBalance({
        address,
        chainId: tenChainIDDecimal,
        token: process.env.NEXT_PUBLIC_ZEN_CONTRACT_ADDRESS as Address,
    });

    const isWrongChain = !chain || Number(chain.id) !== Number(tenChainIDDecimal);

    const handleSwitchChain = async () => {
        console.log('🔄 handleSwitchChain: Starting switch chain process');
        console.log('📝 handleSwitchChain: Current tenToken =', tenToken);
        console.log('📝 handleSwitchChain: Token length =', tenToken?.length);
        
        try {
            // Always get a fresh token when adding/switching to TEN network
            console.log('🎯 handleSwitchChain: Calling joinTestnet() to get fresh token...');
            const freshToken = await joinTestnet();
            console.log('🎯 handleSwitchChain: Received fresh token =', freshToken);
            
            if (!freshToken) {
                console.log('❌ handleSwitchChain: No fresh token received');
                setMissingKeyError(true);
                return;
            }
            
            // Store the fresh token in cookie
            console.log('🍪 handleSwitchChain: Storing fresh token in cookie...');
            await setTenTokenToCookie(freshToken);
            console.log('🍪 handleSwitchChain: Fresh token stored successfully');
            
            const rpcUrl = `${tenGatewayAddress}/v1/?token=${freshToken}`;
            console.log('🌐 handleSwitchChain: About to switch chain with fresh RPC URL =', rpcUrl);
            console.log('🌐 handleSwitchChain: Chain ID =', tenChainIDDecimal);
            
            switchChain({
                chainId: tenChainIDDecimal,
                addEthereumChainParameter: {
                    rpcUrls: [rpcUrl],
                    chainName: tenNetworkName,
                    nativeCurrency: nativeCurrency,
                },
            });
            
            console.log('✅ handleSwitchChain: switchChain call completed');
            
        } catch (error) {
            console.error('❌ handleSwitchChain: Error occurred:', error);
            console.error(
                'Failed to switch to TEN Protocol. Please make sure you have added TEN Protocol to your wallet.',
                error
            );
            toast.error('Failed to switch to TEN Protocol.', {
                duration: 5000,
            });
        }
    };

    const handleOpenChange = () => {
        setMissingKeyError(false);
        switchChainReset();
        onOpenChange(!isOpen);
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Wallet Settings</DialogTitle>
                    <DialogDescription>
                        Manage your wallet connection and network settings
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div className="space-y-2">
                        <h4 className="font-medium">Current Network</h4>
                        <p className="text-sm text-muted-foreground">
                            {chain?.name || 'Unknown Network'} (ID: {chain?.id})
                        </p>

                        {isWrongChain && (
                            <div className="space-y-2">
                                <p className="text-sm text-destructive">
                                    You are on the wrong network. Please switch to TEN Protocol.
                                </p>
                                {missingKeyError && (
                                    <Alert variant="destructive" className="mt-2">
                                        <AlertCircle className="h-4 w-4" />
                                        <AlertTitle>Failed to switch network.</AlertTitle>
                                        <AlertDescription>
                                            TEN token not found. Close modal and try again.
                                        </AlertDescription>
                                    </Alert>
                                )}
                                {switchChainError && (
                                    <Alert variant="destructive" className="mt-2">
                                        <AlertCircle className="h-4 w-4" />
                                        <AlertTitle>Failed to switch network.</AlertTitle>
                                        <AlertDescription>
                                            Please make sure you have added TEN Protocol to your
                                            wallet. Visit the to get onboarded onto the network.
                                        </AlertDescription>
                                    </Alert>
                                )}
                                {switchChainError && <div>SWITCH ERROR</div>}
                                <Button
                                    onClick={handleSwitchChain}
                                    className="w-full"
                                    disabled={isSwitchingChain}
                                >
                                    {isSwitchingChain ? (
                                        <>
                                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                            Switching...
                                        </>
                                    ) : (
                                        'Switch to TEN Protocol'
                                    )}
                                </Button>
                            </div>
                        )}
                    </div>

                    <div className="space-y-2">
                        <h4 className="font-medium">Connected Account</h4>
                        <div className="flex gap-2">
                            {connector?.icon ? (
                                <Image
                                    src={connector.icon.trimStart().trimEnd()}
                                    height={32}
                                    width={32}
                                    alt={connector.name}
                                    className="w-[24px]"
                                    unoptimized
                                />
                            ) : (
                                <Wallet className="h-6 w-6" />
                            )}
                            <p className="text-sm text-muted-foreground break-all">{address}</p>
                        </div>
                    </div>

                    <div className="space-y-3">
                        <h4 className="font-medium">Balances</h4>
                        <div className="space-y-2">
                            <div className="flex items-center">
                                {isLoadingEthBalance ? (
                                    <Skeleton className="w-24 h-5" />
                                ) : (
                                    <span className="font-medium">
                                        {ethBalance?.formatted || '0'} {ethBalance?.symbol || 'ETH'}
                                    </span>
                                )}
                            </div>

                            <div className="flex items-center">
                                {isLoadingZenBalance ? (
                                    <Skeleton className="w-24 h-5" />
                                ) : (
                                    <span className="font-medium">
                                        {zenBalance
                                            ? numeral(
                                                  parseFloat(formatEther(zenBalance.value) ?? 0)
                                              ).format('0.[00]')
                                            : 0}{' '}
                                        ZEN
                                    </span>
                                )}
                            </div>

                            {/*{zenBalance && parseFloat(formatEther(zenBalance.value)) === 0 && (*/}
                            {/*    <ZenToken insufficientBalance={true} />*/}
                            {/*)}*/}
                        </div>
                    </div>

                    <div className="pt-4">
                        <Button
                            onClick={() => {
                                disconnect();
                                onOpenChange(false);
                            }}
                            variant="destructive"
                            className="w-full"
                        >
                            Disconnect Wallet
                        </Button>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}
