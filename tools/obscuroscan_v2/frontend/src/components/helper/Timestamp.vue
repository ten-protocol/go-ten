<template>
  <div>
    {{ formattedTimestamp }}
  </div>
</template>

<script>
export default {
  props: {
    unixTimestampSeconds: {
      type: Number,
      required: true
    }
  },
  computed: {
    formattedTimestamp() {
      // Convert Unix timestamp from seconds to JavaScript epoch milliseconds
      const date = new Date(this.unixTimestampSeconds * 1000);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0'); // months are 0-indexed
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');

      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    }
  }
};
</script>