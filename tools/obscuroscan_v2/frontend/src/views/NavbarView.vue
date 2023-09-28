<template>
  <header>
    <el-container class="align-vertical">
      <el-main>
    <el-row justify="space-around">
      <el-col :span="2" >
        <img src="@/assets/imgs/obscuro_logo_black.png" alt="obscu.ro" class="header-image"  style="max-height: 5vh;"/>
      </el-col>
      <el-col :span="10" :offset="7">
        <nav class="nav-bar">
          <el-radio-group size="default">
            <RouterLink to="/">  <el-button link class="paddedTop"> Home</el-button></RouterLink>
            <RouterLink to="/personal">  <el-button link class="paddedTop">Personal</el-button></RouterLink>
            <el-button link class="paddedTop">
              <el-dropdown>
                <span class="el-dropdown-link">
                  Blockchain
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <RouterLink to="/transactions"> <el-dropdown-item>Transactions</el-dropdown-item></RouterLink>
                    <RouterLink to="/batches"> <el-dropdown-item>Batches</el-dropdown-item></RouterLink>
                    <RouterLink to="/blocks"> <el-dropdown-item>Blocks</el-dropdown-item></RouterLink>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              </el-button>
            <el-button link class="paddedTop" style="margin: 0px">
              <el-dropdown>
                <span class="el-dropdown-link">
                  Resources
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <RouterLink to="/decrypt"> <el-dropdown-item>Decrypt</el-dropdown-item></RouterLink>
                    <RouterLink to="/verified"> <el-dropdown-item>Verified Data</el-dropdown-item></RouterLink>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-button>
          </el-radio-group>
        </nav>
      </el-col>
      <el-col :span="5" >
        <meta-mask-connect-button />
      </el-col>
<!--      <el-col :span="4" >-->
<!--        <search-bar-item />-->
<!--      </el-col>-->
    </el-row>
      </el-main>
    </el-container>
  </header>
</template>

<script>
import MetaMaskConnectButton from '@/components/MetaMaskConnectButton.vue'

export default {
  name: 'NavbarView',
  components: { MetaMaskConnectButton }
}
</script>

<style scoped>
/* Nav bar styles */
.nav-bar {
  border-radius: 8px; /* Optional rounded corners */
}

.header-image {
  max-height: 100%; /* Ensure the image doesn't exceed the header's height */
  object-fit: cover; /* If image's aspect ratio doesn't match, it'll still cover the space without distortion */
}

.paddedTop {
  padding: 10px;
}

</style>
