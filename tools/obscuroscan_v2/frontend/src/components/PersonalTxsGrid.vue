<template>
  <el-table height="250" style="width: 100%" :data="personalTransactionList">
    <el-table-column prop="transactionHash" label="Tx Hash" width="180" >
    </el-table-column>
    <el-table-column prop="status" label="Status" width="180" >
      <template #default="scope">
        <span style="margin-left: 10px">{{ (Number(scope.row.status) === 1) ? "Sucess" : "Failed" }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="gasUsed" label="Fee" width="180" >
      <template #default="scope">
        <span style="margin-left: 10px">{{ Number(scope.row.gasUsed) }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="blockNumber" label="Batch Height" width="180" >
      <template #default="scope">
        <span style="margin-left: 10px">{{ Number(scope.row.blockNumber) }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="blockHash" label="Batch Hash" width="180" >
    </el-table-column>
  </el-table>
</template>

<script>
import {computed, onMounted, onUnmounted} from "vue";
import {usePersonalDataStore} from "@/stores/personalDataStore";

export default {
  name: "PersonalTxsGrid",
  setup() {
    const personalDataStore = usePersonalDataStore()

    // Start polling when the component is mounted
    onMounted(() => {
      personalDataStore.startPolling()
    })

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      personalDataStore.stopPolling()
    })

    return {
      personalTransactionList:  computed(() => personalDataStore.personalTransactionList),
      isAnimating: false
    }
  },
}
</script>

<style scoped>

</style>