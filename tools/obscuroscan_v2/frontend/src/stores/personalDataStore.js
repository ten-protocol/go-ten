import { defineStore } from 'pinia';
import Config from "@/lib/config";
import {useWalletStore} from "@/stores/walletStore";

export const usePersonalDataStore = defineStore({
    id: 'personalDataStore',
    state: () => ({
        loading: false,
        pollingInterval: Config.pollingInterval,
        timer: null,
        walletStarted: false,
        personalTransactionList: null,
    }),
    actions: {
        async fetchPersonalData() {
            this.loading = true;
            try {
                const walletStore = useWalletStore()
                if (walletStore.provider == null) {
                    return
                }

                console.log(this.address)
                const personalTxData = await walletStore.provider.send('eth_getStorageAt', [walletStore.address, null, null])
                this.personalTransactionList = personalTxData.result;
            } catch (error) {
                console.error("Failed to fetch count:", error);
            } finally {
                this.loading = false;
            }
        },

        startPolling() {
            this.stopPolling(); // Ensure previous intervals are cleared
            this.timer = setInterval(async () => {
                await this.fetchPersonalData();
            },  this.pollingInterval);
        },

        stopPolling() {
            if (this.timer) {
                clearInterval(this.timer);
                this.timer = null;
            }
        }

    },
});
