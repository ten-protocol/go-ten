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
        <div className="relative mx-6">
            <div className="relative w-full mb-4 z-10 shape-three">
                <div className="animate-scan-overlay w-full h-full absolute z-2 pointer-events-none" />
                <Image src={imageUrl} height={800} width={800} alt={title} className="w-full" />
            </div>

            <div className="">
                <h3 className=" text-2xl font-bold tracking-tight">{title}</h3>
                <div className="text-sm opacity-80">{description}</div>
            </div>

            <Button variant="secondary" size="sm" className="mt-2" asChild>
                <Link href={url}>
                    Visit
                    <RiExternalLinkLine />
                </Link>
            </Button>
        </div>
    );
}
