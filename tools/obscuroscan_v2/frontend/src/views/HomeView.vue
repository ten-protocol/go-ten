<script>
import { defineComponent } from 'vue'
import RotatingBatchesItem from "@/components/RotatingBatchesItem.vue";
import RecentRollupsTable from "@/components/RecentRollupsTable.vue";
import PublicTxDataGrid from "@/components/PublicTxDataGrid.vue";
import SummaryRow from "@/components/SummaryRow.vue";

export default defineComponent({
  components: {PublicTxDataGrid, RecentRollupsTable, RotatingBatchesItem, SummaryRow }
})
</script>
<template>
  <summary-row />
  <rotating-batches-item/>
  <el-row gutter="12">
    <el-col :span="12">
      <recent-rollups-table/>
    </el-col>
    <el-col :span="12">
      <public-tx-data-grid/>
    </el-col>
  </el-row>
</template>

<style>

</style>
