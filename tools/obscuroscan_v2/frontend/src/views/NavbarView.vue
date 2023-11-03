<template>
  <div class="navbar">
    <div class="logo">
      <img src="@/assets/imgs/obscuro_black.png" alt="Obscuro logo"/>
    </div>
    <div class="spacer"> </div>
    <el-menu
        class="navbar-menu"
        :default-active="$route.path"
        mode="horizontal"
        :router="true"
        :ellipsis="false"
        active-text-color="#5973B8"
    >
      <el-menu-item index="/">Home</el-menu-item>
      <el-menu-item index="/personal">Personal</el-menu-item>
      <el-sub-menu index="blockchain">
        <template #title>Blockchain</template>
        <el-menu-item index="/transactions">Transactions</el-menu-item>
        <el-menu-item index="/batches">Batches</el-menu-item>
        <el-menu-item index="/blocks">Blocks</el-menu-item>
      </el-sub-menu>
      <el-sub-menu index="resources">
        <template #title>Resources</template>
        <el-menu-item index="/decrypt">Decrypt</el-menu-item>
        <el-menu-item index="/data">Verified Data</el-menu-item>
      </el-sub-menu>
    </el-menu>
    <meta-mask-connect-button />
  </div>
</template>

<script>
import MetaMaskConnectButton from '@/components/MetaMaskConnectButton.vue'

export default {
  name: 'NavbarView',
  components: { MetaMaskConnectButton }
}
</script>

<style scoped>
.navbar {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;

  background: white;
  //height: 59px;
  padding-left: 12px;
  padding-right: 8px;
}

.navbar-menu {
  border: none;
  cursor: pointer;
}

.spacer {
  flex-grow: 1;
}

.logo {
  margin-right: 20px;
  max-height: 59px;
}

.logo img {
  max-height: 59px;
}
</style>
