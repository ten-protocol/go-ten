<script>
import { defineComponent } from 'vue'
import PersonalTxsGrid from "@/components/PersonalTxsGrid.vue";

export default defineComponent({
  components: {PersonalTxsGrid}
})
</script>
<template>
  <el-container>
    <el-main>
      <el-card>
       <personal-txs-grid/>
      </el-card>
    </el-main>
  </el-container>
</template>

<style></style>
