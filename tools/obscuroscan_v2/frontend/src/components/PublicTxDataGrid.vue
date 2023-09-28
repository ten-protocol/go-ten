<template>
  <el-card class="fill-width">
    <el-table  height="60vh" :data="publicTransactionsData">
      <el-table-column prop="BatchHeight" label="Batch Height" width="180" />
      <el-table-column prop="Finality" label="Finality" width="180" />
      <el-table-column prop="TransactionHash" label="Tx Hash"  />
    </el-table>
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="currentPage"
      :page-sizes="[10, 20, 30, 40]"
      :page-size="size"
      :page-count="totalPages"
      layout="total, sizes, prev, pager, next"
      :total="publicTransactionsCount"
    ></el-pagination>
  </el-card>
</template>

<script>
import { computed, onMounted, onUnmounted } from 'vue'
import { usePublicDataStore } from '@/stores/publicTxDataStore'

export default {
  name: 'PublicTxDataGrid',
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
      publicTransactionsData: computed(() => publicDataStore.publicTransactionsData),
      publicTransactionsCount: computed(() => publicDataStore.publicTransactionsCount),
      size: computed(() => publicDataStore.size),
      totalPages: computed(() => {
        const publicDataStore = usePublicDataStore()
        if (!publicDataStore.publicTransactionsCount) {
          return 0
        }
        const pages = Math.ceil(publicDataStore.publicTransactionsCount / publicDataStore.size)
        console.log('Recalculated page count - ' + pages)
        return pages
      }),
      currentPage: 0
    }
  },
  methods: {
    // Called when the page size is changed
    handleSizeChange(newSize) {
      const store = usePublicDataStore()
      store.size = newSize
      store.offset = (this.currentPage - 1) * store.size
    },
    // Called when the current page is changed
    handleCurrentChange(newPage) {
      const store = usePublicDataStore()
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
