<template>
  <div class="version-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button @click="handleGoBack" type="text">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div class="app-info">
          <h1 class="page-title">{{ app?.name }} - 版本管理</h1>
          <div class="app-meta">
            <span class="app-akey">{{ app?.akey }}</span>
            <el-tag :type="app?.is_paid ? 'warning' : 'success'" size="small">
              {{ app?.is_paid ? '付费应用' : '免费应用' }}
            </el-tag>
          </div>
        </div>
      </div>
      <el-button type="primary" @click="handleCreateVersion">
        <el-icon><Plus /></el-icon>
        创建版本
      </el-button>
    </div>

    <!-- 数据表格 -->
    <DataTable
      :data="versions"
      :loading="loading"
      :total="total"
      :page="currentPage"
      :page-size="pageSize"
      @page-change="handlePageChange"
    >
      <template #toolbar>
        <div class="table-toolbar">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索版本号或VKey"
            clearable
            style="width: 300px"
            @change="handleSearch"
          />
          <div class="toolbar-actions">
            <el-button @click="refreshVersions">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table-column type="selection" width="55" />
      
      <el-table-column prop="version" label="版本号" width="120">
        <template #default="{ row }">
          <div class="version-cell">
            <span class="version-text">{{ row.version }}</span>
            <el-tag v-if="row.is_latest" type="success" size="small">最新</el-tag>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="vkey" label="版本标识" width="220">
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

      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-space>
            <el-button type="text" size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button type="text" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              type="text"
              size="small"
              @click="handleSetLatest(row)"
              :disabled="row.is_latest"
            >
              设为最新
            </el-button>
            <el-button
              type="text"
              size="small"
              @click="handleDelete(row)"
              style="color: #f56c6c"
            >
              删除
            </el-button>
          </el-space>
        </template>
      </el-table-column>
    </DataTable>

    <!-- 创建/编辑版本对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingVersion ? '编辑版本' : '创建版本'"
      width="500px"
      :close-on-click-modal="false"
      @closed="handleDialogClosed"
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
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleSubmit"
          >
            {{ editingVersion ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 版本详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="版本详情"
      width="600px"
    >
      <div v-if="selectedVersion" class="version-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="版本号">
            {{ selectedVersion.version }}
          </el-descriptions-item>
          <el-descriptions-item label="版本标识">
            <div class="version-key">
              {{ selectedVersion.vkey }}
              <el-button
                type="text"
                size="small"
                @click="copyToClipboard(selectedVersion.vkey)"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="是否最新">
            <el-tag :type="selectedVersion.is_latest ? 'success' : 'info'" size="small">
              {{ selectedVersion.is_latest ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="强制更新">
            <el-tag :type="selectedVersion.is_forced_update ? 'danger' : 'info'" size="small">
              {{ selectedVersion.is_forced_update ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">
            {{ formatDate(selectedVersion.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="版本描述" :span="2">
            {{ selectedVersion.description || '暂无描述' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 确认删除对话框 -->
    <ConfirmDialog
      v-model="showDeleteDialog"
      :message="`确定要删除版本「${selectedVersion?.version}」吗？`"
      sub-message="删除后将无法恢复。"
      :loading="deleting"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { formatDate, copyToClipboard } from '@/utils'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  Plus,
  Refresh,
  CopyDocument
} from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import type { App, Version, CreateVersionRequest, UpdateVersionRequest } from '@/types'
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

// 搜索
const searchKeyword = ref('')

// 对话框状态
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const showDeleteDialog = ref(false)
const submitting = ref(false)
const deleting = ref(false)

// 编辑状态
const editingVersion = ref<Version | null>(null)
const selectedVersion = ref<Version | null>(null)

// 表单引用和数据
const versionFormRef = ref<FormInstance>()
const versionForm = ref<CreateVersionRequest & UpdateVersionRequest>({
  version: '',
  description: '',
  is_latest: false,
  is_forced_update: false
})

// 表单验证规则
const versionRules = {
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' },
    { min: 1, max: 20, message: '版本号长度为1-20个字符', trigger: 'blur' }
  ]
}

// 处理返回操作
const handleGoBack = () => {
  // 优先尝试返回应用管理页面
  router.push('/apps')
}

// 处理分页变化
const handlePageChange = (page: number, size: number) => {
  appStore.fetchVersions(akey.value, { page, size })
}

// 处理搜索
const handleSearch = () => {
  // 这里可以添加搜索逻辑
  console.log('搜索:', searchKeyword.value)
}

// 刷新版本列表
const refreshVersions = () => {
  appStore.fetchVersions(akey.value, { page: currentPage.value, size: pageSize.value })
}

// 查看版本详情
const handleView = (version: Version) => {
  selectedVersion.value = version
  showDetailDialog.value = true
}

// 编辑版本
const handleEdit = (version: Version) => {
  editingVersion.value = version
  versionForm.value = {
    version: version.version,
    description: version.description,
    is_latest: version.is_latest,
    is_forced_update: version.is_forced_update
  }
  showCreateDialog.value = true
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
  showCreateDialog.value = true
}

// 设为最新版本
const handleSetLatest = async (version: Version) => {
  try {
    const updateData = {
      ...version,
      is_latest: true
    }
    await appStore.updateVersion(version.vkey, updateData, akey.value)
  } catch (error) {
    console.error('设置最新版本失败:', error)
  }
}

// 删除版本
const handleDelete = (version: Version) => {
  selectedVersion.value = version
  showDeleteDialog.value = true
}

// 确认删除
const confirmDelete = async () => {
  if (!selectedVersion.value) return

  deleting.value = true
  try {
    await appStore.deleteVersion(selectedVersion.value.vkey, akey.value)
    showDeleteDialog.value = false
    selectedVersion.value = null
  } catch (error) {
    console.error('删除版本失败:', error)
  } finally {
    deleting.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
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

    showCreateDialog.value = false
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

// 对话框关闭处理
const handleDialogClosed = () => {
  editingVersion.value = null
  versionForm.value = {
    version: '',
    description: '',
    is_latest: false,
    is_forced_update: false
  }
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
.version-list {
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

.page-title {
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

.table-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.toolbar-actions {
  display: flex;
  gap: 12px;
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
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.version-description {
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.version-detail {
  padding: 10px 0;
}

.version-key {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: monospace;
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

  .table-toolbar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .app-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>