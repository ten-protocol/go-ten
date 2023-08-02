<template>
  <div>
    <button @click="listMyTransactions">List My Transactions</button>
    <div v-if="transactions">
      <h3>Your Transactions:</h3>
      <ul>
        <li v-for="transaction in transactions" :key="transaction.hash">
          {{ transaction }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ethers } from 'ethers'

export default {
  name: 'ListMyTxsButton',
  data() {
    return {
      transactions: null
    }
  },
  methods: {
    async listMyTransactions() {
      if (typeof window.ethereum !== 'undefined') {
        const provider = new ethers.BrowserProvider(window.ethereum)

        // It will prompt user for account connections if it isnt connected
        const signer = await provider.getSigner()
        console.log('Account:', await signer.getAddress())

        // Get the user's address
        const address = await signer.getAddress()

        const derp = await provider.send('eth_getStorageAt', [null, null, null])

        console.log(derp)
        // NOTE: As previously mentioned, the method `scan_listMyTransactions` does not exist in my knowledge.
        // You may need to integrate with an external service like the Etherscan API.
        // For now, I'll use a mocked method.

        this.transactions = await this.mockedListMyTransactions(address)
      } else {
        alert('Please install MetaMask!')
      }
    },
    // Mocked method - Replace this with actual logic to fetch transactions for a given address
    async mockedListMyTransactions(address) {
      // Mocked response
      return [`Transaction 1 for ${address}`, `Transaction 2 for ${address}`]
    }
  }
}
</script>

<style scoped>
/* Add your styles here */
</style>