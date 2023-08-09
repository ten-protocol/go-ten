import { defineStore } from 'pinia';
import Config from "@/lib/config";
import Poller from "@/lib/poller";

export const useCounterStore = defineStore({
    id: 'counterStore',
    state: () => ({
        totalContractCount: 0,
        totalTransactionCount:0,
        poller: new Poller(() => {
            const store = useCounterStore();
            store.fetch();
        }, Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            try {
                const totContractResp = await fetch( Config.backendServerAddress+'/count/contracts/');
                const totContractData = await totContractResp.json();
                this.totalContractCount = totContractData.count;

                const totTxResp = await fetch( Config.backendServerAddress+'/count/transactions/');
                const totTxData = await totTxResp.json();
                this.totalTransactionCount = totTxData.count;
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
