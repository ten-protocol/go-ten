<template>
  <el-row>
    <el-col :span="4">
      <el-card class="box" shadow="always">
        <p>Ether Price</p>
        <div>
          <div v-if="ethPriceUSD">$ {{ ethPriceUSD }}</div>
          <div v-else>-</div>
        </div>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Nodes</p>
        <p>n/a</p>
      </el-card>
    </el-col>

    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always">
        <p>Latest L2 Batch</p>
        <div>
          <div v-if="latestBatch">{{ latestBatch }}</div>
          <div v-else>-</div>
        </div>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Latest L1 Rollup</p>
        <div>
          <div v-if="latestL1Proof">
            <ShortenedHash :hash="latestL1Proof" />
          </div>
          <div v-else>-</div>
        </div>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always">
        <p>Transactions</p>
        <div>
          <div v-if="totalTransactionCount">{{ totalTransactionCount }}</div>
          <div v-else>-</div>
        </div>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Contracts</p>
        <div>
          <div v-if="totalContractCount">{{ totalContractCount }}</div>
          <div v-else>-</div>
        </div>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always" style="min-height: 100%">
        <p>News from Foundation</p>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
import { useCounterStore } from '@/stores/counterStore'
import { onMounted, onUnmounted } from 'vue'
import { computed } from 'vue'
import { useBatchStore } from '@/stores/batchStore'
import { usePriceStore } from '@/stores/priceStore'
import ShortenedHash from "@/components/helper/ShortenedHash.vue";

export default {
  name: 'SummaryItem',
  components: {ShortenedHash},
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
    }
  }
}
</script>

<style scoped>
.box {
  border-radius: 15px;
}
</style>
