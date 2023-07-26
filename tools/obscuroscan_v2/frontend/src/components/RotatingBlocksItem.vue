<template>
  <el-card class="centered-container">
  <div class="slider-container">
    <div class="slider-content" :class="{ 'slide-out': isAnimating }">
      <el-card v-for="card in displayedCards"  class="card-content">
        <h2>Height: {{ card.number }}</h2>
        <div>Hash: {{ card.hash }}</div>
        <div>Time: {{ card.timestamp }}</div>
        <div>No of Txs: {{ card.hash }}</div>
      </el-card>
    </div>
  </div>
  </el-card>
</template>

<script>
import {computed, onMounted, onUnmounted} from "vue";
import {useBlockStore} from "@/stores/blockStore";

export default {
  name: "RotatingBlocksItem",

  setup() {
    const blocks = useBlockStore()

    // Start polling when the component is mounted
    onMounted(() => {
      blocks.startPolling()
    })

    // Ensure to stop polling when component is destroyed or deactivated
    onUnmounted(() => {
      blocks.stopPolling()
    })

    return {
      displayedCards:  computed(() => blocks.blocks.get()),
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