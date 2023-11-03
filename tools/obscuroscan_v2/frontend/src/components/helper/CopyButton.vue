<template>
  <el-icon class="copy-button" title="Copy to clipboard" @click="copyToClipboard()"><CopyDocument /></el-icon>
</template>

<script>
export default {
  props: ['value'],
  setup(props) {
    function copyToClipboard() {
      navigator.clipboard.writeText(props.value);
    }
    return {
      props,
      copyToClipboard
    }
  }
};
</script>

<style scoped>
.copy-button {
  cursor: pointer;
  opacity: 0.2;
}
.copy-button:hover {
  opacity: 1;
}
</style>
