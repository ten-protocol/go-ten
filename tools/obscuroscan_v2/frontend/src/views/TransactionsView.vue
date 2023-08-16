<template>
  <el-container>
    <el-card class="margin-top">
      <el-header>Transactions</el-header>
      <public-tx-data-grid />
    </el-card>
  </el-container>
</template>

<script>
import PublicTxDataGrid from '@/components/PublicTxDataGrid.vue'

export default {
  name: 'TransactionsView',
  components: { PublicTxDataGrid }
}
</script>

<style scoped>
.margin-top {
  margin-top: 5em;
}
</style>
