<script setup>
const props = defineProps({
  label: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  type: { type: String, default: 'text' },
  modelValue: { type: [String, Number, File, null], default: '' },

  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  readonly: { type: Boolean, default: false },

  maxlength: { type: [String, Number], default: null },
  min: { type: [String, Number], default: null },
  max: { type: [String, Number], default: null },

  accept: { type: String, default: '' },
  numberOnly: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue', 'change'])

function blockInvalidNumber(event) {
  if (!props.numberOnly) return

  if (['e', 'E', '+', '-', ' '].includes(event.key)) {
    event.preventDefault()
  }
}

function handleInput(event) {
  if (props.type === 'file') {
    const file = event.target.files?.[0] || null
    emit('update:modelValue', file)
    emit('change', file)
    return
  }

  emit('update:modelValue', event.target.value)
}
</script>

<template>
  <div style="margin-bottom: 8px;">
    <label>{{ label }}</label>
    <br>

    <input
      :type="type"
      :value="type === 'file' ? undefined : modelValue"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :readonly="readonly"
      :maxlength="maxlength"
      :min="min"
      :max="max"
      :accept="accept"
      @keydown="blockInvalidNumber"
      @input="handleInput"
      @change="type === 'file' ? handleInput($event) : $emit('change', $event)"
    >
  </div>
</template>