// 用户认证相关类型
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: string
  user_info: {
    username: string
  }
}

export interface UserInfo {
  username: string
}

// 应用相关类型
export interface App {
  akey: string
  name: string
  description: string
  is_paid: boolean
  version_count: number
  created_at: string
}

export interface CreateAppRequest {
  name: string
  description: string
  is_paid: boolean
}

export interface UpdateAppRequest {
  name: string
  description: string
  is_paid: boolean
}

// 版本相关类型
export interface Version {
  vkey: string
  version: string
  description: string
  is_latest: boolean
  is_forced_update: boolean
  created_at: string
}

export interface CreateVersionRequest {
  version: string
  description: string
  is_latest: boolean
  is_forced_update: boolean
}

export interface UpdateVersionRequest {
  version: string
  description: string
  is_latest: boolean
  is_forced_update: boolean
}

// 校验相关类型
export interface CheckRequest {
  akey: string
  vkey: string
}

export interface ValidationResponse {
  valid: boolean
  message: string
  app_name?: string
  version?: string
}

export interface UpdateCheckResponse {
  has_update: boolean
  latest_version?: string
  forced_update: boolean
  message: string
}

// 仪表盘相关类型
export interface DashboardStats {
  total_apps: number
  total_versions: number
  recent_apps: App[]
  recent_versions: Version[]
}

export interface Announcement {
  id: number
  title: string
  content: string
  is_active: boolean
  publish_date: string
  url?: string
}

// API 响应通用类型
export interface ApiResponse<T = any> {
  code: number
  data?: T
  message?: string
}

// 分页相关类型
export interface PaginationParams {
  page: number
  size: number
}

export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  size: number
}

// 表格列表项类型
export interface TableColumn {
  prop: string
  label: string
  width?: number
  minWidth?: number
  formatter?: (row: any) => string
}