<template>
  <div>
    <el-tooltip
        class="box-item"
        effect="dark"
        :content="this.hash"
        placement="top-start"
    >
   {{ shortenedHash }}
    </el-tooltip>
  </div>
</template>

<script>
export default {
  props: {
    hash: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      showTooltip: false
    };
  },
  computed: {
    shortenedHash() {
      if (!this.hash) return "";
      const start = this.hash.substring(0, 6);
      const end = this.hash.substring(this.hash.length - 6);
      return `${start}...${end}`;
    }
  }
};
</script>

<style scoped>

</style>
