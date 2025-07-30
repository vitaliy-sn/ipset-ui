<template>
  <div ref="modal" class="modal" v-if="isVisible" @keydown.enter="confirmDelete" tabindex="0">
    <div class="modal-content">
      <span class="close" @click="closeModal">&times;</span>
      <p>Are you sure you want to delete {{ entry }}?</p>
      <div class="button-group">
        <button @click="confirmDelete">Delete</button>
        <button @click="closeModal">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DeleteConfirmationModal',
  props: {
    isVisible: {
      type: Boolean,
      required: true,
    },
    entry: {
      type: String,
      required: true,
    },
  },
  methods: {
    handleKeydown(event) {
      if (event.key === 'Enter' && this.isVisible) {
        this.confirmDelete()
      }
    },
    closeModal() {
      this.$emit('close')
    },
    confirmDelete() {
      this.$emit('confirm')
    },
  },
  mounted() {
    if (this.isVisible) {
      this.$refs.modal.focus()
    }
  },
  watch: {
    isVisible(newVal) {
      if (newVal) {
        this.$nextTick(() => {
          this.$refs.modal.focus()
        })
      }
    },
  },
}
</script>

<style scoped>
.modal {
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  z-index: 1000;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  overflow: auto;
  background-color: rgba(0, 0, 0, 0.5);
  outline: none; /* Remove default outline */
}

.modal-content {
  background-color: #fff;
  padding: 20px;
  border: 1px solid #888;
  width: 300px;
  border-radius: 8px;
  text-align: center;
}

.close {
  color: #aaa;
  float: right;
  font-size: 28px;
  font-weight: bold;
}

.close:hover,
.close:focus {
  color: black;
  text-decoration: none;
  cursor: pointer;
}

.button-group {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}

.button-group button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.button-group button:first-of-type {
  background-color: #dc3545;
  color: white;
}

.button-group button:first-of-type:hover {
  background-color: #c82333;
}

.button-group button:last-of-type {
  background-color: #6c757d;
  color: white;
}

.button-group button:last-of-type:hover {
  background-color: #5a6268;
}
</style>
