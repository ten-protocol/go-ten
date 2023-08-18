import { defineStore } from 'pinia';
import Config from "@/lib/config";
import CachedList from "@/lib/cachedList";
import Poller from "@/lib/poller";


export const useRollupStore = defineStore({
    id: 'rollupStore',
    state: () => ({
        rollups: new CachedList(),
        poller: new Poller(() => {
            const store = useRollupStore();
            store.fetch();
        }, Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            this.loading = true;
            try {
                const response = await fetch( Config.backendServerAddress+'/items/rollup/latest/');
                const data = await response.json();
                this.rollups.add(data.item);
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
