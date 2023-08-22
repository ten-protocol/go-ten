<template>
  <el-container>
    <el-main>
      <el-card>
        <el-header>Official Verified Data</el-header>
        <verified-contracts-data-grid />
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import { defineComponent } from 'vue'
import VerifiedContractsDataGrid from "@/components/VerifiedContractsDataGrid.vue";

export default defineComponent({
  components: {VerifiedContractsDataGrid}
})
</script>
<style></style>
