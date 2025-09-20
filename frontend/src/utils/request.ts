import axios, { type AxiosResponse, type AxiosError } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

// 环境配置函数
function getBaseURL() {
  // 优先使用环境变量配置
  const envBaseURL = import.meta.env.VITE_API_BASE_URL
  if (envBaseURL) {
    return envBaseURL
  }
  
  // 备用方案：根据运行环境自动判断
  // 在开发环境下（Vite dev server），使用相对路径，代理会转发到后端
  // 在正式环境下（embed模式），使用相对路径，直接访问嵌入的后端
  return '/api'
}

// 创建 axios 实例
const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 在开发环境下显示调试信息
if (import.meta.env.DEV) {
  console.log('🔌 前端 API 基本地址:', getBaseURL())
  console.log('🌐 后端代理地址:', import.meta.env.VITE_BACKEND_URL || 'http://localhost:8913')
}

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 添加认证令牌
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data } = response
    
    // 如果响应包含 code 字段，根据 code 处理
    if (typeof data === 'object' && data !== null && 'code' in data) {
      if (data.code === 200) {
        return response
      } else {
        // 非 200 状态码的业务错误
        ElMessage.error(data.message || '请求失败')
        return Promise.reject(new Error(data.message || '请求失败'))
      }
    }
    
    return response
  },
  (error: AxiosError) => {
    // HTTP 状态码错误处理
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          ElMessage.error('认证失效，请重新登录')
          const authStore = useAuthStore()
          authStore.logout()
          router.push('/login')
          break
        case 403:
          ElMessage.error('权限不足')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          const errorMessage = (data as any)?.message || `请求失败 (${status})`
          ElMessage.error(errorMessage)
      }
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export default api