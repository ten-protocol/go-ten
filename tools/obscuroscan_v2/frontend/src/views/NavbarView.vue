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
            <el-radio-button label="top"> <RouterLink to="/">Home</RouterLink></el-radio-button>
            <el-radio-button label="right"><RouterLink to="/personal">Personal</RouterLink></el-radio-button>
            <el-radio-button label="left">
              <el-dropdown>
                <span class="el-dropdown-link">
                  Blockchain
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item><RouterLink to="/transactions">Transactions</RouterLink></el-dropdown-item>
                    <el-dropdown-item><RouterLink to="/batches">Batches</RouterLink></el-dropdown-item>
                    <el-dropdown-item><RouterLink to="/rollups">Rollups</RouterLink></el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              </el-radio-button>
            <el-radio-button label="left"
              ><RouterLink to="/about">Resources</RouterLink></el-radio-button
            >
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
