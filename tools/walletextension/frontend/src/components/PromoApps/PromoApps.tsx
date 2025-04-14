import PromoAppCard from "@/components/PromoApps/PromoAppCard";

export default function PromoApps() {
    return (
        <div className="container mx-auto flex flex-col items-center mt-8">
            <h3 className="text-3xl">Explore some of our apps.</h3>
            <p className="opacity-80">loremipsum lerm ispum .</p>
            <div className="flex space-x-4 mt-4">
                <PromoAppCard title="House of TEN" description="AI Poker on-chain" imageUrl="/assets/promo/houseOfTen.png" url="https://houseof.ten.xyz" />
                <PromoAppCard title="House of TEN" description="AI Poker on-chain" imageUrl="/assets/promo/houseOfTen.png" url="https://houseof.ten.xyz" />
                <PromoAppCard title="House of TEN" description="AI Poker on-chain" imageUrl="/assets/promo/houseOfTen.png" url="https://houseof.ten.xyz" />
                <PromoAppCard title="House of TEN" description="AI Poker on-chain" imageUrl="/assets/promo/houseOfTen.png" url="https://houseof.ten.xyz" />
            </div>
        </div>
    )
}