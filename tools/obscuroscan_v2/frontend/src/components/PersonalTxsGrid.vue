<template>
  <el-card>
  <el-table height="250" style="width: 100%" :data="personalTransactionList">
    <el-table-column prop="blockNumber" label="Batch Height" width="180" >
      <template #default="scope">
        <span style="margin-left: 10px">{{ Number(scope.row.blockNumber) }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="transactionHash" label="Tx Hash" width="250" >
    </el-table-column>
    <el-table-column prop="status" label="Status" width="180" >
      <template #default="scope">
        <span style="margin-left: 10px">{{ (Number(scope.row.status) === 1) ? "Sucess" : "Failed" }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="gasUsed" label="Gas Cost" width="180" >
      <template #default="scope">
        <span style="margin-left: 10px">{{ Number(scope.row.gasUsed) }}</span>
      </template>
    </el-table-column>

    <el-table-column prop="blockHash" label="Batch Hash"  >
    </el-table-column>
  </el-table>
  <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="currentPage"
      :page-sizes="[10, 20, 30, 40]"
      :page-size="size"
      :page-count="totalPages"
      layout="total, sizes, prev, pager, next"
      :total="personalTransactionCount"
  ></el-pagination>
  </el-card>
</template>

<script>
import {computed, onMounted, onUnmounted} from "vue";
import {usePersonalDataStore} from "@/stores/personalDataStore";

export default {
  name: "PersonalTxsGrid",
  setup() {
    const store = usePersonalDataStore()

    // Start polling when the component is mounted
    onMounted(() => {
      store.startPolling()
    })

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      store.stopPolling()
    })

    return {
      personalTransactionList:  computed(() => store.personalTransactionList),
      personalTransactionCount: computed(() => store.personalTransactionCount),
      size: computed(() => store.size),
      totalPages: computed(() => {
        const store = usePersonalDataStore()
        if (!store.personalTransactionCount) {
          return 0
        }
        const pages = Math.ceil(store.personalTransactionCount / store.size)
        console.log('Recalculated page count - ' + pages)
        return pages
      }),
      currentPage: 0
    }
  },
  methods: {
    // Called when the page size is changed
    handleSizeChange(newSize) {
      const store = usePersonalDataStore()
      store.size = newSize
      store.offset = (this.currentPage - 1) * store.size
    },
    // Called when the current page is changed
    handleCurrentChange(newPage) {
      const store = usePersonalDataStore()
      this.currentPage = newPage
      store.offset = (newPage - 1) * store.size
    }
  },
}
</script>

<style scoped></style>
