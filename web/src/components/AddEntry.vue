<template>
  <div class="add-entry">
    <h3>Add New Entry</h3>
    <div class="input-group">
      <div
        class="file-drop"
        :class="{
          empty: !importFile,
          'is-dragover': isDragOver,
          disabled: importInProgress || entry !== '',
        }"
        @click="triggerFileInput"
        @dragover.prevent="onDragOver"
        @dragleave="onDragLeave"
        @drop.prevent="onDrop"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
          <polyline points="14 2 14 8 20 8"></polyline>
        </svg>
        <span v-if="!importFile" class="file-hint">
          <span class="file-hint-main">Select file</span>
          <span class="file-hint-sub">or drag &amp; drop</span>
        </span>
        <template v-else>
          <span class="file-name" :title="importFile.name">{{ importFile.name }}</span>
          <span class="file-size">{{ formatFileSize(importFile.size) }}</span>
          <button class="clear" title="Remove file" @click.stop="clearImportFile">×</button>
        </template>
        <input
          type="file"
          ref="fileInput"
          @change="onFileChange"
          style="display: none"
          :disabled="importInProgress || entry !== ''"
        />
      </div>
    </div>
    <div class="input-group">
      <div class="input-wrapper">
        <input
          v-model="entry"
          type="text"
          placeholder="IP / CIDR / Domain"
          title="Multiple IPs or CIDRs can be added at once — separate with a space or comma."
          :disabled="importInProgress || importFile !== null"
        />
        <button v-if="entry && !importFile" class="clear" title="Clear" @click="clearEntry">×</button>
      </div>
    </div>
    <div class="input-group">
      <div class="input-wrapper">
        <input
          v-model="comment"
          type="text"
          placeholder="Comment (optional)"
          :disabled="importInProgress"
        />
        <button v-if="comment" class="clear" title="Clear" @click="clearComment">×</button>
      </div>
    </div>
    <div class="import-btn-wrapper">
      <button v-if="importFile" :disabled="importInProgress" @click.prevent="importEntries">Import</button>
      <button
        v-else-if="entry && isDomain(entry) && ipAddresses.length === 0"
        @click.prevent="lookupDomain"
      >
        Lookup
      </button>
      <button
        v-else-if="entry && isDomain(entry) && ipAddresses.length > 0"
        @click.prevent="addAllEntries"
      >
        Add All
      </button>
      <button v-else :disabled="!entry || importInProgress" @click.prevent="addEntry(entry)">Add</button>
    </div>
    <div class="ip-list" v-if="ipAddresses.length > 0">
      <div v-for="(ip, index) in ipAddresses" :key="index" class="ip-entry">
        <span>{{ ip }}</span>
        <button @click.prevent="checkIp(ip)">Check</button>
        <button @click.prevent="addEntry(ip, comment)">Add</button>
      </div>
    </div>
    <AppNotification
      v-if="notificationMessage"
      :message="notificationMessage"
      :type="notificationType"
      @close="notificationMessage = ''"
    />
  </div>
</template>

<script>
import { useIpSetStore } from '../stores/ipsetStore'
import AppNotification from './AppNotification.vue'

export default {
  name: 'AddEntry',
  components: {
    AppNotification,
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
      importFile: null, // Selected file for import
      resolvedDomain: '', // Save domain after lookup
      importInProgress: false, // Flag to disable import button
      isDragOver: false, // Drag & drop highlight state
    }
  },
  computed: {
    ipAddresses() {
      return this.ipSetStore.ipAddresses
    },
  },
  methods: {
    clearEntry() {
      this.entry = ''
    },
    clearComment() {
      this.comment = ''
    },
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
            const parsed = String(ip).split(/[\s,]+/).filter(Boolean)
            if (parsed.length > 1) {
              const added =
                typeof response.data?.added === 'number' ? response.data.added : parsed.length
              if (added < parsed.length) {
                const failed = response.data?.failed || []
                const failedPart = failed.length ? ` (failed: ${failed.join(', ')})` : ''
                this.notificationType = 'error'
                this.showNotification(`Added ${added} of ${parsed.length} entries${failedPart}`)
              } else {
                this.notificationType = 'success'
                this.showNotification(`${added} entries successfully added`)
              }
              return
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
      if (files && files.length > 0 && !this.entry) {
        this.setImportFile(files[0])
      } else {
        this.importFile = null
      }
      // Reset so re-selecting the same file fires change again
      event.target.value = ''
    },
    setImportFile(file) {
      this.importFile = file
      // If comment is empty — use file name as comment
      if (!this.comment) {
        this.comment = file.name
      }
    },
    clearImportFile() {
      if (this.importFile && this.comment === this.importFile.name) {
        this.comment = ''
      }
      this.importFile = null
      if (this.$refs.fileInput) {
        this.$refs.fileInput.value = ''
      }
    },
    onDragOver() {
      if (this.importInProgress || this.entry) return
      this.isDragOver = true
    },
    onDragLeave(event) {
      if (event.currentTarget.contains(event.relatedTarget)) return
      this.isDragOver = false
    },
    onDrop(event) {
      if (this.importInProgress || this.entry) return
      this.isDragOver = false
      const file = event.dataTransfer && event.dataTransfer.files[0]
      if (file) {
        this.setImportFile(file)
      }
    },
    formatFileSize(bytes) {
      if (bytes < 1024) return `${bytes} B`
      if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
      return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
    },
    triggerFileInput() {
      if (this.importInProgress || this.entry) return
      this.$refs.fileInput.click()
    },
    async importEntries() {
      if (!this.importFile || !this.selectedSet) return
      this.importInProgress = true // Immediately disable the button
      const formData = new FormData()
      formData.append('file', this.importFile)
      if (this.comment) {
        formData.append('comment', this.comment)
      }
      try {
        const response = await this.ipSetStore.importFromFile(this.selectedSet, formData)
        if (response.status === 200) {
          this.notificationType = 'success'
          this.showNotification('Import successful')
          this.importFile = null
          this.comment = ''
          if (this.$refs.fileInput) {
            this.$refs.fileInput.value = ''
          }
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
      if (this.notificationType === 'success') {
        setTimeout(() => {
          this.notificationMessage = ''
        }, 3000)
      }
    },
    checkIp(ip) {
      this.$emit('check-ip', ip)
    },
  },
  setup() {
    return { ipSetStore: useIpSetStore() }
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
  border-radius: 10px;
  box-shadow: 0 2px 8px #0000000f;
  background-color: #ffffff;
}

.add-entry h3 {
  margin: 0 0 15px;
  font-size: 1.2em;
  color: #333;
  text-align: center;
}

input[type='file'] {
  display: none;
}

.file-drop {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  min-height: 110px;
  border: 1px dashed #ccc;
  border-radius: 4px;
  font-size: 16px;
  background-color: #fff;
  cursor: pointer;
  color: #333;
  transition:
    border-color 0.2s,
    background-color 0.2s;
  user-select: none;
  box-sizing: border-box;
  width: 100%;
}

.file-drop:hover {
  border-color: #28a745;
}

.file-drop.empty {
  flex-direction: column;
  justify-content: center;
  gap: 8px;
}

.file-drop.empty svg {
  width: 24px;
  height: 24px;
}

.file-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  line-height: 1.3;
}

.file-hint-main {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.file-hint-sub {
  font-size: 14px;
  color: #888;
}

.file-drop.is-dragover {
  border-color: #28a745;
  background-color: #f0faf3;
}

.file-drop.disabled {
  background-color: #f5f5f5;
  color: #333;
  opacity: 0.6;
  cursor: not-allowed;
  border-color: #eee;
}

.file-drop svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  color: #666;
}

.file-drop .file-name {
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
}

.file-drop .file-size {
  flex-shrink: 0;
  color: #666;
  font-size: 14px;
}

.file-drop .clear {
  flex-shrink: 0;
  padding: 2px 6px;
  border: none;
  background: none;
  color: #666;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
}

.file-drop .clear:hover {
  background: none;
  color: #333;
}

button[disabled] {
  background: #eee !important;
  color: #888 !important;
  opacity: 0.7;
  cursor: not-allowed !important;
}

.input-group {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.input-group select {
  flex: 1;
  padding: 8px;
  margin-right: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.input-wrapper {
  position: relative;
  flex: 1;
  min-width: 0;
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
