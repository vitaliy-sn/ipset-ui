<template>
  <div class="add-entry">
    <div class="tabs">
      <button :class="{ active: activeTab === 'add' }" @click="activeTab = 'add'">Add</button>
      <button :class="{ active: activeTab === 'import' }" @click="activeTab = 'import'">
        Import
      </button>
    </div>
    <div v-if="activeTab === 'add'">
      <h3>Add New Entry</h3>
      <div class="input-group">
        <input v-model="entry" type="text" placeholder="IP / CIDR / Domain" />
      </div>
      <div class="input-group">
        <input v-model="comment" type="text" placeholder="Comment (optional)" />
      </div>
      <div class="import-btn-wrapper">
        <button v-if="!entry || !isDomain(entry)" @click.prevent="addEntry(entry)">Add</button>
        <button
          v-if="entry && isDomain(entry) && ipAddresses.length === 0"
          @click.prevent="lookupDomain"
        >
          Lookup
        </button>
        <button
          v-if="entry && isDomain(entry) && ipAddresses.length > 0"
          @click.prevent="addAllEntries"
        >
          Add All
        </button>
      </div>
      <div class="ip-list" v-if="ipAddresses.length > 0">
        <div v-for="(ip, index) in ipAddresses" :key="index" class="ip-entry">
          <span>{{ ip }}</span>
          <button @click.prevent="checkIp(ip)">Check</button>
          <button @click.prevent="addEntry(ip, comment)">Add</button>
        </div>
      </div>
      <Notification
        v-if="notificationMessage"
        :message="notificationMessage"
        :type="notificationType"
      />
    </div>
    <div v-else-if="activeTab === 'import'">
      <h3>Import</h3>
      <div class="input-group">
        <div
          class="custom-file-input"
          :class="{ disabled: importInProgress }"
          @click="triggerFileInput"
        >
          <span>{{ importFile ? importFile.name : 'Select file...' }}</span>
        </div>
        <input
          type="file"
          ref="fileInput"
          @change="onFileChange"
          style="display: none"
          :disabled="importInProgress"
        />
      </div>
      <div class="input-group">
        <input v-model="importComment" type="text" placeholder="Comment (optional)" />
      </div>
      <div class="import-btn-wrapper">
        <button :disabled="!importFile || importInProgress" @click="importEntries">Import</button>
      </div>
      <Notification
        v-if="notificationMessage"
        :message="notificationMessage"
        :type="notificationType"
      />
    </div>
  </div>
</template>

<script>
import { watch } from 'vue'
import { useIpSetStore } from '../stores/ipsetStore'
import Notification from './Notification.vue'

export default {
  name: 'AddEntry',
  components: {
    Notification,
  },
  props: {
    selectedSet: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      entry: '',
      comment: '',
      notificationMessage: '',
      notificationType: 'success', // Default type
      activeTab: 'add', // Added state for tab
      importFile: null, // Added for reactive file import
      importComment: '', // For comment during import
      resolvedDomain: '', // Save domain after lookup
      importInProgress: false, // Flag to disable import button
    }
  },
  computed: {
    ipSets() {
      return this.ipSetStore.ipSets
    },
    ipAddresses() {
      return this.ipSetStore.ipAddresses
    },
  },
  methods: {
    async addEntry(ip = this.entry, comment = this.comment) {
      if (this.selectedSet && ip) {
        try {
          // If adding IP from domain resolve and comment is empty — use domain as comment
          let usedComment = comment
          if (this.isDomain(this.resolvedDomain) && this.ipAddresses.includes(ip) && !comment) {
            usedComment = this.resolvedDomain
          }
          const payload = {
            entry: ip,
            comment: usedComment,
          }
          const sendPayload = typeof ip === 'object' && ip.entry ? ip : payload
          const response = await this.ipSetStore.addEntry(this.selectedSet, sendPayload)
          if (response.status === 200) {
            if (ip === this.entry) {
              this.entry = ''
              this.comment = ''
              this.resolvedDomain = ''
            }
            this.notificationType = 'success'
            this.showNotification(`IP address ${ip} successfully added`)
            return
          }
          // Error handling: show the text returned by API if present
          this.notificationType = 'error'
          const errorMsg =
            response.data && response.data.error
              ? response.data.error
              : `Error adding IP address ${ip}`
          this.showNotification(errorMsg)
        } catch (error) {
          console.error('Add entry request failed', error)
          this.notificationType = 'error'
          this.showNotification(error?.response?.data?.error || `Error adding IP address ${ip}`)
        }
      }
    },
    onFileChange(event) {
      const files = event.target.files
      if (files && files.length > 0) {
        this.importFile = files[0]
        // If comment is empty — use file name as comment
        if (!this.importComment) {
          this.importComment = files[0].name
        }
      } else {
        this.importFile = null
      }
    },
    triggerFileInput() {
      if (this.importInProgress) return
      this.$refs.fileInput.click()
    },
    async importEntries() {
      if (!this.importFile || !this.selectedSet) return
      this.importInProgress = true // Immediately disable the button
      const formData = new FormData()
      formData.append('file', this.importFile)
      if (this.importComment) {
        formData.append('comment', this.importComment)
      }
      try {
        const response = await this.ipSetStore.importFromFile(this.selectedSet, formData)
        if (response.status === 200) {
          this.notificationType = 'success'
          this.showNotification('Import successful')
          this.importFile = null
          this.importComment = ''
          this.$refs.fileInput.value = ''
        } else {
          this.notificationType = 'error'
          const errorMsg =
            response.data && response.data.error ? response.data.error : 'Import failed'
          this.showNotification(errorMsg)
        }
      } catch (error) {
        console.error('Import entries request failed', error)
        this.notificationType = 'error'
        this.showNotification(error?.response?.data?.error || 'Import failed')
      } finally {
        this.importInProgress = false
      }
    },
    async lookupDomain() {
      if (this.entry) {
        try {
          await this.ipSetStore.lookupDomain(this.entry)
          // Save domain for use as default comment
          this.resolvedDomain = this.entry
          // If comment is empty — use domain in comment field
          if (!this.comment) {
            this.comment = this.entry
          }
        } catch (error) {
          console.error('Domain lookup request failed', error)
        }
      }
    },
    async addAllEntries() {
      const comment = this.comment || this.resolvedDomain || ''
      const promises = this.ipAddresses.map((ip) =>
        this.ipSetStore.addEntry(this.selectedSet, {
          entry: ip,
          comment,
        }),
      )
      let successCount = 0

      try {
        const results = await Promise.allSettled(promises)
        results.forEach((result) => {
          if (result.status === 'fulfilled' && result.value.status === 200) {
            successCount++
          }
        })

        this.ipSetStore.ipAddresses = [] // Clear the IP addresses list after adding all

        if (successCount > 0) {
          this.notificationType = 'success'
          this.showNotification('All IP addresses successfully added')
        } else {
          this.notificationType = 'error'
          this.showNotification('Error adding all IP addresses')
        }
      } catch (error) {
        console.error('Add all entries request failed', error)
        this.notificationType = 'error'
        this.showNotification('Error adding all IP addresses')
      }
    },
    isDomain(entry) {
      // Simple regex to check if the entry is a domain name
      const domainPattern = /^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
      return domainPattern.test(entry)
    },
    showNotification(message) {
      this.notificationMessage = message
      setTimeout(() => {
        this.notificationMessage = ''
      }, 3000)
    },
    checkIp(ip) {
      this.$emit('check-ip', ip)
    },
  },
  setup() {
    const ipSetStore = useIpSetStore()

    // Watcher for entry field
    watch(
      () => ipSetStore.entry,
      (newEntry, oldEntry) => {
        if (oldEntry && ipSetStore.isDomain(oldEntry) && newEntry !== oldEntry) {
          ipSetStore.ipAddresses = []
        }
      },
    )

    return { ipSetStore }
  },
  watch: {
    entry(newEntry, oldEntry) {
      if (this.isDomain(oldEntry) && newEntry !== oldEntry) {
        this.ipSetStore.ipAddresses = []
      }
    },
  },
}
</script>

<style scoped>
.add-entry {
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #ffffff;
}

.tabs {
  display: flex;
  margin-bottom: 20px;
}

.tabs button {
  flex: 1;
  padding: 10px 0;
  border: none;
  background: #efefef;
  color: #333;
  font-size: 16px;
  cursor: pointer;
  border-bottom: 2px solid #ccc;
  transition:
    background 0.2s,
    border-bottom 0.2s;
}

.tabs button.active {
  background: #fdfdfd;
  border-bottom: 2px solid #28a745;
  font-weight: bold;
}

input[type='file'] {
  display: none;
}

.custom-file-input {
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
  background-color: #fff;
  margin-bottom: 0;
  margin-right: 0;
  width: 100%;
  cursor: pointer;
  color: #888;
  transition: border-color 0.2s;
  user-select: none;
  box-sizing: border-box;
}
.custom-file-input:hover {
  border-color: #28a745;
  color: #333;
}
.custom-file-input.disabled {
  background-color: #f5f5f5;
  color: #bbb;
  cursor: not-allowed;
  border-color: #eee;
}

button[disabled] {
  background: #ccc !important;
  cursor: not-allowed !important;
}

.input-group {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.input-group select,
.input-group input {
  padding: 8px;
  margin-right: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.input-group input[type='text'] {
  width: 100%;
  margin-right: 0;
  box-sizing: border-box;
}

.input-group select {
  flex: 1;
}

.input-group input {
  flex: 2;
}

.input-group button {
  padding: 8px 16px;
  border: none;
  background-color: #28a745;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.input-group button:hover {
  background-color: #218838;
}

.ip-list {
  max-height: 200px;
  overflow-y: auto;
  margin-top: 20px;
}

.ip-entry {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-bottom: 1px solid #ccc;
}

.ip-entry span {
  flex: 1;
  text-align: left;
}

.ip-entry button {
  padding: 4px 8px;
  margin: 0 2px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  background-color: #28a745;
  color: white;
}

.ip-entry button:hover {
  background-color: #218838;
}
.import-btn-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.import-btn-wrapper button {
  padding: 8px 16px;
  border: none;
  background-color: #28a745;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}
.import-btn-wrapper button:hover {
  background-color: #218838;
}
</style>
