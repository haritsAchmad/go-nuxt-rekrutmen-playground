<script setup>
const props = defineProps({
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  modelValue: {
    type: [String, Number],
    default: ''
  },
  numberOnly: {
    type: Boolean,
    default: false
  }
})

defineEmits(['update:modelValue'])

function blockInvalidNumber(event) {
  if (!props.numberOnly) {
    return
  }

  if (['e', 'E', '+', '-', ' '].includes(event.key)) {
    event.preventDefault()
  }
}
</script>

<template>
  <div style="margin-bottom: 8px;">
    <label>{{ label }}</label>
    <br>

    <input
  :type="type"
  :value="modelValue"
  :placeholder="placeholder"
  @keydown="blockInvalidNumber"
  @input="$emit('update:modelValue', $event.target.value)"
>
  </div>
</template>