<script setup>
const form = reactive({
  judul: '',
  unit: '',
  status: 'aktif'
})

const editId = ref(null)

const { data, pending, error, refresh } = await useFetch('http://localhost:8080/api/lowongan')

const lowonganList = computed(() => {
  return data.value?.data || []
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
</script>

<template>
  <main style="padding: 24px; font-family: Arial, sans-serif;">
    <h1>Lowongan Rekrutmen</h1>

    <form @submit.prevent="submitLowongan" style="margin-bottom: 24px;">
      <div style="margin-bottom: 8px;">
        <label>Judul Lowongan</label><br>
        <input v-model="form.judul" type="text" placeholder="Contoh: Backend Developer">
      </div>

      <div style="margin-bottom: 8px;">
        <label>Unit</label><br>
        <input v-model="form.unit" type="text" placeholder="Contoh: Direktorat SDM">
      </div>

      	<button type="submit">
  	{{ editId ? 'Simpan Perubahan' : 'Tambah Lowongan' }}
	</button>

	<button v-if="editId" type="button" @click="resetForm">
  	Batal Edit
	</button>
    </form>

    <p v-if="pending">Loading...</p>
    <p v-else-if="error">Gagal ambil data lowongan</p>

    <table v-else border="1" cellpadding="8" cellspacing="0">
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
  	<button type="button" @click="editLowongan(lowongan)">
    	  Edit
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
  </main>
</template>