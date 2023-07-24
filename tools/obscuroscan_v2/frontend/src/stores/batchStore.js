import { defineStore } from 'pinia';
import Config from "@/lib/config";

export const useBatchStore = defineStore({
    id: 'batchStore',
    state: () => ({
        latestBatch: null,
        latestL1Proof: null,
        loading: false,
        pollingInterval: 5000,  // 5 seconds
        timer: null,
    }),
    actions: {
        async fetchCount() {
            this.loading = true;
            try {
                let response = await fetch( Config.backendServerAddress+'/items/batch/latest/');
                let data = await response.json();
                this.latestBatch = data.item.Number;
                this.latestL1Proof = data.item.L1Proof;

                console.log("Fetched "+this.latestBatch);
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
