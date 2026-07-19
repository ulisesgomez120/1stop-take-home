<script setup lang="ts">
import type { Device } from '../api/devices'
import { usePreferencesStore } from '../stores/preferences'
import { useSelectionStore } from '../stores/selection'
import { apiURL } from '../api/client'

const props = defineProps<{
  devices: Device[]
}>()

const preferencesStore = usePreferencesStore()
const selection = useSelectionStore()

const sortableColumns = [
  { key: 'display_name', label: 'Name' },
  { key: 'active_state', label: 'Active State' },
  { key: 'drive_status', label: 'Drive Status' },
  { key: 'visible', label: 'Visible' },
] as const

function isHidden(deviceId: string): boolean {
  return preferencesStore.preferences?.hidden_device_ids.includes(deviceId) ?? false
}

function iconUrl(deviceId: string): string | null {
  const path = preferencesStore.preferences?.device_icons[deviceId]
  return path ? apiURL(path) : null
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
          </td>
          <td>{{ device.display_name }}</td>
          <td>{{ device.active_state }}</td>
          <td>{{ device.drive_status }}</td>
          <td>
            <input
              type="checkbox"
              class="visible-toggle"
              :checked="!isHidden(device.device_id)"
              :aria-label="`Show ${device.display_name} on the map`"
              @click.stop
              @change="toggleHidden(device.device_id)"
            />
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
}

.visible-toggle {
  width: 16px;
  height: 16px;
  accent-color: var(--accent);
  cursor: pointer;
}
</style>
