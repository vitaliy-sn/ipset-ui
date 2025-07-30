<template>
  <transition name="fade">
    <div v-if="isVisible" :class="['notification', type]">
      <span class="message">{{ message }}</span>
    </div>
  </transition>
</template>

<script>
export default {
  name: 'Notification',
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
  data() {
    return {
      isVisible: true,
    }
  },
  mounted() {
    setTimeout(() => {
      this.isVisible = false
    }, 3000)
  },
}
</script>

<style scoped>
.notification {
  position: fixed;
  top: 20px;
  left: 20px;
  display: flex;
  align-items: center;
  padding: 16px 24px;
  background-color: #323232;
  border-radius: 4px;
  box-shadow:
    0 2px 4px rgba(0, 0, 0, 0.2),
    0 4px 8px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  min-width: 300px;
  max-width: 400px;
}

.notification.success .message {
  color: #4caf50;
}

.notification.error .message {
  color: #f44336;
}

.message {
  flex: 1;
  font-size: 16px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 2s;
}

.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>
