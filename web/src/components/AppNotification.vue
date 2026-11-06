<template>
  <transition name="fade">
    <div v-if="isVisible" :class="['notification', type]">
      <div class="header">
        <span class="title">{{ type === 'error' ? 'Error' : 'Success' }}</span>
        <button class="close" aria-label="Close" @click="close">&times;</button>
      </div>
      <div class="message">{{ message }}</div>
    </div>
  </transition>
</template>

<script>
export default {
  name: 'AppNotification',
  props: {
    message: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      required: true,
      validator: (value) => ['success', 'error'].includes(value),
    },
  },
  emits: ['close'],
  data() {
    return {
      isVisible: true,
      hideTimer: null,
    }
  },
  mounted() {
    if (this.type === 'success') {
      this.hideTimer = setTimeout(() => {
        this.isVisible = false
      }, 3000)
    }
  },
  beforeUnmount() {
    if (this.hideTimer) {
      clearTimeout(this.hideTimer)
      this.hideTimer = null
    }
  },
  methods: {
    close() {
      this.$emit('close')
    },
  },
}
</script>

<style scoped>
.notification {
  position: fixed;
  top: 20px;
  left: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  padding: 16px 24px;
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.18);
  z-index: 1000;
  min-width: 300px;
  max-width: 400px;
  max-height: calc(100vh - 40px);
  max-height: calc(100dvh - 40px);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  flex: 0 0 auto;
}

.title {
  font-size: 16px;
  font-weight: 600;
}

.notification.success .title {
  color: #28a745;
}

.notification.error .title {
  color: #d32f2f;
}

.message {
  font-size: 16px;
  color: #333;
  white-space: pre-line;
  overflow-y: auto;
  min-height: 0;
}

.close {
  padding: 4px;
  border: none;
  background: transparent;
  color: #9e9e9e;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
}

.close:hover {
  color: #333333;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 2s;
}

.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>
