<template>
  <div class="app-detail">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button @click="handleGoBack" type="text">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div class="app-info">
          <h1 class="app-name">{{ app?.name }}</h1>
          <div class="app-meta">
            <span class="app-akey">{{ app?.akey }}</span>
            <el-tag :type="app?.is_paid ? 'warning' : 'success'" size="small">
              {{ app?.is_paid ? '付费应用' : '免费应用' }}
            </el-tag>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <el-button @click="handleEdit">
          <el-icon><Edit /></el-icon>
          编辑应用
        </el-button>
        <el-button type="primary" @click="handleCreateVersion">
          <el-icon><Plus /></el-icon>
          创建版本
        </el-button>
      </div>
    </div>

    <!-- 应用描述 -->
    <el-card v-if="app?.description" class="description-card">
      <div class="description-content">
        <h3>应用描述</h3>
        <p>{{ app.description }}</p>
      </div>
    </el-card>

    <!-- 版本列表 -->
    <el-card class="versions-card">
      <template #header>
        <div class="card-header">
          <span>版本列表</span>
          <el-button type="text" @click="refreshVersions">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <DataTable
        :data="versions"
        :loading="loading"
        :total="total"
        :page="currentPage"
        :page-size="pageSize"
        @page-change="handlePageChange"
      >
        <el-table-column prop="version" label="版本号" width="120">
          <template #default="{ row }">
            <div class="version-cell">
              <span class="version-text">{{ row.version }}</span>
              <el-tag v-if="row.is_latest" type="success" size="small">最新</el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="vkey" label="版本标识" width="200">
          <template #default="{ row }">
            <div class="vkey-cell">
              <span class="vkey-text">{{ row.vkey }}</span>
              <el-button
                type="text"
                size="small"
                @click="copyToClipboard(row.vkey)"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="description" label="版本描述" min-width="200">
          <template #default="{ row }">
            <div class="version-description">
              {{ row.description || '暂无描述' }}
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="is_forced_update" label="强制更新" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_forced_update ? 'danger' : 'info'" size="small">
              {{ row.is_forced_update ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-space>
              <el-button type="text" size="small" @click="handleEditVersion(row)">
                编辑
              </el-button>
              <el-button
                type="text"
                size="small"
                @click="handleDeleteVersion(row)"
                style="color: #f56c6c"
              >
                删除
              </el-button>
            </el-space>
          </template>
        </el-table-column>
      </DataTable>
    </el-card>

    <!-- 编辑应用对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑应用"
      width="500px"
      :close-on-click-modal="false"
      @closed="handleEditDialogClosed"
    >
      <el-form
        ref="appFormRef"
        :model="appForm"
        :rules="appRules"
        label-width="80px"
      >
        <el-form-item label="应用名称" prop="name">
          <el-input v-model="appForm.name" placeholder="请输入应用名称" />
        </el-form-item>

        <el-form-item label="应用描述" prop="description">
          <el-input
            v-model="appForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入应用描述（可选）"
          />
        </el-form-item>

        <el-form-item label="应用类型" prop="is_paid">
          <el-radio-group v-model="appForm.is_paid">
            <el-radio :label="false">免费应用</el-radio>
            <el-radio :label="true">付费应用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleUpdateApp"
          >
            更新
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 创建/编辑版本对话框 -->
    <el-dialog
      v-model="showVersionDialog"
      :title="editingVersion ? '编辑版本' : '创建版本'"
      width="500px"
      :close-on-click-modal="false"
      @closed="handleVersionDialogClosed"
    >
      <el-form
        ref="versionFormRef"
        :model="versionForm"
        :rules="versionRules"
        label-width="80px"
      >
        <el-form-item label="版本号" prop="version">
          <el-input
            v-model="versionForm.version"
            placeholder="请输入版本号，如：1.0.0"
          />
        </el-form-item>

        <el-form-item label="版本描述" prop="description">
          <el-input
            v-model="versionForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入版本描述（可选）"
          />
        </el-form-item>

        <el-form-item label="最新版本">
          <el-switch
            v-model="versionForm.is_latest"
            active-text="是"
            inactive-text="否"
          />
          <div class="form-help">设置为最新版本后，其他版本的"最新"标记将被取消</div>
        </el-form-item>

        <el-form-item label="强制更新">
          <el-switch
            v-model="versionForm.is_forced_update"
            active-text="是"
            inactive-text="否"
          />
          <div class="form-help">强制更新将要求用户必须升级到此版本</div>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showVersionDialog = false">取消</el-button>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleSubmitVersion"
          >
            {{ editingVersion ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 确认删除版本对话框 -->
    <ConfirmDialog
      v-model="showDeleteVersionDialog"
      :message="`确定要删除版本「${selectedVersion?.version}」吗？`"
      sub-message="删除后将无法恢复。"
      :loading="deleting"
      @confirm="confirmDeleteVersion"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { formatDate, copyToClipboard } from '@/utils'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Edit,
  Plus,
  Refresh,
  CopyDocument
} from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import type { App, Version, UpdateAppRequest, CreateVersionRequest, UpdateVersionRequest } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

// 路由参数
const akey = computed(() => route.params.akey as string)

// 数据状态
const app = ref<App | null>(null)
const loading = computed(() => appStore.loading)
const versions = computed(() => appStore.versions)
const total = computed(() => appStore.total)
const currentPage = computed(() => appStore.currentPage)
const pageSize = computed(() => appStore.pageSize)

// 对话框状态
const showEditDialog = ref(false)
const showVersionDialog = ref(false)
const showDeleteVersionDialog = ref(false)
const submitting = ref(false)
const deleting = ref(false)

// 编辑状态
const editingVersion = ref<Version | null>(null)
const selectedVersion = ref<Version | null>(null)

// 表单引用和数据
const appFormRef = ref<FormInstance>()
const versionFormRef = ref<FormInstance>()

const appForm = ref<UpdateAppRequest>({
  name: '',
  description: '',
  is_paid: false
})

const versionForm = ref<CreateVersionRequest & UpdateVersionRequest>({
  version: '',
  description: '',
  is_latest: false,
  is_forced_update: false
})

// 表单验证规则
const appRules = {
  name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' },
    { min: 2, max: 50, message: '应用名称长度为2-50个字符', trigger: 'blur' }
  ]
}

const versionRules = {
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' },
    { min: 1, max: 20, message: '版本号长度为1-20个字符', trigger: 'blur' }
  ]
}

// 处理返回操作
const handleGoBack = () => {
  // 返回应用管理页面
  router.push('/apps')
}

// 处理分页变化
const handlePageChange = (page: number, size: number) => {
  appStore.fetchVersions(akey.value, { page, size })
}

// 刷新版本列表
const refreshVersions = () => {
  appStore.fetchVersions(akey.value, { page: currentPage.value, size: pageSize.value })
}

// 编辑应用
const handleEdit = () => {
  if (!app.value) return
  
  appForm.value = {
    name: app.value.name,
    description: app.value.description,
    is_paid: app.value.is_paid
  }
  showEditDialog.value = true
}

// 更新应用
const handleUpdateApp = async () => {
  if (!appFormRef.value || !app.value) return

  try {
    await appFormRef.value.validate()
    submitting.value = true

    await appStore.updateApp(app.value.akey, appForm.value)
    
    // 更新本地应用信息
    app.value = { ...app.value, ...appForm.value }
    showEditDialog.value = false
  } catch (error) {
    console.error('更新应用失败:', error)
  } finally {
    submitting.value = false
  }
}

// 创建版本
const handleCreateVersion = () => {
  editingVersion.value = null
  versionForm.value = {
    version: '',
    description: '',
    is_latest: false,
    is_forced_update: false
  }
  showVersionDialog.value = true
}

// 编辑版本
const handleEditVersion = (version: Version) => {
  editingVersion.value = version
  versionForm.value = {
    version: version.version,
    description: version.description,
    is_latest: version.is_latest,
    is_forced_update: version.is_forced_update
  }
  showVersionDialog.value = true
}

// 提交版本表单
const handleSubmitVersion = async () => {
  if (!versionFormRef.value) return

  try {
    await versionFormRef.value.validate()
    submitting.value = true

    if (editingVersion.value) {
      // 更新版本
      await appStore.updateVersion(editingVersion.value.vkey, versionForm.value, akey.value)
    } else {
      // 创建版本
      await appStore.createVersion(akey.value, versionForm.value)
    }

    showVersionDialog.value = false
  } catch (error) {
    console.error('提交版本失败:', error)
  } finally {
    submitting.value = false
  }
}

// 删除版本
const handleDeleteVersion = (version: Version) => {
  selectedVersion.value = version
  showDeleteVersionDialog.value = true
}

// 确认删除版本
const confirmDeleteVersion = async () => {
  if (!selectedVersion.value) return

  deleting.value = true
  try {
    await appStore.deleteVersion(selectedVersion.value.vkey, akey.value)
    showDeleteVersionDialog.value = false
    selectedVersion.value = null
  } catch (error) {
    console.error('删除版本失败:', error)
  } finally {
    deleting.value = false
  }
}

// 对话框关闭处理
const handleEditDialogClosed = () => {
  appFormRef.value?.resetFields()
}

const handleVersionDialogClosed = () => {
  editingVersion.value = null
  versionFormRef.value?.resetFields()
}

// 复制到剪贴板
const handleCopy = async (text: string) => {
  const success = await copyToClipboard(text)
  if (success) {
    ElMessage.success('已复制到剪贴板')
  } else {
    ElMessage.error('复制失败')
  }
}

// 从应用列表中找到当前应用
const findCurrentApp = () => {
  const currentApp = appStore.apps.find(item => item.akey === akey.value)
  if (currentApp) {
    app.value = currentApp
  }
}

// 监听路由变化
watch(
  () => route.params.akey,
  (newAkey) => {
    if (newAkey) {
      findCurrentApp()
      appStore.fetchVersions(newAkey as string, { page: 1, size: 10 })
    }
  },
  { immediate: true }
)

// 初始化数据
onMounted(() => {
  findCurrentApp()
  
  // 如果没有找到应用信息，可能需要重新获取应用列表
  if (!app.value) {
    appStore.fetchApps({ page: 1, size: 100 }).then(() => {
      findCurrentApp()
    })
  }
  
  // 获取版本列表
  appStore.fetchVersions(akey.value, { page: 1, size: 10 })
})
</script>

<style scoped>
.app-detail {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
}

.header-left {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.app-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.app-name {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.app-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-akey {
  font-family: monospace;
  font-size: 14px;
  background: #f0f2f5;
  padding: 4px 8px;
  border-radius: 4px;
  color: #606266;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.description-card {
  margin-bottom: 24px;
}

.description-content h3 {
  margin: 0 0 12px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.description-content p {
  margin: 0;
  color: #606266;
  line-height: 1.6;
}

.versions-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.version-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.version-text {
  font-weight: 500;
  color: #303133;
}

.vkey-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.vkey-text {
  font-family: monospace;
  font-size: 12px;
  color: #606266;
}

.version-description {
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.form-help {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.dialog-footer {
  text-align: right;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .header-left {
    flex-direction: column;
    gap: 12px;
  }

  .header-actions {
    justify-content: flex-start;
  }

  .app-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>