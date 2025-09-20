<template>
  <div class="app-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1 class="page-title">应用管理</h1>
        <p class="page-subtitle">管理您的应用和版本信息</p>
      </div>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        创建应用
      </el-button>
    </div>

    <!-- 数据表格 -->
    <DataTable
      :data="apps"
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
            placeholder="搜索应用名称或AKey"
            clearable
            style="width: 300px"
            @change="handleSearch"
          />
        </div>
      </template>

      <el-table-column type="selection" width="55" />
      
      <el-table-column prop="name" label="应用名称" min-width="150">
        <template #default="{ row }">
          <div class="app-name-cell">
            <div class="app-name">{{ row.name }}</div>
            <div class="app-akey">{{ row.akey }}</div>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="description" label="描述" min-width="200">
        <template #default="{ row }">
          <div class="app-description">
            {{ row.description || '暂无描述' }}
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="is_paid" label="类型" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_paid ? 'warning' : 'success'" size="small">
            {{ row.is_paid ? '付费' : '免费' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="version_count" label="版本数" width="80">
        <template #default="{ row }">
          <el-button
            type="text"
            @click="handleViewVersions(row)"
            :disabled="row.version_count === 0"
          >
            {{ row.version_count }}
          </el-button>
        </template>
      </el-table-column>

      <el-table-column prop="created_at" label="创建时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button type="text" size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-dropdown @command="(command: string) => handleAction(command, row)" trigger="click">
              <el-button type="text" size="small">
                更多
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-dropdown-item>
                  <el-dropdown-item command="versions">
                    <el-icon><DocumentCopy /></el-icon>
                    版本管理
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    <el-icon><Delete /></el-icon>
                    <span style="color: #f56c6c">删除</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-table-column>
    </DataTable>

    <!-- 创建/编辑应用对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingApp ? '编辑应用' : '创建应用'"
      width="500px"
      :close-on-click-modal="false"
      @closed="handleDialogClosed"
    >
      <el-form
        ref="appFormRef"
        :model="appForm"
        :rules="appRules"
        label-width="80px"
      >
        <el-form-item label="应用名称" prop="name">
          <el-input
            v-model="appForm.name"
            placeholder="请输入应用名称"
          />
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
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleSubmit"
          >
            {{ editingApp ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 应用详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="应用详情"
      width="600px"
    >
      <div v-if="selectedApp" class="app-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="应用名称">
            {{ selectedApp.name }}
          </el-descriptions-item>
          <el-descriptions-item label="应用标识">
            <div class="app-key">
              {{ selectedApp.akey }}
              <el-button
                type="text"
                size="small"
                @click="copyToClipboard(selectedApp.akey)"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="应用类型">
            <el-tag :type="selectedApp.is_paid ? 'warning' : 'success'">
              {{ selectedApp.is_paid ? '付费应用' : '免费应用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="版本数量">
            {{ selectedApp.version_count }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">
            {{ formatDate(selectedApp.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="应用描述" :span="2">
            {{ selectedApp.description || '暂无描述' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 确认删除对话框 -->
    <ConfirmDialog
      v-model="showDeleteDialog"
      :message="`确定要删除应用「${selectedApp?.name}」吗？`"
      sub-message="删除后将无法恢复，且会同时删除该应用下的所有版本。"
      :loading="deleting"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { formatDate, copyToClipboard } from '@/utils'
import { ElMessage } from 'element-plus'
import { Plus, CopyDocument, ArrowDown, Edit, DocumentCopy, Delete } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import type { App, CreateAppRequest, UpdateAppRequest } from '@/types'
import DataTable from '@/components/common/DataTable.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const router = useRouter()
const appStore = useAppStore()

// 数据状态
const loading = computed(() => appStore.loading)
const apps = computed(() => appStore.apps)
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
const editingApp = ref<App | null>(null)
const selectedApp = ref<App | null>(null)

// 表单引用和数据
const appFormRef = ref<FormInstance>()
const appForm = ref<CreateAppRequest & UpdateAppRequest>({
  name: '',
  description: '',
  is_paid: false
})

// 表单验证规则
const appRules = {
  name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' },
    { min: 2, max: 50, message: '应用名称长度为2-50个字符', trigger: 'blur' }
  ]
}

// 处理分页变化
const handlePageChange = (page: number, size: number) => {
  appStore.fetchApps({ page, size })
}

// 处理搜索
const handleSearch = () => {
  // 这里可以添加搜索逻辑
  console.log('搜索:', searchKeyword.value)
}

// 查看应用详情
const handleView = (app: App) => {
  selectedApp.value = app
  showDetailDialog.value = true
}

// 编辑应用
const handleEdit = (app: App) => {
  editingApp.value = app
  appForm.value = {
    name: app.name,
    description: app.description,
    is_paid: app.is_paid
  }
  showCreateDialog.value = true
}

// 查看版本列表
const handleViewVersions = (app: App) => {
  router.push(`/apps/${app.akey}/versions`)
}

// 删除应用
const handleDelete = (app: App) => {
  selectedApp.value = app
  showDeleteDialog.value = true
}

// 确认删除
const confirmDelete = async () => {
  if (!selectedApp.value) return

  deleting.value = true
  try {
    await appStore.deleteApp(selectedApp.value.akey)
    showDeleteDialog.value = false
    selectedApp.value = null
  } catch (error) {
    console.error('删除应用失败:', error)
  } finally {
    deleting.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!appFormRef.value) return

  try {
    await appFormRef.value.validate()
    submitting.value = true

    if (editingApp.value) {
      // 更新应用
      await appStore.updateApp(editingApp.value.akey, appForm.value)
    } else {
      // 创建应用
      await appStore.createApp(appForm.value)
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
  editingApp.value = null
  appForm.value = {
    name: '',
    description: '',
    is_paid: false
  }
  appFormRef.value?.resetFields()
}

// 处理下拉菜单操作
const handleAction = (command: string, app: App) => {
  switch (command) {
    case 'edit':
      handleEdit(app)
      break
    case 'versions':
      handleViewVersions(app)
      break
    case 'delete':
      handleDelete(app)
      break
  }
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

// 初始化数据
onMounted(() => {
  appStore.fetchApps({ page: 1, size: 10 })
})
</script>

<style scoped>
.app-list {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
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

.table-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.app-name {
  font-weight: 500;
  color: #303133;
}

.app-akey {
  font-family: monospace;
  font-size: 12px;
  color: #909399;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  display: inline-block;
}

.app-description {
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.app-detail {
  padding: 10px 0;
}

.app-key {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: monospace;
}

.dialog-footer {
  text-align: right;
}

/* 操作按钮样式 */
.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: flex-start;
}

.action-buttons .el-button {
  margin: 0 !important;
  padding: 4px 8px;
}

/* 下拉菜单项样式 */
:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
}

:deep(.el-dropdown-menu__item .el-icon) {
  font-size: 14px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .table-toolbar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  /* 移动端操作按钮优化 */
  :deep(.el-table__cell) {
    padding: 8px 4px;
  }
  
  :deep(.el-button--small) {
    padding: 4px 8px;
    font-size: 12px;
  }
  
  /* 操作列在移动端的优化 */
  .action-buttons {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .action-buttons .el-button {
    margin: 0 !important;
    min-width: auto;
  }

  /* 下拉菜单项样式 */
  :deep(.el-dropdown-menu__item) {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}
</style>