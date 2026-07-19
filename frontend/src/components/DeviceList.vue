<script setup lang="ts">
import { ref } from 'vue'
import type { Device } from '../api/devices'
import { usePreferencesStore } from '../stores/preferences'
import { useSelectionStore } from '../stores/selection'
import { apiURL } from '../api/client'
import { uploadDeviceIcon } from '../api/icons'

const props = defineProps<{
  devices: Device[]
}>()

const preferencesStore = usePreferencesStore()
const selection = useSelectionStore()
const uploadErrors = ref<Record<string, string>>({})

const sortableColumns = [
  { key: 'display_name', label: 'Name' },
  { key: 'active_state', label: 'Active State' },
  { key: 'drive_status', label: 'Drive Status' },
] as const

function isHidden(deviceId: string): boolean {
  return preferencesStore.preferences?.hidden_device_ids.includes(deviceId) ?? false
}

function iconUrl(deviceId: string): string | null {
  const path = preferencesStore.preferences?.device_icons[deviceId]
  return path ? apiURL(path) : null
}

async function onIconChange(deviceId: string, event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  uploadErrors.value = { ...uploadErrors.value, [deviceId]: '' }
  try {
    const url = await uploadDeviceIcon(deviceId, file)
    const icons = { ...(preferencesStore.preferences?.device_icons ?? {}), [deviceId]: url }
    await preferencesStore.update({ device_icons: icons })
  } catch (err) {
    uploadErrors.value = {
      ...uploadErrors.value,
      [deviceId]: err instanceof Error ? err.message : String(err),
    }
  } finally {
    input.value = ''
  }
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
  <div class="device-list-panel">
    <p v-if="preferencesStore.loading" class="prefs-status">Loading preferences…</p>
    <p v-else-if="preferencesStore.error" class="prefs-status prefs-error">
      Preferences error: {{ preferencesStore.error }}
    </p>

    <p v-if="props.devices.length === 0" class="empty-state">
      Waiting for devices… If this persists, the device feed may be unavailable.
    </p>

    <table v-else class="device-list">
      <thead>
        <tr>
          <th>Icon</th>
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
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="device in props.devices"
          :key="device.device_id"
          :class="{ 'is-hidden': isHidden(device.device_id) }"
          class="device-row"
          @click="selection.openDetail(device.device_id)"
        >
          <td>
            <img
              v-if="iconUrl(device.device_id)"
              :src="iconUrl(device.device_id)!"
              alt=""
              class="device-icon"
            />
            <input
              type="file"
              accept="image/png,image/jpeg,image/webp"
              @click.stop
              @change="onIconChange(device.device_id, $event)"
            />
            <div v-if="uploadErrors[device.device_id]" class="upload-error">
              {{ uploadErrors[device.device_id] }}
            </div>
          </td>
          <td>{{ device.display_name }}</td>
          <td>{{ device.active_state }}</td>
          <td>{{ device.drive_status }}</td>
          <td>
            <button type="button" @click.stop="toggleHidden(device.device_id)">
              {{ isHidden(device.device_id) ? 'Show' : 'Hide' }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.device-list-panel {
  padding: 0 16px 16px;
}

.prefs-status {
  padding: 8px 0;
  font-size: 14px;
}

.prefs-error {
  color: var(--danger, #c0392b);
}

.empty-state {
  padding: 16px 0;
  font-size: 15px;
}

.device-list {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 14px;
}

.device-list th,
.device-list td {
  padding: 8px 10px;
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

.device-row {
  cursor: pointer;
}

.device-row:hover {
  background: var(--accent-bg);
}

.device-list tr.is-hidden {
  opacity: 0.5;
}

.device-icon {
  display: block;
  width: 32px;
  height: 32px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 4px;
}

.device-list input[type='file'] {
  max-width: 130px;
  font-size: 12px;
}

.upload-error {
  color: var(--danger, #c0392b);
  font-size: 0.8em;
}
</style>
