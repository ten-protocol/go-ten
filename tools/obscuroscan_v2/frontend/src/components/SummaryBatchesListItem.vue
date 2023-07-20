<template>
  <el-card>
    <p>Batches</p>
    <el-row :gutter="12">
      <el-col :span="8">
        <el-card shadow="always">
          <el-container>
            <el-header>Height</el-header>
            <el-main>
              <p>No of Txs</p>
              <p>Time</p>
              <p>Batch Hash</p>
            </el-main>
          </el-container>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="always">
          <el-container>
            <el-header>Height</el-header>
            <el-main>
              <p>No of Txs</p>
              <p>Time</p>
              <p>Batch Hash</p>
            </el-main>
          </el-container>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="always">
          <el-container>
            <el-header>Height</el-header>
            <el-main>
              <p>No of Txs</p>
              <p>Time</p>
              <p>Batch Hash</p>
            </el-main>
          </el-container>
        </el-card>
      </el-col>
    </el-row>
  </el-card>
</template>

<script>
export default {
  name: "SummaryBatchesListItem"
}
</script>

<style scoped>

</style>