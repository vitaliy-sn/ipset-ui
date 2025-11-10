<template>
  <div v-if="visible" class="modal">
    <div class="modal-content">
      <p>
        <slot name="message">
          {{ message }}
        </slot>
      </p>
      <div class="modal-actions">
        <button @click="onConfirm" :disabled="loading" class="confirm-btn">Yes</button>
        <button @click="onCancel" :disabled="loading" class="cancel-btn">No</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  visible: { type: Boolean, required: true },
  message: { type: String, default: '' },

  loading: { type: Boolean, default: false },
  onConfirm: { type: Function, required: true },
  onCancel: { type: Function, required: true },
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
.modal-actions {
  display: flex;
  flex-direction: row;
  gap: 18px;
  justify-content: center;
  width: 100%;
  margin-top: 8px;
}
.confirm-btn {
  background: #28a745;
  color: #fff;
  border: none;
  padding: 8px 18px;
  border-radius: 6px;
  font-size: 1em;
  cursor: pointer;
  font-weight: 500;
  width: 110px;
  transition: background 0.2s;
}
.confirm-btn:hover {
  background: #218838;
}
.cancel-btn {
  background: #e0e0e0;
  color: #333;
  border: none;
  padding: 8px 18px;
  border-radius: 6px;
  font-size: 1em;
  cursor: pointer;
  font-weight: 500;
  width: 110px;
  transition: background 0.2s;
}
.cancel-btn:hover {
  background: #bdbdbd;
}
</style>
