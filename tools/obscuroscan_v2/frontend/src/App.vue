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

.bot-bar {
  background:black;
}

.nav-bar-header {
  margin: 0;
  padding: 0;
  filter:drop-shadow(0px 1vh 1vh #E1E1E1);
}

.main-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
