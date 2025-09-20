import { defineStore } from 'pinia'
import { appApi, versionApi } from '@/api'
import type { App, CreateAppRequest, UpdateAppRequest, Version, CreateVersionRequest, UpdateVersionRequest, PaginationParams } from '@/types'
import { ElMessage } from 'element-plus'

export const useAppStore = defineStore('app', {
  state: () => ({
    apps: [] as App[],
    currentApp: null as App | null,
    versions: [] as Version[],
    loading: false,
    total: 0,
    currentPage: 1,
    pageSize: 10
  }),

  actions: {
    // 获取应用列表
    async fetchApps(params: PaginationParams) {
      this.loading = true
      try {
        const response = await appApi.getList(params)
        const { list, total } = response.data.data!
        this.apps = list
        this.total = total
        this.currentPage = params.page
        this.pageSize = params.size
      } catch (error) {
        console.error('获取应用列表失败:', error)
      } finally {
        this.loading = false
      }
    },

    // 创建应用
    async createApp(data: CreateAppRequest) {
      try {
        const response = await appApi.create(data)
        ElMessage.success('应用创建成功')
        // 刷新列表
        await this.fetchApps({ page: this.currentPage, size: this.pageSize })
        return response.data.data
      } catch (error) {
        console.error('创建应用失败:', error)
        throw error
      }
    },

    // 更新应用
    async updateApp(akey: string, data: UpdateAppRequest) {
      try {
        await appApi.update(akey, data)
        ElMessage.success('应用更新成功')
        // 刷新列表
        await this.fetchApps({ page: this.currentPage, size: this.pageSize })
      } catch (error) {
        console.error('更新应用失败:', error)
        throw error
      }
    },

    // 删除应用
    async deleteApp(akey: string) {
      try {
        await appApi.delete(akey)
        ElMessage.success('应用删除成功')
        // 刷新列表
        await this.fetchApps({ page: this.currentPage, size: this.pageSize })
      } catch (error) {
        console.error('删除应用失败:', error)
        throw error
      }
    },

    // 设置当前应用
    setCurrentApp(app: App) {
      this.currentApp = app
    },

    // 获取版本列表
    async fetchVersions(akey: string, params: PaginationParams) {
      this.loading = true
      try {
        const response = await versionApi.getList(akey, params)
        const { list, total } = response.data.data!
        this.versions = list
        this.total = total
        this.currentPage = params.page
        this.pageSize = params.size
      } catch (error) {
        console.error('获取版本列表失败:', error)
      } finally {
        this.loading = false
      }
    },

    // 创建版本
    async createVersion(akey: string, data: CreateVersionRequest) {
      try {
        const response = await versionApi.create(akey, data)
        ElMessage.success('版本创建成功')
        // 刷新列表
        await this.fetchVersions(akey, { page: this.currentPage, size: this.pageSize })
        return response.data.data
      } catch (error) {
        console.error('创建版本失败:', error)
        throw error
      }
    },

    // 更新版本
    async updateVersion(vkey: string, data: UpdateVersionRequest, akey: string) {
      try {
        await versionApi.update(vkey, data)
        ElMessage.success('版本更新成功')
        // 刷新列表
        await this.fetchVersions(akey, { page: this.currentPage, size: this.pageSize })
      } catch (error) {
        console.error('更新版本失败:', error)
        throw error
      }
    },

    // 删除版本
    async deleteVersion(vkey: string, akey: string) {
      try {
        await versionApi.delete(vkey)
        ElMessage.success('版本删除成功')
        // 刷新列表
        await this.fetchVersions(akey, { page: this.currentPage, size: this.pageSize })
      } catch (error) {
        console.error('删除版本失败:', error)
        throw error
      }
    }
  }
})