<template>
  <div class="dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1 class="page-title">仪表盘</h1>
      <p class="page-subtitle">查看系统概览和最新动态</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
        <div class="stat-card">
          <div class="stat-icon apps">
            <el-icon size="32"><Grid /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats?.total_apps || 0 }}</div>
            <div class="stat-label">应用总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
        <div class="stat-card">
          <div class="stat-icon versions">
            <el-icon size="32"><DocumentCopy /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats?.total_versions || 0 }}</div>
            <div class="stat-label">版本总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
        <div class="stat-card">
          <div class="stat-icon recent">
            <el-icon size="32"><Clock /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ recentAppsCount }}</div>
            <div class="stat-label">最近应用</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
        <div class="stat-card">
          <div class="stat-icon active">
            <el-icon size="32"><Check /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ recentVersionsCount }}</div>
            <div class="stat-label">最近版本</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 内容区域 -->
    <el-row :gutter="20">
      <!-- 最近应用 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>最近应用</span>
              <el-button type="text" @click="$router.push('/apps')">
                查看全部 <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
          </template>
          
          <div v-if="loading" class="loading-container">
            <el-skeleton :rows="3" animated />
          </div>
          
          <div v-else-if="!stats?.recent_apps?.length" class="empty-container">
            <el-empty description="暂无应用数据" :image-size="80" />
          </div>
          
          <div v-else class="recent-list">
            <div
              v-for="app in stats.recent_apps"
              :key="app.akey"
              class="recent-item"
              @click="handleAppClick(app)"
            >
              <div class="recent-info">
                <div class="recent-title">{{ app.name }}</div>
                <div class="recent-meta">
                  <span class="recent-akey">{{ app.akey }}</span>
                  <span class="recent-time">{{ formatDate(app.created_at) }}</span>
                </div>
                <div v-if="app.description" class="recent-desc">{{ app.description }}</div>
              </div>
              <div class="recent-badge">
                <el-tag v-if="app.is_paid" type="warning" size="small">付费</el-tag>
                <el-tag type="info" size="small">{{ app.version_count }} 个版本</el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 最近版本 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>最近版本</span>
            </div>
          </template>
          
          <div v-if="loading" class="loading-container">
            <el-skeleton :rows="3" animated />
          </div>
          
          <div v-else-if="!stats?.recent_versions?.length" class="empty-container">
            <el-empty description="暂无版本数据" :image-size="80" />
          </div>
          
          <div v-else class="recent-list">
            <div
              v-for="version in stats.recent_versions"
              :key="version.vkey"
              class="recent-item"
            >
              <div class="recent-info">
                <div class="recent-title">版本 {{ version.version }}</div>
                <div class="recent-meta">
                  <span class="recent-vkey">{{ version.vkey }}</span>
                  <span class="recent-time">{{ formatDate(version.created_at) }}</span>
                </div>
                <div v-if="version.description" class="recent-desc">{{ version.description }}</div>
              </div>
              <div class="recent-badge">
                <el-tag v-if="version.is_latest" type="success" size="small">最新</el-tag>
                <el-tag v-if="version.is_forced_update" type="danger" size="small">强制</el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统信息区域 -->
    <el-card class="dashboard-card system-info-card">
      <template #header>
        <div class="card-header">
          <span>系统信息</span>
        </div>
      </template>
      
      <div class="system-info">
        <div class="info-item">
          <span class="info-label">后端版本:</span>
          <el-tag v-if="systemInfo.backendVersion" type="success" size="small">
            v{{ systemInfo.backendVersion }}
          </el-tag>
          <span v-else class="info-value">加载中...</span>
        </div>
        <div class="info-item">
          <span class="info-label">最后更新:</span>
          <span class="info-value">{{ systemInfo.lastUpdate || '未知' }}</span>
        </div>
      </div>
    </el-card>

    <!-- 公告区域 -->
    <el-card v-if="announcements?.length" class="dashboard-card">
      <template #header>
        <div class="card-header">
          <span>系统公告</span>
        </div>
      </template>
      
      <div class="announcements">
        <div
          v-for="announcement in announcements"
          :key="announcement.id"
          class="announcement-item"
        >
          <div class="announcement-content">
            <h4 class="announcement-title">{{ announcement.title }}</h4>
            <p class="announcement-text">{{ announcement.content }}</p>
            <div class="announcement-meta">
              <span>{{ formatDate(announcement.publish_date) }}</span>
              <a
                v-if="announcement.url"
                :href="announcement.url"
                target="_blank"
                class="announcement-link"
              >
                查看详情 <el-icon><Link /></el-icon>
              </a>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDashboardStore } from '@/stores/dashboard'
import { NetworkChecker } from '@/utils/network'
import { formatDate } from '@/utils'
import type { App } from '@/types'
import {
  Grid,
  DocumentCopy,
  Clock,
  Check,
  ArrowRight,
  Link
} from '@element-plus/icons-vue'

const router = useRouter()
const dashboardStore = useDashboardStore()

// 数据状态
const loading = ref(false)
const stats = computed(() => dashboardStore.stats)
const announcements = computed(() => dashboardStore.announcements)

// 系统信息
const systemInfo = ref({
  backendVersion: '',
  lastUpdate: ''
})

// 计算属性
const recentAppsCount = computed(() => stats.value?.recent_apps?.length || 0)
const recentVersionsCount = computed(() => stats.value?.recent_versions?.length || 0)

// 处理应用点击
const handleAppClick = (app: App) => {
  router.push(`/apps/${app.akey}`)
}

// 初始化数据
const initData = async () => {
  loading.value = true
  try {
    await Promise.all([
      dashboardStore.fetchStats(),
      dashboardStore.fetchAnnouncements()
    ])
    
    // 获取后端版本信息
    try {
      const networkChecker = NetworkChecker.getInstance()
      const versionResult = await networkChecker.getBackendVersion()
      if (versionResult.version) {
        systemInfo.value.backendVersion = versionResult.version
      }
      systemInfo.value.lastUpdate = new Date().toLocaleString()
    } catch (error) {
      console.warn('获取系统信息失败:', error)
    }
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initData()
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.page-subtitle {
  color: #909399;
  margin: 0;
}

.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  color: #fff;
}

.stat-icon.apps {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.versions {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.recent {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.active {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.dashboard-card {
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.system-info {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-label {
  font-weight: 500;
  color: #606266;
  min-width: 80px;
}

.info-value {
  color: #303133;
  font-family: monospace;
  font-size: 14px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.loading-container,
.empty-container {
  padding: 20px;
}

.recent-list {
  max-height: 300px;
  overflow-y: auto;
}

.recent-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;
}

.recent-item:hover {
  background-color: #f8f9fa;
}

.recent-item:last-child {
  border-bottom: none;
}

.recent-info {
  flex: 1;
  min-width: 0;
}

.recent-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.recent-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.recent-akey,
.recent-vkey {
  font-family: monospace;
  font-size: 12px;
  background: #f0f2f5;
  padding: 2px 6px;
  border-radius: 4px;
  color: #606266;
}

.recent-time {
  font-size: 12px;
  color: #909399;
}

.recent-desc {
  font-size: 14px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.recent-badge {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex-shrink: 0;
}

.announcements {
  max-height: 400px;
  overflow-y: auto;
}

.announcement-item {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.announcement-item:last-child {
  border-bottom: none;
}

.announcement-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin: 0 0 8px 0;
}

.announcement-text {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
  margin: 0 0 8px 0;
}

.announcement-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
}

.announcement-link {
  color: #409eff;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 4px;
}

.announcement-link:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .stats-row {
    margin-bottom: 16px;
  }
  
  .stat-card {
    padding: 16px;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    margin-right: 12px;
  }
  
  .stat-value {
    font-size: 20px;
  }
}
</style>", "original_text": "", "replace_all": false}]