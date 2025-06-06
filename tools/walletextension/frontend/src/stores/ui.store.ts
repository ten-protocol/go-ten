import { create } from 'zustand';

type UiStore = {
    tenToken?: string
    authEvents: number
    isConnectionModalOpen: boolean;
    isSettingsModalOpen: boolean;
    setConnectionModal: (open: boolean) => void;
    setSettingsModal: (open: boolean) => void;
    setTenToken: (tenToken: string) => void;
    incrementAuthEvents: () => void;
};

export const useUiStore = create<UiStore>()((set, get) => ({
    authEvents: 0,
    isConnectionModalOpen: false,
    setConnectionModal: (open) => set({ isConnectionModalOpen: open }),
    isSettingsModalOpen: false,
    setSettingsModal: (open) => set({ isSettingsModalOpen: open }),
    setTenToken: (tenToken) => set({ tenToken }),
    incrementAuthEvents: () => {set({authEvents: get().authEvents + 1})}
}));
