<template>
  <el-card class="fill-width">
    <el-table height="250" :data="blocksListing">
      <el-table-column label="Height" width="180">
        <template #default="scope">
          <span style="margin-left: 10px">{{ Number(scope.row.blockHeader.number) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="blockHeader.hash" label="Hash" width="180">
        <template #default="scope">
          <ShortenedHash :hash="scope.row.blockHeader.hash" />
        </template>
      </el-table-column>
      <el-table-column label="Time" width="180">
        <template #default="scope">
          <Timestamp :unixTimestampSeconds="Number(scope.row.blockHeader.timestamp)" />
        </template>
      </el-table-column>
      <el-table-column label="RollupHash" width="180">
        <template #default="scope">
          <span v-if="scope.row.rollupHash !== '0x0000000000000000000000000000000000000000000000000000000000000000'" style="margin-left: 10px">
            {{ scope.row.rollupHash }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="blockHeader.parentHash" label="Parent"/>

    </el-table>
    <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="size"
        :page-count="totalPages"
        layout="total, sizes, prev, pager, next"
        :total="blocksListingCount"
    ></el-pagination>
  </el-card>
</template>

<script>
import {computed, onMounted, onUnmounted} from 'vue'
import {useBlockStore} from "@/stores/blockStore";
import ShortenedHash from "@/components/helper/ShortenedHash.vue";
import Timestamp from "@/components/helper/Timestamp.vue";

export default {
  name: 'BlocksDataGrid',
  components: {Timestamp, ShortenedHash},
  setup() {
    const store = useBlockStore()

    // Start polling when the component is mounted
    onMounted(() => {
      store.startPolling()
    })

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      store.stopPolling()
    })

    return {
      blocksListing: computed(() => store.blocksListing),
      blocksListingCount: computed(() => store.blocksListingCount),
      size: computed(() => store.size),
      totalPages: computed(() => {
        const store = useBlockStore()
        if (!store.blocksListingCount) {
          return 0
        }
        const pages = Math.ceil(store.blocksListingCount / store.size)
        console.log('Recalculated page count - ' + pages)
        return pages
      }),
      currentPage: 0
    }
  },
  methods: {
    // Called when the page size is changed
    handleSizeChange(newSize) {
      const store = useBlockStore()
      store.size = newSize
      store.offset = (this.currentPage - 1) * store.size
    },
    // Called when the current page is changed
    handleCurrentChange(newPage) {
      const store = useBlockStore()
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
