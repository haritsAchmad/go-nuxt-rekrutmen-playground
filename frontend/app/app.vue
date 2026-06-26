<script setup>
const form = reactive({
  judul: '',
  unit: '',
  status: 'aktif',
jumlah: ''
})

const filter = reactive({
  keyword: '',
  status: ''
})

const paginationMeta = computed(() => {
  return data.value?.data?.meta || {
    page: 1,
    limit: 10,
    total: 0,
    total_page: 1
  }
})

const editId = ref(null)
const selectedIds = ref([])
const selectedLowongan = ref(null)

const bulkLoading = ref(false)
const bulkMessage = ref('')
const bulkError = ref('')

const page = ref(1)
const limit = ref(10)

const isAllSelected = computed(() => {
  return lowonganList.value.length > 0 &&
    selectedIds.value.length === lowonganList.value.length
})

const apiUrl = computed(() => {
  const params = new URLSearchParams()

  if (filter.keyword) {
    params.append('keyword', filter.keyword)
  }

  if (filter.status) {
    params.append('status', filter.status)
  }

  params.append('page', page.value)
  params.append('limit', limit.value)

  return `http://localhost:8080/api/lowongan?${params.toString()}`
})

const { data, pending, error, refresh } = await useFetch(apiUrl)

async function resetFilter() {
  filter.keyword = ''
  filter.status = ''
  page.value = 1
  selectedIds.value = []
  await refresh()
}

async function applyFilter() {
  page.value = 1
  selectedIds.value = []
  await refresh()
}

async function showDetail(id) {
  const result = await $fetch(`http://localhost:8080/api/lowongan/detail?id=${id}`)
  selectedLowongan.value = result.data
}

const lowonganList = computed(() => {
  return data.value?.data?.data || []
})

async function submitLowongan() {
  if (editId.value) {
    await $fetch(`http://localhost:8080/api/lowongan?id=${editId.value}`, {
      method: 'PUT',
      body: {
        judul: form.judul,
        unit: form.unit,
        status: form.status || 'aktif'
      }
    })
  } else {
    await $fetch('http://localhost:8080/api/lowongan', {
      method: 'POST',
      body: {
        judul: form.judul,
        unit: form.unit
      }
    })
  }

  resetForm()
  await refresh()
}

async function deleteLowongan(id) {
  const yakin = confirm('Yakin mau hapus lowongan ini?')

  if (!yakin) {
    return
  }

  await $fetch(`http://localhost:8080/api/lowongan?id=${id}`, {
    method: 'DELETE'
  })

  await refresh()
}

async function toggleStatus(lowongan) {
  const nextStatus = lowongan.status === 'aktif' ? 'nonaktif' : 'aktif'

  await $fetch(`http://localhost:8080/api/lowongan/status?id=${lowongan.id}`, {
    method: 'PUT',
    body: {
      status: nextStatus
    }
  })

  await refresh()
}

async function bulkUpdateStatus(status) {
  if (selectedIds.value.length === 0) {
    bulkError.value = 'Pilih minimal satu lowongan'
    bulkMessage.value = ''
    return
  }

  bulkLoading.value = true
  bulkError.value = ''
  bulkMessage.value = ''

  try {
    await $fetch('http://localhost:8080/api/lowongan/bulk-status', {
      method: 'PUT',
      body: {
        ids: selectedIds.value,
        status: status
      }
    })

    selectedIds.value = []
    bulkMessage.value = 'Status lowongan terpilih berhasil diubah'
    await refresh()
  } catch (error) {
    bulkError.value = 'Gagal mengubah status lowongan'
  } finally {
    bulkLoading.value = false
  }
}

async function bulkDelete() {
  if (selectedIds.value.length === 0) {
    bulkError.value = 'Pilih minimal satu lowongan'
    bulkMessage.value = ''
    return
  }

  const yakin = confirm('Yakin mau hapus lowongan yang dipilih?')

  if (!yakin) {
    return
  }

  bulkLoading.value = true
  bulkError.value = ''
  bulkMessage.value = ''

  try {
    await $fetch('http://localhost:8080/api/lowongan/bulk-delete', {
      method: 'DELETE',
      body: {
        ids: selectedIds.value
      }
    })

    selectedIds.value = []
    bulkMessage.value = 'Lowongan terpilih berhasil dihapus'
    await refresh()
  } catch (error) {
    bulkError.value = 'Gagal menghapus lowongan'
  } finally {
    bulkLoading.value = false
  }
}

async function nextPage() {
  if (page.value >= paginationMeta.value.total_page) {
    return
  }

  page.value++
  selectedIds.value = []
  await refresh()
}

async function previousPage() {
  if (page.value <= 1) {
    return
  }

  page.value--
  selectedIds.value = []
  await refresh()
}

function editLowongan(lowongan) {
  editId.value = lowongan.id
  form.judul = lowongan.judul
  form.unit = lowongan.unit
  form.status = lowongan.status
}

function resetForm() {
  editId.value = null
  form.judul = ''
  form.unit = ''
  form.status = 'aktif'
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value = []
    return
  }

  selectedIds.value = lowonganList.value.map((lowongan) => lowongan.id)
}
</script>

<template>
  <main style="padding: 24px; font-family: Arial, sans-serif;">
    <h1>Lowongan Rekrutmen</h1>

    <form @submit.prevent="submitLowongan" style="margin-bottom: 24px;">
      <BaseInput
  v-model="form.judul"
  label="Judul Lowongan"
  required
  maxlength="100"
/>

      <BaseInput
    label="Unit"
    placeholder="Contoh: Direktorat SDM"
    v-model="form.unit"
/>

<BaseInput
  v-model="form.jumlah"
  label="Jumlah Dibutuhkan"
  type="number"
  :min="1"
  :max="100"
  placeholder="Contoh: 3"
  number-only
/>

      	<BaseButton
  type="submit"
  :label="editId ? 'Simpan Perubahan' : 'Tambah Lowongan'"
  :color="editId ? 'blue' : 'green'"
/>

	<button v-if="editId" type="button" @click="resetForm">
  	Batal Edit
	</button>
    </form>

    <div style="margin-bottom: 16px;">
  <h3>Filter Lowongan</h3>

  <input
    v-model="filter.keyword"
    type="text"
    placeholder="Cari judul/unit"
  >

  <select v-model="filter.status">
    <option value="">Semua Status</option>
    <option value="aktif">Aktif</option>
    <option value="nonaktif">Nonaktif</option>
  </select>

  <button type="button" @click="applyFilter">
  Cari
</button>

  <button type="button" @click="resetFilter">
    Reset
  </button>
</div>

<div style="margin-bottom: 16px;">
  <h3>Bulk Action</h3>

  <p>Terpilih: {{ selectedIds.length }} lowongan</p>

<p v-if="bulkMessage" style="color: green;">
  {{ bulkMessage }}
</p>

<p v-if="bulkError" style="color: red;">
  {{ bulkError }}
</p>

  <button type="button" @click="bulkUpdateStatus('aktif')">
    Aktifkan Terpilih
  </button>

  <button type="button" @click="bulkUpdateStatus('nonaktif')">
    Nonaktifkan Terpilih
  </button>

  <button type="button" @click="bulkDelete">
    Hapus Terpilih
  </button>
</div>

<p v-if="pending">Loading...</p>
<p v-else-if="error">Gagal ambil data lowongan</p>

<table v-if="!pending && !error" border="1" cellpadding="8" cellspacing="0">

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

  <tbody>
    <tr v-for="lowongan in lowonganList" :key="lowongan.id">
<td>
  <input
    v-model="selectedIds"
    type="checkbox"
    :value="lowongan.id"
  >
</td>
      <td>{{ lowongan.id }}</td>
      <td>{{ lowongan.judul }}</td>
      <td>{{ lowongan.unit }}</td>
      <td>{{ lowongan.status }}</td>
      <td>
        <button type="button" @click="editLowongan(lowongan)">
          Edit
        </button>

	<button type="button" @click="showDetail(lowongan.id)">
  	  Detail
	</button>

        <button type="button" @click="toggleStatus(lowongan)">
          {{ lowongan.status === 'aktif' ? 'Nonaktifkan' : 'Aktifkan' }}
        </button>

        <button type="button" @click="deleteLowongan(lowongan.id)">
          Hapus
        </button>
      </td>
    </tr>
  </tbody>
</table>

<div style="margin-top: 16px;">
  <button
    type="button"
    :disabled="page <= 1 || pending"
    @click="previousPage"
  >
    Previous
  </button>

  <span style="margin: 0 12px;">
    Page {{ paginationMeta.page }} of {{ paginationMeta.total_page }}
    — Total {{ paginationMeta.total }} data
  </span>

  <button
    type="button"
    :disabled="page >= paginationMeta.total_page || pending"
    @click="nextPage"
  >
    Next
  </button>
</div>

<div v-if="selectedLowongan" style="margin-top: 16px;">
  <h3>Detail Lowongan</h3>
  <p>ID: {{ selectedLowongan.id }}</p>
  <p>Judul: {{ selectedLowongan.judul }}</p>
  <p>Unit: {{ selectedLowongan.unit }}</p>
  <p>Status: {{ selectedLowongan.status }}</p>

  <button type="button" @click="selectedLowongan = null">
    Tutup Detail
  </button>
</div>

  </main>
</template>