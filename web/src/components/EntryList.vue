<template>
  <div class="container">
    <div class="input-group">
      <input
        v-model="filter"
        @keyup.enter="searchIpSet"
        type="text"
        placeholder="Empty / IP / CIDR"
      />
      <button @click="searchIpSet">Search</button>
    </div>
    <div class="results">
      <div v-if="ipEntries.length === 0">
        {{ errorMessage || 'Not Found' }}
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
            <button class="whois" @click="whois(entry.Entry)">whois</button>
            <button class="delete" @click="confirmDelete(entry.Entry)">delete</button>
          </div>
        </div>
      </div>
    </div>
    <WhoisModal
      :isVisible="isWhoisModalVisible"
      :whoisData="whoisData"
      @close="isWhoisModalVisible = false"
    />
    <DeleteConfirmationModal
      :isVisible="isDeleteModalVisible"
      :entry="entryToDelete"
      @close="isDeleteModalVisible = false"
      @confirm="deleteEntryConfirmed"
    />
    <Notification
      v-if="notificationMessage"
      :message="notificationMessage"
      :type="notificationType"
    />
  </div>
</template>

<script>
import { ref } from 'vue'
import WhoisModal from './WhoisModal.vue'
import DeleteConfirmationModal from './DeleteConfirmationModal.vue'
import Notification from './Notification.vue'
import axios from '../axios'

export default {
  name: 'EntryList',
  components: {
    WhoisModal,
    DeleteConfirmationModal,
    Notification,
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
  },
  setup(props, { emit }) {
    const filter = ref('')
    const isWhoisModalVisible = ref(false)
    const whoisData = ref('')
    const isDeleteModalVisible = ref(false)
    const entryToDelete = ref('')
    const notificationMessage = ref('')
    const notificationType = ref('success')

    const searchIpSet = () => {
      if (props.selectedSet) {
        emit('search-ip-set', filter.value)
      } else {
        emit('update:errorMessage', 'Please select an IP set')
      }
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

    const deleteEntryConfirmed = async () => {
      try {
        const message = await emit('delete-entry', entryToDelete.value)
        isDeleteModalVisible.value = false
        searchIpSet() // Refresh the list after deletion
        showNotification(message, 'success')
      } catch (error) {
        console.error('Delete request failed', error)
        showNotification('Error deleting entry', 'error')
      }
    }

    const showNotification = (message, type = 'success') => {
      notificationMessage.value = message
      notificationType.value = type
      setTimeout(() => {
        notificationMessage.value = ''
      }, 2000)
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
      whois,
      confirmDelete,
      deleteEntryConfirmed,
      showNotification,
    }
  },
}
</script>

<style scoped>
.container {
  width: 400px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #ffffff;
}

.input-group {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.input-group input {
  flex: 1;
  padding: 8px;
  margin-right: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.input-group button {
  padding: 8px 16px;
  border: none;
  background-color: #007bff;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.input-group button:hover {
  background-color: #0056b3;
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
  margin-left: 10px;
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.results .entry .buttons .whois {
  background-color: #007bff;
  color: white;
}

.results .entry .buttons .whois:hover {
  background-color: #0056b3;
}

.results .entry .buttons .delete {
  background-color: #dc3545;
  color: white;
}

.results .entry .buttons .delete:hover {
  background-color: #c82333;
}
</style>
