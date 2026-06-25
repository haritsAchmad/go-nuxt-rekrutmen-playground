<template>
  <thead>
    <tr>
      <th>
        <input
          type="checkbox"
          :checked="isAllSelected"
          @change="$emit('toggle-select-all')"
        >
      </th>
      <th>ID</th>
      <th>Judul</th>
      <th>Unit</th>
      <th>Status</th>
      <th>Aksi</th>
    </tr>
  </thead>
</template>

<script setup>
defineProps({
  isAllSelected: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle-select-all'])
</script>