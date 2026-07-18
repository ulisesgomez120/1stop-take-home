<script setup lang="ts">
import type { Device } from '../api/devices'
import { usePreferencesStore } from '../stores/preferences'

const props = defineProps<{
  devices: Device[]
}>()

const preferencesStore = usePreferencesStore()

const sortableColumns = [
  { key: 'display_name', label: 'Name' },
  { key: 'active_state', label: 'Active State' },
  { key: 'drive_status', label: 'Drive Status' },
] as const

function isHidden(deviceId: string): boolean {
  return preferencesStore.preferences?.hidden_device_ids.includes(deviceId) ?? false
}

function toggleHidden(deviceId: string) {
  const hidden = preferencesStore.preferences?.hidden_device_ids ?? []
  const next = hidden.includes(deviceId)
    ? hidden.filter((id) => id !== deviceId)
    : [...hidden, deviceId]
  preferencesStore.update({ hidden_device_ids: next })
}

function sortBy(key: string) {
  const prefs = preferencesStore.preferences
  if (!prefs) return
  const dir = prefs.sort_by === key && prefs.sort_dir === 'asc' ? 'desc' : 'asc'
  preferencesStore.update({ sort_by: key, sort_dir: dir })
}
</script>

<template>
  <table class="device-list">
    <thead>
      <tr>
        <th
          v-for="col in sortableColumns"
          :key="col.key"
          class="sortable"
          @click="sortBy(col.key)"
        >
          {{ col.label }}
          <span v-if="preferencesStore.preferences?.sort_by === col.key">
            {{ preferencesStore.preferences.sort_dir === 'asc' ? '▲' : '▼' }}
          </span>
        </th>
        <th>Position</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="device in props.devices"
        :key="device.device_id"
        :class="{ 'is-hidden': isHidden(device.device_id) }"
      >
        <td>{{ device.display_name }}</td>
        <td>{{ device.active_state }}</td>
        <td>{{ device.drive_status }}</td>
        <td>{{ device.lat.toFixed(5) }}, {{ device.lng.toFixed(5) }}</td>
        <td>
          <button type="button" @click="toggleHidden(device.device_id)">
            {{ isHidden(device.device_id) ? 'Show' : 'Hide' }}
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>
.device-list {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.device-list th,
.device-list td {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
}

.device-list th {
  color: var(--text-h);
  font-weight: 500;
}

.device-list th.sortable {
  cursor: pointer;
  user-select: none;
}

.device-list tr.is-hidden {
  opacity: 0.5;
}
</style>
