import {defineStore} from 'pinia';
import Config from "@/lib/config";

export const usePublicDataStore = defineStore({
    id: 'publicDataStore',
    state: () => ({
        loading: false,
        pollingInterval: Config.pollingInterval,
        timer: null,
        publicTransactionsData: null,
    }),
    actions: {
        async fetch() {
            this.loading = true;
            try {
                let response = await fetch( Config.backendServerAddress+'/items/transactions/');

                let data  = await response.json();
                this.publicTransactionsData = data.result;
            } catch (error) {
                console.error("Failed to fetch count:", error);
            } finally {
                this.loading = false;
            }
        },

        startPolling() {
            this.stopPolling(); // Ensure previous intervals are cleared
            this.timer = setInterval(async () => {
                await this.fetch();
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
