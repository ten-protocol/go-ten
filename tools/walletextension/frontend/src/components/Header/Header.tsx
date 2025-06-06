import Link from 'next/link';
import Image from 'next/image';
import NetworkStatus from '@/components/NetworkStatus/NetworkStatus';
import Social from '@/components/Social/Social';
import ConnectWalletButton from '@/components/ConnectWallet/ConnectWalletButton';

export default function Header() {
    return (
        <header>
            <nav className="bg-[rgba(255,255,255,.01)] backdrop-blur-3xl text-white p-4 fixed z-30 w-screen border-b border-[rgba(255,255,255,4%)] top-0 left-0">
                <div className="px-4 flex justify-between items-center relative">
                    <Link href="/" className="text-xl font-bold">
                        <Image src="/assets/logo.svg" height={42} width={140} alt="TEN Protocol" />
                    </Link>

                    <div className="flex gap-x-6">
                        <NetworkStatus />
                        <ConnectWalletButton />
                    </div>

                    <div className="absolute right-4 top-[200%]">
                        <Social />
                    </div>
                </div>
            </nav>
        </header>
    );
}
