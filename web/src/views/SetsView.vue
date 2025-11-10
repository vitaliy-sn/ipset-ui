<template>
  <div class="main-container">
    <div class="header-wrapper">
      <div class="header-container">
        <AppTabs />
      </div>
    </div>
    <div class="content-wrapper">
      <SetList @toggle-backup="toggleBackup" />
      <SetBackup v-if="selectedSet" :setName="selectedSet" key="selectedSet" />
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import SetList from '../components/SetList.vue'
import SetBackup from '../components/SetBackup.vue'
import AppTabs from '../components/AppTabs.vue'

export default {
  components: {
    SetList,
    SetBackup,
    AppTabs,
  },
  setup() {
    const router = useRouter()

    const selectedSet = ref('')

    const goHome = () => {
      router.push('/')
    }

    const toggleBackup = (set) => {
      if (selectedSet.value === set) {
        selectedSet.value = ''
      } else {
        selectedSet.value = set
      }
    }

    return {
      goHome,
      selectedSet,
      toggleBackup,
    }
  },
}
</script>

<style scoped>
@import '../assets/common.css';
</style>
