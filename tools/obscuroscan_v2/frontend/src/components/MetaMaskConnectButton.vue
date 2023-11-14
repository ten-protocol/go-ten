<template>
  <el-button @click="connectMetamask" size="large" >
    <img src="@/assets/imgs/icon_metamask.png" alt="Connect with MetaMask" class="metamask-icon" />
    <div v-if="walletStore.provider">Connected!</div>
    <div v-else>Connect with MetaMask</div>
  </el-button>
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
      connectMetamask,
      walletStore
    };
  }
}
</script>

<style scoped>
.metamask-icon {
  width: 24px;       /* Set desired width */
  height: 24px;      /* Set desired height */
  object-fit: cover; /* Ensure image content is not distorted */
  margin-right: 8px; /* Optional space between the icon and the text */
}
</style>