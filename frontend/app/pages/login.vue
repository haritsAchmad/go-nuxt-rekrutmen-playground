<script setup>
definePageMeta({
  layout: false
})

const authStore = useAuthStore()

const loginForm = reactive({
  email: 'admin@example.com',
  password: 'admin123'
})

onMounted(() => {
  authStore.restoreSession()

  if (authStore.isLoggedIn) {
    navigateTo('/')
  }
})

async function login() {
  try {
    await authStore.login(loginForm.email, loginForm.password)
    await navigateTo('/')
  } catch (error) {
    // error message sudah disimpan di authStore
  }
}
</script>

<template>
  <main class="relative min-h-screen overflow-hidden bg-slate-950 px-4 py-10 text-slate-100">
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_top_left,_rgba(99,102,241,0.35),_transparent_35%),radial-gradient(circle_at_bottom_right,_rgba(20,184,166,0.25),_transparent_35%)]"></div>
    <div class="absolute left-1/2 top-20 h-72 w-72 -translate-x-1/2 rounded-full bg-indigo-500/20 blur-3xl"></div>

    <section class="relative mx-auto grid min-h-[calc(100vh-5rem)] max-w-6xl items-center gap-10 lg:grid-cols-[1.1fr_0.9fr]">
      <div class="hidden lg:block">
        <p class="mb-4 inline-flex rounded-full border border-white/10 bg-white/10 px-4 py-2 text-sm font-medium text-indigo-100 backdrop-blur">
          Rekrutmen Playground
        </p>
        <h1 class="max-w-2xl text-5xl font-black tracking-tight text-white">
          Kelola lowongan tanpa tampilan abu-abu menyedihkan.
        </h1>
        <p class="mt-5 max-w-xl text-lg leading-8 text-slate-300">
          Playground buat latihan Nuxt + Go: login, CRUD lowongan, bulk action, filter, dan pagination dalam satu dashboard ringan.
        </p>

        <div class="mt-8 grid max-w-xl grid-cols-3 gap-3">
          <div class="rounded-2xl border border-white/10 bg-white/10 p-4 backdrop-blur">
            <p class="text-2xl font-bold text-white">Nuxt</p>
            <p class="text-sm text-slate-300">Frontend</p>
          </div>
          <div class="rounded-2xl border border-white/10 bg-white/10 p-4 backdrop-blur">
            <p class="text-2xl font-bold text-white">Pinia</p>
            <p class="text-sm text-slate-300">State</p>
          </div>
          <div class="rounded-2xl border border-white/10 bg-white/10 p-4 backdrop-blur">
            <p class="text-2xl font-bold text-white">Go</p>
            <p class="text-sm text-slate-300">Backend</p>
          </div>
        </div>
      </div>

      <form
        class="mx-auto w-full max-w-md rounded-[2rem] border border-white/15 bg-white/95 p-8 text-slate-900 shadow-2xl shadow-indigo-950/30 backdrop-blur"
        @submit.prevent="login"
      >
        <div class="mb-8">
          <p class="mb-3 inline-flex rounded-full bg-indigo-50 px-3 py-1 text-xs font-bold uppercase tracking-wide text-indigo-700">
            Welcome back
          </p>
          <h1 class="text-3xl font-black tracking-tight text-slate-950">Login</h1>
          <p class="mt-2 text-sm text-slate-500">Masuk ke Rekrutmen Playground.</p>
        </div>

        <div class="space-y-5">
          <div>
            <label class="mb-2 block text-sm font-semibold text-slate-700">Email</label>
            <input
              v-model="loginForm.email"
              class="w-full"
              type="email"
              required
            >
          </div>

          <div>
            <label class="mb-2 block text-sm font-semibold text-slate-700">Password</label>
            <input
              v-model="loginForm.password"
              class="w-full"
              type="password"
              required
            >
          </div>
        </div>

        <p v-if="authStore.error" class="mt-5 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
          {{ authStore.error }}
        </p>

        <button type="submit" :disabled="authStore.loading" class="primary-button mt-6 w-full">
          {{ authStore.loading ? 'Memproses...' : 'Masuk Dashboard' }}
        </button>
      </form>
    </section>
  </main>
</template>
