<script setup>
definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()

const filter = reactive({
  keyword: '',
  status: ''
})

const selectedLowongan = ref(null)
const page = ref(1)
const limit = ref(10)

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
  server: false,
  watch: false,
  headers: computed(() => authStore.authHeaders)
})

const lowonganList = computed(() => data.value?.data?.data || [])

const paginationMeta = computed(() => {
  return data.value?.data?.meta || {
    page: 1,
    limit: 10,
    total: 0,
    total_page: 1
  }
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

const canManageLowongan = computed(() => {
  return ['admin', 'superadmin'].includes(authStore.role)
})

async function logout() {
  authStore.logout()
  data.value = null
  selectedLowongan.value = null
  await navigateTo('/login')
}

function handleUnauthorized(error) {
  if (authStore.handleUnauthorized(error)) {
    data.value = null
    selectedLowongan.value = null
    navigateTo('/login')
    return true
  }

  return false
}

async function applyFilter() {
  page.value = 1
  await refresh()
}

async function resetFilter() {
  filter.keyword = ''
  filter.status = ''
  page.value = 1
  await refresh()
}

async function changeLimit() {
  page.value = 1
  await refresh()
}

async function showDetail(id) {
  try {
    const result = await $fetch(`http://localhost:8080/api/lowongan/detail?id=${id}`, {
      headers: authStore.authHeaders
    })

    selectedLowongan.value = result.data
  } catch (error) {
    handleUnauthorized(error)
  }
}

async function nextPage() {
  if (page.value >= paginationMeta.value.total_page) {
    return
  }

  page.value++
  await refresh()
}

async function previousPage() {
  if (page.value <= 1) {
    return
  }

  page.value--
  await refresh()
}
</script>

<template>
  <main style="padding: 24px; font-family: Arial, sans-serif;">
    <header style="display: flex; justify-content: space-between; align-items: center; gap: 16px; margin-bottom: 24px;">
      <div>
        <h1 style="margin-bottom: 4px;">Lowongan Rekrutmen</h1>
        <p style="margin-top: 0; color: #666;">
          Login sebagai {{ authStore.user?.name || 'User' }} ({{ authStore.role || '-' }})
        </p>
      </div>

      <button type="button" @click="logout">
        Logout
      </button>
    </header>

    <p v-if="canManageLowongan" style="padding: 12px; background: #eef6ff; border: 1px solid #bfdbfe;">
      Role ini boleh mengelola lowongan. Tombol CRUD akan kita pisah lagi bertahap.
    </p>

    <p v-else style="padding: 12px; background: #f5f5f5; border: 1px solid #ddd;">
      Role viewer hanya dapat melihat daftar dan detail lowongan.
    </p>

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

      <select
        v-model="limit"
        :disabled="pending"
        style="margin-left: 12px;"
        @change="changeLimit"
      >
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
          <th>ID</th>
          <th>Judul</th>
          <th>Unit</th>
          <th>Status</th>
          <th>Aksi</th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="lowongan in lowonganList" :key="lowongan.id">
          <td>{{ lowongan.id }}</td>
          <td>{{ lowongan.judul }}</td>
          <td>{{ lowongan.unit }}</td>
          <td>{{ lowongan.status }}</td>
          <td>
            <button type="button" @click="showDetail(lowongan.id)">
              Detail
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div style="margin-top: 16px;">
      <span>
        Menampilkan {{ showingFrom }} - {{ showingTo }} dari {{ paginationMeta.total }} data
      </span>
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
      <p>Status: {{ selectedLowongan.status }}</p>

      <button type="button" @click="selectedLowongan = null">
        Tutup Detail
      </button>
    </div>
  </main>
</template>
