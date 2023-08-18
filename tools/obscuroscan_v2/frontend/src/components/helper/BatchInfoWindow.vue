<template>
  <el-drawer v-model="isVisible" :direction="'rtl'">
    <router-link :to="`/decrypt?encryptedString=${encodeURIComponent(displayedData.EncryptedTxBlob)}`">
      <el-button>Decrypt Transactions</el-button>
    </router-link>
    <p>&nbsp;</p>
    <el-card>
    <vue-json-pretty :data="displayedData"></vue-json-pretty>
    </el-card>
  </el-drawer>
</template>

<script>
import {ref} from "vue";
import {useBatchStore} from "@/stores/batchStore";
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';

export default {
  name: 'BatchInfoWindow',
  components: {
    VueJsonPretty,
  },
  setup() {
    return {
      isVisible: ref(false),
      displayedData: null
    }
  },
  methods: {
    async displayData(hash){
      const store = useBatchStore();
      this.displayedData = await store.getByHash(hash);
      this.isVisible = true;
    }
  }
}
</script>

<style scoped>
/* Add styling if required */
</style>