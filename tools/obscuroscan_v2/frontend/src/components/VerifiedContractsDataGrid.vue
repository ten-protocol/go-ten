<template>
    <el-card style="width: 100%">
      <el-header>Verified Sequencer Data</el-header>
      <el-table :height="'25vh'" style="width: 100%" :data="sequencerData">
        <el-table-column prop="name" label="Contract Name" width="180" />
        <el-table-column prop="confirmed" label="Confirmed" width="180" />
        <el-table-column prop="address" label="Contract Address">
        </el-table-column>
      </el-table>
    </el-card>
  <p>&nbsp;</p>
    <el-card>
      <el-header>Verified Contracts</el-header>
    <el-table :height="'30vh'" style="width: 100%" :data="contracts">
      <el-table-column prop="name" label="Contract Name" width="180" />
      <el-table-column prop="confirmed" label="Confirmed" width="180" />
      <el-table-column prop="address" label="Contract Address">
      </el-table-column>
    </el-table>
    </el-card>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useVerifiedContractStore } from "@/stores/verifiedContractStore";

export default {
  name: "VerifiedContractsDataGrid",
  components: {  },
  setup() {
    const contracts = ref([]);
    const sequencerData = ref([]);

    async function loadContracts() {
      const store = useVerifiedContractStore();

      await store.update();
      contracts.value = store.contracts;
      sequencerData.value = store.sequencerData;
    }

    onMounted(() => {
      loadContracts();
    });

    return {
      contracts,
      sequencerData
    };
  }
}
</script>

<style scoped></style>