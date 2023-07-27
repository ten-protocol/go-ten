import { defineStore } from 'pinia';
import Config from "@/lib/config";

export const usePriceStore = defineStore({
    id: 'priceStore',
    state: () => ({
        ethPriceUSD: null,
        loading: false,
        pollingInterval: 60*Config.pollingInterval,
        timer: null,
    }),
    actions: {
        async fetchCount() {
            this.loading = true;
            try {
                const response = await fetch( 'https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd');
                const data = await response.json();
                this.ethPriceUSD = data.ethereum.usd;

                console.log("Fetched "+this.ethPriceUSD);
            } catch (error) {
                console.error("Failed to fetch count:", error);
            } finally {
                this.loading = false;
            }
        },

        startPolling() {
            this.stopPolling(); // Ensure previous intervals are cleared
            this.timer = setInterval(async () => {
                await this.fetchCount();
            }, this.pollingInterval);
        },

        stopPolling() {
            if (this.timer) {
                clearInterval(this.timer);
                this.timer = null;
            }
        }

    },
});
