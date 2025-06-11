import PrimaryCard from '@/components/PrimaryCard/PrimaryCard';
import SupportedWallets from '@/components/SupportedWallets/SupportedWallets';
import ConnectWalletButton from '@/components/ConnectWallet/ConnectWalletButton';

export default function DisconnectedWallet() {
    return (
        <PrimaryCard>
            <div className="p-4 flex flex-col justify-center items-center mb-6">
                <h3 className="text-xl mb-4">
                    Three clicks to setup encrypted communication between MetaMask and TEN.
                </h3>
                <ol className="list-decimal list-inside space-y-2 mb-6">
                    <li>Hit Connect to TEN and start your journey</li>
                    <li>Allow MetaMask to switch networks to the TEN Testnet</li>
                    <li>Sign the Signature Request (this is not a transaction)</li>
                </ol>
                <div className="flex justify-center">
                    <ConnectWalletButton />
                </div>
            </div>
            <SupportedWallets />
        </PrimaryCard>
    );
}
