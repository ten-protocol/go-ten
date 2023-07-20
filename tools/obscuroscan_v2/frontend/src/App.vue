<template>
  <div id="app">
    <el-container>
      <el-header class="nav-bar-header">
        <NavbarView />
      </el-header>
      <el-main class="main-content" >
        <RouterView style="min-height: 100vh"/>
      </el-main>
      <el-footer class="bot-bar">
        <bot-bar-item></bot-bar-item>
      </el-footer>
    </el-container>
  </div>
</template>

<script>
import NavbarView from '@/views/NavbarView.vue'
import BotBarItem from "@/components/BotBarItem.vue";

export default {
  name: 'App',
  components: {BotBarItem, NavbarView },
  data() {
    return {
      // data properties go here
    }
  },
  methods: {
    // methods go here
  },
  computed: {
    // computed properties go here
  }
}
</script>

<style scoped>
.nav-bar-header {
  background:black;
  min-height: 40vh;
}

.main-content {
  margin-top: -35vh;
}

.bot-bar {
  background:black;
}
</style>
