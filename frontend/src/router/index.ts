import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

// 布局组件
import DefaultLayout from '@/layouts/DefaultLayout.vue'

// 页面组件
import LoginView from '@/views/auth/LoginView.vue'
import DashboardView from '@/views/dashboard/DashboardView.vue'
import AppListView from '@/views/apps/AppListView.vue'
import AppDetailView from '@/views/apps/AppDetailView.vue'
import VersionListView from '@/views/versions/VersionListView.vue'
import ToolsView from '@/views/tools/ToolsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: LoginView,
      meta: {
        requiresAuth: false,
        title: '登录'
      }
    },
    {
      path: '/',
      component: DefaultLayout,
      meta: {
        requiresAuth: true
      },
      children: [
        {
          path: '',
          redirect: '/dashboard'
        },
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: DashboardView,
          meta: {
            title: '仪表盘'
          }
        },
        {
          path: 'apps',
          name: 'AppList',
          component: AppListView,
          meta: {
            title: '应用管理'
          }
        },
        {
          path: 'apps/:akey',
          name: 'AppDetail',
          component: AppDetailView,
          meta: {
            title: '应用详情'
          }
        },
        {
          path: 'apps/:akey/versions',
          name: 'VersionList',
          component: VersionListView,
          meta: {
            title: '版本管理'
          }
        },
        {
          path: 'tools',
          name: 'Tools',
          component: ToolsView,
          meta: {
            title: '校验工具'
          }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      redirect: '/dashboard'
    }
  ]
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - VerKeyOSS`
  } else {
    document.title = 'VerKeyOSS - 应用版本管理系统'
  }
  
  // 检查是否需要认证
  if (to.meta.requiresAuth !== false) {
    // 需要认证的路由
    if (!authStore.token) {
      // 没有token，跳转到登录页
      ElMessage.warning('请先登录')
      next('/login')
      return
    }
    
    // 有token但没有用户信息，尝试获取用户信息
    if (!authStore.userInfo) {
      try {
        await authStore.fetchUserInfo()
      } catch (error) {
        // 获取用户信息失败，可能token已过期
        authStore.logout()
        next('/login')
        return
      }
    }
  } else {
    // 不需要认证的路由（如登录页）
    if (to.path === '/login' && authStore.isAuthenticated) {
      // 已经登录，重定向到仪表盘
      next('/dashboard')
      return
    }
  }
  
  next()
})

export default router
