import {useAccount} from "wagmi";
import ConnectedAccounts from "@/components/ConnectedAccounts/ConnectedAccounts";
import PrimaryCard from "@/components/PrimaryCard/PrimaryCard";
import {Button} from "@/components/ui/button";

export default function WalletConnected() {
    const { connector } = useAccount();


    return (
    <PrimaryCard>
        <div className="text-center max-w-[600px]">
            <h3 className="text-2xl font-bold">Wallet Connected</h3>
            <p className="mb-8 opacity-70">Manage the accounts you have connected to the TEN Gateway. You can revoke access to your accounts at any time and request new tokens from the TEN Discord.</p>
            <ConnectedAccounts/>

            {connector && (
                <p className="text-sm mt-1 opacity-70">Connected with: {connector.name}</p>
            )}

            <Button className="mt-4">GET TOKENS</Button>
        </div>
    </PrimaryCard>
    )
}