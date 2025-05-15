import PromoAppCard from '@/components/PromoApps/PromoAppCard';

export default function PromoApps() {
    return (
        <div className="container mx-auto flex flex-col items-center mt-8">
            <h3 className="text-3xl">Explore some of our apps.</h3>
            <p className="opacity-80 text-center">
                Discover exciting decentralized apps built on TEN. Ready for you to dive in and
                explore!
            </p>
            <div className="grid grid-cols-1 md:grid-cols-3 space-4 space-y-6 mt-8">
                <PromoAppCard
                    title="House of TEN"
                    description="An Onchain poker tournament played by frontier AI models."
                    imageUrl="/assets/promo/houseOfTen.png"
                    url="https://houseof.ten.xyz"
                />
                <PromoAppCard
                    title="Battleships"
                    description="Sink ships, win ZEN! Play a new vartiation of Battleships."
                    imageUrl="/assets/promo/bs.png"
                    url="https://battleships.ten.xyz"
                />
                <PromoAppCard
                    title="TENZEN"
                    description="Play to hit zero! Every extra zero boosts your prize!"
                    imageUrl="/assets/promo/tenzen.png"
                    url="https://tenzen.ten.xyz"
                />
            </div>
        </div>
    );
}
