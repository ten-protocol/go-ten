import { defineStore } from 'pinia';
import Config from "@/lib/config";

export const useCounterStore = defineStore({
    id: 'counterStore',
    state: () => ({
        count: 0,
        loading: false,
        pollingInterval: 5000,  // 5 seconds
        timer: null,
    }),
    actions: {
        async fetchCount() {
            this.loading = true;
            try {
                console.log(Config.backendServerAddress)
                const response = await fetch( Config.backendServerAddress+'/count/contracts/');
                const data = await response.json();
                this.count = data.count;
                console.log("Fetched "+this.count);
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
    getters: {
        getCount: (state) => state.count,
    },
});
