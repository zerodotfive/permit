<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { message, theme } from 'ant-design-vue'
import { fmt } from '../lib.js'

// Page with a header, filter bar and an auto-refreshing table.
// Slots: #filters (extra controls), #cell="{ column, record, text }" (custom cells;
// when it renders nothing the default cell is shown). Columns may set `date: true`.
const props = defineProps({
  title: String,
  subtitle: String,
  columns: Array,
  load: Function, // async () => rows
  rowKey: { default: 'id' },
  parentKey: String, // row field with the parent's key: rows are shown nested under their parent
  filterKeys: { default: () => ['id', 'name'] }, // row fields matched by the filter box
  filterPlaceholder: { default: 'Filter by ID or name' },
  creatable: Boolean,
  refreshSec: { default: 10 },
})
defineEmits(['create'])
const { token } = theme.useToken() // colors come from the theme (src/mount.js)

const rows = ref([])
const loading = ref(false)
const q = ref('')
const updated = ref('')
const matched = computed(() => {
  const s = q.value.toLowerCase()
  return rows.value.filter((r) => !s || props.filterKeys.map((k) => r[k]).some((v) => v?.toLowerCase().includes(s)))
})
// Table rows: the matched rows, nested under their parent when the parent is shown too
// (a row whose parent is missing or filtered out stays at the top level)
const shown = computed(() => {
  if (!props.parentKey) return matched.value
  const byKey = new Map(matched.value.map((r) => [r[props.rowKey], { ...r }]))
  const roots = []
  for (const r of byKey.values()) {
    const parent = byKey.get(r[props.parentKey])
    if (parent && parent !== r) (parent.children ??= []).push(r)
    else roots.push(r)
  }
  return roots
})

// Expanded parent rows: all collapsed by default, the user's choice survives auto-refresh
const expanded = ref([])

async function refresh() {
  loading.value = true
  try {
    rows.value = await props.load()
    updated.value = new Date().toLocaleTimeString('en-GB')
  } catch (e) {
    message.error(e.message)
  } finally {
    loading.value = false
  }
}

let timer
onMounted(() => {
  refresh()
  timer = setInterval(refresh, props.refreshSec * 1000)
})
onUnmounted(() => clearInterval(timer))
defineExpose({ refresh })
</script>

<template>
  <div :style="{ padding: '16px 24px', borderBottom: '1px solid ' + token.colorBorderSecondary }">
    <span style="font-size: 18px">{{ title }}</span>
    <a-typography-text type="secondary" style="margin-left: 12px">{{ subtitle }}</a-typography-text>
  </div>
  <div style="padding: 20px 24px">
    <a-flex gap="small" align="center" style="margin-bottom: 16px">
      <a-input v-model:value="q" :placeholder="filterPlaceholder" allow-clear style="width: 280px">
        <template #prefix><SearchOutlined /></template>
      </a-input>
      <slot name="filters" />
      <span style="flex: 1" />
      <a-typography-text type="secondary">auto-refresh every {{ refreshSec }} s · updated {{ updated }}</a-typography-text>
      <span>{{ matched.length }} of {{ rows.length }}</span>
      <a-button v-if="creatable" type="primary" @click="$emit('create')">
        <template #icon><PlusOutlined /></template>Create
      </a-button>
      <a-button :loading="loading" @click="refresh">
        <template #icon><ReloadOutlined /></template>Refresh
      </a-button>
    </a-flex>
    <a-table v-model:expanded-row-keys="expanded" :data-source="shown" :columns="columns" :row-key="rowKey" :pagination="false">
      <template #bodyCell="s">
        <slot name="cell" v-bind="s">
          <a-typography-text v-if="s.column.key === 'id'" code copyable>{{ s.text }}</a-typography-text>
          <template v-else-if="s.column.dataIndex">{{ s.column.date ? fmt(s.text) : s.text }}</template>
        </slot>
      </template>
    </a-table>
  </div>
</template>
