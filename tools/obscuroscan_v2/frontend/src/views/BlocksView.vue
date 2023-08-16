<template>
  <el-container>
    <el-card class="margin-top">
      <el-header>Blocks</el-header>
      <blocks-data-grid />
    </el-card>
  </el-container>
</template>

<script>
import BlocksDataGrid from "@/components/BlocksDataGrid.vue";

export default {
  name: "BlocksView",
  components: {BlocksDataGrid}
}
</script>

<style scoped>
.margin-top {
  margin-top: 5em;
}
</style>