<template>
  <div class="tabs">
    <button
      v-for="tab in tabList"
      :key="tab.value"
      :class="{ active: activeTab === tab.value }"
      @click="switchTab(tab)"
    >
      {{ tab.label }}
    </button>
  </div>
</template>

<script>
export default {
  name: 'AppTabs',
  emits: ['tab-change'],
  data() {
    return {
      tabList: [
        { value: 'entries', label: 'Entries', route: '/entries' },
        { value: 'sets', label: 'Sets', route: '/sets' },
      ],
      activeTab: 'entries',
    }
  },
  methods: {
    switchTab(tab) {
      if (this.activeTab !== tab.value) {
        this.activeTab = tab.value
        this.$emit('tab-change', tab.value)
        this.$router.push(tab.route)
      }
    },
    syncTabWithRoute() {
      const currentRoute = this.$route.path
      const found = this.tabList.find((tab) => tab.route === currentRoute)
      if (found) {
        this.activeTab = found.value
        this.$emit('tab-change', found.value)
      }
    },
  },
  mounted() {
    this.syncTabWithRoute()
    this.$watch(
      () => this.$route.path,
      () => {
        this.syncTabWithRoute()
      },
    )
  },
}
</script>

<style scoped>
.tabs {
  display: flex;
  /*margin-bottom: 20px;*/
}

.tabs button {
  flex: 1;
  min-width: 120px;
  padding: 10px 20px;
  border: none;
  background: #efefef;
  color: #333;
  font-size: 16px;
  cursor: pointer;
  border-bottom: 2px solid #ccc;
  transition:
    background 0.2s,
    border-bottom 0.2s;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.tabs button.active {
  background: #fdfdfd;
  border-bottom: 2px solid #28a745;
  font-weight: bold;
}
</style>
