<script setup>
definePageMeta({
  middleware: 'auth',
  layout: 'app'
})

const authStore = useAuthStore()
const authHeaders = computed(() => authStore.authHeaders)

const statusOptions = ['baru', 'diproses', 'diterima', 'ditolak']

const form = reactive({
  lowonganId: '',
  nama: '',
  email: '',
  noHp: '',
  status: 'baru'
})

const filter = reactive({
  keyword: '',
  status: '',
  sort: 'newest',
  lowonganId: ''
})

const editId = ref(null)
const selectedIds = ref([])
const selectedPelamar = ref(null)
const isFormModalOpen = ref(false)

const bulkLoading = ref(false)
const bulkMessage = ref('')
const bulkError = ref('')
const actionMessage = ref('')
const actionError = ref('')

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

  if (filter.lowonganId) {
    params.append('lowonganId', filter.lowonganId)
  }

  params.append('sort', filter.sort)
  params.append('page', page.value)
  params.append('limit', limit.value)

  return `http://localhost:8080/api/pelamar?${params.toString()}`
})

const modalTitle = computed(() => editId.value ? 'Edit Pelamar' : 'Tambah Pelamar')
const modalDescription = computed(() => editId.value ? 'Perbarui data pelamar yang dipilih.' : 'Tambahkan pelamar baru untuk lowongan tertentu.')

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

const pelamarList = computed(() => {
  const items = data.value?.data?.data || []

  return items.map(normalizePelamar)
})

const isAllSelected = computed(() => {
  return pelamarList.value.length > 0 && selectedIds.value.length === pelamarList.value.length
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

function normalizePelamar(pelamar) {
  if (!pelamar) {
    return null
  }

  return {
    id: pelamar.id ?? pelamar.ID,
    lowonganId: pelamar.lowonganId ?? pelamar.lowongan_id ?? pelamar.LowonganID ?? '',
    lowonganJudul: pelamar.lowonganJudul ?? pelamar.lowongan_judul ?? pelamar.LowonganJudul ?? '',
    nama: pelamar.nama ?? pelamar.Nama ?? '',
    email: pelamar.email ?? pelamar.Email ?? '',
    noHp: pelamar.noHp ?? pelamar.no_hp ?? pelamar.NoHP ?? '',
    status: pelamar.status ?? pelamar.Status ?? '',
    createdAt: pelamar.createdAt ?? pelamar.created_at ?? pelamar.CreatedAt ?? ''
  }
}

function formatDate(value) {
  if (!value || String(value).startsWith('0001-01-01')) {
    return '-'
  }

  return String(value).slice(0, 10)
}

function statusClass(status) {
  const classes = {
    baru: 'bg-sky-100 text-sky-700 ring-sky-200',
    diproses: 'bg-amber-100 text-amber-700 ring-amber-200',
    diterima: 'bg-emerald-100 text-emerald-700 ring-emerald-200',
    ditolak: 'bg-rose-100 text-rose-700 ring-rose-200'
  }

  return classes[status] || 'bg-slate-100 text-slate-600 ring-slate-200'
}

function clearPageState() {
  selectedIds.value = []
  selectedPelamar.value = null
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

function handleUnauthorized(error) {
  if (authStore.handleUnauthorized(error)) {
    clearPageState()
    data.value = null
    navigateTo('/login')
    return true
  }

  return false
}

function buildPelamarPayload(includeStatus = false) {
  const payload = {
    lowonganId: Number(form.lowonganId),
    nama: form.nama,
    email: form.email,
    noHp: form.noHp
  }

  if (includeStatus) {
    payload.status = form.status || 'baru'
  }

  return payload
}

function validateForm() {
  if (!form.lowonganId || Number(form.lowonganId) <= 0) {
    setActionError('ID lowongan wajib diisi')
    return false
  }

  if (!form.nama.trim() || !form.email.trim()) {
    setActionError('Nama dan email wajib diisi')
    return false
  }

  actionError.value = ''
  return true
}

async function resetFilter() {
  filter.keyword = ''
  filter.status = ''
  filter.sort = 'newest'
  filter.lowonganId = ''
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
    const result = await $fetch(`http://localhost:8080/api/pelamar/detail?id=${id}`, {
      headers: authHeaders.value
    })
    selectedPelamar.value = normalizePelamar(result.data)
    setActionMessage('Detail pelamar berhasil dimuat')
  } catch (error) {
    if (!handleUnauthorized(error)) {
      setActionError('Gagal memuat detail pelamar')
    }
  }
}

async function submitPelamar() {
  if (!validateForm()) {
    return
  }

  const isEditing = Boolean(editId.value)
  const yakin = confirm(isEditing ? 'Yakin mau simpan perubahan pelamar ini?' : 'Yakin mau menambahkan pelamar baru?')

  if (!yakin) {
    return
  }

  try {
    if (editId.value) {
      await $fetch(`http://localhost:8080/api/pelamar?id=${editId.value}`, {
        method: 'PUT',
        headers: authHeaders.value,
        body: buildPelamarPayload(true)
      })
    } else {
      await $fetch('http://localhost:8080/api/pelamar', {
        method: 'POST',
        headers: authHeaders.value,
        body: buildPelamarPayload()
      })
    }

    closeFormModal()
    await refresh()
    setActionMessage(isEditing ? 'Pelamar berhasil diperbarui' : 'Pelamar berhasil ditambahkan')
  } catch (error) {
    if (!handleUnauthorized(error)) {
      setActionError(editId.value ? 'Gagal memperbarui pelamar' : 'Gagal menambahkan pelamar')
    }
  }
}

async function deletePelamar(id) {
  const yakin = confirm('Yakin mau hapus pelamar ini?')

  if (!yakin) {
    return
  }

  try {
    await $fetch(`http://localhost:8080/api/pelamar?id=${id}`, {
      method: 'DELETE',
      headers: authHeaders.value
    })

    await refresh()
    setActionMessage('Pelamar berhasil dihapus')
  } catch (error) {
    if (!handleUnauthorized(error)) {
      setActionError('Gagal menghapus pelamar')
    }
  }
}

async function bulkUpdateStatus(status) {
  if (selectedIds.value.length === 0) {
    bulkError.value = 'Pilih minimal satu pelamar'
    bulkMessage.value = ''
    return
  }

  bulkLoading.value = true
  bulkError.value = ''
  bulkMessage.value = ''

  try {
    await $fetch('http://localhost:8080/api/pelamar/bulk-status', {
      method: 'PUT',
      headers: authHeaders.value,
      body: {
        ids: selectedIds.value,
        status: status
      }
    })

    selectedIds.value = []
    bulkMessage.value = 'Status pelamar terpilih berhasil diubah'
    await refresh()
  } catch (error) {
    if (!handleUnauthorized(error)) {
      bulkError.value = 'Gagal mengubah status pelamar'
    }
  } finally {
    bulkLoading.value = false
  }
}

async function bulkDelete() {
  if (selectedIds.value.length === 0) {
    bulkError.value = 'Pilih minimal satu pelamar'
    bulkMessage.value = ''
    return
  }

  const yakin = confirm('Yakin mau hapus pelamar yang dipilih?')

  if (!yakin) {
    return
  }

  bulkLoading.value = true
  bulkError.value = ''
  bulkMessage.value = ''

  try {
    await $fetch('http://localhost:8080/api/pelamar/bulk-delete', {
      method: 'DELETE',
      headers: authHeaders.value,
      body: {
        ids: selectedIds.value
      }
    })

    selectedIds.value = []
    bulkMessage.value = 'Pelamar terpilih berhasil dihapus'
    await refresh()
  } catch (error) {
    if (!handleUnauthorized(error)) {
      bulkError.value = 'Gagal menghapus pelamar'
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

function openCreateModal() {
  resetForm()
  actionError.value = ''
  isFormModalOpen.value = true
}

function editPelamar(pelamar) {
  editId.value = pelamar.id
  form.lowonganId = pelamar.lowonganId
  form.nama = pelamar.nama
  form.email = pelamar.email
  form.noHp = pelamar.noHp
  form.status = pelamar.status || 'baru'
  actionError.value = ''
  isFormModalOpen.value = true
}

function closeFormModal() {
  isFormModalOpen.value = false
  resetForm()
}

function resetForm() {
  editId.value = null
  form.lowonganId = ''
  form.nama = ''
  form.email = ''
  form.noHp = ''
  form.status = 'baru'
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value = []
    return
  }

  selectedIds.value = pelamarList.value.map((pelamar) => pelamar.id)
}
</script>

<template>
  <div class="space-y-6">
    <section class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">Rekrutmen</p>
        <h1 class="mt-1 text-3xl font-black tracking-tight text-slate-950">Pelamar</h1>
        <p class="mt-1 text-sm text-slate-500">Kelola data pelamar dan status proses rekrutmen.</p>
      </div>
    </section>

    <div v-if="actionMessage" class="rounded-2xl border border-emerald-200 bg-emerald-50 px-5 py-4 text-sm font-semibold text-emerald-700">
      {{ actionMessage }}
    </div>

    <div v-if="actionError" class="rounded-2xl border border-rose-200 bg-rose-50 px-5 py-4 text-sm font-semibold text-rose-700">
      {{ actionError }}
    </div>

    <section class="dashboard-card space-y-5">
      <div>
        <h2 class="mt-1 text-2xl font-black text-slate-950">Daftar Pelamar</h2>
        <p class="mt-1 text-sm text-slate-500">
          Menampilkan {{ showingFrom }} - {{ showingTo }} dari {{ paginationMeta.total }} data
        </p>
      </div>

      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div class="grid gap-3 sm:grid-cols-[104px_156px_140px_140px_minmax(220px,1fr)_auto_auto] lg:flex-1">
          <select v-model="limit" :disabled="pending" class="w-full" @change="changeLimit">
            <option :value="5">5</option>
            <option :value="10">10</option>
            <option :value="25">25</option>
            <option :value="50">50</option>
          </select>

          <select v-model="filter.status" class="w-full" @change="applyFilter">
            <option value="">Semua Status</option>
            <option v-for="status in statusOptions" :key="status" :value="status">{{ status }}</option>
          </select>

          <select v-model="filter.sort" class="w-full" @change="applyFilter">
            <option value="newest">Terbaru</option>
            <option value="oldest">Terlama</option>
          </select>

          <input v-model="filter.lowonganId" class="w-full" type="number" min="1" placeholder="ID lowongan" @keyup.enter="applyFilter">
          <input v-model="filter.keyword" class="w-full" type="text" placeholder="Cari nama, email, HP, atau lowongan" @keyup.enter="applyFilter">

          <button type="button" class="muted-button" @click="applyFilter">
            Cari
          </button>

          <button type="button" class="muted-button" @click="resetFilter">
            Reset
          </button>
        </div>

        <button type="button" class="primary-button whitespace-nowrap" @click="openCreateModal">
          + Tambah Pelamar
        </button>
      </div>

      <Transition name="fade">
        <div v-if="selectedIds.length > 0" class="rounded-2xl border border-indigo-100 bg-indigo-50/80 p-4 shadow-sm">
          <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p class="text-sm font-bold text-indigo-800">{{ selectedIds.length }} pelamar dipilih</p>
              <p class="text-xs text-indigo-600">Pilih aksi cepat untuk data yang sudah dicentang.</p>
            </div>

            <div class="flex flex-wrap gap-2">
              <button v-for="status in statusOptions" :key="status" type="button" class="muted-button" :disabled="bulkLoading" @click="bulkUpdateStatus(status)">
                {{ status }}
              </button>

              <button type="button" class="danger-button" :disabled="bulkLoading" @click="bulkDelete">
                Hapus
              </button>

              <button type="button" class="muted-button" :disabled="bulkLoading" @click="selectedIds = []">
                Batal Pilih
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <div v-if="bulkMessage" class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-semibold text-emerald-700">
        {{ bulkMessage }}
      </div>

      <div v-if="bulkError" class="rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">
        {{ bulkError }}
      </div>

      <div v-if="bulkLoading" class="rounded-2xl border border-indigo-200 bg-indigo-50 px-4 py-3 text-sm font-semibold text-indigo-700">
        Memproses aksi pilihan...
      </div>

      <p v-if="pending" class="soft-panel text-sm font-semibold text-slate-600">Loading data pelamar...</p>
      <p v-else-if="error" class="rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">Gagal ambil data pelamar</p>

      <div v-if="!pending && !error" class="overflow-hidden rounded-2xl border border-slate-200">
        <div class="overflow-x-auto">
          <table class="w-full min-w-[1050px] table-fixed divide-y divide-slate-200 text-left text-sm">
            <colgroup>
              <col class="w-14">
              <col class="w-[24%]">
              <col class="w-[26%]">
              <col class="w-32">
              <col class="w-32">
              <col class="w-32">
              <col class="w-36">
            </colgroup>

            <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
              <tr>
                <th class="px-4 py-3">
                  <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll">
                </th>
                <th class="px-4 py-3">Pelamar</th>
                <th class="px-4 py-3">Lowongan</th>
                <th class="px-4 py-3 text-center">No HP</th>
                <th class="px-4 py-3 text-center">Status</th>
                <th class="px-4 py-3 text-center">Daftar</th>
                <th class="px-4 py-3 text-center">Aksi</th>
              </tr>
            </thead>

            <tbody class="divide-y divide-slate-100 bg-white">
              <tr v-for="pelamar in pelamarList" :key="pelamar.id" class="hover:bg-indigo-50/40">
                <td class="px-4 py-4">
                  <input v-model="selectedIds" type="checkbox" :value="pelamar.id">
                </td>
                <td class="px-4 py-4">
                  <p class="truncate font-bold text-slate-900" :title="pelamar.nama">{{ pelamar.nama }}</p>
                  <p class="mt-1 truncate text-xs text-slate-500" :title="pelamar.email">{{ pelamar.email }}</p>
                </td>
                <td class="px-4 py-4">
                  <p class="truncate font-semibold text-slate-700" :title="pelamar.lowonganJudul || 'Lowongan belum ditemukan'">
                    {{ pelamar.lowonganJudul || 'Lowongan belum ditemukan' }}
                  </p>
                  <p class="mt-1 text-xs text-slate-400">ID Lowongan: #{{ pelamar.lowonganId }}</p>
                </td>
                <td class="whitespace-nowrap px-4 py-4 text-center text-slate-600">{{ pelamar.noHp || '-' }}</td>
                <td class="px-4 py-4 text-center">
                  <span :class="statusClass(pelamar.status)" class="inline-flex rounded-full px-3 py-1 text-xs font-bold ring-1">
                    {{ pelamar.status }}
                  </span>
                </td>
                <td class="whitespace-nowrap px-4 py-4 text-center text-slate-600">{{ formatDate(pelamar.createdAt) }}</td>
                <td class="px-4 py-4">
                  <div class="flex items-center justify-center gap-1.5">
                    <button type="button" class="muted-button grid h-9 w-9 place-items-center p-0 text-slate-500 hover:border-indigo-200 hover:bg-indigo-50 hover:text-indigo-600" title="Lihat detail" aria-label="Lihat detail pelamar" @click="showDetail(pelamar.id)">
                      <svg class="h-4.5 w-4.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                        <path d="M2.5 12s3.5-6 9.5-6 9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z" />
                        <circle cx="12" cy="12" r="2.5" />
                      </svg>
                    </button>

                    <button type="button" class="muted-button grid h-9 w-9 place-items-center p-0 text-slate-500 hover:border-amber-200 hover:bg-amber-50 hover:text-amber-600" title="Edit pelamar" aria-label="Edit pelamar" @click="editPelamar(pelamar)">
                      <svg class="h-4.5 w-4.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                        <path d="M12 20h9" />
                        <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L8 18l-4 1 1-4Z" />
                      </svg>
                    </button>

                    <button type="button" class="grid h-9 w-9 place-items-center rounded-xl border border-rose-100 bg-white p-0 text-rose-500 shadow-sm transition hover:bg-rose-50 hover:text-rose-700" title="Hapus pelamar" aria-label="Hapus pelamar" @click="deletePelamar(pelamar.id)">
                      <svg class="h-4.5 w-4.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                        <path d="M3 6h18" />
                        <path d="M8 6V4h8v2" />
                        <path d="m19 6-1 14H6L5 6" />
                        <path d="M10 11v5M14 11v5" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>

              <tr v-if="pelamarList.length === 0">
                <td colspan="7" class="px-4 py-10 text-center text-sm font-semibold text-slate-500">
                  Belum ada data pelamar yang cocok.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="flex flex-col gap-3 border-t border-slate-100 pt-4 sm:flex-row sm:items-center sm:justify-between">
        <span class="text-sm font-medium text-slate-500">
          Menampilkan {{ showingFrom }} - {{ showingTo }} dari {{ paginationMeta.total }} data
        </span>

        <div class="flex items-center gap-2">
          <button type="button" class="muted-button" :disabled="page <= 1 || pending" @click="previousPage">
            Sebelumnya
          </button>
          <span class="rounded-xl bg-slate-100 px-4 py-2 text-sm font-bold text-slate-700">
            {{ paginationMeta.page }} / {{ paginationMeta.total_page }}
          </span>
          <button type="button" class="muted-button" :disabled="page >= paginationMeta.total_page || pending" @click="nextPage">
            Berikutnya
          </button>
        </div>
      </div>
    </section>

    <section v-if="selectedPelamar" class="dashboard-card">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">Detail Pelamar</p>
          <h2 class="mt-1 text-2xl font-black text-slate-950">{{ selectedPelamar.nama }}</h2>
          <p class="mt-2 text-sm text-slate-500">ID: #{{ selectedPelamar.id }}</p>
        </div>

        <button type="button" class="muted-button" @click="selectedPelamar = null">
          Tutup Detail
        </button>
      </div>

      <div class="mt-6 grid gap-4 md:grid-cols-4">
        <div class="soft-panel">
          <p class="text-xs font-bold uppercase tracking-wide text-slate-400">Email</p>
          <p class="mt-2 break-words text-sm font-bold text-slate-800">{{ selectedPelamar.email }}</p>
        </div>
        <div class="soft-panel">
          <p class="text-xs font-bold uppercase tracking-wide text-slate-400">No HP</p>
          <p class="mt-2 text-sm font-bold text-slate-800">{{ selectedPelamar.noHp || '-' }}</p>
        </div>
        <div class="soft-panel">
          <p class="text-xs font-bold uppercase tracking-wide text-slate-400">Status</p>
          <p class="mt-2 text-sm font-bold text-slate-800">{{ selectedPelamar.status }}</p>
        </div>
        <div class="soft-panel">
          <p class="text-xs font-bold uppercase tracking-wide text-slate-400">Tanggal Daftar</p>
          <p class="mt-2 text-sm font-bold text-slate-800">{{ formatDate(selectedPelamar.createdAt) }}</p>
        </div>
      </div>

      <div class="mt-4 soft-panel">
        <p class="text-xs font-bold uppercase tracking-wide text-slate-400">Lowongan</p>
        <p class="mt-2 text-sm font-bold text-slate-800">{{ selectedPelamar.lowonganJudul || '-' }}</p>
        <p class="mt-1 text-xs text-slate-500">ID Lowongan: #{{ selectedPelamar.lowonganId }}</p>
      </div>
    </section>

    <Transition name="fade">
      <div v-if="isFormModalOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 px-4 py-6">
        <div class="w-full max-w-2xl rounded-3xl bg-white p-6 shadow-2xl">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">Form Pelamar</p>
              <h2 class="mt-1 text-2xl font-black text-slate-950">{{ modalTitle }}</h2>
              <p class="mt-1 text-sm text-slate-500">{{ modalDescription }}</p>
            </div>
            <button type="button" class="muted-button" @click="closeFormModal">✕</button>
          </div>

          <form class="mt-6 grid gap-4" @submit.prevent="submitPelamar">
            <label class="grid gap-2 text-sm font-bold text-slate-700">
              ID Lowongan
              <input v-model="form.lowonganId" type="number" min="1" placeholder="Masukkan ID lowongan" required>
            </label>

            <label class="grid gap-2 text-sm font-bold text-slate-700">
              Nama
              <input v-model="form.nama" type="text" placeholder="Masukkan nama pelamar" required>
            </label>

            <label class="grid gap-2 text-sm font-bold text-slate-700">
              Email
              <input v-model="form.email" type="email" placeholder="Masukkan email pelamar" required>
            </label>

            <label class="grid gap-2 text-sm font-bold text-slate-700">
              No HP
              <input v-model="form.noHp" type="text" placeholder="Masukkan nomor HP">
            </label>

            <label v-if="editId" class="grid gap-2 text-sm font-bold text-slate-700">
              Status
              <select v-model="form.status">
                <option v-for="status in statusOptions" :key="status" :value="status">{{ status }}</option>
              </select>
            </label>

            <div class="mt-2 flex justify-end gap-2">
              <button type="button" class="muted-button" @click="closeFormModal">Batal</button>
              <button type="submit" class="primary-button">Simpan</button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
  </div>
</template>
