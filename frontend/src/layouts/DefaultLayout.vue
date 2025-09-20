<template>
  <div class="app-container">
    <!-- 移动端遮罩层 -->
    <div 
      v-if="isMobile && showSidebar" 
      class="sidebar-overlay"
      @click="closeSidebar"
    ></div>
    
    <!-- 左侧边栏 -->
    <div :class="['app-sidebar', { 'mobile-hidden': isMobile && !showSidebar }]">
      <!-- 侧边栏标题区域 -->
      <div class="sidebar-header">
        <div class="logo-area">
          <h1 class="logo-title">VerKeyOSS</h1>
          <span class="logo-subtitle">应用版本管理系统</span>
        </div>
      </div>

      <!-- 菜单区域 -->
      <div class="sidebar-menu-container">
        <el-menu
          :default-active="$route.path"
          :unique-opened="true"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
          class="sidebar-menu"
        >
          <el-menu-item index="/dashboard">
            <el-icon><Odometer /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>
          
          <el-menu-item index="/apps">
            <el-icon><Grid /></el-icon>
            <template #title>应用管理</template>
          </el-menu-item>
          
          <el-menu-item index="/tools">
            <el-icon><Tools /></el-icon>
            <template #title>校验工具</template>
          </el-menu-item>
        </el-menu>
      </div>
    </div>

    <!-- 右侧内容区域 -->
    <div :class="['app-content', { 'mobile-full': isMobile }]">
      <!-- 顶部导航栏 -->
      <div class="app-header">
        <div class="header-content">
          <div class="header-left">
            <!-- 移动端菜单按钮 -->
            <el-button 
              v-if="isMobile"
              class="mobile-menu-btn"
              @click="toggleSidebar"
              type="text"
            >
              <el-icon><Menu /></el-icon>
            </el-button>
          </div>
          
          <div class="header-actions">
            <!-- 版本信息 -->
            <div v-if="backendVersion" class="version-info">
              <el-tag size="small" type="info">v{{ backendVersion }}</el-tag>
            </div>
            
            <el-dropdown @command="handleCommand">
              <span class="el-dropdown-link">
                <el-icon><User /></el-icon>
                {{ userInfo?.username }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人设置</el-dropdown-item>
                  <el-dropdown-item command="password">修改密码</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>

      <!-- 主内容区域 -->
      <div class="app-main">
        <router-view />
      </div>
    </div>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="passwordDialogVisible"
      title="修改密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="80px"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            show-password
            placeholder="请输入原密码"
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            show-password
            placeholder="请输入新密码"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="passwordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleChangePassword">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { NetworkChecker } from '@/utils/network'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  User,
  ArrowDown,
  Odometer,
  Grid,
  Tools,
  Menu
} from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

// 移动端侧边栏状态
const isMobile = ref(false)
const showSidebar = ref(false)

// 后端版本信息
const backendVersion = ref<string>('')

// 用户信息
const userInfo = computed(() => authStore.userInfo)

// 修改密码对话框
const passwordDialogVisible = ref(false)
const passwordFormRef = ref<FormInstance>()
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码表单验证规则
const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 监听窗口尺寸变化
const handleResize = () => {
  const mobile = window.innerWidth < 768
  isMobile.value = mobile
  
  // 小屏幕时隐藏侧边栏，大屏幕时显示侣边栏
  if (mobile) {
    showSidebar.value = false
  } else {
    showSidebar.value = true
  }
}

// 切换侧边栏显示状态
const toggleSidebar = () => {
  showSidebar.value = !showSidebar.value
}

// 点击遮罩层关闭侧边栏
const closeSidebar = () => {
  if (isMobile.value) {
    showSidebar.value = false
  }
}

// 处理下拉菜单命令
const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      // 暂时不实现个人设置
      ElMessage.info('个人设置功能开发中')
      break
    case 'password':
      passwordDialogVisible.value = true
      break
    case 'logout':
      handleLogout()
      break
  }
}

// 处理登出
const handleLogout = () => {
  ElMessageBox.confirm(
    '确定要退出登录吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    authStore.logout()
  }).catch(() => {
    // 取消登出
  })
}

// 处理修改密码
const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    await authStore.changePassword(
      passwordForm.value.oldPassword,
      passwordForm.value.newPassword
    )
    
    passwordDialogVisible.value = false
    passwordForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (error) {
    console.error('修改密码失败:', error)
  }
}

// 初始化
onMounted(async () => {
  if (!authStore.userInfo) {
    await authStore.fetchUserInfo()
  }
  
  // 获取后端版本信息
  try {
    const networkChecker = NetworkChecker.getInstance()
    const versionResult = await networkChecker.getBackendVersion()
    if (versionResult.version) {
      backendVersion.value = versionResult.version
    }
  } catch (error) {
    console.warn('获取后端版本信息失败:', error)
  }
  
  // 初始化屏幕尺寸检测
  handleResize()
  window.addEventListener('resize', handleResize)
})

// 清理事件监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
/* 防止横向滚动条 */
* {
  box-sizing: border-box;
}

.app-container {
  height: 100vh;
  width: 100vw;
  position: relative;
  overflow: hidden;
}

/* 移动端遮罩层 */
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

/* 左侧边栏样式 */
.app-sidebar {
  width: 200px;
  background: #304156;
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  z-index: 1000;
  overflow: hidden;
  transition: transform 0.3s ease;
}

.app-sidebar.mobile-hidden {
  transform: translateX(-100%);
}

/* 侧边栏标题区域 */
.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  min-height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-area {
  text-align: center;
}

.logo-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #409eff;
  line-height: 1.2;
}

.logo-subtitle {
  font-size: 12px;
  color: #bfcbd9;
  margin-top: 4px;
  display: block;
}

/* 菜单容器 */
.sidebar-menu-container {
  flex: 1;
  overflow-y: auto;
}

.sidebar-menu {
  border: none;
  background: transparent;
}

/* 修复菜单项对齐问题 */
.sidebar-menu .el-menu-item {
  display: flex;
  align-items: center;
  padding-left: 20px;
  height: 56px;
}

.sidebar-menu .el-menu-item .el-icon {
  margin-right: 8px;
  font-size: 16px;
}

/* 右侧内容区域 */
.app-content {
  position: fixed;
  right: 0;
  top: 0;
  height: 100vh;
  width: calc(100vw - 200px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-content.mobile-full {
  width: 100vw;
  left: 0;
  right: auto;
}

/* 顶部导航栏 */
.app-header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  height: 60px;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  z-index: 999;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 20px;
}

.header-left {
  flex: 1;
  display: flex;
  align-items: center;
}

.mobile-menu-btn {
  color: #409eff;
  font-size: 18px;
  padding: 8px;
  margin-right: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mobile-menu-btn:hover {
  background-color: #f5f7fa;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.version-info {
  display: flex;
  align-items: center;
}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 主内容区域 */
.app-main {
  flex: 1;
  background: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
  overflow-x: hidden;
  min-width: 0;
}

.dialog-footer {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-sidebar {
    z-index: 1001;
  }
  
  .sidebar-header {
    padding: 15px;
    min-height: 60px;
  }
  
  .logo-title {
    font-size: 18px;
  }
  
  .logo-subtitle {
    font-size: 11px;
  }
  
  .app-header {
    height: 50px;
  }
  
  .header-content {
    padding: 0 15px;
  }
  
  .app-main {
    padding: 15px;
  }
}

@media (max-width: 480px) {
  .app-sidebar {
    width: 180px;
  }
  
  .sidebar-header {
    padding: 8px;
    min-height: 50px;
  }
  
  .logo-title {
    font-size: 16px;
  }
  
  .logo-subtitle {
    font-size: 10px;
  }
  
  .app-header {
    height: 45px;
  }
  
  .header-content {
    padding: 0 8px;
  }
  
  .app-main {
    padding: 8px;
  }
}
</style>