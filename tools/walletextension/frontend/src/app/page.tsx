'use client';
import { useAccount } from 'wagmi';
import WalletConnected from '@/components/WalletConnected/WalletConnected';
import DisconnectedWallet from '@/components/DisconnectedWallet/DisconnectedWallet';
import Header from '@/components/Header/Header';
import Footer from '@/components/Footer/Footer';
import PromoApps from '@/components/PromoApps/PromoApps';
import { tenChainIDDecimal } from '@/lib/constants';
import ConnectWalletModal from '@/components/ConnectWallet/ConnectWalletModal';
import { useUiStore } from '@/stores/ui.store';
import { useShallow } from 'zustand/react/shallow';
import { useEffect, useState } from 'react';
import WalletSettingsModal from '@/components/ConnectWallet/WalletSettingsModal';
import { useLocalStorage } from 'usehooks-ts';

export default function Home() {
    const [tenToken] = useLocalStorage<string>('ten_token', '');
    const { isConnected, chainId } = useAccount();

    const [isConnectionModalOpen, isSettingsModalOpen, setConnectionModal, setSettingsModal] =
        useUiStore(
            useShallow((state) => [
                state.isConnectionModalOpen,
                state.isSettingsModalOpen,
                state.setConnectionModal,
                state.setSettingsModal,
            ])
        );
    const isOnTen = chainId === tenChainIDDecimal;
    const [isWalletReady, setIsWalletReady] = useState(false);

    useEffect(() => {
        const timer = setTimeout(() => {
            setIsWalletReady(true);
        }, 500);

        return () => clearTimeout(timer);
    }, []);

    return (
        <div className="min-h-screen flex flex-col items-center justify-center p-8 overflow-y-hidden relative w-full">
            <div className="fixed inset-0 pointer-events-none z-40 opacity-[.03] grain-overlay" />
            <Header />

            <main className="flex flex-col items-center justify-center flex-1 gap-8 mt-16">
                <div className="text-center mb-12 mt-8">
                    <h1 className="text-[3rem] font-bold -mb-1">Welcome to the TEN Gateway!</h1>
                    <h2 className="opacity-80 text-lg">
                        Your portal into the universe of encrypted Ethereum on TEN Protocol.
                    </h2>
                </div>

                {isWalletReady && isConnected && isOnTen && tenToken !== '' ? (
                    <WalletConnected />
                ) : isWalletReady ? (
                    <DisconnectedWallet />
                ) : (
                    <div className="flex justify-center items-center h-64">
                        <div className="animate-pulse text-lg opacity-70">Loading wallet...</div>
                    </div>
                )}
                <PromoApps />
            </main>

            <Footer />
            <div className="bg-after-glow" />
            <WalletSettingsModal isOpen={isSettingsModalOpen} onOpenChange={setSettingsModal} />
            <ConnectWalletModal isOpen={isConnectionModalOpen} onOpenChange={setConnectionModal} />
        </div>
    );
}
