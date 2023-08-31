<template>
  <el-container>
    <el-main>
      <el-card>
        <el-header>
          <h2>Static Keys</h2>
        </el-header>
        <el-card style="text-align: center">
          <h3>Decrypting transaction blobs is only possible on testnet, where the rollup encryption key is long-lived and well-known. </h3>
          <h3>On mainnet, rollups will use rotating keys that are not known to anyone - or anything - other than the Obscuro enclaves.</h3>
        </el-card>
        <p>&nbsp;</p>
        <h3> Current Static Encryption Key: bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c</h3>
      </el-card>
      <p>&nbsp;</p>
      <el-card>
        <el-main>
          <el-text class="mx-1" size="large">Encrypted Rollup</el-text>
          <p>&nbsp;</p>
          <el-input
              v-model="inputEncryptedRollup"
              :rows="3"
              type="textarea"
              placeholder="Please input encrypted rollup"
          />
          <p>&nbsp;</p>
          <el-button @click="decryptData">Decrypt</el-button>
          <p v-if="error" style="color: red;">{{ error }}</p>
          <p>&nbsp;</p>
          <el-text class="mx-1" size="large">Decrypted Rollup</el-text>
          <vue-json-pretty :data="outputDecryptedRollup"></vue-json-pretty>
        </el-main>
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import Config from "@/lib/config";
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';

export default defineComponent({
  components: {
    VueJsonPretty,
  },
  setup() {
    const inputEncryptedRollup = ref(null);
    const outputDecryptedRollup = ref(null);
    const error = ref(null);

    const route = useRoute(); // Use the useRoute hook
    if(route.query.encryptedString) {
      inputEncryptedRollup.value = decodeURIComponent(route.query.encryptedString);
    }

    async function decryptData() {
      try {
        const response = await axios.post(Config.backendServerAddress+`/actions/decryptTxBlob/`, { StrData: inputEncryptedRollup.value });
        outputDecryptedRollup.value = response.data.result;
      } catch (err) {
        error.value = "There was an issue with the decrypt operation.";
      }
    }

    return {
      inputEncryptedRollup,
      outputDecryptedRollup,
      error,
      decryptData
    }
  }
})
</script>
<style></style>