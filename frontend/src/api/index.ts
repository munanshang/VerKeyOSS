import request from '@/utils/request'
import type {
  LoginRequest,
  LoginResponse,
  UserInfo,
  App,
  CreateAppRequest,
  UpdateAppRequest,
  Version,
  CreateVersionRequest,
  UpdateVersionRequest,
  CheckRequest,
  ValidationResponse,
  UpdateCheckResponse,
  DashboardStats,
  Announcement,
  ApiResponse,
  PaginatedResponse,
  PaginationParams
} from '@/types'

// 认证相关 API
export const authApi = {
  // 登录
  login: (data: LoginRequest) =>
    request.post<ApiResponse<LoginResponse>>('/auth/login', data),
  
  // 获取用户信息
  getUserInfo: () =>
    request.get<ApiResponse<UserInfo>>('/auth/user-info'),
  
  // 修改密码
  changePassword: (data: { old_password: string; new_password: string }) =>
    request.put<ApiResponse>('/auth/password', data)
}

// 应用管理相关 API
export const appApi = {
  // 创建应用
  create: (data: CreateAppRequest) =>
    request.post<ApiResponse<App>>('/app', data),
  
  // 获取应用列表
  getList: (params: PaginationParams) =>
    request.get<ApiResponse<PaginatedResponse<App>>>('/app', { params }),
  
  // 更新应用
  update: (akey: string, data: UpdateAppRequest) =>
    request.put<ApiResponse>(`/app/${akey}`, data),
  
  // 删除应用
  delete: (akey: string) =>
    request.delete<ApiResponse>(`/app/${akey}`)
}

// 版本管理相关 API
export const versionApi = {
  // 创建版本
  create: (akey: string, data: CreateVersionRequest) =>
    request.post<ApiResponse<Version>>(`/app/${akey}/versions`, data),
  
  // 获取版本列表
  getList: (akey: string, params: PaginationParams) =>
    request.get<ApiResponse<PaginatedResponse<Version>>>(`/app/${akey}/versions`, { params }),
  
  // 更新版本
  update: (vkey: string, data: UpdateVersionRequest) =>
    request.put<ApiResponse>(`/versions/${vkey}`, data),
  
  // 删除版本
  delete: (vkey: string) =>
    request.delete<ApiResponse>(`/versions/${vkey}`)
}

// 校验相关 API
export const checkApi = {
  // 校验 AKey 和 VKey
  validate: (data: CheckRequest) =>
    request.post<ApiResponse<ValidationResponse>>('/check/validate', data),
  
  // 检查更新
  checkUpdate: (data: CheckRequest) =>
    request.post<ApiResponse<UpdateCheckResponse>>('/check/update', data)
}

// 仪表盘相关 API
export const dashboardApi = {
  // 获取仪表盘统计数据
  getStats: () =>
    request.get<ApiResponse<DashboardStats>>('/dashboard/stats'),
  
  // 获取公告列表
  getAnnouncements: () =>
    request.get<ApiResponse<Announcement[]>>('/dashboard/announcements')
}