import { supportedWallets } from '@/lib/supportedWallets';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import Image from 'next/image';

export default function SupportedWallets() {
    return (
        <div className="text-center">
            <h4 className="text-2xl">Recommended Wallets</h4>
            <p className="opacity-70">
                These popular wallets are confirmed to work smoothly with the Gateway.
            </p>

            <div className="flex justify-center gap-6 mt-6">
                {supportedWallets
                    .map((wallet) => (
                        <div key={wallet.name}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <Image
                                        src={wallet.logo}
                                        height={48}
                                        width={48}
                                        alt={wallet.name}
                                        className="w-[32px]"
                                    />
                                </TooltipTrigger>
                                <TooltipContent>
                                    <p>{wallet.name}</p>
                                </TooltipContent>
                            </Tooltip>
                        </div>
                    ))}
            </div>
        </div>
    );
}
