<script setup>
defineProps({
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  modelValue: {
    type: String,
    default: ''
  },
type: {
  type: String,
  default: 'text'
}
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div style="margin-bottom: 8px;">
    <label>{{ label }}</label>
    <br>

    <input
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      @input="$emit('update:modelValue', $event.target.value)"
    >
  </div>
</template>