<template>
  <el-card class="fill-width">
    <BatchInfoWindow ref="batchInfoWindowRef" />
    <el-table height="60vh" :data="batchesData" @row-click="toggleWindow" style="cursor: pointer">
      <el-table-column prop="number" label="Height" width="180" />
      <el-table-column prop="hash" label="Hash" width="250">
        <template #default="scope">
          <ShortenedHash :hash="scope.row.hash" />
        </template>
      </el-table-column>
      <el-table-column prop="timestamp" label="Time" width="180">
        <template #default="scope">
          <Timestamp :unixTimestampSeconds="Number(scope.row.timestamp)" />
        </template>
      </el-table-column>
      <el-table-column
        label="No. Transactions"
        :formatter="
          function (row) {
            return row.txHashes ? row.txHashes.length : 0
          }
        "
        width="180"
      />
      <el-table-column prop="l1Proof" label="L1 block" />
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
import { computed, onMounted } from 'vue'
import { useBatchStore } from '@/stores/batchStore'
import ShortenedHash from '@/components/helper/ShortenedHash.vue'
import Timestamp from '@/components/helper/Timestamp.vue'
import BatchInfoWindow from '@/components/helper/BatchInfoWindow.vue'

export default {
  name: 'BatchesDataGrid',
  components: { BatchInfoWindow, Timestamp, ShortenedHash },
  setup() {
    const store = useBatchStore()

    // Reload batch data onMount
    onMounted(() => {
      store.fetch()
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
      currentPage: 0,
      isWindowVisible: false
    }
  },
  methods: {
    // Called when the page size is changed
    handleSizeChange(newSize) {
      const store = useBatchStore()
      store.size = newSize
      store.offset = (this.currentPage - 1) * store.size
      // reload data
      store.fetch()
    },
    // Called when the current page is changed
    handleCurrentChange(newPage) {
      const store = useBatchStore()
      this.currentPage = newPage
      store.offset = (newPage - 1) * store.size
      // reload data
      store.fetch()
    },
    toggleWindow(data) {
      this.$refs.batchInfoWindowRef.displayData(data.hash)
    }
  }
}
</script>

<style scoped>
.fill-width {
  width: 100%;
}
</style>
