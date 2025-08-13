'use client';

import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { AlertCircle, Copy, Check } from 'lucide-react';
import { Connector } from 'wagmi';
import { supportedWallets } from '@/lib/supportedWallets';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import useConnectToTenChain from '@/hooks/useConnectToTenChain';
import ConnectWalletListItem from '@/components/ConnectWallet/ConnectWalletListItem';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import SupportedWallets from '@/components/SupportedWallets/SupportedWallets';
import ConnectWalletProgressAnimation from '@/components/ConnectWallet/ConnectWalletProgressAnimation';
import EncryptedTextAnimation from '@/components/EncryptedTextAnimation/EncryptedTextAnimation';
import { motion } from 'framer-motion';
import Link from 'next/link';
import { RiExternalLinkLine } from 'react-icons/ri';
import { useMemo, useState } from 'react';
import { useTokenFromCookie } from '@/hooks/useTokenFromCookie';

type Props = {
    isOpen: boolean;
    onOpenChange: (open: boolean) => void;
};

export default function ConnectWalletModal({ isOpen, onOpenChange }: Props) {
    const { connectors, connectToTen, step, reset, error } = useConnectToTenChain();
    const [tenToken] = useTokenFromCookie();
    const [isCopied, setIsCopied] = useState(false);
    const rpcUrl = `${process.env.NEXT_PUBLIC_GATEWAY_URL}/v1/?token=${tenToken}`;

    const handleChange = () => {
        reset();
        onOpenChange(false);
    };

    const stepMessages = [
        'Choose a wallet to connect to TEN Protocol',
        'Confirm connection',
        'Add TEN to wallet',
        'Sign signature & authenticate account',
        'Onboarding completed',
    ];

    const usableWallets: Connector[] = [];
    const unsupportedWallets: Connector[] = [];

    connectors.forEach((connector) => {
        const matchingSupportedWallet = supportedWallets.find(
            (wallet) => wallet.name === connector.name
        );

        if (matchingSupportedWallet) {
            usableWallets.push(connector);
        } else {
            unsupportedWallets.push(connector);
        }
    });

    const unsupportedWalletList = useMemo(
        () =>
            unsupportedWallets
                .filter((connector) => connector.icon)
                .map((connector) => {
                    return (
                        <ConnectWalletListItem
                            key={connector.id}
                            onClick={() => {}}
                            connector={connector}
                            supported={false}
                        />
                    );
                }),
        [connectors]
    );

    const usableWalletList = useMemo(
        () =>
            usableWallets.map((connector) => {
                return (
                    <ConnectWalletListItem
                        key={connector.id}
                        onClick={() => connectToTen(connector)}
                        connector={connector}
                        supported={true}
                    />
                );
            }),
        [connectors]
    );

    const WalletList = () => {
        return (
            <div className="grid gap-4 py-4">
                {usableWalletList.length === 0 ? (
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-center">
                                You have no wallets compatible with TEN Protocol.
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <SupportedWallets />
                        </CardContent>
                    </Card>
                ) : (
                    <div className="flex flex-col gap-4">
                        {usableWalletList.length > 0 && <h4>Your Supported Wallets</h4>}
                        {usableWalletList}
                    </div>
                )}

                {unsupportedWalletList.length > 0 && <h4>Your Unsupported Wallets</h4>}
                {unsupportedWalletList}
            </div>
        );
    };

    const copyToClipboard = () => {
        navigator.clipboard.writeText(rpcUrl);
        setIsCopied(true);
        setTimeout(() => setIsCopied(false), 2000);
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Connect Wallet</DialogTitle>
                    <DialogDescription>{stepMessages[step]}</DialogDescription>
                </DialogHeader>
                {step === 0 && <WalletList />}
                {step >= 1 && (
                    <div>
                        <div className="py-4 px-40">
                            <motion.div
                                initial={{ opacity: 0, scale: 0.6 }}
                                animate={{ opacity: 1, scale: 1 }}
                                transition={{ duration: 0.6, ease: 'easeOut' }}
                            >
                                <ConnectWalletProgressAnimation
                                    progress={(100 / 4) * step}
                                    error={!!error}
                                />
                            </motion.div>
                        </div>
                        <h3 className="text-center mb-4">
                            <EncryptedTextAnimation
                                text={stepMessages[step].toUpperCase()}
                                hover={false}
                            />
                        </h3>
                    </div>
                )}
                {step >= 4 && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ duration: 0.6, delay: 1, ease: 'easeOut' }}
                        className="w-full mb-2"
                    >
                        <div className="">
                            <p className="text-center text-sm text-red-300 mb-2">
                                Keep this token safe and private - do not share it with anyone
                            </p>
                            <div className="flex items-center justify-center gap-2 bg-white/5 rounded p-1 mb-4">
                                <code className="text-xs">{rpcUrl}</code>
                                <Button variant="ghost" size="icon" onClick={copyToClipboard}>
                                    {isCopied ? (
                                        <Check className="h-4 w-4" />
                                    ) : (
                                        <Copy className="h-4 w-4" />
                                    )}
                                </Button>
                            </div>
                        </div>

                        <h4 className="text-center font-medium mb-2">Next Steps</h4>
                        <p className="text-center text-muted-foreground mb-4">
                            Visit the TEN Protocol faucet to get test tokens or close this dialog to
                            continue.
                        </p>
                        <div className="flex justify-center gap-4">
                            <Link href="https://faucet.ten.xyz" target="_blank">
                                <Button variant="default">
                                    VISIT FAUCET <RiExternalLinkLine />
                                </Button>
                            </Link>
                            <Button variant="outline" onClick={handleChange}>
                                CLOSE
                            </Button>
                        </div>
                    </motion.div>
                )}
                {error && (
                    <div className="flex flex-col justify-center items-center gap-4">
                        <Alert variant="destructive" className="mt-2">
                            <AlertCircle className="h-4 w-4" />
                            <AlertTitle>{error.name}</AlertTitle>
                            <AlertDescription>{error.message}</AlertDescription>
                        </Alert>{' '}
                        <Button
                            variant="outline"
                            onClick={handleChange}
                            size="sm"
                            className="mx-auto"
                        >
                            CLOSE
                        </Button>
                    </div>
                )}
            </DialogContent>
        </Dialog>
    );
}
