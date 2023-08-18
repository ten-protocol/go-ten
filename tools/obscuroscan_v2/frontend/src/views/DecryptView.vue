<template>
  <el-container>
    <el-main>
      <el-card>
        <el-header>
          <h2>Static Keys</h2>
        </el-header>
        <h3> Current public Key: TBD</h3>
        <h3> Current private Key: TBD</h3>
      </el-card>
      <p>&nbsp;</p>
      <el-card>
       <el-main>
         <el-text class="mx-1" size="large">Encrypted Rollup</el-text>
         <p>&nbsp;</p>
         <el-input
             v-model="textarea"
             :rows="3"
             type="textarea"
             placeholder="Please input encrypted rollup"
         />
         <p>&nbsp;</p>
         <el-text class="mx-1" size="large">Decrypted Rollup</el-text>
         <el-input
             v-model="textarea"
             :rows="3"
             type="textarea"
             disabled
             placeholder="Output"
         />
       </el-main>
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  components: {}
})
</script>
<style></style>
