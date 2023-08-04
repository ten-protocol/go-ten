<template>
  <button @click="connectMetamask">Connect to Metamask</button>
</template>

<script>
import detectEthereumProvider from '@metamask/detect-provider';
import {useWalletStore} from "@/stores/walletStore";

export default {
  name: 'MetaMaskConnectButton',
  setup() {
    const walletStore = useWalletStore();

    async function connectMetamask() {
      const provider = await detectEthereumProvider();

      if (provider) {
        // From now on, this should always be true:
        // provider === window.ethereum
        startApp(provider); // initialize your app with the provider
      } else {
        console.log('Please install MetaMask!');
      }
    }

    async function startApp(provider) {
     // Request account access if needed
      const accounts = await provider.request({ method: 'eth_requestAccounts' });

      // Set provider and address in the store
      walletStore.setProvider(provider);
      walletStore.setAddress(accounts[0]);
    }

    return {
      connectMetamask
    };
  }
}
</script>

<style scoped>
/* Add your styles here */
</style>