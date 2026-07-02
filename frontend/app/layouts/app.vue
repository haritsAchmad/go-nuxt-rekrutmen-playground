<script setup>
const authStore = useAuthStore()
const route = useRoute()
const isSidebarOpen = ref(false)
const isProfileMenuOpen = ref(false)
const profileMenu = ref(null)

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

const navigation = [
  { label: 'Dashboard', to: '/' },
  { label: 'Lowongan', to: '/lowongan' }
]

const pageTitle = computed(() => {
  const titles = {
    '/': 'Dashboard',
    '/lowongan': 'Kelola Lowongan',
    '/profil': 'Profil Saya'
  }

  return titles[route.path] || 'Sistem Rekrutmen'
})

watch(
  () => route.path,
  () => {
    isSidebarOpen.value = false
    isProfileMenuOpen.value = false
  }
)

function handleOutsideClick(event) {
  if (profileMenu.value && !profileMenu.value.contains(event.target)) {
    isProfileMenuOpen.value = false
  }
}

function handleEscape(event) {
  if (event.key === 'Escape') {
    isProfileMenuOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
  document.addEventListener('keydown', handleEscape)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleOutsideClick)
  document.removeEventListener('keydown', handleEscape)
})

async function logout() {
  isProfileMenuOpen.value = false
  authStore.logout()
  await navigateTo('/login')
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 text-slate-900">
    <div
      v-if="isSidebarOpen"
      class="fixed inset-0 z-40 bg-slate-950/50 lg:hidden"
      @click="isSidebarOpen = false"
    />

    <aside
      class="fixed inset-y-0 left-0 z-50 flex w-72 flex-col bg-slate-950 px-5 py-6 text-white shadow-2xl transition-transform lg:translate-x-0"
      :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <div class="flex items-center justify-between gap-3 px-2">
        <NuxtLink to="/" class="flex min-w-0 items-center gap-3">
          <span class="grid h-11 w-11 shrink-0 place-items-center rounded-2xl bg-gradient-to-br from-indigo-500 to-cyan-400 text-lg font-black shadow-lg shadow-indigo-950/50">
            RP
          </span>
          <span class="min-w-0">
            <span class="block truncate text-xs font-bold uppercase tracking-[0.2em] text-indigo-300">Rekrutmen</span>
            <span class="mt-1 block truncate text-xl font-black">Playground</span>
          </span>
        </NuxtLink>

        <button
          type="button"
          class="grid h-9 w-9 place-items-center p-0 text-slate-300 hover:bg-white/10 lg:hidden"
          aria-label="Tutup menu"
          @click="isSidebarOpen = false"
        >
          ✕
        </button>
      </div>

      <nav class="mt-9 space-y-2">
        <NuxtLink
          v-for="item in navigation"
          :key="item.to"
          :to="item.to"
          class="block rounded-xl px-4 py-3 text-sm font-bold transition"
          :class="route.path === item.to
            ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-950/40'
            : 'text-slate-300 hover:bg-white/10 hover:text-white'"
        >
          {{ item.label }}
        </NuxtLink>
      </nav>

      <div class="mt-auto rounded-2xl border border-white/10 bg-white/5 p-4">
        <p class="truncate text-sm font-bold">{{ currentUser?.name || 'User' }}</p>
        <p class="mt-1 truncate text-xs text-slate-400">{{ currentUser?.role || '-' }}</p>
      </div>
    </aside>

    <div class="lg:pl-72">
      <header class="sticky top-0 z-30 border-b border-slate-200/80 bg-white/90 backdrop-blur">
        <div class="flex h-18 items-center justify-between gap-4 px-4 sm:px-6 lg:px-8">
          <div class="flex min-w-0 items-center gap-3">
            <button
              type="button"
              class="muted-button grid h-10 w-10 place-items-center p-0 lg:hidden"
              aria-label="Buka menu"
              @click="isSidebarOpen = true"
            >
              ☰
            </button>
            <div class="min-w-0">
              <p class="truncate text-sm font-black text-slate-900">
                {{ pageTitle }}
              </p>
              <p class="truncate text-xs text-slate-500">Sistem Rekrutmen</p>
            </div>
          </div>

          <div ref="profileMenu" class="relative">
            <button
              type="button"
              class="flex items-center gap-3 rounded-2xl border border-slate-200 bg-white px-2 py-1.5 text-left shadow-sm hover:bg-slate-50"
              :aria-expanded="isProfileMenuOpen"
              aria-haspopup="menu"
              @click.stop="isProfileMenuOpen = !isProfileMenuOpen"
            >
              <span class="grid h-9 w-9 shrink-0 place-items-center rounded-xl bg-gradient-to-br from-indigo-600 to-cyan-500 text-xs font-black text-white">
                {{ userInitials }}
              </span>
              <span class="hidden min-w-0 sm:block">
                <span class="block max-w-36 truncate text-sm font-bold text-slate-800">{{ currentUser?.name || 'User' }}</span>
                <span class="block max-w-36 truncate text-xs text-slate-500">{{ currentUser?.role || '-' }}</span>
              </span>
              <span
                class="mr-1 text-xs text-slate-400 transition-transform"
                :class="isProfileMenuOpen ? 'rotate-180' : ''"
              >
                ▼
              </span>
            </button>

            <Transition name="profile-menu">
              <div
                v-if="isProfileMenuOpen"
                class="absolute right-0 top-full z-50 mt-2 w-56 overflow-hidden rounded-2xl border border-slate-200 bg-white p-2 shadow-xl shadow-slate-300/60"
                role="menu"
              >
                <div class="border-b border-slate-100 px-3 py-2 sm:hidden">
                  <p class="truncate text-sm font-bold text-slate-800">{{ currentUser?.name || 'User' }}</p>
                  <p class="truncate text-xs text-slate-500">{{ currentUser?.role || '-' }}</p>
                </div>

                <NuxtLink
                  to="/profil"
                  class="mt-1 block rounded-xl px-3 py-2.5 text-sm font-semibold text-slate-700 hover:bg-slate-100"
                  role="menuitem"
                >
                  Lihat Profil
                </NuxtLink>
                <button
                  type="button"
                  class="mt-1 w-full px-3 py-2.5 text-left text-sm font-semibold text-rose-600 hover:bg-rose-50"
                  role="menuitem"
                  @click="logout"
                >
                  Logout
                </button>
              </div>
            </Transition>
          </div>
        </div>
      </header>

      <main class="px-4 py-6 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.profile-menu-enter-active,
.profile-menu-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.profile-menu-enter-from,
.profile-menu-leave-to {
  opacity: 0;
  transform: translateY(-0.375rem);
}
</style>
