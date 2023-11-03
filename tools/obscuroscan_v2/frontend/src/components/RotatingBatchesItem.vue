<template>
  <el-row class="batches-row">
    <el-button class="all-batches-btn obs-link-button" @click="$router.push('/batches')"><el-icon><ArrowLeftBold /></el-icon> View older batches</el-button>
    <el-card v-for="batch in batchList" :key="batch.hash" class="batch-card">
      <template #header>
        <div class="card-header" >Batch <span class="batch-height">{{ batch.number }}</span></div>
      </template>
      <p class="prop-label">Batch Hash:</p>
      <el-row class="hash prop-val"><ShortenedHash :hash="batch.hash" /><CopyButton :value="batch.hash"/></el-row>
      <p class="prop-label">L1 Block:</p>
      <el-row class="hash prop-val"><ShortenedHash :hash="batch.l1Proof" /><CopyButton :value="batch.l1Proof"/></el-row>
      <p class="prop-label">Tx Count:</p>
      <p class="prop-val">{{ batch.txHashes ? batch.txHashes.length : 0}}</p>
      <p class="prop-label">Timestamp:</p>
      <p class="prop-val timestamp"><Timestamp :unixTimestampSeconds="Number(batch.timestamp)" /></p>
    </el-card>
  </el-row>
</template>

<script>
import {useBatchStore} from "@/stores/batchStore";
import {computed, onMounted, onUnmounted} from "vue";
import Timestamp from "@/components/helper/Timestamp.vue";
import ShortenedHash from "@/components/helper/ShortenedHash.vue";
import CopyButton from "@/components/helper/CopyButton.vue";

export default {
  name: "RotatingBatchesItem",
  components: {CopyButton, ShortenedHash, Timestamp},

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
      batchList:  computed(() => batch.batches.get()),
      isAnimating: false
    }
  },
}
</script>

<style scoped>
.card-header {
  text-align: right;
}
.batch-height {
  font-weight: bold;
  color: var(--obs-tertiary)
}

.batches-row {
  align-content: center;
  height: fit-content;
}

.all-batches-btn {
  border-radius: 8px;
  margin: auto 0;
  flex-grow: 1;
  text-align: center;
}

.batch-card {
  width: 14rem;
  margin: 12px;
  border-radius: 8px;
}

.batch-card.el-card__header {
  border-bottom-color: var(--obs-secondary) !important;
}
.prop-label {
  font-size: 0.6rem;
  color: #A3A3A3;
}
.hash {
  gap: 0.5rem;
  justify-items: end;
  align-items: center;
}

.prop-val {
  font-size: 0.75rem;
}

</style>