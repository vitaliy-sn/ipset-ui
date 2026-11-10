<template>
  <div class="progress">
    <div class="progress-bar">
      <div class="progress-fill" :style="{ width: percent + '%', backgroundColor: color }"></div>
    </div>
    <span class="progress-text">{{ processed }} / {{ total }} ({{ percent }}%)</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  processed: { type: Number, default: 0 },
  total: { type: Number, default: 0 },
  color: { type: String, default: '#28a745' },
})

const percent = computed(() => {
  if (props.total <= 0) return 0
  return Math.min(100, Math.round((props.processed / props.total) * 100))
})
</script>

<style scoped>
.progress {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}
.progress-bar {
  flex: 1;
  height: 8px;
  background-color: #eee;
  border-radius: 4px;
  overflow: hidden;
}
.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.2s;
}
.progress-text {
  flex-shrink: 0;
  font-size: 14px;
  color: #666;
  white-space: nowrap;
}
</style>
