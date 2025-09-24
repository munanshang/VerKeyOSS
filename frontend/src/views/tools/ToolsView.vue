<template>
  <div class="tools">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1 class="page-title">校验工具</h1>
      <p class="page-subtitle">验证 AKey 和 VKey 的有效性，检查应用更新</p>
    </div>

    <el-row :gutter="20">
      <!-- AKey/VKey 校验工具 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card class="tool-card">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><Lock /></el-icon>
              <span>AKey/VKey 校验</span>
            </div>
          </template>

          <el-form
            ref="validateFormRef"
            :model="validateForm"
            :rules="validateRules"
            label-width="80px"
            @submit.prevent="handleValidate"
          >
            <el-form-item label="AKey" prop="akey">
              <el-input
                v-model="validateForm.akey"
                placeholder="请输入应用标识 AKey"
                clearable
              />
            </el-form-item>

            <el-form-item label="VKey" prop="vkey">
              <el-input
                v-model="validateForm.vkey"
                placeholder="请输入版本标识 VKey"
                clearable
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                :loading="validateLoading"
                @click="handleValidate"
                style="width: 100%"
              >
                <el-icon><Search /></el-icon>
                开始校验
              </el-button>
            </el-form-item>
          </el-form>

          <!-- 校验结果 -->
          <div v-if="validateResult" class="result-container">
            <el-divider content-position="left">校验结果</el-divider>
            <div class="result-content">
              <div class="result-status">
                <el-icon
                  :class="['status-icon', validateResult.valid ? 'success' : 'error']"
                  size="20"
                >
                  <component :is="validateResult.valid ? 'CircleCheck' : 'CircleClose'" />
                </el-icon>
                <span class="status-text">
                  {{ validateResult.valid ? '校验通过' : '校验失败' }}
                </span>
              </div>
              <div class="result-message">
                {{ validateResult.message }}
              </div>
              <!-- 显示应用名和版本号 -->
              <div v-if="validateResult.valid && validateResult.app_name" class="app-info">
                <div class="info-item">
                  <span class="info-label">应用名称：</span>
                  <span class="info-value">{{ validateResult.app_name }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">版本号：</span>
                  <span class="info-value">{{ validateResult.version }}</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 更新检测工具 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card class="tool-card">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><Refresh /></el-icon>
              <span>更新检测</span>
            </div>
          </template>

          <el-form
            ref="updateFormRef"
            :model="updateForm"
            :rules="updateRules"
            label-width="80px"
            @submit.prevent="handleCheckUpdate"
          >
            <el-form-item label="AKey" prop="akey">
              <el-input
                v-model="updateForm.akey"
                placeholder="请输入应用标识 AKey"
                clearable
              />
            </el-form-item>

            <el-form-item label="VKey" prop="vkey">
              <el-input
                v-model="updateForm.vkey"
                placeholder="请输入当前版本标识 VKey"
                clearable
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="success"
                :loading="updateLoading"
                @click="handleCheckUpdate"
                style="width: 100%"
              >
                <el-icon><Upload /></el-icon>
                检查更新
              </el-button>
            </el-form-item>
          </el-form>

          <!-- 更新检测结果 -->
          <div v-if="updateResult" class="result-container">
            <el-divider content-position="left">检测结果</el-divider>
            <div class="result-content">
              <div class="result-status">
                <el-icon
                  :class="['status-icon', updateResult.has_update ? 'warning' : 'success']"
                  size="20"
                >
                  <component :is="updateResult.has_update ? 'Warning' : 'CircleCheck'" />
                </el-icon>
                <span class="status-text">
                  {{ updateResult.has_update ? '有可用更新' : '已是最新版本' }}
                </span>
              </div>
              
              <div v-if="updateResult.has_update" class="update-info">
                <div class="update-item">
                  <span class="update-label">最新版本：</span>
                  <span class="update-value">{{ updateResult.latest_version }}</span>
                </div>
                <div v-if="updateResult.forced_update" class="update-item">
                  <el-tag type="danger" size="small">强制更新</el-tag>
                </div>
              </div>
              
              <div class="result-message">
                {{ updateResult.message }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 使用说明 -->
    <el-card class="help-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><QuestionFilled /></el-icon>
          <span>使用说明</span>
        </div>
      </template>

      <el-collapse v-model="activeHelp">
        <el-collapse-item title="什么是 AKey 和 VKey？" name="1">
          <div class="help-content">
            <p><strong>AKey（应用标识）：</strong>每个应用的唯一标识符，用于区分不同的应用。</p>
            <p><strong>VKey（版本标识）：</strong>每个版本的唯一标识符，用于区分同一应用的不同版本。</p>
          </div>
        </el-collapse-item>

        <el-collapse-item title="如何获取 AKey 和 VKey？" name="2">
          <div class="help-content">
            <p>1. 在<router-link to="/apps">应用管理</router-link>页面可以查看所有应用的 AKey</p>
            <p>2. 点击进入应用详情，可以查看该应用下所有版本的 VKey</p>
            <p>3. 每个标识符都支持一键复制功能</p>
          </div>
        </el-collapse-item>

        <el-collapse-item title="校验工具的作用" name="3">
          <div class="help-content">
            <p><strong>AKey/VKey 校验：</strong>验证给定的 AKey 和 VKey 组合是否有效，确保应用版本的合法性。</p>
            <p><strong>更新检测：</strong>检查指定版本是否有可用的更新，如果有更新会显示最新版本信息和是否强制更新。</p>
          </div>
        </el-collapse-item>

        <el-collapse-item title="API 集成示例" name="4">
          <div class="help-content">
            <p>您可以通过以下 API 端点进行集成：</p>
            <div class="api-example">
              <h4>校验接口</h4>
              <code>POST /api/check/validate</code>
              <pre>{
  "akey": "应用标识",
  "vkey": "版本标识"
}</pre>

              <h4>更新检测接口</h4>
              <code>POST /api/check/update</code>
              <pre>{
  "akey": "应用标识", 
  "vkey": "当前版本标识"
}</pre>
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { checkApi } from '@/api'
import { ElMessage } from 'element-plus'
import {
  Lock,
  Refresh,
  Search,
  Upload,
  QuestionFilled,
  CircleCheck,
  CircleClose,
  Warning
} from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import type { CheckRequest, ValidationResponse, UpdateCheckResponse } from '@/types'

// 表单引用
const validateFormRef = ref<FormInstance>()
const updateFormRef = ref<FormInstance>()

// 校验表单
const validateForm = ref<CheckRequest>({
  akey: '',
  vkey: ''
})

// 更新检测表单
const updateForm = ref<CheckRequest>({
  akey: '',
  vkey: ''
})

// 表单验证规则
const validateRules = {
  akey: [
    { required: true, message: '请输入 AKey', trigger: 'blur' }
  ],
  vkey: [
    { required: true, message: '请输入 VKey', trigger: 'blur' }
  ]
}

const updateRules = {
  akey: [
    { required: true, message: '请输入 AKey', trigger: 'blur' }
  ],
  vkey: [
    { required: true, message: '请输入 VKey', trigger: 'blur' }
  ]
}

// 加载状态
const validateLoading = ref(false)
const updateLoading = ref(false)

// 结果
const validateResult = ref<ValidationResponse | null>(null)
const updateResult = ref<UpdateCheckResponse | null>(null)

// 帮助面板
const activeHelp = ref(['1'])

// 处理校验
const handleValidate = async () => {
  if (!validateFormRef.value) return

  try {
    await validateFormRef.value.validate()
    validateLoading.value = true
    validateResult.value = null

    const response = await checkApi.validate(validateForm.value)
    validateResult.value = response.data.data!
    
    if (validateResult.value.valid) {
      ElMessage.success('校验通过')
    } else {
      ElMessage.warning('校验失败')
    }
  } catch (error) {
    console.error('校验失败:', error)
    validateResult.value = {
      valid: false,
      message: '校验请求失败，请检查网络连接或联系管理员'
    }
  } finally {
    validateLoading.value = false
  }
}

// 处理更新检测
const handleCheckUpdate = async () => {
  if (!updateFormRef.value) return

  try {
    await updateFormRef.value.validate()
    updateLoading.value = true
    updateResult.value = null

    const response = await checkApi.checkUpdate(updateForm.value)
    updateResult.value = response.data.data!
    
    if (updateResult.value.has_update) {
      ElMessage.info('检测到可用更新')
    } else {
      ElMessage.success('已是最新版本')
    }
  } catch (error) {
    console.error('更新检测失败:', error)
    updateResult.value = {
      has_update: false,
      forced_update: false,
      message: '更新检测失败，请检查网络连接或联系管理员'
    }
  } finally {
    updateLoading.value = false
  }
}
</script>

<style scoped>
.tools {
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

.tool-card {
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.header-icon {
  color: #409eff;
}

.result-container {
  margin-top: 20px;
}

.result-content {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 6px;
}

.result-status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.status-icon.success {
  color: #67c23a;
}

.status-icon.error {
  color: #f56c6c;
}

.status-icon.warning {
  color: #e6a23c;
}

.status-text {
  font-weight: 500;
  font-size: 16px;
}

.result-message {
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
}

.update-info {
  margin: 12px 0;
}

.update-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.update-label {
  font-weight: 500;
}

.app-info {
  margin-top: 12px;
  padding: 12px;
  background: #ffffff;
  border-radius: 4px;
  border-left: 4px solid #409eff;
}

.info-item {
  display: flex;
  margin-bottom: 8px;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-label {
  font-weight: 500;
  color: #606266;
  width: 80px;
}

.info-value {
  color: #303133;
  flex: 1;
}

.help-card {
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.help-content {
  line-height: 1.6;
}

.help-content p {
  margin: 8px 0;
  color: #606266;
}

.help-content strong {
  color: #303133;
}

.help-content a {
  color: #409eff;
  text-decoration: none;
}

.help-content a:hover {
  text-decoration: underline;
}

.api-example {
  margin-top: 16px;
}

.api-example h4 {
  margin: 16px 0 8px 0;
  color: #303133;
  font-size: 14px;
}

.api-example code {
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: monospace;
  color: #e6a23c;
  font-size: 13px;
}

.api-example pre {
  background: #2d3748;
  color: #e2e8f0;
  padding: 12px;
  border-radius: 6px;
  font-family: monospace;
  font-size: 13px;
  margin: 8px 0;
  overflow-x: auto;
}

:deep(.el-divider__text) {
  font-weight: 500;
  color: #303133;
}

:deep(.el-collapse-item__header) {
  font-weight: 500;
}

@media (max-width: 768px) {
  .result-status {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .update-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}
</style>