<template>
  <div class="sets-container">
    <div class="input-group">
      <div class="input-wrapper">
        <input v-model="newSetName" type="text" placeholder="Enter new set name" />
        <button v-if="newSetName" class="clear" title="Clear" @click="clearName">×</button>
      </div>
      <button @click="addSet">Create</button>
    </div>
    <div class="results">
      <div v-if="ipSetStore.ipSets.length === 0">No sets available.</div>
      <div v-else>
        <div v-for="set in ipSetStore.ipSets" :key="set" class="entry">
          <span>{{ set }}</span>
          <div class="buttons">
            <button class="backups" title="backups" aria-label="backups" @click="toggleBackup(set)">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
                <path d="M3 5v14a9 3 0 0 0 18 0V5"></path>
                <path d="M3 12a9 3 0 0 0 18 0"></path>
              </svg>
            </button>
            <button class="delete" title="Delete" aria-label="Delete" @click="confirmDelete(set)">
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
    <ConfirmModal
      :visible="isDeleteModalVisible"
      :message="`Are you sure you want to delete ${setToDelete}?`"
      :onConfirm="deleteSetConfirmed"
      :onCancel="
        () => {
          isDeleteModalVisible = false
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
import { ref, onMounted } from 'vue'
import { useIpSetStore } from '../stores/ipsetStore'
import ConfirmModal from './ConfirmModal.vue'
import AppNotification from './AppNotification.vue'

export default {
  components: {
    ConfirmModal,
    AppNotification,
  },
  setup(_, { emit }) {
    const ipSetStore = useIpSetStore()
    const newSetName = ref('')
    const isDeleteModalVisible = ref(false)
    const setToDelete = ref('')
    const notificationMessage = ref('')
    const notificationType = ref('success')

    const clearName = () => {
      newSetName.value = ''
    }

    const addSet = async () => {
      if (newSetName.value.trim()) {
        try {
          await ipSetStore.addSet(newSetName.value.trim())
          newSetName.value = ''
          showNotification('Set added successfully', 'success')
        } catch {
          showNotification('Failed to add set', 'error')
        }
      }
    }

    const confirmDelete = (set) => {
      setToDelete.value = set
      isDeleteModalVisible.value = true
    }

    const deleteSetConfirmed = async () => {
      try {
        await ipSetStore.deleteSet(setToDelete.value)
        showNotification('Set deleted successfully', 'success')
      } catch (e) {
        const msg = e?.response?.data?.error || e.message || 'Failed to delete set'
        showNotification(msg, 'error')
      } finally {
        isDeleteModalVisible.value = false
      }
    }

    const toggleBackup = (set) => {
      emit('toggle-backup', set)
    }

    const showNotification = (message, type = 'success') => {
      notificationMessage.value = message
      notificationType.value = type
      if (type === 'success') {
        setTimeout(() => {
          notificationMessage.value = ''
        }, 3000)
      }
    }

    onMounted(async () => {
      await ipSetStore.fetchIpSets()
    })

    return {
      newSetName,
      isDeleteModalVisible,
      setToDelete,
      notificationMessage,
      notificationType,
      clearName,
      addSet,
      confirmDelete,
      deleteSetConfirmed,
      toggleBackup,
      ipSetStore,
    }
  },
}
</script>

<style scoped>
.sets-container {
  width: 400px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 10px;
  box-shadow: 0 2px 8px #0000000f;
  background-color: #ffffff;
}

.input-group {
  display: flex;
  margin-bottom: 20px;
}

.input-wrapper {
  position: relative;
  flex: 1;
  min-width: 0;
  margin-right: 10px;
}

.input-wrapper input {
  box-sizing: border-box;
  width: 100%;
  padding: 8px 32px 8px 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.input-wrapper .clear {
  position: absolute;
  top: 50%;
  right: 6px;
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

.input-group button {
  padding: 8px 16px;
  border: none;
  background-color: #6c757d;
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

.results .entry .buttons .backups:hover {
  color: #28a745;
}

.results .entry .buttons .delete:hover {
  color: #dc3545;
}
</style>
