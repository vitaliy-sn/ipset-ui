<template>
  <div class="sidebar">
    <h3>IP sets</h3>
    <ul class="set-list">
      <li
        v-for="(set, index) in ipSets"
        :key="set"
        :class="{
          selected: set === selectedSet,
          'first-item': index === 0,
          'last-item': index === ipSets.length - 1,
        }"
        @click="selectSet(set)"
      >
        {{ set }}
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'HomeSetList',
  props: {
    ipSets: {
      type: Array,
      required: true,
    },
    selectedSet: {
      type: String,
      required: true,
    },
  },
  methods: {
    selectSet(set) {
      this.$emit('select-set', set)
    },
  },
}
</script>

<style scoped>
.sidebar {
  width: 200px;
  margin-right: 20px;
  padding: 15px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #ffffff;
}

.sidebar h3 {
  margin-top: 0;
  font-size: 1.2em;
  color: #333;
  text-align: center;
  padding: 20px;
}

.set-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.set-list li {
  padding: 8px 20px;
  margin: 0;
  cursor: pointer;
  transition:
    background-color 0.3s,
    color 0.3s,
    border-color 0.3s;
  border-bottom: 1px solid #ccc;
  background-color: whitesmoke;
}

.set-list li:first-child {
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
}

.set-list li:last-child {
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
  border-bottom: none;
}

.set-list li:hover {
  background-color: #f0f0f0;
  border-color: #bbb;
}

.set-list li.selected {
  background-color: #dcdcdc;
  border-color: #dcdcdc;
}
</style>
