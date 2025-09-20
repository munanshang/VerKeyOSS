import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import App from './App.vue'
import router from './router'
import '@/assets/styles/global.css'
import { NetworkChecker } from './utils/network'

const app = createApp(App)

// 注册 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus, {
  locale: zhCn,
})

// 在开发环境下检测网络连接
if (import.meta.env.DEV) {
  const networkChecker = NetworkChecker.getInstance()
  // 延迟检测，等待应用初始化完成
  setTimeout(() => {
    networkChecker.showNetworkStatus()
  }, 1000)
}

app.mount('#app')
