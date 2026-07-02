<script setup>
definePageMeta({
  middleware: 'auth',
  layout: 'app'
})

const authStore = useAuthStore()
const currentUser = computed(() => authStore.user)
</script>

<template>
  <div class="space-y-6">
    <section class="overflow-hidden rounded-[2rem] bg-slate-950 shadow-2xl shadow-slate-300">
      <div class="relative px-6 py-10 sm:px-8">
        <div class="absolute inset-0 bg-[radial-gradient(circle_at_top_left,_rgba(99,102,241,0.45),_transparent_38%),radial-gradient(circle_at_bottom_right,_rgba(20,184,166,0.3),_transparent_34%)]" />

        <div class="relative">
          <p class="mb-3 text-sm font-bold uppercase tracking-[0.2em] text-indigo-200">
            Dashboard Rekrutmen
          </p>
          <h1 class="text-3xl font-black tracking-tight text-white sm:text-4xl">
            Selamat datang, {{ currentUser?.name || 'User' }}!
          </h1>
          <p class="mt-3 max-w-2xl text-sm leading-6 text-slate-300 sm:text-base">
            Kelola proses rekrutmen dari satu tempat. Untuk sekarang, kamu bisa mulai dari pengelolaan data lowongan.
          </p>
        </div>
      </div>
    </section>
    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <NuxtLink
        to="/lowongan"
        class="group dashboard-card block transition hover:-translate-y-1 hover:border-indigo-200 hover:shadow-indigo-100"
      >
        <div class="flex items-start justify-between gap-4">
          <div>
            <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">Menu utama</p>
            <h2 class="mt-2 text-xl font-black text-slate-950">Kelola Lowongan</h2>
            <p class="mt-2 text-sm leading-6 text-slate-500">
              Tambah, edit, lihat detail, dan atur status lowongan rekrutmen.
            </p>
          </div>
          <span class="grid h-11 w-11 shrink-0 place-items-center rounded-2xl bg-indigo-100 text-xl text-indigo-700 transition group-hover:bg-indigo-600 group-hover:text-white">
            →
          </span>
        </div>
      </NuxtLink>

      <div class="dashboard-card border-dashed">
        <p class="text-sm font-bold uppercase tracking-wide text-slate-400">Segera hadir</p>
        <h2 class="mt-2 text-xl font-black text-slate-700">Modul Berikutnya</h2>
        <p class="mt-2 text-sm leading-6 text-slate-500">
          Area ini siap dipakai saat menu rekrutmen baru mulai ditambahkan.
        </p>
      </div>
    </section>
  </div>
</template>
