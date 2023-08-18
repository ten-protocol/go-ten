<template>
  <el-card class="fill-width">
    <el-table height="250" :data="batchesData">
      <el-table-column prop="sequencerOrderNo" label="Height" width="180"/>
      <el-table-column prop="hash" label="Hash" width="250">
        <template #default="scope">
          <ShortenedHash :hash="scope.row.hash" />
        </template>
      </el-table-column>
      <el-table-column prop="timestamp" label="Time"  width="180">
        <template #default="scope">
          <Timestamp :unixTimestampSeconds="Number(scope.row.timestamp)" />
        </template>
      </el-table-column>
      <el-table-column prop="l1Proof" label="Included in Rollup"/>
    </el-table>
    <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="size"
        :page-count="totalPages"
        layout="total, sizes, prev, pager, next"
        :total="batchListingCount"
    ></el-pagination>
  </el-card>
</template>

<script>
import {computed, onMounted, onUnmounted} from 'vue'
import {useBatchStore} from "@/stores/batchStore";
import ShortenedHash from "@/components/helper/ShortenedHash.vue";
import Timestamp from "@/components/helper/Timestamp.vue";

export default {
  name: 'BatchesDataGrid',
  components: {Timestamp, ShortenedHash},
  setup() {
    const store = useBatchStore()

    // Start polling when the component is mounted
    onMounted(() => {
      store.startPolling()
    })

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      store.stopPolling()
    })

    return {
      batchesData: computed(() => store.batchListing),
      batchListingCount: computed(() => store.batchListingCount),
      size: computed(() => store.size),
      totalPages: computed(() => {
        const store = useBatchStore()
        if (!store.batchListingCount) {
          return 0
        }
        const pages = Math.ceil(store.batchListingCount / store.size)
        console.log('Recalculated page count - ' + pages)
        return pages
      }),
      currentPage: 0
    }
  },
  methods: {
    // Called when the page size is changed
    handleSizeChange(newSize) {
      const store = useBatchStore()
      store.size = newSize
      store.offset = (this.currentPage - 1) * store.size
    },
    // Called when the current page is changed
    handleCurrentChange(newPage) {
      const store = useBatchStore()
      this.currentPage = newPage
      store.offset = (newPage - 1) * store.size
    }
  },
}
</script>

<style scoped>
.fill-width {
  width: 100%;
}
</style>
