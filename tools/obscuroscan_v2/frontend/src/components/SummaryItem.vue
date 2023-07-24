<template>
    <button @click="fetchCount">Reload Count</button>
  <el-row>
    <el-col :span="4">
      <el-card class="box" shadow="always">
        <p>Ether Price</p>
        <p>$123</p>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Nodes</p>
        <p>1000</p>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always">
        <p>Latest L2 Batch</p>
        <p>0x123412312423</p>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Latest L1 Rollup</p>
        <p>0x123412312423</p>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always">
        <p>Transactions</p>
        <div>
          <div>{{ totalTransactionCount }}</div>
        </div>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Contracts</p>
        <div>
          <div>{{ totalContractCount }}</div>
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
import { useCounterStore } from "@/stores/counterStore";
import { onMounted, onUnmounted } from 'vue';
import { computed } from 'vue'


export default {
  name: 'SummaryItem',
  setup() {
    const counter = useCounterStore()

    // Start polling when the component is mounted
    onMounted(() => {
      counter.startPolling();
    });

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      counter.stopPolling();
    });


    console.log(counter)
    return {
      totalContractCount: computed(() => counter.totalContractCount),
      totalTransactionCount: computed(() => counter.totalTransactionCount),
      loading: counter.loading,
      fetchCount: counter.fetchCount
    }
  }
}
</script>

<style scoped>
.box {
  border-radius: 15px;
}
</style>
