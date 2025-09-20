import { defineStore } from 'pinia'
import { dashboardApi } from '@/api'
import type { DashboardStats, Announcement } from '@/types'

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    stats: null as DashboardStats | null,
    announcements: [] as Announcement[],
    loading: false
  }),

  actions: {
    // 获取仪表盘统计数据
    async fetchStats() {
      this.loading = true
      try {
        const response = await dashboardApi.getStats()
        this.stats = response.data.data!
      } catch (error) {
        console.error('获取仪表盘数据失败:', error)
      } finally {
        this.loading = false
      }
    },

    // 获取公告列表
    async fetchAnnouncements() {
      try {
        const response = await dashboardApi.getAnnouncements()
        this.announcements = response.data.data!
      } catch (error) {
        console.error('获取公告列表失败:', error)
      }
    }
  }
})