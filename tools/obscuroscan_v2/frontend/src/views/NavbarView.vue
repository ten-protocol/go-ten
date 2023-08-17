<template>
  <header>
    <el-row class="align-vertical" justify="space-between">
      <el-col :span="6">
        <img src="@/assets/imgs/obscuroscan_logo.png" alt="obscu.ro" class="header-image" />
      </el-col>
      <el-col :span="4">
        <meta-mask-connect-button />
      </el-col>
      <el-col :span="8">
        <nav class="nav-bar">
          <el-radio-group size="default">
            <RouterLink to="/">  <el-radio-button label="top"> Home</el-radio-button></RouterLink>
            <RouterLink to="/personal">  <el-radio-button label="right">Personal</el-radio-button></RouterLink>
            <el-radio-button label="left">
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
              </el-radio-button>
            <el-radio-button label="left">
              <el-dropdown>
                <span class="el-dropdown-link">
                  Resources
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <RouterLink to="/decrypt"> <el-dropdown-item>Decrypt</el-dropdown-item></RouterLink>
                    <RouterLink to="/verified"> <el-dropdown-item>Verified Contracts</el-dropdown-item></RouterLink>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-radio-button>
          </el-radio-group>
        </nav>
      </el-col>
      <el-col :span="6">
        <search-bar-item />
      </el-col>
    </el-row>
    <el-row>
      <div class="white-bar"></div>
    </el-row>
  </header>
</template>

<script>
import SearchBarItem from '@/components/SearchBarItem.vue'
import MetaMaskConnectButton from '@/components/MetaMaskConnectButton.vue'

export default {
  name: 'NavbarView',
  components: { MetaMaskConnectButton, SearchBarItem }
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

.white-bar {
  background-color: white;
  height: 1px;
  width: 98%;
  margin-top: 5px; /* Adds a small margin at the top */
  align-self: center; /* Centers the bar in the flex layout */
}

.align-vertical {
  display: flex;
  align-items: center;
}
</style>
