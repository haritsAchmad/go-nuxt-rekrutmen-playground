export default defineNuxtRouteMiddleware(() => {
  const authStore = useAuthStore()

  authStore.restoreSession()

  if (!authStore.isLoggedIn) {
    return navigateTo('/login')
  }
})
