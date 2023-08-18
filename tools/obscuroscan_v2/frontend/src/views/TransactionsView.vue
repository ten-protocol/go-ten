<template>
  <el-container>
    <el-main>
      <el-card>
        <el-header>Transactions</el-header>
        <public-tx-data-grid />
      </el-card>
    </el-main>
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
