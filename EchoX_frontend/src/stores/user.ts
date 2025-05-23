import { defineStore } from 'pinia'
import axios from 'axios'

interface User {
  id: number
  username: string
  email: string
}

interface UserState {
  token: string | null
  user: User | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem('token'),
    user: null
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    currentUser: (state) => state.user
  },

  actions: {
    async register(username: string, password: string, email: string) {
      try {
        await axios.post('http://localhost:8080/api/v1/auth/register', {
          username,
          password,
          email
        })
        return true
      } catch (error) {
        console.error('Registration failed:', error)
        return false
      }
    },

    async login(username: string, password: string) {
      try {
        const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
          username,
          password
        })

        const { token, user } = response.data
        this.token = token
        this.user = user
        localStorage.setItem('token', token)
        return true
      } catch (error) {
        console.error('Login failed:', error)
        return false
      }
    },

    async fetchProfile() {
      try {
        const response = await axios.get('http://localhost:8080/api/v1/users/profile', {
          headers: {
            Authorization: `Bearer ${this.token}`
          }
        })
        this.user = response.data
      } catch (error) {
        console.error('Failed to fetch profile:', error)
        this.logout()
      }
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
    }
  }
}) 