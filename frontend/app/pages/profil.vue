<script setup>
definePageMeta({
  middleware: 'auth',
  layout: 'app'
})

const authStore = useAuthStore()
const currentUser = computed(() => authStore.user)

const userInitials = computed(() => {
  const name = currentUser.value?.name?.trim()

  if (!name) {
    return 'U'
  }

  return name
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part.charAt(0))
    .join('')
    .toUpperCase()
})
</script>

<template>
  <div class="space-y-6">
    <section>
      <p class="text-sm font-bold uppercase tracking-wide text-indigo-600">Akun</p>
      <h1 class="mt-1 text-3xl font-black tracking-tight text-slate-950">Profil Saya</h1>
      <p class="mt-1 text-sm text-slate-500">Informasi akun yang sedang digunakan.</p>
    </section>

    <section class="dashboard-card">
      <div class="flex flex-col gap-5 sm:flex-row sm:items-center">
        <div class="grid h-24 w-24 shrink-0 place-items-center rounded-3xl bg-gradient-to-br from-indigo-600 to-cyan-500 text-3xl font-black text-white shadow-xl shadow-indigo-200">
          {{ userInitials }}
        </div>

        <div class="min-w-0">
          <h2 class="truncate text-2xl font-black text-slate-950">{{ currentUser?.name || 'User' }}</h2>
          <p class="mt-1 truncate text-sm text-slate-500">{{ currentUser?.email || '-' }}</p>
          <span class="mt-3 inline-flex rounded-full bg-indigo-100 px-3 py-1 text-xs font-bold uppercase tracking-wide text-indigo-700">
            {{ currentUser?.role || '-' }}
          </span>
        </div>
      </div>

      <dl class="mt-8 grid gap-4 border-t border-slate-100 pt-6 sm:grid-cols-2">
        <div class="soft-panel">
          <dt class="text-xs font-bold uppercase tracking-wide text-slate-400">Nama</dt>
          <dd class="mt-2 font-semibold text-slate-800">{{ currentUser?.name || '-' }}</dd>
        </div>
        <div class="soft-panel">
          <dt class="text-xs font-bold uppercase tracking-wide text-slate-400">Email</dt>
          <dd class="mt-2 break-all font-semibold text-slate-800">{{ currentUser?.email || '-' }}</dd>
        </div>
        <div class="soft-panel">
          <dt class="text-xs font-bold uppercase tracking-wide text-slate-400">Role</dt>
          <dd class="mt-2 font-semibold capitalize text-slate-800">{{ currentUser?.role || '-' }}</dd>
        </div>
        <div class="soft-panel">
          <dt class="text-xs font-bold uppercase tracking-wide text-slate-400">User ID</dt>
          <dd class="mt-2 font-semibold text-slate-800">{{ currentUser?.id || '-' }}</dd>
        </div>
      </dl>
    </section>
  </div>
</template>
