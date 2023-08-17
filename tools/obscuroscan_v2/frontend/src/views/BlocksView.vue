<template>
  <el-container>
    <el-main>
      <el-card>
        <el-header>Blocks</el-header>
        <blocks-data-grid />
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import BlocksDataGrid from '@/components/BlocksDataGrid.vue'

export default {
  name: 'BlocksView',
  components: { BlocksDataGrid }
}
</script>

<style scoped></style>
