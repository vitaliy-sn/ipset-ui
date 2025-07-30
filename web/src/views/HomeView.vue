<template>
  <div class="main-container">
    <div class="content-wrapper">
      <HomeSetList :ipSets="ipSetStore.ipSets" :selectedSet="selectedSet" @select-set="selectSet" />
      <EntryList
        :selectedSet="selectedSet"
        :ipEntries="ipSetStore.ipEntries"
        :errorMessage="ipSetStore.errorMessage"
        @search-ip-set="searchIpSet"
        @delete-entry="deleteEntry"
      />
      <div class="add-entry-container">
        <AddEntry :selectedSet="selectedSet" @check-ip="checkIp" />
      </div>
    </div>
    <button class="floating-button" @click="goToSets">Manage IP Sets</button>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useIpSetStore } from '../stores/ipsetStore'
import AddEntry from '../components/AddEntry.vue'
import HomeSetList from '../components/HomeSetList.vue'
import EntryList from '../components/EntryList.vue'

export default {
  components: {
    AddEntry,
    HomeSetList,
    EntryList,
  },
  setup() {
    const ipSetStore = useIpSetStore()
    const selectedSet = ref('')
    const router = useRouter()

    const searchIpSet = (filter = '') => {
      if (selectedSet.value) {
        ipSetStore.searchIpSet(selectedSet.value, filter)
      } else {
        ipSetStore.errorMessage = 'Please select an IP set'
      }
    }

    const deleteEntry = async (entry) => {
      try {
        const message = await ipSetStore.deleteEntry(selectedSet.value, entry)
        searchIpSet('') // Refresh the list after deletion
        return message
      } catch (error) {
        console.error('Delete request failed', error)
        throw error
      }
    }

    const selectSet = (set) => {
      selectedSet.value = set
      ipSetStore.ipEntries = [] // Clear the current entries
      ipSetStore.errorMessage = '' // Clear any error message
    }

    const checkIp = (ip) => {
      searchIpSet(ip)
    }

    const goToSets = () => {
      router.push('/sets')
    }

    onMounted(async () => {
      await ipSetStore.fetchIpSets()
      if (ipSetStore.ipSets.length > 0) {
        selectedSet.value = ipSetStore.ipSets[0]
      }
    })

    return {
      selectedSet,
      ipSetStore,
      searchIpSet,
      deleteEntry,
      selectSet,
      checkIp,
      goToSets,
    }
  },
}
</script>

<style scoped>
.main-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
}

.content-wrapper {
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.add-entry-container {
  margin-left: 20px;
  width: 360px;
}

.floating-button {
  position: fixed;
  top: 10px;
  right: 20px;
  padding: 10px 20px;
  border-radius: 4px;
  background-color: #28a745;
  color: white;
  border: none;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.floating-button:hover {
  background-color: #218838;
}
</style>
