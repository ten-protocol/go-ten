<template>
  <el-card shadow="never" style=" border-radius: 20px; ">
  <el-row justify="space-evenly">
    <el-col :span="4" >
      <el-card class="box" shadow="never">
        <p class="header-text">
          <img class="icon" src="@/assets/imgs/icon_ethereum.png"/>
          Ether Price
        </p>

        <div>
          <div v-if="ethPriceUSD">$ {{ ethPriceUSD }}</div>
          <div v-else v-loading=true element-loading-background="#F4F6FF">loading...</div>
        </div>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="never">
        <p class="header-text">
          <img class="icon" src="@/assets/imgs/icon_nodes.png"/>
          Nodes
        </p>
        <div>Coming soon</div>
      </el-card>
    </el-col>

    <el-col :span="4" :offset="3">
      <el-card class="box" shadow="never">
        <p class="header-text">
          <img class="icon" src="@/assets/imgs/icon_l2_batch.png"/>
          Latest L2 Batch
        </p>
        <div>
          <div v-if="latestBatch">{{ latestBatch }}</div>
          <div v-else v-loading=true element-loading-background="#F4F6FF">loading...</div>
        </div>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="never">
        <p class="header-text">
          <img class="icon" src="@/assets/imgs/icon_l1_rollup.png"/>
          Latest L1 Rollup
        </p>
        <div>
          <div v-if="latestL1Proof">
            <ShortenedHash :hash="latestL1Proof" />
          </div>
          <div v-else v-loading=true element-loading-background="#F4F6FF">loading...</div>
        </div>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="3">
      <el-card class="box" shadow="never">
        <p class="header-text">
          <img class="icon" src="@/assets/imgs/icon_transactions.png"/>
          Transactions</p>
        <div>
          <div v-if="totalTransactionCount">{{ totalTransactionCount }}</div>
          <div v-else v-loading=true element-loading-background="#F4F6FF">loading...</div>
        </div>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="never">
        <p class="header-text">
          <img class="icon" src="@/assets/imgs/icon_contracts.png"/>
          Contracts</p>
        <div>
          <div v-if="totalContractCount">{{ totalContractCount }}</div>
          <div v-else v-loading=true element-loading-background="#F4F6FF">loading...</div>
        </div>
      </el-card>
    </el-col>
  </el-row>
  </el-card>
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
  background: #F4F6FF;
}

.icon {
  height: 24px;      /* Set desired height */
  object-fit: cover; /* Ensure image content is not distorted */
  margin-right: 8px; /* Optional space between the icon and the text */
}

.header-text {
  color: #5973B8;
  font-weight: bold;
}

.el-icon-loading:before {
  color: red;
}
</style>
