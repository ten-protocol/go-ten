'use client';

import { Button } from '@/components/ui/button';
import { useAccount, useBalance } from 'wagmi';
import { Loader2 } from 'lucide-react';
import { tenChainIDDecimal, tenNetworkName } from '@/lib/constants';
import { useUiStore } from '@/stores/ui.store';
import { shallow } from 'zustand/shallow';
import { TbPrompt } from 'react-icons/tb';
import Image from 'next/image';
import { useLocalStorage } from 'usehooks-ts';

type Props = {
    className?: string;
    onChainChange?: (chainId: number, isCorrect: boolean) => void;
};

export default function ConnectWalletButton({ className }: Props) {
    const { address, isConnected, chain, connector } = useAccount();
    const [setConnectionModal, setSettingsModal] = useUiStore(
        (state) => [state.setConnectionModal, state.setSettingsModal],
        shallow
    );
    const [tenToken] = useLocalStorage('ten_token', '');

    const isWrongChain = !chain || Number(chain.id) !== Number(tenChainIDDecimal);

    const { data: ethBalance, isLoading: isLoadingEthBalance } = useBalance({
        address,
        chainId: tenChainIDDecimal,
        query: {
            enabled: isConnected && !isWrongChain,
        },
    });

    return (
        <div className={`w-full max-w-md flex justify-center ${className || ''}`}>
            <div className="w-full flex justify-center">
                {!isConnected || !tenToken ? (
                    <>
                        <Button
                            onClick={() => setConnectionModal(true)}
                            className="bg-primary hover:bg-primary/90"
                        >
                            <TbPrompt className="text-2xl size-" />
                            <span className="hidden md:block">Connect to {tenNetworkName}</span>
                            <span className="md:hidden">Connect</span>
                        </Button>
                    </>
                ) : (
                    <>
                        {isWrongChain ? (
                            <div className="flex gap-2 items-start">
                                <Button
                                    size="sm"
                                    className="bg-destructive hover:bg-destructive/90"
                                    onClick={() => setSettingsModal(true)}
                                >
                                    SWITCH to {tenNetworkName}
                                </Button>
                            </div>
                        ) : (
                            <Button
                                className="bg-primary hover:bg-primary/90 flex items-center justify-center lg:py-2  gap-x-2"
                                onClick={() => setSettingsModal(true)}
                            >
                                {connector?.icon && (
                                    <Image
                                        src={connector.icon.trimStart().trimEnd()}
                                        height={32}
                                        width={32}
                                        alt={connector.name}
                                        className="w-[24px]"
                                        unoptimized
                                    />
                                )}
                                <div className="flex flex-col">
                                    <span className="text-xs">
                                        {address?.slice(0, 6)}...{address?.slice(-4)}
                                    </span>
                                    <div className="hidden md:flex gap-2 items-center text-xs text-primary-foreground/80">
                                        {isLoadingEthBalance ? (
                                            <Loader2 className="h-3 w-3 animate-spin" />
                                        ) : (
                                            <span>
                                                {ethBalance?.formatted?.slice(0, 6) || '0'}{' '}
                                                {ethBalance?.symbol}
                                            </span>
                                        )}
                                        {/*|*/}
                                        {/*{isLoadingZenBalance ? (*/}
                                        {/*    <Loader2 className="h-3 w-3 animate-spin" />*/}
                                        {/*) : (*/}
                                        {/*    <span>*/}
                                        {/*        {zenBalance &&*/}
                                        {/*            numeral(*/}
                                        {/*                parseFloat(formatEther(zenBalance.value))*/}
                                        {/*            ).format('0.[00]')}{' '}*/}
                                        {/*        ZEN*/}
                                        {/*    </span>*/}
                                        {/*)}*/}
                                    </div>
                                </div>
                            </Button>
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
