<script setup lang="ts">
import { onMounted } from 'vue'
import { useDevicesStore } from './stores/devices'
import { usePreferencesStore } from './stores/preferences'

const devicesStore = useDevicesStore()
const preferencesStore = usePreferencesStore()

onMounted(async () => {
  await devicesStore.fetchDevices()
  console.log('devices', devicesStore.devices)
  await preferencesStore.fetchPreferences()
  console.log('preferences', preferencesStore.preferences)
})
</script>

<template>
  <div>
    <h1>OneStepGPS Scaffold</h1>
    <p v-if="devicesStore.loading">Loading devices...</p>
    <p v-else-if="devicesStore.error">Error: {{ devicesStore.error }}</p>
    <p v-else>Fetched {{ devicesStore.devices.length }} devices (see console).</p>
  </div>
</template>
