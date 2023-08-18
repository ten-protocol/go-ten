import {defineStore} from 'pinia';
import Config from "@/lib/config";
import Poller from "@/lib/poller";

export const usePublicDataStore = defineStore({
    id: 'publicDataStore',
    state: () => ({
        publicTransactionsData: null,
        publicTransactionsCount: null,
        offset: 0,
        size: 10,
        poller: new Poller(() => {
            const store = usePublicDataStore();
            store.fetch();
        }, Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            try {
                const response = await fetch( Config.backendServerAddress+`/items/transactions/?offset=${this.offset}&size=${this.size}`);

                const data  = await response.json();
                this.publicTransactionsData = data.result.TransactionsData;
                this.publicTransactionsCount = data.result.Total;
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
