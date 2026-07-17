<script setup lang="ts">
import { onMounted } from 'vue'
import { useDevicesStore } from './stores/devices'
import { usePreferencesStore } from './stores/preferences'

const devicesStore = useDevicesStore()
const preferencesStore = usePreferencesStore()

onMounted(async () => {
  devicesStore.connect()
  await preferencesStore.fetchPreferences()
  console.log('preferences', preferencesStore.preferences)
})
</script>

<template>
  <div>
    <h1>OneStepGPS Scaffold</h1>
    <p>Connection: {{ devicesStore.connectionState }}</p>
    <p>{{ devicesStore.devices.length }} devices (see console).</p>
  </div>
</template>
