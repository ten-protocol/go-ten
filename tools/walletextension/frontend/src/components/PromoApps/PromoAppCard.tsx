import Link from 'next/link';
import { Button } from '@/components/ui/button';
import Image from 'next/image';
import { RiExternalLinkLine } from 'react-icons/ri';

type Props = {
    title: string;
    description: string;
    imageUrl: string;
    url: string;
};

export default function PromoAppCard({ imageUrl, title, description, url }: Props) {
    return (
        <div className="relative min-w-0 w-full overflow-hidden">
            <div className="relative w-full mb-4 z-10 shape-three aspect-video">
                <div className="animate-scan-overlay w-full h-full absolute z-2 pointer-events-none" />
                <Image src={imageUrl} height={400} width={400} alt={title} className="w-full h-full object-cover" />
            </div>

            <div className="min-w-0 w-full">
                <h3 className="text-sm sm:text-lg md:text-xl lg:text-2xl font-bold tracking-tight break-all overflow-wrap-anywhere">{title}</h3>
                <div className="text-xs sm:text-sm opacity-80 break-all overflow-wrap-anywhere">{description}</div>
            </div>

            <Button variant="secondary" size="sm" className="mt-2 text-xs sm:text-sm min-w-0 px-1 sm:px-2 md:px-3 gap-0.5 sm:gap-1 md:gap-1.5 h-6 sm:h-7 md:h-8" asChild>
                <Link href={url}>
                    Visit
                    <RiExternalLinkLine />
                </Link>
            </Button>
        </div>
    );
}
