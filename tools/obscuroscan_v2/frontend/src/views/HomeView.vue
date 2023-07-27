<script>
import { defineComponent } from 'vue'
import SummaryItem from '@/components/SummaryItem.vue'
import RotatingBatchesItem from "@/components/RotatingBatchesItem.vue";
import RotatingBlocksItem from "@/components/RotatingRollupsItem.vue";

export default defineComponent({
  components: {RotatingBlocksItem, RotatingBatchesItem, SummaryItem }
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
      <el-card>
        <el-table height="250" style="width: 100%">
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
