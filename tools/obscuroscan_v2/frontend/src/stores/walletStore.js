import { defineStore } from 'pinia';

export const useWalletStore = defineStore({
    id: 'wallet',
    state: () => ({
        provider: null,
        address: null
    }),
    actions: {
        setProvider(provider) {
            this.provider = provider;
        },
        setAddress(address) {
            this.address = address;
        }
    }
});