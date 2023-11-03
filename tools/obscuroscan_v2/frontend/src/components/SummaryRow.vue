<template>
  <el-row justify="space-evenly" class="stats-row">
    <SummaryItem label="Ether Price" :icon="ethIcon">
      <div v-if="ethPriceUSD">$ {{ ethPriceUSD }}</div>
    </SummaryItem>
    <el-divider direction="vertical" class="stat-divider"></el-divider>

    <SummaryItem label="Latest Batch" :icon="batchIcon">
      <div v-if="latestBatch">{{ latestBatch }}</div>
    </SummaryItem>
    <el-divider direction="vertical" class="stat-divider"></el-divider>

    <SummaryItem label="Latest Rollup" :icon="rollupIcon">
      <div v-if="latestL1Proof">
        <ShortenedHash :hash="latestL1Proof"/>
      </div>
    </SummaryItem>
    <el-divider direction="vertical" class="stat-divider"></el-divider>

    <SummaryItem label="Transactions" :icon="txIcon">
      <div v-if="totalTransactionCount">{{ totalTransactionCount }}</div>
    </SummaryItem>
    <el-divider direction="vertical" class="stat-divider"></el-divider>

    <SummaryItem label="Contracts" :icon="contractsIcon">
      <div v-if="totalContractCount">{{ totalContractCount }}</div>
    </SummaryItem>
  </el-row>
</template>

<script>
import { useCounterStore } from '@/stores/counterStore'
import { onMounted, onUnmounted } from 'vue'
import { computed } from 'vue'
import { useBatchStore } from '@/stores/batchStore'
import { usePriceStore } from '@/stores/priceStore'
import ShortenedHash from "@/components/helper/ShortenedHash.vue";
import SummaryItem from "@/components/SummaryItem.vue";

import ethIcon from "@/assets/imgs/icon_ethereum.png";
import batchIcon from "@/assets/imgs/icon_l2_batch.png";
import rollupIcon from "@/assets/imgs/icon_l1_rollup.png";
import txIcon from "@/assets/imgs/icon_transactions.png";
import contractsIcon from "@/assets/imgs/icon_contracts.png";

export default {
  name: 'SummaryRow',
  components: {SummaryItem, ShortenedHash},
  setup() {
    const counter = useCounterStore()
    const batch = useBatchStore()
    const price = usePriceStore()

    // Start polling when the component is mounted
    onMounted(() => {
      counter.startPolling()
      batch.startPolling()
      price.startPolling()
    })

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      counter.stopPolling()
      batch.stopPolling()
      price.stopPolling()
    })

    return {
      totalContractCount: computed(() => counter.totalContractCount),
      totalTransactionCount: computed(() => counter.totalTransactionCount),
      latestBatch: computed(() => batch.latestBatch),
      latestL1Proof: computed(() => batch.latestL1Proof),
      ethPriceUSD: computed(() => price.ethPriceUSD),
      ethIcon,
      batchIcon,
      rollupIcon,
      txIcon,
      contractsIcon
    }
  }
}
</script>

<style scoped>
.stats-row {
  display: flex;
  flex-direction: row;
  border-radius: 20px;
  background-color: white;
}

.stat-divider {
  height: 80px;
  margin-top: 10px;
  margin-bottom: 10px;
}
</style>
