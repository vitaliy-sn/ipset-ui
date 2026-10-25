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
        :initialFilter="initialFilter"
        @search-ip-set="searchIpSet"
      />
      <div class="add-entry-container">
        <AddEntry :selectedSet="selectedSet" @check-ip="checkIp" />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
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
    const route = useRoute()

    // Seed the search box from the URL synchronously so a page reload restores it
    const initialFilter = ref(typeof route.query.filter === 'string' ? route.query.filter : '')

    const searchIpSet = (filter = '') => {
      if (selectedSet.value) {
        ipSetStore.searchIpSet(selectedSet.value, filter)
        router.replace({ query: { set: selectedSet.value, ...(filter ? { filter } : {}) } })
      } else {
        ipSetStore.errorMessage = 'Please select an IP set'
      }
    }

    const selectSet = (set) => {
      selectedSet.value = set
      ipSetStore.ipEntries = [] // Clear the current entries
      ipSetStore.errorMessage = '' // Clear any error message
      ipSetStore.lastFilter = '' // Reset the applied search filter
      router.replace({ query: set ? { set } : {} })
    }

    const checkIp = (ip) => {
      searchIpSet(ip)
    }

    const goToSets = () => {
      router.push('/sets')
    }

    onMounted(async () => {
      const requestedSet = typeof route.query.set === 'string' ? route.query.set : ''
      const requestedFilter = initialFilter.value
      // A warm store means we came back from another tab; keep the loaded list as-is
      const wasWarm = ipSetStore.ipSets.length > 0
      await ipSetStore.fetchIpSets()
      const validSet = requestedSet && ipSetStore.ipSets.includes(requestedSet)
      if (validSet) {
        selectedSet.value = requestedSet
        // Only re-run the search on a cold start (a real page reload)
        if (requestedFilter && !wasWarm) {
          searchIpSet(requestedFilter)
        }
      } else {
        // The set from the URL no longer exists: fall back and drop the stale filter
        initialFilter.value = ''
        selectSet(ipSetStore.ipSets[0] || '')
      }
    })

    return {
      ipSetStore,
      selectedSet,
      initialFilter,
      searchIpSet,
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
