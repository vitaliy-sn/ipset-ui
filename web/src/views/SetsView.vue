<template>
  <div class="main-container">
    <div class="content-wrapper">
      <SetList @toggle-backup="toggleBackup" />
      <SetBackup v-if="selectedSet" :setName="selectedSet" key="selectedSet" />
    </div>
    <button class="floating-button" @click="goHome">Back</button>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import SetList from '../components/SetList.vue'
import SetBackup from '../components/SetBackup.vue'

export default {
  components: {
    SetList,
    SetBackup,
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
.main-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
}

.content-wrapper {
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.floating-button {
  position: fixed;
  top: 10px;
  right: 20px;
  padding: 10px 20px;
  border-radius: 4px;
  background-color: #28a745;
  color: white;
  border: none;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.floating-button:hover {
  background-color: #218838;
}
</style>
