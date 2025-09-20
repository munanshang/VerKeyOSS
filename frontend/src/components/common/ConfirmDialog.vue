<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="400px"
    :close-on-click-modal="false"
    @closed="handleClosed"
  >
    <div class="confirm-content">
      <el-icon class="warning-icon" color="#E6A23C" size="24">
        <WarningFilled />
      </el-icon>
      <div class="confirm-text">
        <p class="confirm-message">{{ message }}</p>
        <p v-if="subMessage" class="confirm-sub-message">{{ subMessage }}</p>
      </div>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">{{ cancelText }}</el-button>
        <el-button 
          type="danger" 
          @click="handleConfirm"
          :loading="loading"
        >
          {{ confirmText }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { WarningFilled } from '@element-plus/icons-vue'

interface Props {
  modelValue: boolean
  title?: string
  message: string
  subMessage?: string
  confirmText?: string
  cancelText?: string
  loading?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  title: '确认删除',
  confirmText: '确定',
  cancelText: '取消',
  loading: false
})

const emit = defineEmits<Emits>()

const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  visible.value = false
  emit('cancel')
}

const handleClosed = () => {
  emit('cancel')
}
</script>

<style scoped>
.confirm-content {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 10px 0;
}

.warning-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.confirm-text {
  flex: 1;
}

.confirm-message {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
  line-height: 1.4;
}

.confirm-sub-message {
  margin: 0;
  font-size: 14px;
  color: #909399;
  line-height: 1.4;
}

.dialog-footer {
  text-align: right;
}
</style>