type AuthUser = {
  id: number
  name: string
  email: string
  role: string
}

type LoginResponse = {
  success: boolean
  message: string
  data: {
    token: string
    user: AuthUser
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null as string | null,
    user: null as AuthUser | null,
    loading: false,
    error: ''
  }),

  getters: {
    isLoggedIn: (state) => Boolean(state.token),
    role: (state) => state.user?.role || '',
    authHeaders: (state) => {
      if (!state.token) {
        return {}
      }

      return {
        Authorization: `Bearer ${state.token}`
      }
    }
  },

  actions: {
    restoreSession() {
      const tokenCookie = useCookie<string | null>('auth_token')
      const userCookie = useCookie<AuthUser | null>('auth_user')

      this.token = tokenCookie.value || null
      this.user = userCookie.value || null
    },

    async login(email: string, password: string) {
      this.loading = true
      this.error = ''

      try {
        const result = await $fetch<LoginResponse>('http://localhost:8080/api/auth/login', {
          method: 'POST',
          body: {
            email,
            password
          }
        })

        this.setSession(result.data.token, result.data.user)
      } catch (error) {
        this.error = 'Login gagal. Cek email dan password.'
        throw error
      } finally {
        this.loading = false
      }
    },

    setSession(token: string, user: AuthUser) {
      const tokenCookie = useCookie<string | null>('auth_token', {
        maxAge: 60 * 60 * 24
      })
      const userCookie = useCookie<AuthUser | null>('auth_user', {
        maxAge: 60 * 60 * 24
      })

      tokenCookie.value = token
      userCookie.value = user

      this.token = token
      this.user = user
    },

    logout(message = '') {
      const tokenCookie = useCookie<string | null>('auth_token')
      const userCookie = useCookie<AuthUser | null>('auth_user')

      tokenCookie.value = null
      userCookie.value = null

      this.token = null
      this.user = null
      this.error = message
    },

    handleUnauthorized(error: any) {
      if (error?.statusCode === 401 || error?.response?.status === 401) {
        this.logout('Sesi login habis. Silakan login ulang.')
        return true
      }

      return false
    }
  }
})
