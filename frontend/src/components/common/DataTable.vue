<template>
  <div class="data-table">
    <!-- 表格工具栏 -->
    <div v-if="$slots.toolbar" class="table-toolbar">
      <slot name="toolbar" />
    </div>

    <!-- 表格 -->
    <el-table
      :data="data"
      :loading="loading"
      v-bind="$attrs"
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <slot />
    </el-table>

    <!-- 分页 -->
    <div v-if="showPagination" class="table-pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="currentPageSize"
        :page-sizes="pageSizes"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  data: any[]
  loading?: boolean
  total?: number
  page?: number
  pageSize?: number
  pageSizes?: number[]
  showPagination?: boolean
}

interface Emits {
  (e: 'page-change', page: number, pageSize: number): void
  (e: 'selection-change', selection: any[]): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  total: 0,
  page: 1,
  pageSize: 10,
  pageSizes: () => [10, 20, 50, 100],
  showPagination: true
})

const emit = defineEmits<Emits>()

const currentPage = ref(props.page)
const currentPageSize = ref(props.pageSize)

// 监听外部传入的页码变化
watch(() => props.page, (newVal) => {
  currentPage.value = newVal
})

watch(() => props.pageSize, (newVal) => {
  currentPageSize.value = newVal
})

// 处理页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  emit('page-change', page, currentPageSize.value)
}

// 处理每页条数变化
const handleSizeChange = (size: number) => {
  currentPageSize.value = size
  currentPage.value = 1 // 重置到第一页
  emit('page-change', 1, size)
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  emit('selection-change', selection)
}
</script>

<style scoped>
.data-table {
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.table-toolbar {
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
}

.table-pagination {
  padding: 16px 20px;
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #e4e7ed;
}
</style>