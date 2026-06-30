<script setup>
definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()

const currentUser = computed(() => authStore.user)
const authHeaders = computed(() => authStore.authHeaders)

const form = reactive({
  judul: '',
  unit: '',
  tanggalBuka: '',
  tanggalTutup: '',
  deskripsi: '',
  status: 'aktif'
})

const filter = reactive({
  keyword: '',
  status: ''
})

const editId = ref(null)
const selectedIds = ref([])
const selectedLowongan = ref(null)

const bulkLoading = ref(false)
const bulkMessage = ref('')
const bulkError = ref('')
const actionMessage = ref('')
const actionError = ref('')

const page = ref(1)
const limit = ref(10)
const today = new Date().toISOString().slice(0, 10)

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

const { data, pending, error, refresh } = await useFetch(apiUrl, {
  immediate: false,
  server: false,
  watch: false,
  headers: authHeaders
})

const paginationMeta = computed(() => {
  return data.value?.data?.meta || {
    page: 1,
    limit: 10,
    total: 0,
    total_page: 1
  }
})

const lowonganList = computed(() => {
  const items = data.value?.data?.data || []

  return items.map((lowongan) => ({
    id: lowongan.id ?? lowongan.ID,
    judul: lowongan.judul ?? lowongan.Judul ?? '',
    unit: lowongan.unit ?? lowongan.Unit ?? '',
    tanggalBuka: lowongan.tanggalBuka ?? lowongan.tanggal_buka ?? lowongan.TanggalBuka ?? '',
    tanggalTutup: lowongan.tanggalTutup ?? lowongan.tanggal_tutup ?? lowongan.TanggalTutup ?? '',
    deskripsi: lowongan.deskripsi ?? lowongan.Deskripsi ?? '',
    status: lowongan.status ?? lowongan.Status ?? ''
  }))
})

const isAllSelected = computed(() => {
  return lowonganList.value.length > 0 &&
    selectedIds.value.length === lowonganList.value.length
})

const showingFrom = computed(() => {
  if (paginationMeta.value.total === 0) {
    return 0
  }

  return ((paginationMeta.value.page - 1) * paginationMeta.value.limit) + 1
})

const showingTo = computed(() => {
  const end = paginationMeta.value.page * paginationMeta.value.limit

  if (end > paginationMeta.value.total) {
    return paginationMeta.value.total
  }

  return end
})

onMounted(async () => {
  authStore.restoreSession()

  if (authStore.isLoggedIn) {
    await refresh()
  }
})

function clearPageState() {
  selectedIds.value = []
  selectedLowongan.value = null
  bulkMessage.value = ''
  bulkError.value = ''
  actionMessage.value = ''
  actionError.value = ''
}

function setActionMessage(message) {
  actionMessage.value = message
  actionError.value = ''
}

function setActionError(message) {
  actionMessage.value = ''
  actionError.value = message
}

function buildLowonganPayload(includeStatus = false) {
  const payload = {
    judul: form.judul,
    unit: form.unit,
    tanggalBuka: form.tanggalBuka,
    tanggalTutup: form.tanggalTutup,
    deskripsi: form.deskripsi
  }

  if (includeStatus) {
    payload.status = form.status || 'aktif'
  }

  return payload
}

function validateDateRange() {
  if (form.tanggalBuka && form.tanggalTutup && form.tanggalBuka > form.tanggalTutup) {
    setActionError('Tanggal buka tidak boleh melewati tanggal tutup')
    return false
  }

  if (actionError.value === 'Tanggal buka tidak boleh melewati tanggal tutup') {
    actionError.value = ''
  }

  return true
}

function normalizeLowongan(lowongan) {
  if (!lowongan) {
    return null
  }

  return {
    id: lowongan.id ?? lowongan.ID,
    judul: lowongan.judul ?? lowongan.Judul ?? '',
    unit: lowongan.unit ?? lowongan.Unit ?? '',
    tanggalBuka: lowongan.tanggalBuka ?? lowongan.tanggal_buka ?? lowongan.TanggalBuka ?? '',
    tanggalTutup: lowongan.tanggalTutup ?? lowongan.tanggal_tutup ?? lowongan.TanggalTutup ?? '',
    deskripsi: lowongan.deskripsi ?? lowongan.Deskripsi ?? '',
    status: lowongan.status ?? lowongan.Status ?? ''
  }
}

function formatDate(value) {
  if (!value || String(value).startsWith('0001-01-01')) {
    return '-'
  }

  return String(value).slice(0, 10)
}

function dateForEdit(value) {
  const formattedDate = formatDate(value)

  if (formattedDate === '-') {
    return today
  }

  return formattedDate
}

async function logout() {
  authStore.logout()
  clearPageState()
  data.value = null
  await navigateTo('/login')
}

function handleUnauthorized(error) {
  if (authStore.handleUnauthorized(error)) {
    clearPageState()
    data.value = null
    navigateTo('/login')
    return true
  }

  return false
}

async function resetFilter() {
  filter.keyword = ''
  filter.status = ''
  page.value = 1
  selectedIds.value = []
  await refresh()
  setActionMessage('Filter berhasil direset')
}

async function applyFilter() {
  page.value = 1
  selectedIds.value = []
  await refresh()
  setActionMessage('Filter berhasil diterapkan')
}

async function changeLimit() {
  page.value = 1
  selectedIds.value = []
  await refresh()
  setActionMessage('Jumlah data per halaman berhasil diubah')
}

async function showDetail(id) {
  try {
    const result = await $fetch(`http://localhost:8080/api/lowongan/detail?id=${id}`, {
      headers: authHeaders.value
    })
    selectedLowongan.value = normalizeLowongan(result.data)
    setActionMessage('Detail lowongan berhasil dimuat')
  } catch (error) {
    if (!handleUnauthorized(error)) {
      setActionError('Gagal memuat detail lowongan')
    }
  }
}

async function submitLowongan() {
  if (!validateDateRange()) {
    return
  }

  try {
    const isEditing = Boolean(editId.value)

    if (editId.value) {
      await $fetch(`http://localhost:8080/api/lowongan?id=${editId.value}`, {
        method: 'PUT',
        headers: authHeaders.value,
        body: buildLowonganPayload(true)
      })
    } else {
      await $fetch('http://localhost:8080/api/lowongan', {
        method: 'POST',
        headers: authHeaders.value,
        body: buildLowonganPayload()
      })
    }

    resetForm()
    await refresh()
    setActionMessage(isEditing ? 'Lowongan berhasil diperbarui' : 'Lowongan berhasil ditambahkan')
  } catch (error) {
    if (!handleUnauthorized(error)) {
      setActionError(editId.value ? 'Gagal memperbarui lowongan' : 'Gagal menambahkan lowongan')
    }
  }
}

async function deleteLowongan(id) {
  const yakin = confirm('Yakin mau hapus lowongan ini?')

  if (!yakin) {
    return
  }

  try {
    await $fetch(`http://localhost:8080/api/lowongan?id=${id}`, {
      method: 'DELETE',
      headers: authHeaders.value
    })

    await refresh()
    setActionMessage('Lowongan berhasil dihapus')
  } catch (error) {
    if (!handleUnauthorized(error)) {
      setActionError('Gagal menghapus lowongan')
    }
  }
}

async function toggleStatus(lowongan) {
  const nextStatus = lowongan.status === 'aktif' ? 'nonaktif' : 'aktif'

  try {
    await $fetch(`http://localhost:8080/api/lowongan/status?id=${lowongan.id}`, {
      method: 'PUT',
      headers: authHeaders.value,
      body: {
        status: nextStatus
      }
    })

    await refresh()
    setActionMessage('Status lowongan berhasil diubah')
  } catch (error) {
    if (!handleUnauthorized(error)) {
      setActionError('Gagal mengubah status lowongan')
    }
  }
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
      headers: authHeaders.value,
      body: {
        ids: selectedIds.value,
        status: status
      }
    })

    selectedIds.value = []
    bulkMessage.value = 'Status lowongan terpilih berhasil diubah'
    await refresh()
  } catch (error) {
    if (!handleUnauthorized(error)) {
      bulkError.value = 'Gagal mengubah status lowongan'
    }
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
      headers: authHeaders.value,
      body: {
        ids: selectedIds.value
      }
    })

    selectedIds.value = []
    bulkMessage.value = 'Lowongan terpilih berhasil dihapus'
    await refresh()
  } catch (error) {
    if (!handleUnauthorized(error)) {
      bulkError.value = 'Gagal menghapus lowongan'
    }
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
  form.tanggalBuka = dateForEdit(lowongan.tanggalBuka)
  form.tanggalTutup = dateForEdit(lowongan.tanggalTutup)
  form.deskripsi = lowongan.deskripsi
  form.status = lowongan.status
}

function resetForm() {
  editId.value = null
  form.judul = ''
  form.unit = ''
  form.tanggalBuka = ''
  form.tanggalTutup = ''
  form.deskripsi = ''
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
    <header style="display: flex; justify-content: space-between; align-items: center; gap: 16px; margin-bottom: 24px;">
      <div>
        <h1 style="margin-bottom: 4px;">Lowongan Rekrutmen</h1>
        <p style="margin-top: 0; color: #666;">
          Login sebagai {{ currentUser?.name || 'User' }} ({{ currentUser?.role || '-' }})
        </p>
      </div>

      <button type="button" @click="logout">
        Logout
      </button>
    </header>

    <p v-if="actionMessage"
      style="padding: 10px 12px; color: #166534; background: #dcfce7; border: 1px solid #86efac; margin-bottom: 16px;">
      {{ actionMessage }}
    </p>

    <p v-if="actionError"
      style="padding: 10px 12px; color: #991b1b; background: #fee2e2; border: 1px solid #fca5a5; margin-bottom: 16px;">
      {{ actionError }}
    </p>

    <form @submit.prevent="submitLowongan" style="margin-bottom: 24px;">
      <BaseInput v-model="form.judul" label="Judul Lowongan" required maxlength="100" />

      <BaseInput v-model="form.unit" label="Unit" placeholder="Contoh: Direktorat SDM" />

      <BaseInput
        v-model="form.tanggalBuka"
        label="Tanggal Buka"
        type="date"
        :max="form.tanggalTutup || null"
        @change="validateDateRange"
      />

      <BaseInput
        v-model="form.tanggalTutup"
        label="Tanggal Tutup"
        type="date"
        :min="form.tanggalBuka || null"
        @change="validateDateRange"
      />

      <div style="margin-bottom: 12px;">
        <label>Deskripsi</label>
        <textarea
          v-model="form.deskripsi"
          rows="4"
          placeholder="Tulis ringkasan pekerjaan, tanggung jawab, atau kualifikasi utama"
          style="width: 100%; padding: 8px; margin-top: 4px; box-sizing: border-box;"
        ></textarea>
      </div>

      <div v-if="editId" style="margin-bottom: 12px;">
        <label>Status</label>
        <select v-model="form.status" style="display: block; padding: 8px; margin-top: 4px;">
          <option value="aktif">Aktif</option>
          <option value="nonaktif">Nonaktif</option>
        </select>
      </div>

      <BaseButton type="submit" :label="editId ? 'Simpan Perubahan' : 'Buka Lowongan'"
        :color="editId ? 'blue' : 'green'" />

      <button v-if="editId" type="button" @click="resetForm">
        Batal Edit
      </button>
    </form>

    <div style="margin-bottom: 16px;">
      <h3>Filter Lowongan</h3>

      <input v-model="filter.keyword" type="text" placeholder="Cari judul/unit">

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

      <p v-if="bulkLoading">
        Memproses bulk action...
      </p>

      <button type="button" :disabled="selectedIds.length === 0 || bulkLoading" @click="bulkUpdateStatus('aktif')">
        Aktifkan Terpilih
      </button>

      <button type="button" :disabled="selectedIds.length === 0 || bulkLoading" @click="bulkUpdateStatus('nonaktif')">
        Nonaktifkan Terpilih
      </button>

      <button type="button" :disabled="selectedIds.length === 0 || bulkLoading" @click="bulkDelete">
        Hapus Terpilih
      </button>

      <select v-model="limit" :disabled="pending" style="margin-left: 12px;" @change="changeLimit">
        <option :value="5">5 / page</option>
        <option :value="10">10 / page</option>
        <option :value="25">25 / page</option>
        <option :value="50">50 / page</option>
      </select>
    </div>

    <p v-if="pending">Loading...</p>
    <p v-else-if="error">Gagal ambil data lowongan</p>

    <table v-if="!pending && !error" border="1" cellpadding="8" cellspacing="0">
      <thead>
        <tr>
          <th>
            <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll">
          </th>
          <th>ID</th>
          <th>Judul</th>
          <th>Unit</th>
          <th>Tanggal Buka</th>
          <th>Tanggal Tutup</th>
          <th>Status</th>
          <th>Deskripsi</th>
          <th>Aksi</th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="lowongan in lowonganList" :key="lowongan.id">
          <td>
            <input v-model="selectedIds" type="checkbox" :value="lowongan.id">
          </td>
          <td>{{ lowongan.id }}</td>
          <td>{{ lowongan.judul }}</td>
          <td>{{ lowongan.unit }}</td>
          <td>{{ formatDate(lowongan.tanggalBuka) }}</td>
          <td>{{ formatDate(lowongan.tanggalTutup) }}</td>
          <td>{{ lowongan.status }}</td>
          <td>{{ lowongan.deskripsi || '-' }}</td>
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
      <div style="margin-bottom: 8px;">
        <span>
          Menampilkan {{ showingFrom }} - {{ showingTo }}
          dari {{ paginationMeta.total }} data
        </span>
      </div>
    </div>

    <div style="margin-top: 12px;">
      <button type="button" :disabled="page <= 1 || pending" @click="previousPage">
        Sebelumnya
      </button>
      <span style="margin: 0 8px;">
        Halaman {{ paginationMeta.page }} / {{ paginationMeta.total_page }}
      </span>
      <button type="button" :disabled="page >= paginationMeta.total_page || pending" @click="nextPage">
        Berikutnya
      </button>
    </div>

    <div v-if="selectedLowongan" style="margin-top: 16px;">
      <h3>Detail Lowongan</h3>
      <p>ID: {{ selectedLowongan.id }}</p>
      <p>Judul: {{ selectedLowongan.judul }}</p>
      <p>Unit: {{ selectedLowongan.unit }}</p>
      <p>Tanggal Buka: {{ formatDate(selectedLowongan.tanggalBuka) }}</p>
      <p>Tanggal Tutup: {{ formatDate(selectedLowongan.tanggalTutup) }}</p>
      <p>Status: {{ selectedLowongan.status }}</p>
      <p>Deskripsi: {{ selectedLowongan.deskripsi || '-' }}</p>

      <button type="button" @click="selectedLowongan = null">
        Tutup Detail
      </button>
    </div>
  </main>
</template>
