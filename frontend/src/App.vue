<script setup lang="ts">
import { onMounted } from 'vue'
import { useDevicesStore } from './stores/devices'
import { usePreferencesStore } from './stores/preferences'
import { useVisibleDevices } from './composables/useVisibleDevices'

const devicesStore = useDevicesStore()
const preferencesStore = usePreferencesStore()
const visibleDevices = useVisibleDevices()

onMounted(async () => {
  devicesStore.connect()
  await preferencesStore.fetchPreferences()
})
</script>

<template>
  <div>
    <h1>OneStepGPS Scaffold</h1>
    <p>Connection: {{ devicesStore.connectionState }}</p>
    <p>{{ devicesStore.devices.length }} total, {{ visibleDevices.length }} visible (see console).</p>
  </div>
</template>
