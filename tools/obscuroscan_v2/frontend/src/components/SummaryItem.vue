<template>
  <el-row>
    <el-col :span="4">
      <el-card class="box" shadow="always">
        <p>Ether Price</p>
        <p>$123</p>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Nodes</p>
        <p>1000</p>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always">
        <p>Latest L2 Batch</p>
        <p>0x123412312423</p>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Latest L1 Rollup</p>
        <p>0x123412312423</p>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always">
        <p>Transactions</p>
        <p>bazilion</p>
      </el-card>
      <p>&nbsp;</p>

      <el-card class="box" shadow="always">
        <p>Contracts</p>
        <p>bazilion</p>
      </el-card>
    </el-col>
    <el-col :span="4" :offset="2">
      <el-card class="box" shadow="always" style="min-height: 100%">
        <p>News from Foundation</p>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
export default {
  name: "SummaryItem"
}
</script>

<style scoped>
.box {
  border-radius: 15px;
}
</style>