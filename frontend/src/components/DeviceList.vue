<script setup lang="ts">
import type { Device } from '../api/devices'

defineProps<{
  devices: Device[]
}>()
</script>

<template>
  <table class="device-list">
    <thead>
      <tr>
        <th>Name</th>
        <th>Position</th>
        <th>Active State</th>
        <th>Drive Status</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="device in devices" :key="device.device_id">
        <td>{{ device.display_name }}</td>
        <td>{{ device.lat.toFixed(5) }}, {{ device.lng.toFixed(5) }}</td>
        <td>{{ device.active_state }}</td>
        <td>{{ device.drive_status }}</td>
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
</style>
