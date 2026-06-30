<script setup>
const authStore = useAuthStore()

const loginForm = reactive({
  email: 'admin@example.com',
  password: ''
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
  <main style="min-height: 100vh; display: grid; place-items: center; font-family: Arial, sans-serif; background: #f5f5f5;">
    <form
      style="width: 360px; padding: 24px; border: 1px solid #ddd; border-radius: 8px; background: white;"
      @submit.prevent="login"
    >
      <h1 style="margin-top: 0;">Login</h1>
      <p style="color: #666;">Masuk ke Rekrutmen Playground</p>

      <div style="margin-bottom: 12px;">
        <label>Email</label>
        <input
          v-model="loginForm.email"
          type="email"
          required
          style="width: 100%; padding: 8px; margin-top: 4px; box-sizing: border-box;"
        >
      </div>

      <div style="margin-bottom: 12px;">
        <label>Password</label>
        <input
          v-model="loginForm.password"
          type="password"
          required
          style="width: 100%; padding: 8px; margin-top: 4px; box-sizing: border-box;"
        >
      </div>

      <p v-if="authStore.error" style="color: red;">{{ authStore.error }}</p>

      <button type="submit" :disabled="authStore.loading" style="width: 100%; padding: 10px;">
        {{ authStore.loading ? 'Memproses...' : 'Login' }}
      </button>
    </form>
  </main>
</template>
