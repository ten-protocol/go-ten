<template>
  <el-container>
    <el-main>
      <el-card>
        <el-header>Batches</el-header>
        <batches-data-grid />
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import BatchesDataGrid from '@/components/BatchDataGrid.vue'

export default {
  name: 'BatchesView',
  components: { BatchesDataGrid }
}
</script>

<style scoped></style>
