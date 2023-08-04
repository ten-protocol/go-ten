<script>
import { defineComponent } from 'vue'
import SummaryItem from '@/components/SummaryItem.vue'
import RotatingBatchesItem from "@/components/RotatingBatchesItem.vue";
import RotatingBlocksItem from "@/components/RotatingRollupsItem.vue";
import PublicTxDataGrid from "@/components/PublicTxDataGrid.vue";

export default defineComponent({
  components: {PublicTxDataGrid, RotatingBlocksItem, RotatingBatchesItem, SummaryItem }
})
</script>
<template>
  <el-container>
    <el-main>
      <summary-item />
      <p>&nbsp;</p>
      <rotating-batches-item/>
      <p>&nbsp;</p>
     <rotating-blocks-item/>
      <p>&nbsp;</p>
      <public-tx-data-grid/>
    </el-main>
  </el-container>
</template>

<style></style>
