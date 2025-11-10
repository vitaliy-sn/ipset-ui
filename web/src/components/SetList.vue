<template>
  <div class="sets-container">
    <div class="input-group">
      <input v-model="newSetName" type="text" placeholder="Enter new set name" />
      <button @click="addSet">Add</button>
    </div>
    <div class="results">
      <div v-if="ipSetStore.ipSets.length === 0">No sets available.</div>
      <div v-else>
        <div v-for="set in ipSetStore.ipSets" :key="set" class="entry">
          <span>{{ set }}</span>
          <div class="buttons">
            <button class="backups" @click="toggleBackup(set)">backups</button>
            <button class="delete" @click="confirmDelete(set)">delete</button>
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
    <Notification
      v-if="notificationMessage"
      :message="notificationMessage"
      :type="notificationType"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useIpSetStore } from '../stores/ipsetStore'
import ConfirmModal from './ConfirmModal.vue'
import Notification from './Notification.vue'

export default {
  components: {
    ConfirmModal,
    Notification,
  },
  setup(_, { emit }) {
    const ipSetStore = useIpSetStore()
    const newSetName = ref('')
    const isDeleteModalVisible = ref(false)
    const setToDelete = ref('')
    const notificationMessage = ref('')
    const notificationType = ref('success')

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
      setTimeout(() => {
        notificationMessage.value = ''
      }, 3000)
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

.results .entry .buttons .backups {
  background-color: #28a745;
  color: white;
}

.results .entry .buttons .backups:hover {
  background-color: #218838;
}

.results .entry .buttons .delete {
  background-color: #dc3545;
  color: white;
}

.results .entry .buttons .delete:hover {
  background-color: #c82333;
}
</style>
