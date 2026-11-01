<script setup>
import { ref } from 'vue'
import { message, theme } from 'ant-design-vue'

// Dialog with a highlighted link and a copy button, open while mounted (use v-if)
const props = defineProps({ title: String, url: String })
defineEmits(['close'])
const { token } = theme.useToken() // colors come from the theme (src/mount.js)
const copied = ref(false)

async function copy() {
  try {
    await navigator.clipboard.writeText(props.url)
  } catch {
    // clipboard API is unavailable on non-secure origins: fall back to execCommand
    const t = document.createElement('textarea')
    t.value = props.url
    document.body.appendChild(t)
    t.select()
    const ok = document.execCommand('copy')
    t.remove()
    if (!ok) return message.error('Copy failed')
  }
  copied.value = true
  message.success('Link copied')
  setTimeout(() => (copied.value = false), 2000)
}
</script>

<template>
  <a-modal open :title="title" width="640px" :footer="null" @cancel="$emit('close')">
    <a-flex gap="small" align="center">
      <span
        :style="{
          flex: 1,
          minWidth: 0,
          padding: '8px 12px',
          fontFamily: 'monospace',
          userSelect: 'all',
          wordBreak: 'break-all',
          background: token.colorPrimaryBg,
          border: '1px solid ' + token.colorPrimaryBorder,
          borderRadius: token.borderRadius + 'px',
        }"
      >{{ url }}</span>
      <a-button type="primary" @click="copy">
        <template #icon><CheckOutlined v-if="copied" /><CopyOutlined v-else /></template>
        {{ copied ? 'Copied' : 'Copy' }}
      </a-button>
    </a-flex>
  </a-modal>
</template>
