<template>
  <div class="backup-container">
    <div class="input-group">
      <input
        v-model="backupNamePart"
        :readonly="!isEditing"
        @focus="startEditing"
        @blur="stopEditing"
        maxlength="30"
        placeholder="Enter backup name"
      />
      <button @click="createBackup">Backup</button>
    </div>
    <div class="results">
      <div v-if="!ipSetStore.backups || ipSetStore.backups.length === 0">
        No backups available for {{ setName }}.
      </div>
      <div v-else>
        <div v-for="backup in ipSetStore.backups" :key="backup" class="entry">
          <span>{{ backup }}</span>
          <div class="buttons">
            <button @click="confirmRestore(backup)">Restore</button>
            <button @click="confirmDelete(backup)">Delete</button>
          </div>
        </div>
      </div>
    </div>
    <ConfirmModal
      :visible="isRestoreModalVisible"
      :message="`Are you sure you want to restore the backup ${backupToRestore}?`"
      :onConfirm="restoreBackupConfirmed"
      :onCancel="
        () => {
          isRestoreModalVisible = false
        }
      "
    />
    <ConfirmModal
      :visible="isDeleteModalVisible"
      :message="`Are you sure you want to delete ${backupToDelete}?`"
      :onConfirm="deleteBackupConfirmed"
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
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { useIpSetStore } from '../stores/ipsetStore'
import ConfirmModal from './ConfirmModal.vue'
import Notification from './Notification.vue'

export default {
  components: {
    ConfirmModal,
    Notification,
  },
  props: {
    setName: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const ipSetStore = useIpSetStore()
    const backupNamePart = ref(
      new Date().toISOString().slice(0, 16).replace('T', '-').replace(':', '-'),
    )
    const isEditing = ref(false)
    const isRestoreModalVisible = ref(false)
    const backupToRestore = ref('')
    const isDeleteModalVisible = ref(false)
    const backupToDelete = ref('')
    const notificationMessage = ref('')
    const notificationType = ref('success')

    const fetchBackups = async () => {
      await ipSetStore.fetchBackups(props.setName)
    }

    const createBackup = async () => {
      try {
        await ipSetStore.createBackup(props.setName, backupNamePart.value)
        backupNamePart.value = new Date()
          .toISOString()
          .slice(0, 16)
          .replace('T', '-')
          .replace(':', '-')
        isEditing.value = false
        await fetchBackups() // Refresh the list after creating a backup
        showNotification('Backup created successfully', 'success')
      } catch (error) {
        console.error('Failed to create backup', error)
        showNotification('Failed to create backup', 'error')
      }
    }

    const confirmRestore = (backup) => {
      backupToRestore.value = backup
      isRestoreModalVisible.value = true
    }

    const restoreBackupConfirmed = async () => {
      try {
        await ipSetStore.restoreBackup(props.setName, backupToRestore.value)
        isRestoreModalVisible.value = false
        showNotification('Backup restored successfully', 'success')
      } catch (error) {
        console.error('Failed to restore backup', error)
        showNotification('Failed to restore backup', 'error')
      }
    }

    const confirmDelete = (backup) => {
      backupToDelete.value = backup
      isDeleteModalVisible.value = true
    }

    const deleteBackupConfirmed = async () => {
      try {
        await ipSetStore.deleteBackup(props.setName, backupToDelete.value)
        isDeleteModalVisible.value = false
        await fetchBackups() // Refresh the list after deletion
        showNotification('Backup deleted successfully', 'success')
      } catch (error) {
        console.error('Failed to delete backup', error)
        showNotification('Failed to delete backup', 'error')
      }
    }

    const startEditing = () => {
      isEditing.value = true
      backupNamePart.value = ''
    }

    const stopEditing = () => {
      if (!backupNamePart.value) {
        backupNamePart.value = new Date()
          .toISOString()
          .slice(0, 16)
          .replace('T', '-')
          .replace(':', '-')
        isEditing.value = false
      }
    }

    const showNotification = (message, type = 'success') => {
      notificationMessage.value = message
      notificationType.value = type
      setTimeout(() => {
        notificationMessage.value = ''
      }, 3000)
    }

    watch(() => props.setName, fetchBackups)

    onMounted(fetchBackups)

    onBeforeUnmount(() => {
      ipSetStore.backups = []
    })

    return {
      ipSetStore,
      backupNamePart,
      isEditing,
      isRestoreModalVisible,
      backupToRestore,
      isDeleteModalVisible,
      backupToDelete,
      notificationMessage,
      notificationType,
      fetchBackups,
      createBackup,
      confirmRestore,
      restoreBackupConfirmed,
      confirmDelete,
      deleteBackupConfirmed,
      startEditing,
      stopEditing,
      showNotification,
    }
  },
}
</script>

<style scoped>
.backup-container {
  width: 400px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 10px;
  box-shadow: 0 2px 8px #0000000f;
  background-color: #ffffff;
  margin-left: 20px;
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
  color: #888;
}

.input-group input:focus {
  color: #000;
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

.results .entry .buttons .restore {
  background-color: #28a745;
  color: white;
}

.results .entry .buttons .restore:hover {
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
