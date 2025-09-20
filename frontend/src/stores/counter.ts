import { defineStore } from 'pinia'
import { authApi } from '@/api'
import type { LoginRequest, UserInfo } from '@/types'
import { ElMessage } from 'element-plus'
import router from '@/router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('verkeyoss_token') || '',
    userInfo: null as UserInfo | null,
    isLoggedIn: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && state.isLoggedIn
  },

  actions: {
    // 登录
    async login(loginData: LoginRequest) {
      try {
        const response = await authApi.login(loginData)
        const { token, user_info } = response.data.data!
        
        this.token = token
        this.userInfo = user_info
        this.isLoggedIn = true
        
        // 存储到本地
        localStorage.setItem('verkeyoss_token', token)
        
        ElMessage.success('登录成功')
        router.push('/dashboard')
      } catch (error) {
        console.error('登录失败:', error)
        throw error
      }
    },

    // 登出
    logout() {
      this.token = ''
      this.userInfo = null
      this.isLoggedIn = false
      
      localStorage.removeItem('verkeyoss_token')
      router.push('/login')
      ElMessage.success('已退出登录')
    },

    // 获取用户信息
    async fetchUserInfo() {
      try {
        const response = await authApi.getUserInfo()
        this.userInfo = response.data.data!
        this.isLoggedIn = true
      } catch (error) {
        console.error('获取用户信息失败:', error)
        this.logout()
      }
    },

    // 初始化认证状态
    async initAuth() {
      if (this.token) {
        try {
          await this.fetchUserInfo()
        } catch (error) {
          this.logout()
        }
      }
    },

    // 修改密码
    async changePassword(oldPassword: string, newPassword: string) {
      try {
        await authApi.changePassword({
          old_password: oldPassword,
          new_password: newPassword
        })
        ElMessage.success('密码修改成功')
      } catch (error) {
        console.error('修改密码失败:', error)
        throw error
      }
    }
  }
})
