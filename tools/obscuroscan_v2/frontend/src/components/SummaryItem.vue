<template>
    <el-col :span="4" class="stats-item">
          <p class="header-text">
            <img class="icon" :src="props.icon"/>
            {{ props.label }}
          </p>

      <div class="stat-value">
        <slot>
          <!-- placeholder if no content provided -->
          <div v-loading=true></div>
        </slot>
      </div>
    </el-col>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'SummaryItem',
  props: ['label', 'icon'],
  setup(props) {
    return {
      props,
    }
  }
}
</script>

<style scoped>
.icon {
  height: 24px;      /* Set desired height */
  object-fit: cover; /* Ensure image content is not distorted */
  margin-right: 8px; /* Optional space between the icon and the text */
}
.header-text {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.stat-value {
  text-align: right;
  width: fit-content;
  margin: auto;
}
</style>
