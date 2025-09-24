import PromoAppCard from '@/components/PromoApps/PromoAppCard';
import {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious,
    CarouselDots,
} from '@/components/ui/carousel';
import {PROMO_APPS} from "@/lib/constants";
import Autoplay from "embla-carousel-autoplay";
import { useRef } from 'react';


export default function PromoApps() {
    const plugin = useRef(
        Autoplay({ 
            delay: 2000, 
            stopOnInteraction: true,
            stopOnMouseEnter: false // We'll handle this manually
        })
    );
    
    return (
        <div className="container flex flex-col items-center mt-8">
            <h3 className="text-3xl">Explore some of our apps.</h3>
            <p className="opacity-80 text-center">
                Discover exciting decentralized apps built on TEN. Ready for you to dive in and
                explore!
            </p>
            <div className="w-full mt-8">
                <Carousel
                    opts={{
                        align: 'start',
                        loop: true,
                    }}
                    plugins={[plugin.current]}
                    className="w-full"
                    onMouseEnter={() => plugin.current.stop()}
                    onMouseLeave={() => plugin.current.play()}
                >
                    <CarouselContent className="-ml-1 sm:-ml-2 md:-ml-4 min-w-0">
                        {PROMO_APPS.map((app, index) => (
                            <CarouselItem key={index} className="pl-1 sm:pl-4 md:pl-8 basis-[85%] sm:basis-1/2 lg:basis-1/2 xl:basis-1/3 min-w-0">
                                <PromoAppCard
                                    title={app.title}
                                    description={app.description}
                                    imageUrl={app.imageUrl}
                                    url={app.url}
                                />
                            </CarouselItem>
                        ))}
                    </CarouselContent>
                    <div className="hidden md:block w-[calc(100%-64px)] ml-[42px] h-16 absolute top-1/3">
                        <CarouselPrevious />
                        <CarouselNext />
                    </div>

                    <CarouselDots />
                </Carousel>
            </div>
        </div>
    );
}


//md:pl-4 basis-[100%] sm:basis-[80%] md:basis-[60%] lg:basis-[45%] xl:basis-[35%] 2xl:basis-[28%]