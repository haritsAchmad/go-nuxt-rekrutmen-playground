<script setup>
const props = defineProps({
  label: String,
  type: {
    type: String,
    default: 'button'
  },
  color: {
    type: String,
    default: 'default'
  }
})

const buttonClass = computed(() => {
  return `btn btn-${props.color}`
})
</script>

<template>
  <button :type="type" :class="buttonClass">
    {{ label }}
  </button>
</template>

<style scoped>
.btn {
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid #ccc;
}

.btn-green {
  background: #16a34a;
  color: white;
}

.btn-blue {
  background: #2563eb;
  color: white;
}

.btn-red {
  background: #dc2626;
  color: white;
}

.btn-default {
  background: #e5e7eb;
  color: black;
}
</style>