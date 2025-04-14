import Image from 'next/image';
import Link from 'next/link';


export default function Footer() {
    return (
        <footer className="relative container py-12 mx-auto px-4 mt-24">
            <div className="grid grid-cols-12 gap-4 md:grid-cols-6">
                <div className="col-span-12 md:col-span-2 flex flex-col justify-between">
                    <Link href="/" className="text-xl font-bold mb-4">
                        <Image src="/assets/logo.svg" height={42} width={140} alt="TEN Protocol" />
                    </Link>
                </div>

                <div className="col-span-12 md:col-span-4 flex items-start justify-end gap-12">

                </div>
            </div>

            <div className="flex flex-col gap-4 md:flex-row justify-between border-t border-[rgab(255,255,255,.2)] pt-4 mt-4">
                <div className="flex gap-4">
                    <Link href="/privacy-policy" className="text-sm">
                        Privacy Policy
                    </Link>
                    <Link href="/tos" className="text-sm">
                        Terms of Service
                    </Link>
                </div>
                <p className="text-sm">© 2025 TEN Protocol. All right reserved.</p>
            </div>
        </footer>
    );
}
