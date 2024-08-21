import { create } from "zustand";
import { IContractState } from "../types";

const useContractStore = create<IContractState>((set) => ({
  bridgeContract: null,
  managementContract: null,
  messageBusContract: null,
  wallet: null,
  messageBusAddress: "",
  setContractState: (state) => set((prevState) => ({ ...prevState, ...state })),
}));

export default useContractStore;
