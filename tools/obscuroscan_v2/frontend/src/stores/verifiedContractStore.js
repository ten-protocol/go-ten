import {defineStore} from 'pinia';
import Config from "@/lib/config";

export const useVerifiedContractStore = defineStore({
    id: 'VerifiedContractStore',
    state: () => ({
        contracts: null,
        sequencerData: null
    }),
    actions: {
        async update() {
            try {
                const response = await fetch( Config.backendServerAddress+`/info/obscuro/`);

                const data  = await response.json();
                this.contracts =
                    [
                        {
                            "name": "Management Contract",
                            "address": data.item.ManagementContractAddress,
                            "confirmed": true
                        },
                        {
                            "name": "Message Bus Contract",
                            "address": data.item.MessageBusAddress,
                            "confirmed": true
                        }
                ];
                this.sequencerData = [
                    {
                        "name": "Sequencer ID",
                        "address": data.item.SequencerID,
                        "confirmed": true,
                    },
                    {
                        "name": "L1 Start Hash",
                        "address": data.item.L1StartHash,
                        "confirmed": true,
                    }
                ]
            } catch (error) {
                console.error("Failed to fetch item:", error);
            }
        },
    },
});
