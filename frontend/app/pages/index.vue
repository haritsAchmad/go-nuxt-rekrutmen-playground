<script setup>
definePageMeta({ middleware: 'auth' })

const authStore = useAuthStore()
const currentUser = computed(() => authStore.user)
const authHeaders = computed(() => authStore.authHeaders)

const form = reactive({ judul: '', unit: '', tanggalBuka: '', tanggalTutup: '', deskripsi: '', status: 'aktif' })
const filter = reactive({ keyword: '', status: '' })

const editId = ref(null)
const selectedIds = ref([])
const selectedLowongan = ref(null)
const isFormModalOpen = ref(false)
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
  if (filter.keyword) params.append('keyword', filter.keyword)
  if (filter.status) params.append('status', filter.status)
  params.append('page', page.value)
  params.append('limit', limit.value)
  return `http://localhost:8080/api/lowongan?${params.toString()}`
})

const modalTitle = computed(() => editId.value ? 'Edit Lowongan' : 'Tambah Lowongan')
const modalDescription = computed(() => editId.value ? 'Perbarui data lowongan yang dipilih.' : 'Buat lowongan baru dari form singkat ini.')

const { data, pending, error, refresh } = await useFetch(apiUrl, {
  immediate: false,
  server: false,
  watch: false,
  headers: authHeaders
})

const paginationMeta = computed(() => data.value?.data?.meta || { page: 1, limit: 10, total: 0, total_page: 1 })

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

const isAllSelected = computed(() => lowonganList.value.length > 0 && selectedIds.value.length === lowonganList.value.length)
const showingFrom = computed(() => paginationMeta.value.total === 0 ? 0 : ((paginationMeta.value.page - 1) * paginationMeta.value.limit) + 1)
const showingTo = computed(() => Math.min(paginationMeta.value.page * paginationMeta.value.limit, paginationMeta.value.total))

onMounted(async () => {
  authStore.restoreSession()
  if (authStore.isLoggedIn) await refresh()
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
  if (includeStatus) payload.status = form.status || 'aktif'
  return payload
}

function validateDateRange() {
  if (form.tanggalBuka && form.tanggalTutup && form.tanggalBuka > form.tanggalTutup) {
    setActionError('Tanggal buka tidak boleh melewati tanggal tutup')
    return false
  }
  if (actionError.value === 'Tanggal buka tidak boleh melewati tanggal tutup') actionError.value = ''
  return true
}

function normalizeLowongan(lowongan) {
  if (!lowongan) return null
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
  if (!value || String(value).startsWith('0001-01-01')) return '-'
  return String(value).slice(0, 10)
}

function dateForEdit(value) {
  const formattedDate = formatDate(value)
  return formattedDate === '-' ? today : formattedDate
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
    const result = await $fetch(`http://localhost:8080/api/lowongan/detail?id=${id}`, { headers: authHeaders.value })
    selectedLowongan.value = normalizeLowongan(result.data)
    setActionMessage('Detail lowongan berhasil dimuat')
  } catch (error) {
    if (!handleUnauthorized(error)) setActionError('Gagal memuat detail lowongan')
  }
}

async function submitLowongan() {
  if (!validateDateRange()) return
  const isEditing = Boolean(editId.value)
  const yakin = confirm(isEditing ? 'Yakin mau simpan perubahan lowongan ini?' : 'Yakin mau menambahkan lowongan baru?')
  if (!yakin) return

  try {
    if (editId.value) {
      await $fetch(`http://localhost:8080/api/lowongan?id=${editId.value}`, { method: 'PUT', headers: authHeaders.value, body: buildLowonganPayload(true) })
    } else {
      await $fetch('http://localhost:8080/api/lowongan', { method: 'POST', headers: authHeaders.value, body: buildLowonganPayload() })
    }
    closeFormModal()
    await refresh()
    setActionMessage(isEditing ? 'Lowongan berhasil diperbarui' : 'Lowongan berhasil ditambahkan')
  } catch (error) {
    if (!handleUnauthorized(error)) setActionError(editId.value ? 'Gagal memperbarui lowongan' : 'Gagal menambahkan lowongan')
  }
}

async function deleteLowongan(id) {
  if (!confirm('Yakin mau hapus lowongan ini?')) return
  try {
    await $fetch(`http://localhost:8080/api/lowongan?id=${id}`, { method: 'DELETE', headers: authHeaders.value })
    await refresh()
    setActionMessage('Lowongan berhasil dihapus')
  } catch (error) {
    if (!handleUnauthorized(error)) setActionError('Gagal menghapus lowongan')
  }
}

async function toggleStatus(lowongan) {
  const nextStatus = lowongan.status === 'aktif' ? 'nonaktif' : 'aktif'
  try {
    await $fetch(`http://localhost:8080/api/lowongan/status?id=${lowongan.id}`, { method: 'PUT', headers: authHeaders.value, body: { status: nextStatus } })
    await refresh()
    setActionMessage('Status lowongan berhasil diubah')
  } catch (error) {
    if (!handleUnauthorized(error)) setActionError('Gagal mengubah status lowongan')
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
    await $fetch('http://localhost:8080/api/lowongan/bulk-status', { method: 'PUT', headers: authHeaders.value, body: { ids: selectedIds.value, status } })
    selectedIds.value = []
    bulkMessage.value = 'Status lowongan terpilih berhasil diubah'
    await refresh()
  } catch (error) {
    if (!handleUnauthorized(error)) bulkError.value = 'Gagal mengubah status lowongan'
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
  if (!confirm('Yakin mau hapus lowongan yang dipilih?')) return
  bulkLoading.value = true
  bulkError.value = ''
  bulkMessage.value = ''
  try {
    await $fetch('http://localhost:8080/api/lowongan/bulk-delete', { method: 'DELETE', headers: authHeaders.value, body: { ids: selectedIds.value } })
    selectedIds.value = []
    bulkMessage.value = 'Lowongan terpilih berhasil dihapus'
    await refresh()
  } catch (error) {
    if (!handleUnauthorized(error)) bulkError.value = 'Gagal menghapus lowongan'
  } finally {
    bulkLoading.value = false
  }
}

async function nextPage() {
  if (page.value >= paginationMeta.value.total_page) return
  page.value++
  selectedIds.value = []
  await refresh()
}

async function previousPage() {
  if (page.value <= 1) return
  page.value--
  selectedIds.value = []
  await refresh()
}

function openCreateModal() {
  resetForm()
  actionError.value = ''
  isFormModalOpen.value = true
}

function editLowongan(lowongan) {
  editId.value = lowongan.id
  form.judul = lowongan.judul
  form.unit = lowongan.unit
  form.tanggalBuka = dateForEdit(lowongan.tanggalBuka)
  form.tanggalTutup = dateForEdit(lowongan.tanggalTutup)
  form.deskripsi = lowongan.deskripsi
  form.status = lowongan.status
  actionError.value = ''
  isFormModalOpen.value = true
}

function closeFormModal() {
  isFormModalOpen.value = false
  resetForm()
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
  selectedIds.value = isAllSelected.value ? [] : lowonganList.value.map((lowongan) => lowongan.id)
}
</script>

<template>
  <main class="min-h-screen bg-slate-100 px-4 py-6 text-slate-900 sm:px-6 lg:px-8">
    <div class="mx-auto max-w-7xl space-y-6">
      <header class="overflow-hidden rounded-[2rem] bg-slate-950 shadow-2xl shadow-slate-300">
        <div class="relative px-6 py-8 sm:px-8">
          <div class="absolute inset-0 bg-[radial-gradient(circle_at_top_left,_rgba(99,102,241,0.4),_transparent_35%),radial-gradient(circle_at_bottom_right,_rgba(20,184,166,0.28),_transparent_32%)]"></div>
          <div class="relative flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p class="mb-3 inline-flex rounded-full border border-white/10 bg-white/10 px-4 py-2 text-sm font-semibold text-indigo-100 backdrop-blur">Rekrutmen Playground</p>
              <h1 class="text-3xl font-black tracking-tight text-white sm:text-4xl">Lowongan Rekrutmen</h1>
              <p class="mt-2 max-w-2xl text-sm leading-6 text-slate-300 sm:text-base">Login sebagai {{ currentUser?.name || 'User' }} dengan role {{ currentUser?.role || '-' }}.</p>
            </div>

            <button type="button" class="border border-white/10 bg-white/10 text-white hover:bg-white/20" @click="logout">Logout</button>
          </div>
        </div>
      </header>

      <div v-if="actionMessage" class="rounded-2xl border border-emerald-200 bg-emerald-50 px-5 py-4 text-sm font-semibold text-emerald-700">{{ actionMessage }}</div>
      <div v-if="actionError" class="rounded-2xl border border-rose-200 bg-rose-50 px-5 py-4 text-sm font-semibold text-rose-700">{{ actionError }}</div>

      <section class="dashboard-card space-y-5">
        <div>
          <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">Data lowongan</p>
          <h2 class="mt-1 text-2xl font-black text-slate-950">Daftar Lowongan</h2>
          <p class="mt-1 text-sm text-slate-500">Menampilkan {{ showingFrom }} - {{ showingTo }} dari {{ paginationMeta.total }} data</p>
        </div>

        <div class="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
          <div class="grid w-full gap-3 sm:grid-cols-[96px_150px_minmax(220px,1fr)_auto_auto] xl:max-w-4xl">
            <select v-model="limit" :disabled="pending" class="w-full" @change="changeLimit">
              <option :value="5">5</option>
              <option :value="10">10</option>
              <option :value="25">25</option>
              <option :value="50">50</option>
            </select>

            <select v-model="filter.status" class="w-full" @change="applyFilter">
              <option value="">Semua Status</option>
              <option value="aktif">Aktif</option>
              <option value="nonaktif">Nonaktif</option>
            </select>

            <input v-model="filter.keyword" class="w-full" type="text" placeholder="Cari judul atau unit" @keyup.enter="applyFilter">

            <button type="button" class="muted-button" @click="applyFilter">Cari</button>
            <button type="button" class="muted-button" @click="resetFilter">Reset</button>
          </div>

          <button type="button" class="primary-button whitespace-nowrap xl:ml-auto" @click="openCreateModal">+ Tambah Lowongan</button>
        </div>

        <Transition name="fade">
          <div v-if="selectedIds.length > 0" class="rounded-2xl border border-indigo-100 bg-indigo-50/80 p-4 shadow-sm">
            <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <p class="text-sm font-bold text-indigo-800">{{ selectedIds.length }} lowongan dipilih</p>
                <p class="text-xs text-indigo-600">Pilih aksi cepat untuk data yang sudah dicentang.</p>
              </div>
              <div class="flex flex-wrap gap-2">
                <button type="button" class="success-button" :disabled="bulkLoading" @click="bulkUpdateStatus('aktif')">Aktifkan</button>
                <button type="button" class="muted-button" :disabled="bulkLoading" @click="bulkUpdateStatus('nonaktif')">Nonaktifkan</button>
                <button type="button" class="danger-button" :disabled="bulkLoading" @click="bulkDelete">Hapus</button>
                <button type="button" class="muted-button" :disabled="bulkLoading" @click="selectedIds = []">Batal Pilih</button>
              </div>
            </div>
          </div>
        </Transition>

        <div v-if="bulkMessage" class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-semibold text-emerald-700">{{ bulkMessage }}</div>
        <div v-if="bulkError" class="rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">{{ bulkError }}</div>
        <div v-if="bulkLoading" class="rounded-2xl border border-indigo-200 bg-indigo-50 px-4 py-3 text-sm font-semibold text-indigo-700">Memproses aksi pilihan...</div>

        <p v-if="pending" class="soft-panel text-sm font-semibold text-slate-600">Loading data lowongan...</p>
        <p v-else-if="error" class="rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">Gagal ambil data lowongan</p>

        <div v-if="!pending && !error" class="overflow-hidden rounded-2xl border border-slate-200">
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-slate-200 text-left text-sm">
              <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
                <tr>
                  <th class="px-4 py-3"><input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll"></th>
                  <th class="px-4 py-3">ID</th>
                  <th class="px-4 py-3">Judul</th>
                  <th class="px-4 py-3">Unit</th>
                  <th class="px-4 py-3">Buka</th>
                  <th class="px-4 py-3">Tutup</th>
                  <th class="px-4 py-3">Status</th>
                  <th class="px-4 py-3">Deskripsi</th>
                  <th class="px-4 py-3">Aksi</th>
                </tr>
              </thead>

              <tbody class="divide-y divide-slate-100 bg-white">
                <tr v-for="lowongan in lowonganList" :key="lowongan.id" class="hover:bg-indigo-50/40">
                  <td class="px-4 py-4"><input v-model="selectedIds" type="checkbox" :value="lowongan.id"></td>
                  <td class="px-4 py-4 font-semibold text-slate-500">#{{ lowongan.id }}</td>
                  <td class="px-4 py-4 font-bold text-slate-900">{{ lowongan.judul }}</td>
                  <td class="px-4 py-4 text-slate-600">{{ lowongan.unit }}</td>
                  <td class="px-4 py-4 text-slate-600">{{ formatDate(lowongan.tanggalBuka) }}</td>
                  <td class="px-4 py-4 text-slate-600">{{ formatDate(lowongan.tanggalTutup) }}</td>
                  <td class="px-4 py-4"><span :class="lowongan.status === 'aktif' ? 'bg-emerald-100 text-emerald-700 ring-emerald-200' : 'bg-slate-100 text-slate-600 ring-slate-200'" class="inline-flex rounded-full px-3 py-1 text-xs font-bold ring-1">{{ lowongan.status }}</span></td>
                  <td class="max-w-xs px-4 py-4 text-slate-600"><p class="line-clamp-2">{{ lowongan.deskripsi || '-' }}</p></td>
                  <td class="px-4 py-4">
                    <div class="flex flex-wrap gap-2">
                      <button type="button" class="muted-button px-3 py-1.5" @click="editLowongan(lowongan)">Edit</button>
                      <button type="button" class="muted-button px-3 py-1.5" @click="showDetail(lowongan.id)">Detail</button>
                      <button type="button" class="primary-button px-3 py-1.5" @click="toggleStatus(lowongan)">{{ lowongan.status === 'aktif' ? 'Nonaktifkan' : 'Aktifkan' }}</button>
                      <button type="button" class="danger-button px-3 py-1.5" @click="deleteLowongan(lowongan.id)">Hapus</button>
                    </div>
                  </td>
                </tr>

                <tr v-if="lowonganList.length === 0">
                  <td colspan="9" class="px-4 py-10 text-center text-sm font-semibold text-slate-500">Belum ada data lowongan yang cocok.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="flex flex-col gap-3 border-t border-slate-100 pt-4 sm:flex-row sm:items-center sm:justify-between">
          <span class="text-sm font-medium text-slate-500">Menampilkan {{ showingFrom }} - {{ showingTo }} dari {{ paginationMeta.total }} data</span>
          <div class="flex items-center gap-2">
            <button type="button" class="muted-button" :disabled="page <= 1 || pending" @click="previousPage">Sebelumnya</button>
            <span class="rounded-xl bg-slate-100 px-4 py-2 text-sm font-bold text-slate-700">{{ paginationMeta.page }} / {{ paginationMeta.total_page }}</span>
            <button type="button" class="muted-button" :disabled="page >= paginationMeta.total_page || pending" @click="nextPage">Berikutnya</button>
          </div>
        </div>
      </section>

      <section v-if="selectedLowongan" class="dashboard-card">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">Detail Lowongan</p>
            <h2 class="mt-1 text-2xl font-black text-slate-950">{{ selectedLowongan.judul }}</h2>
            <p class="mt-2 text-sm text-slate-500">ID: #{{ selectedLowongan.id }}</p>
          </div>
          <button type="button" class="muted-button" @click="selectedLowongan = null">Tutup Detail</button>
        </div>

        <div class="mt-6 grid gap-4 md:grid-cols-4">
          <div class="soft-panel"><p class="text-xs font-bold uppercase tracking-wide text-slate-400">Unit</p><p class="mt-1 font-semibold text-slate-800">{{ selectedLowongan.unit || '-' }}</p></div>
          <div class="soft-panel"><p class="text-xs font-bold uppercase tracking-wide text-slate-400">Tanggal Buka</p><p class="mt-1 font-semibold text-slate-800">{{ formatDate(selectedLowongan.tanggalBuka) }}</p></div>
          <div class="soft-panel"><p class="text-xs font-bold uppercase tracking-wide text-slate-400">Tanggal Tutup</p><p class="mt-1 font-semibold text-slate-800">{{ formatDate(selectedLowongan.tanggalTutup) }}</p></div>
          <div class="soft-panel"><p class="text-xs font-bold uppercase tracking-wide text-slate-400">Status</p><p class="mt-1 font-semibold text-slate-800">{{ selectedLowongan.status || '-' }}</p></div>
        </div>

        <div class="soft-panel mt-4"><p class="text-xs font-bold uppercase tracking-wide text-slate-400">Deskripsi</p><p class="mt-2 leading-7 text-slate-700">{{ selectedLowongan.deskripsi || '-' }}</p></div>
      </section>
    </div>

    <Teleport to="body">
      <Transition name="fade">
        <div v-if="isFormModalOpen" class="fixed inset-0 z-50 grid place-items-center bg-slate-950/60 px-4 py-6 backdrop-blur-sm" @click.self="closeFormModal">
          <form class="w-full max-w-xl rounded-[2rem] border border-white/70 bg-white p-6 shadow-2xl shadow-slate-950/30" @submit.prevent="submitLowongan">
            <div class="mb-5 flex items-start justify-between gap-4">
              <div>
                <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">{{ editId ? 'Mode update' : 'Mode create' }}</p>
                <h2 class="mt-1 text-2xl font-black text-slate-950">{{ modalTitle }}</h2>
                <p class="mt-1 text-sm text-slate-500">{{ modalDescription }}</p>
              </div>
              <button type="button" class="muted-button px-3 py-1.5" @click="closeFormModal">✕</button>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div class="sm:col-span-2"><label class="mb-2 block text-sm font-semibold text-slate-700">Judul Lowongan</label><input v-model="form.judul" class="w-full" type="text" required maxlength="100" placeholder="Contoh: Frontend Developer"></div>
              <div class="sm:col-span-2"><label class="mb-2 block text-sm font-semibold text-slate-700">Unit</label><input v-model="form.unit" class="w-full" type="text" placeholder="Contoh: Direktorat SDM"></div>
              <div><label class="mb-2 block text-sm font-semibold text-slate-700">Tanggal Buka</label><input v-model="form.tanggalBuka" class="w-full" type="date" :max="form.tanggalTutup || null" @change="validateDateRange"></div>
              <div><label class="mb-2 block text-sm font-semibold text-slate-700">Tanggal Tutup</label><input v-model="form.tanggalTutup" class="w-full" type="date" :min="form.tanggalBuka || null" @change="validateDateRange"></div>
              <div v-if="editId" class="sm:col-span-2"><label class="mb-2 block text-sm font-semibold text-slate-700">Status</label><select v-model="form.status" class="w-full"><option value="aktif">Aktif</option><option value="nonaktif">Nonaktif</option></select></div>
              <div class="sm:col-span-2"><label class="mb-2 block text-sm font-semibold text-slate-700">Deskripsi</label><textarea v-model="form.deskripsi" rows="4" placeholder="Tulis ringkasan pekerjaan, tanggung jawab, atau kualifikasi utama"></textarea></div>
            </div>

            <div class="mt-6 flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
              <button type="button" class="muted-button" @click="closeFormModal">Batal</button>
              <button type="submit" :class="editId ? 'primary-button' : 'success-button'">{{ editId ? 'Simpan Perubahan' : 'Simpan Lowongan' }}</button>
            </div>
          </form>
        </div>
      </Transition>
    </Teleport>
  </main>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
