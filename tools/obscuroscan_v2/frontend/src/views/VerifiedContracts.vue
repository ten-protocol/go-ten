<template>
  <el-container>
    <el-main>
      <el-card>
        <el-header>Personal Transactions</el-header>
       <personal-txs-grid/>
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import { defineComponent } from 'vue'
import PersonalTxsGrid from "@/components/PersonalTxsGrid.vue";

export default defineComponent({
  components: {PersonalTxsGrid}
})
</script>
<style></style>
