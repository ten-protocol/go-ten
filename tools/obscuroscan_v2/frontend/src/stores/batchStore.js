import { defineStore } from 'pinia';
import Config from "@/lib/config";
import CachedList from "@/lib/cachedList";
import Poller from "@/lib/poller";

export const useBatchStore = defineStore({
    id: 'batchStore',
    state: () => ({
        latestBatch: null,
        latestL1Proof: null,
        batches: new CachedList(),
        poller: new Poller(() => {
            const store = useBatchStore();
            store.fetch();
        }, Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            try {
                let response = await fetch( Config.backendServerAddress+'/items/batch/latest/');
                let data = await response.json();
                this.latestBatch = data.item.number;
                this.latestL1Proof = data.item.l1Proof;

                this.batches.add(data.item);
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
