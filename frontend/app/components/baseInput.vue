<script setup>
defineProps({
  label: {
    type: String,
    default: ''
  },

  placeholder: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <div style="margin-bottom: 8px;">
    <label>{{ label }}</label>
    <br>

    <input
      type="text"
      :placeholder="placeholder"
    >
  </div>
</template>