<template>
  <div class="main-container">
    <div class="header-wrapper">
      <div class="header-container">
        <AppTabs />
      </div>
    </div>
    <div class="content-wrapper">
      <SetListSidebar
        :ipSets="ipSetStore.ipSets"
        :selectedSet="selectedSet"
        @select-set="selectSet"
      />
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
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useIpSetStore } from '../stores/ipsetStore'
import AddEntry from '../components/AddEntry.vue'
import SetListSidebar from '../components/SetListSidebar.vue'
import EntryList from '../components/EntryList.vue'
import AppTabs from '../components/AppTabs.vue'

export default {
  components: {
    AddEntry,
    SetListSidebar,
    EntryList,
    AppTabs,
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
      ipSetStore,
      selectedSet,
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
@import '../assets/common.css';

.add-entry-container {
  width: 360px;
}
</style>
