<script>
import { defineComponent } from 'vue'
import SummaryItem from '@/components/SummaryItem.vue'
import SummaryBatchesListItem from '@/components/SummaryBatchesListItem.vue'
import SummaryRollupsListItem from '@/components/SummaryRollupsListItem.vue'

export default defineComponent({
  components: { SummaryRollupsListItem, SummaryBatchesListItem, SummaryItem }
})
</script>
<template>
  <el-container>
    <el-main>
      <summary-item />
      <p>&nbsp;</p>
      <summary-batches-list-item />
      <p>&nbsp;</p>
      <summary-rollups-list-item />
      <p>&nbsp;</p>
      <el-card>
        <el-table :data="tableData" height="250" style="width: 100%">
          <el-table-column prop="date" label="Tx Hash" width="180" />
          <el-table-column prop="name" label="Result" width="180" />
          <el-table-column prop="name" label="Fee" width="180" />
          <el-table-column prop="name" label="Batch" width="180" />
          <el-table-column prop="name" label="Time" width="180" />
        </el-table>
      </el-card>
    </el-main>
  </el-container>
</template>

<style></style>
