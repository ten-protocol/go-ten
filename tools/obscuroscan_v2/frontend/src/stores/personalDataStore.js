import { defineStore } from 'pinia';
import Config from "@/lib/config";
import {useWalletStore} from "@/stores/walletStore";
import Poller from "@/lib/poller";

export const usePersonalDataStore = defineStore({
    id: 'personalDataStore',
    state: () => ({
        walletStarted: false,
        personalTransactionList: null,
        personalTransactionCount: null,
        offset: 0,
        size: 10,
        poller: new Poller(() => {
            const store = usePersonalDataStore();
            store.fetch();
        }, Config.pollingInterval)
    }),
    actions: {
        async fetch() {
            try {
                const walletStore = useWalletStore()
                if (walletStore.provider == null) {
                    return
                }

                console.log(this.address)
                const requestPayload = {
                    "address": walletStore.address,
                    "pagination": {"offset": this.offset, "size": this.size},
                }
                const personalTxData = await walletStore.provider.send('eth_getStorageAt', ["listPersonalTransactions", requestPayload, null])
                this.personalTransactionList = personalTxData.result.Receipts;
                this.personalTransactionCount = personalTxData.result.Total;
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
