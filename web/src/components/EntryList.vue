<template>
  <div class="container">
    <div class="input-group">
      <div class="input-wrapper">
        <input
          v-model="filter"
          :class="{ 'has-filter': filter }"
          @keyup.enter="searchIpSet"
          type="text"
          placeholder="Empty / IP / CIDR"
        />
        <button v-if="filter" class="clear" title="Clear" @click="clearFilter">×</button>
        <button class="search-icon" title="Search" @click="searchIpSet">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
      </div>
    </div>
    <div v-if="canDeleteAll" class="bulk-actions">
      <span class="count">{{ ipEntries.length }} entries found</span>
      <button class="delete-all" @click="confirmDeleteAll">Delete all</button>
    </div>
    <div class="results">
      <div v-if="ipEntries.length === 0 && (errorMessage || hasSearched)">
        {{ errorMessage || 'Entries Not Found' }}
      </div>
      <div v-else>
        <div v-for="(entry, index) in ipEntries" :key="index" class="entry">
          <span>
            {{ entry.Entry }}
            <template v-if="entry.Comment">
              <span class="comment">({{ entry.Comment }})</span>
            </template>
          </span>
          <div class="buttons">
            <button class="whois" title="whois" @click="confirmWhois(entry.Entry)">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="10"></circle>
                <path d="M12 17v-6.5"></path>
                <path d="M12 7.5h.01"></path>
              </svg>
            </button>
            <button class="delete" title="delete" @click="confirmDelete(entry.Entry)">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M3 6h18"></path>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"></path>
                <path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
    <WhoisModal
      :isVisible="isWhoisModalVisible"
      :whoisData="whoisData"
      @close="isWhoisModalVisible = false"
    />
    <ConfirmModal
      :visible="isDeleteModalVisible"
      :message="`Are you sure you want to delete ${entryToDelete}?`"
      :onConfirm="deleteEntryConfirmed"
      :onCancel="
        () => {
          isDeleteModalVisible = false
        }
      "
    />
    <ConfirmModal
      :visible="isWhoisConfirmVisible"
      :message="whoisConfirmText"
      :onConfirm="whoisConfirmed"
      :onCancel="cancelWhois"
    />
    <ConfirmModal
      :visible="isDeleteAllModalVisible"
      :message="deleteAllConfirmText"
      :onConfirm="deleteAllConfirmed"
      :onCancel="
        () => {
          isDeleteAllModalVisible = false
        }
      "
    />
    <AppNotification
      v-if="notificationMessage"
      :message="notificationMessage"
      :type="notificationType"
      @close="notificationMessage = ''"
    />
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import WhoisModal from './WhoisModal.vue'
import ConfirmModal from './ConfirmModal.vue'
import AppNotification from './AppNotification.vue'
import axios from '../axios'
import { useIpSetStore } from '../stores/ipsetStore'

export default {
  name: 'EntryList',
  components: {
    WhoisModal,
    ConfirmModal,
    AppNotification,
  },
  props: {
    selectedSet: {
      type: String,
      required: true,
    },
    ipEntries: {
      type: Array,
      required: true,
    },
    errorMessage: {
      type: String,
      required: false,
    },
    initialFilter: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const ipSetStore = useIpSetStore()
    const filter = ref(props.initialFilter)
    // Keep the search box in sync when the restored filter is invalidated (e.g. set removed)
    watch(
      () => props.initialFilter,
      (value) => {
        filter.value = value
      },
    )
    const isWhoisModalVisible = ref(false)
    const whoisData = ref('')
    const isDeleteModalVisible = ref(false)
    const entryToDelete = ref('')
    const notificationMessage = ref('')
    const notificationType = ref('success')

    // Для подтверждения Whois
    const isWhoisConfirmVisible = ref(false)
    const whoisConfirmText = ref(
      'Frequent Whois requests may result in your IP address being temporarily blocked. Are you sure you want to proceed?',
    )
    const entryToWhois = ref('')

    const searchIpSet = () => {
      if (props.selectedSet) {
        emit('search-ip-set', filter.value)
      } else {
        emit('update:errorMessage', 'Please select an IP set')
      }
    }

    const clearFilter = () => {
      filter.value = ''
      emit('clear-filter')
    }

    const confirmWhois = (entry) => {
      entryToWhois.value = entry
      isWhoisConfirmVisible.value = true
    }

    const whoisConfirmed = async () => {
      isWhoisConfirmVisible.value = false
      await whois(entryToWhois.value)
    }

    const cancelWhois = () => {
      isWhoisConfirmVisible.value = false
      entryToWhois.value = ''
    }

    const whois = async (entry) => {
      try {
        const response = await axios.post('/whois', { object: entry })
        whoisData.value = response.data
        isWhoisModalVisible.value = true
      } catch (error) {
        console.error('Whois request failed', error)
      }
    }

    const confirmDelete = (entry) => {
      entryToDelete.value = entry
      isDeleteModalVisible.value = true
    }

    const isDeleteAllModalVisible = ref(false)
    const canDeleteAll = computed(() => props.ipEntries.length > 0 && ipSetStore.lastFilter !== '')
    const hasSearched = computed(() => ipSetStore.hasSearched)
    const deleteAllConfirmText = computed(
      () =>
        `Are you sure you want to delete all ${props.ipEntries.length} entries ` +
        `matching "${ipSetStore.lastFilter}"? This action cannot be undone.`,
    )

    const confirmDeleteAll = () => {
      if (!canDeleteAll.value) return
      isDeleteAllModalVisible.value = true
    }

    const deleteAllConfirmed = async () => {
      isDeleteAllModalVisible.value = false
      try {
        const message = await ipSetStore.deleteEntriesByFilter(
          props.selectedSet,
          ipSetStore.lastFilter,
        )
        // Refresh with the filter that is currently applied to the list
        await ipSetStore.searchIpSet(props.selectedSet, ipSetStore.lastFilter)
        showNotification(message, 'success')
      } catch (error) {
        console.error('Bulk delete request failed', error)
        showNotification('Error deleting entries', 'error')
      }
    }

    const deleteEntryConfirmed = async () => {
      isDeleteModalVisible.value = false
      try {
        const message = await ipSetStore.deleteEntry(props.selectedSet, entryToDelete.value)
        // Refresh with the filter that is currently applied to the list
        await ipSetStore.searchIpSet(props.selectedSet, ipSetStore.lastFilter)
        showNotification(message, 'success')
      } catch (error) {
        console.error('Delete request failed', error)
        showNotification('Error deleting entry', 'error')
      }
    }

    const showNotification = (message, type = 'success') => {
      notificationMessage.value = message
      notificationType.value = type
      if (type === 'success') {
        setTimeout(() => {
          notificationMessage.value = ''
        }, 2000)
      }
    }

    return {
      filter,
      isWhoisModalVisible,
      whoisData,
      isDeleteModalVisible,
      entryToDelete,
      notificationMessage,
      notificationType,
      searchIpSet,
      clearFilter,
      whois,
      confirmDelete,
      deleteEntryConfirmed,
      isDeleteAllModalVisible,
      canDeleteAll,
      hasSearched,
      deleteAllConfirmText,
      confirmDeleteAll,
      deleteAllConfirmed,
      showNotification,
      // Новые для подтверждения Whois
      isWhoisConfirmVisible,
      whoisConfirmText,
      confirmWhois,
      whoisConfirmed,
      cancelWhois,
    }
  },
}
</script>

<style scoped>
.container {
  width: 400px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 10px;
  box-shadow: 0 2px 8px #0000000f;
  background-color: #ffffff;
}

.input-group {
  margin-bottom: 20px;
}

.input-wrapper {
  position: relative;
}

.input-wrapper input {
  box-sizing: border-box;
  width: 100%;
  padding: 8px 36px 8px 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.input-wrapper input.has-filter {
  padding-right: 64px;
}

.input-wrapper .search-icon {
  position: absolute;
  top: 50%;
  right: 6px;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  padding: 0;
  border: none;
  background: none;
  color: #666;
  border-radius: 4px;
  cursor: pointer;
  transition: color 0.2s;
}

.input-wrapper .search-icon svg {
  width: 16px;
  height: 16px;
}

.input-wrapper .search-icon:hover {
  color: #007bff;
}

.input-wrapper .clear {
  position: absolute;
  top: 50%;
  right: 38px;
  transform: translateY(-50%);
  padding: 2px 6px;
  border: none;
  background: none;
  color: #666;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
}

.input-wrapper .clear:hover {
  background: none;
  color: #333;
}

.bulk-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.bulk-actions .count {
  color: #666;
  font-size: 14px;
}

.bulk-actions .delete-all {
  padding: 4px 12px;
  border: none;
  border-radius: 4px;
  background-color: #6c757d;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.bulk-actions .delete-all:hover {
  background-color: #5a6268;
}

.results .entry {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-bottom: 1px solid #ccc;
  font-size: 16px;
}

.results .entry span {
  flex: 1;
  text-align: left;
}

.results .entry .buttons {
  display: flex;
  justify-content: flex-end;
}

.results .entry .buttons button {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 10px;
  padding: 4px 6px;
  border: none;
  border-radius: 4px;
  background: none;
  color: #666;
  cursor: pointer;
  transition: color 0.2s;
}

.results .entry .buttons button svg {
  width: 16px;
  height: 16px;
}

.results .entry .buttons .whois:hover {
  color: #007bff;
}

.results .entry .buttons .delete:hover {
  color: #dc3545;
}
</style>
