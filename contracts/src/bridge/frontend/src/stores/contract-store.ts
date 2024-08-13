import { create } from "zustand";
import { ethers } from "ethers";

interface ContractState {
  bridgeContract: ethers.Contract | null;
  managementContract: ethers.Contract | null;
  messageBusContract: ethers.Contract | null;
  wallet: ethers.Wallet | null;
  messageBusAddress: string;
  setContractState: (state: Partial<ContractState>) => void;
}

const useContractStore = create<ContractState>((set) => ({
  bridgeContract: null,
  managementContract: null,
  messageBusContract: null,
  wallet: null,
  messageBusAddress: "",
  setContractState: (state) => set((prevState) => ({ ...prevState, ...state })),
}));

export default useContractStore;
