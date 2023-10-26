<template>
  <el-row>
    <el-col :span="4" class="socials">
      <img src="@/assets/imgs/github_logo.png" alt="Github Logo" class="header-image" />

      <img src="@/assets/imgs/twitter_logo.png" alt="Twitter Logo" class="header-image" />
    </el-col>
    <el-col :span="20" class="version-column">
      <div class="version-display">version: {{ version }}</div>
    </el-col>
  </el-row>
</template>

<script>
import Config from "@/lib/config";

export default {
  name: "BotBarItem",
  computed: {
    version() {
      return Config.version;
    }
  }
}
</script>

<style scoped>
.version-column {
  text-align: right;
}

.socials {
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
    align-items: center;
}
</style>