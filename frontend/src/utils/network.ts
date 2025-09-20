// 网络连接检测工具
export class NetworkChecker {
  private static instance: NetworkChecker
  private backendUrl: string
  
  private constructor() {
    this.backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8913'
  }
  
  public static getInstance(): NetworkChecker {
    if (!NetworkChecker.instance) {
      NetworkChecker.instance = new NetworkChecker()
    }
    return NetworkChecker.instance
  }
  
  // 获取后端版本信息
  async getBackendVersion(): Promise<{ version?: string; error?: string }> {
    try {
      const url = import.meta.env.DEV 
        ? `${this.backendUrl}/api/check/health`
        : '/api/check/health'
      
      const response = await fetch(url, {
        method: 'GET',
        signal: AbortSignal.timeout(5000)
      })
      
      if (response.ok) {
        const data = await response.json()
        return { version: data.data?.version || 'unknown' }
      } else {
        return { error: `获取版本信息失败 (状态码: ${response.status})` }
      }
    } catch (error: any) {
      return { error: `获取版本信息失败: ${error.message}` }
    }
  }

  // 检测后端连接状态
  async checkBackendConnection(): Promise<{ connected: boolean; message: string; latency?: number }> {
    const startTime = Date.now()
    
    try {
      // 在开发环境下，检测后端是否运行
      if (import.meta.env.DEV) {
        const response = await fetch(`${this.backendUrl}/api/check/health`, {
          method: 'GET',
          mode: 'cors',
          signal: AbortSignal.timeout(5000) // 5秒超时
        })
        
        const latency = Date.now() - startTime
        
        if (response.ok) {
          return {
            connected: true,
            message: `后端连接正常 (${latency}ms)`,
            latency
          }
        } else {
          return {
            connected: false,
            message: `后端响应异常 (状态码: ${response.status})`
          }
        }
      } else {
        // 正式环境下，检测内置的API
        const response = await fetch('/api/check/health', {
          method: 'GET',
          signal: AbortSignal.timeout(5000)
        })
        
        const latency = Date.now() - startTime
        
        if (response.ok) {
          return {
            connected: true,
            message: `服务连接正常 (${latency}ms)`,
            latency
          }
        } else {
          return {
            connected: false,
            message: `服务响应异常 (状态码: ${response.status})`
          }
        }
      }
    } catch (error: any) {
      const latency = Date.now() - startTime
      
      if (error.name === 'TimeoutError') {
        return {
          connected: false,
          message: `连接超时 (>${latency}ms)`
        }
      } else if (error.name === 'TypeError' && error.message.includes('fetch')) {
        return {
          connected: false,
          message: '网络连接失败，请检查网络或后端服务'
        }
      } else {
        return {
          connected: false,
          message: `连接错误: ${error.message}`
        }
      }
    }
  }
  
  // 显示网络状态
  async showNetworkStatus(): Promise<void> {
    if (import.meta.env.DEV) {
      console.group('🌐 网络连接状态检测')
      
      const status = await this.checkBackendConnection()
      
      if (status.connected) {
        console.log(`✅ ${status.message}`)
      } else {
        console.warn(`❌ ${status.message}`)
        console.log('💡 解决建议：')
        console.log('   1. 确保后端服务已启动 (verkeyoss.exe)')
        console.log('   2. 检查后端端口配置 (默认: 8913)')
        console.log('   3. 确认没有防火墙阻止连接')
      }
      
      console.groupEnd()
    }
  }
}