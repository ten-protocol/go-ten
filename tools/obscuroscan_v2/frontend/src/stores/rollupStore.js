import { defineStore } from 'pinia';
import Config from "@/lib/config";
import CachedList from "@/lib/cachedList";


export const useRollupStore = defineStore({
    id: 'rollupStore',
    state: () => ({
        rollups: new CachedList(),
        loading: false,
        pollingInterval: Config.pollingInterval,
        timer: null,
    }),
    actions: {
        async fetchCount() {
            this.loading = true;
            try {
                let response = await fetch( Config.backendServerAddress+'/items/rollup/latest/');
                let data = await response.json();
                this.rollups.add(data.item);

                console.log("Fetched "+data.item.L1ProofNumber);
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
