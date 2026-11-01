<script setup>
import { ref } from 'vue'
import { A, api, json, list, col, sure } from '../lib.js'
import { colors } from '../mount.js'

const page = ref()
const create = ref(null)
const created = ref(null) // URL of the permit that was just created (the API returns it only once)

const columns = [
  col('ID', 'id'),
  col('Comment', 'comment'),
  col('Embed site', 'embed_site'),
  col('Target site', 'target_site'),
  col('Target location', 'target_location'),
  col('Token valid until', 'token_use_not_after', { date: true }),
  col('Session valid until', 'session_token_use_not_after', { date: true }),
  { title: 'Status', key: 'status' },
  { title: '', key: 'actions', align: 'right' },
]
const fields = [
  // all fields are required by the API except target_location
  { key: 'comment', label: 'Comment', required: true },
  { key: 'embed_site', label: 'Embed site', required: true },
  { key: 'target_site', label: 'Target site', required: true },
  { key: 'target_location', label: 'Target location (optional)' },
  { key: 'token_lifetime_minutes', label: 'Token lifetime, minutes', type: 'number', max: 1440, required: true },
  { key: 'session_lifetime_minutes', label: 'Session lifetime, minutes', type: 'number', max: 1440, required: true },
]
const load = () => list(`${A}/permits`, 'permits')

// [label, tag color]: by the session expiry when there is a session, otherwise by the token expiry
// (session_token_use_not_after is null until the token has been used)
const status = (p) => (new Date(p.session_token_use_not_after ?? p.token_use_not_after) > new Date() ? ['valid', colors.ok] : ['expired', colors.warn])

// target_location is optional: leave it out of the request when empty
async function submit({ target_location, ...body }) {
  const r = await json('POST', `${A}/permits`, target_location?.trim() ? { ...body, target_location: target_location.trim() } : body)
  created.value = { title: 'Permit created', url: r.url }
}
</script>

<template>
  <DataPage
      ref="page"
      title="Permits"
      subtitle="Temporary access permits"
      creatable
      parent-key="parent_id"
      :columns="columns"
      :load="load"
      :filter-keys="['id', 'parent_id', 'comment', 'embed_site', 'target_site', 'target_location']"
      filter-placeholder="Filter by ID, comment or site"
      @create="create = { comment: '', embed_site: '', target_site: '', target_location: '', token_lifetime_minutes: 60, session_lifetime_minutes: 5 }"
  >
    <template #cell="{ column, record }">
      <a-tag v-if="column.key === 'status'" :color="status(record)[1]">{{ status(record)[0] }}</a-tag>
      <a v-else-if="column.key === 'actions' && !record.parent_id" @click="sure(`Delete permit ${record.id}?`, () => api('DELETE', `${A}/permits/${record.id}`), page)">
        <DeleteOutlined />
      </a>
    </template>
  </DataPage>
  <FormModal v-if="create" title="Create permit" :fields="fields" :init="create" :submit="submit" @close="create = null; page.refresh()" />
  <LinkModal v-if="created" v-bind="created" @close="created = null" />
</template>
