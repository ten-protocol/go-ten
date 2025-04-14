import ConnectedAccount from "@/components/ConnectedAccounts/ConnectedAccount";
import {useAccount} from "wagmi";
import {useEffect, useState} from "react";
import {Address} from "viem";

export default function ConnectedAccounts() {
    const { address, connector } = useAccount();
    const [accounts, setAccounts] = useState<Address[]>([]);

    useEffect(() => {
        const fetchAccounts = async () => {
            if (connector) {
                try {
                    const walletAccounts = await connector.getAccounts();
                    setAccounts([...walletAccounts]);
                } catch (error) {
                    console.error("Error fetching accounts:", error);
                    setAccounts(address ? [address] : []);
                }
            }
        };

        fetchAccounts();
    }, [connector, address]);


    return (
        <div>
            <div className="flex justify-between">
                <p className="pl-3">Wallet</p>
                <p className="pr-3">Authenticated</p>
            </div>
            {accounts.map((acc, index) => (
                <ConnectedAccount address={acc} key={index} active={acc === address}/>
            ))}
        </div>
    )

}