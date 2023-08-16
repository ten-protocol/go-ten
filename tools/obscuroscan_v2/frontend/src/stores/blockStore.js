import { defineStore } from 'pinia';
import Config from "@/lib/config";
import Poller from "@/lib/poller";


export const useBlockStore = defineStore({
    id: 'blockStore',
    state: () => ({
        blocksListing: null,
        blocksListingCount: null,
        offset: 0,
        size: 10,

        poller: new Poller(() => {
            const store = useBlockStore();
            store.fetch();
        }, Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            this.loading = true;
            try {
                // fetch data listing
                const responseList = await fetch( Config.backendServerAddress+`/items/blocks/?offset=${this.offset}&size=${this.size}`);
                const dataList  = await responseList.json();
                this.blocksListing = dataList.result.BlocksData;
                this.blocksListingCount = dataList.result.Total;
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
