import { defineStore } from 'pinia';
import Config from "@/lib/config";
import Poller from "@/lib/poller";

export const usePriceStore = defineStore({
    id: 'priceStore',
    state: () => ({
        ethPriceUSD: null,
        poller: new Poller(() => {
            const store = usePriceStore();
            store.fetch();
        }, 60*Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            try {
                const response = await fetch( 'https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd');
                const data = await response.json();
                this.ethPriceUSD = data.ethereum.usd;

                console.log("Fetched "+this.ethPriceUSD);
            } catch (error) {
                console.error("Failed to fetch count:", error);
            }
        },

        startPolling() {
            this.poller.start();
        },

        stopPolling() {
            this.poller.stop();
        }
    },
});
