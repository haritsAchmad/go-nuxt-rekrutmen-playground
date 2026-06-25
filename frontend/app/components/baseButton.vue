<script setup>
defineProps({
  label: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'button'
  }
})
</script>

<template>
  <button :type="type">
    {{ label }}
  </button>
</template>