<template>
  <el-container>
    <el-card class="margin-top">
      <el-header>Batches</el-header>
      <batches-data-grid />
    </el-card>
  </el-container>
</template>

<script>
import BatchesDataGrid from "@/components/BatchDataGrid.vue";

export default {
name: "BatchesView",
  components: {BatchesDataGrid}
}
</script>

<style scoped>
.margin-top {
  margin-top: 5em;
}
</style>