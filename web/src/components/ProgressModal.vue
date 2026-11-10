<template>
  <div v-if="visible" class="modal">
    <div class="modal-content">
      <p>{{ message }}</p>
      <ProgressBar v-if="total > 0" :processed="processed" :total="total" :color="color" />
    </div>
  </div>
</template>

<script setup>
import ProgressBar from './ProgressBar.vue'

defineProps({
  visible: { type: Boolean, required: true },
  message: { type: String, default: '' },
  processed: { type: Number, default: 0 },
  total: { type: Number, default: 0 },
  color: { type: String, default: '#dc3545' },
})
</script>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal-content {
  background: #fff;
  padding: 24px;
  border-radius: 10px;
  min-width: 260px;
  max-width: 400px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.18);
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
  word-break: break-word;
}
</style>
