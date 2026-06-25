<script setup>

</script>

<template>
<thead>
    <tr>
  <th>
    <input
      type="checkbox"
      :checked="isAllSelected"
      @change="toggleSelectAll"
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