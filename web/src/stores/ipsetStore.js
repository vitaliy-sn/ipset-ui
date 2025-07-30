import { defineStore } from 'pinia'
import axios from '../axios'

export const useIpSetStore = defineStore('ipset', {
  state: () => ({
    ipEntries: [],
    errorMessage: '',
    ipSets: [],
    ipAddresses: [],
    backups: [],
  }),
  actions: {
    async searchIpSet(setName, filter) {
      try {
        const response = await axios.post(`/ipsets/${setName}/entries/search`, { filter })
        this.ipEntries = response.data
        this.errorMessage = ''
      } catch (error) {
        console.error('Search IP set request failed', error)
        this.ipEntries = []
        this.errorMessage = 'Not Found'
      }
    },
    async fetchIpSets() {
      try {
        const response = await axios.get('/ipsets')
        this.ipSets = response.data
      } catch (error) {
        console.error('Fetch IP sets request failed', error)
        this.ipSets = []
        this.errorMessage = 'Failed to fetch IP sets'
      }
    },
    async addEntry(setName, entry) {
      try {
        // If entry is a string, wrap it; if it's an object, send as is
        const payload = typeof entry === 'object' ? entry : { entry }
        const response = await axios.post(`/ipsets/${setName}/entries`, payload)
        return response
      } catch (error) {
        console.error('Add entry request failed', error)
        throw error
      }
    },
    async deleteEntry(setName, entry) {
      try {
        const response = await axios.delete(`/ipsets/${setName}/entries`, { data: { entry } })
        return response.data.message
      } catch (error) {
        console.error('Delete entry request failed', error)
        throw error
      }
    },
    async lookupDomain(domain) {
      try {
        const response = await axios.post('/dns-lookup', { domain })
        this.ipAddresses = response.data
      } catch (error) {
        console.error('Domain lookup request failed', error)
        this.ipAddresses = []
      }
    },
    async addSet(setName) {
      try {
        const response = await axios.post('/ipsets', { setName })
        this.ipSets.push(setName)
        return response
      } catch (error) {
        console.error('Add set request failed', error)
        throw error
      }
    },
    async deleteSet(setName) {
      try {
        const response = await axios.delete(`/ipsets/${setName}`)
        this.ipSets = this.ipSets.filter((set) => set !== setName)
        return response.data.message
      } catch (error) {
        console.error('Delete set request failed', error)
        throw error
      }
    },
    async fetchBackups(setName) {
      try {
        const response = await axios.get(`/ipsets/${setName}/backups`)
        this.backups = response.data
      } catch (error) {
        console.error('Failed to fetch backups', error)
        this.backups = []
      }
    },
    async createBackup(setName, fileNamePart) {
      try {
        const response = await axios.post(`/ipsets/${setName}/save`, { setName, fileNamePart })
        if (!this.backups) {
          this.backups = []
        }
        this.backups.push(response.data)
        return response
      } catch (error) {
        console.error('Failed to create backup', error)
        throw error
      }
    },
    async restoreBackup(setName, fileNamePart) {
      try {
        const response = await axios.post(`/ipsets/${setName}/restore`, { setName, fileNamePart })
        return response
      } catch (error) {
        console.error('Failed to restore backup', error)
        throw error
      }
    },
    async deleteBackup(setName, fileNamePart) {
      try {
        const response = await axios.delete(`/ipsets/${setName}/backups`, {
          data: { fileNamePart },
        })
        this.backups = this.backups.filter((backup) => backup !== fileNamePart)
        return response
      } catch (error) {
        console.error('Failed to delete backup', error)
        throw error
      }
    },
    async importFromFile(setName, formData) {
      try {
        const response = await axios.post(`/ipsets/${setName}/entries/import`, formData, {
          headers: { 'Content-Type': 'multipart/form-data' },
        })
        return response
      } catch (error) {
        console.error('Import from file failed', error)
        throw error
      }
    },
  },
})
