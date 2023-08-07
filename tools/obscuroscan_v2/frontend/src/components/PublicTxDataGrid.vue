<template>
  <el-card>
    <el-table height="250" :data="publicTransactionsData" style="width: 100%">
      <el-table-column prop="TransactionHash" label="Tx Hash" width="180" />
      <el-table-column prop="BatchHeight" label="BatchHeight" width="180" />
      <el-table-column prop="Finality" label="Finality" width="180" />
    </el-table>
  </el-card>
</template>

<script>

import {computed, onMounted, onUnmounted} from "vue";
import {usePublicDataStore} from "@/stores/publicTxDataStore";

export default {
  name: "PublicTxDataGrid",
    setup() {
      const publicDataStore = usePublicDataStore()

      // Start polling when the component is mounted
      onMounted(() => {
        publicDataStore.startPolling()
      })

      // Ensure to stop polling when component is destroyed or deactivated
      onUnmounted(() => {
        publicDataStore.stopPolling()
      })

      return {
        publicTransactionsData:  computed(() => publicDataStore.publicTransactionsData),
        isAnimating: false
      }
    },
  }
</script>

<style scoped>

</style>