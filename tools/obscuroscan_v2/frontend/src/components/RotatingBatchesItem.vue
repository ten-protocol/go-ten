<template>
  <el-card class="centered-container">
  <div class="slider-container">
    <div class="slider-content" :class="{ 'slide-out': isAnimating }">
      <el-card v-for="card in displayedCards" :key="card.hash" class="card-content">
        <h2>Height: {{ card.Number }}</h2>
        <div>Hash: {{ card.hash }}</div>
        <div>Time: {{ card.timestamp }}</div>
        <div>No of Txs: {{ card.hash }}</div>
      </el-card>
    </div>
  </div>
  </el-card>
</template>

<script>
import {useBatchStore} from "@/stores/batchStore";
import {computed, onMounted, onUnmounted} from "vue";

export default {
  name: "RotatingBatchesItem",

  setup() {
    const batch = useBatchStore()

    // Start polling when the component is mounted
    onMounted(() => {
      batch.startPolling()
    })

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      batch.stopPolling()
    })

    return {
      displayedCards:  computed(() => batch.batches.get()),
    }
  },
}
</script>

<style scoped>

.centered-container {
  display: flex;
  justify-content: center;  /* Center children horizontally */
  align-items: center;      /* Center children vertically */
}

.slider-container {
  overflow-x: auto;
  width: 1000px;  /* Adjust based on your design */
  white-space: nowrap;
}

.slider-content {
  display: flex;
}

.card-content {
  display: inline-block;
  width: 200px;  /* If 5 cards need to be displayed at a time, and container width is 1000px, then each card can be approximately 200px wide */
  margin-right: 20px;  /* Adjust as needed */
}


.slide-out {
  animation: slideOut 1s forwards; /* Animation to slide the cards to the left */
}

@keyframes slideOut {
  0% { transform: translateX(0%); }
  100% { transform: translateX(-20%); } /* Assuming each card takes 20% width */
}

</style>