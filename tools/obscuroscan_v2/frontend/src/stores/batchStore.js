import { defineStore } from 'pinia';
import Config from "@/lib/config";
import CachedList from "@/lib/cachedList";
import Poller from "@/lib/poller";

export const useBatchStore = defineStore({
    id: 'batchStore',
    state: () => ({
        latestBatch: null,
        latestL1Proof: null,

        batchListing: null,
        batchListingCount: null,
        offset: 0,
        size: 10,

        batches: new CachedList(),
        poller: new Poller(() => {
            const store = useBatchStore();
            store.fetch();
        }, Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            try {
                // fetch the latest batch
                let response = await fetch( Config.backendServerAddress+'/items/batch/latest/');
                let data = await response.json();
                this.latestBatch = data.item.number;
                this.latestL1Proof = data.item.l1Proof;

                this.batches.add(data.item);

                // fetch data listing
                let responseList = await fetch( Config.backendServerAddress+`/items/batches/?offset=${this.offset}&size=${this.size}`);
                let dataList  = await responseList.json();
                this.batchListing = dataList.result.BatchData;
                this.batchListingCount = dataList.result.Total;

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
