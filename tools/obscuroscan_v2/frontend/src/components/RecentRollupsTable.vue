<template>
  <el-card class="fill-width obs-container-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span>Recent Rollups</span>
        <el-button class="obs-link-button" @click="$router.push('/blocks')">View all rollups</el-button>
      </div>
    </template>
    <el-table :data="displayedRollups">
      <el-table-column prop="Timestamp" label="Timestamp" width="180" />
      <el-table-column prop="RollupHash" label="Rollup Hash" width="200" />
      <el-table-column prop="LastBatchHeight" label="Up to Batch" width="200"/>
    </el-table>
  </el-card>
</template>

<script>
import {computed, onMounted, onUnmounted} from "vue";
import {useRollupStore} from "@/stores/rollupStore";
import ShortenedHash from "@/components/helper/ShortenedHash.vue";
import {usePublicDataStore} from "@/stores/publicTxDataStore";

export default {
  name: "RecentRollupsTable",
  components: {ShortenedHash},

  setup() {
    const rollupsStore = useRollupStore()

    return {
      displayedRollups:  computed(() => rollupsStore.rollups.get())
    }
  }
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>