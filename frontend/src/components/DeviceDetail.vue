<script setup lang="ts">
import { computed } from 'vue'
import { useDevicesStore } from '../stores/devices'
import { useSelectionStore } from '../stores/selection'

const devicesStore = useDevicesStore()
const selection = useSelectionStore()

const device = computed(
  () => devicesStore.devices.find((d) => d.device_id === selection.activeDeviceId) ?? null,
)

function formatLastUpdate(iso: string): string {
  return new Date(iso).toLocaleString(undefined, {
    dateStyle: 'medium',
    timeStyle: 'medium',
  })
}
</script>

<template>
  <div class="device-detail">
    <button type="button" class="back-button" @click="selection.closeDetail()">
      ← Back to list
    </button>

    <template v-if="device">
      <h2 class="detail-name">{{ device.display_name }}</h2>
      <dl class="detail-fields">
        <dt>Active state</dt>
        <dd>{{ device.active_state }}</dd>
        <dt>Drive status</dt>
        <dd>{{ device.drive_status }}</dd>
        <dt>Heading</dt>
        <dd>{{ Math.round(device.heading) }}°</dd>
        <dt>Position</dt>
        <dd>{{ device.lat.toFixed(5) }}, {{ device.lng.toFixed(5) }}</dd>
        <dt>Last update</dt>
        <dd>{{ formatLastUpdate(device.dt_tracker) }}</dd>
      </dl>
    </template>
    <p v-else class="detail-missing">Device not found in the current feed.</p>
  </div>
</template>

<style scoped>
.device-detail {
  padding: 16px;
}

.back-button {
  border: 1px solid var(--border);
  border-radius: 6px;
  background: transparent;
  color: var(--text);
  font: 14px/1.4 var(--sans);
  padding: 6px 12px;
  cursor: pointer;
}

.back-button:hover {
  color: var(--text-h);
  border-color: var(--accent-border);
}

.detail-name {
  margin: 16px 0 8px;
}

.detail-fields {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 8px 16px;
  margin: 0;
  font-size: 15px;
}

.detail-fields dt {
  color: var(--text-h);
  font-weight: 500;
}

.detail-fields dd {
  margin: 0;
}

.detail-missing {
  margin-top: 16px;
  font-size: 15px;
}
</style>
