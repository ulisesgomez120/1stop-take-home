<script setup lang="ts">
import { onMounted } from 'vue'
import { useDevicesStore } from './stores/devices'
import { usePreferencesStore } from './stores/preferences'
import { useSortedDevices } from './composables/useSortedDevices'
import DeviceList from './components/DeviceList.vue'

const devicesStore = useDevicesStore()
const preferencesStore = usePreferencesStore()
const sortedDevices = useSortedDevices()

onMounted(async () => {
  devicesStore.connect()
  await preferencesStore.fetchPreferences()
})
</script>

<template>
  <div>
    <h1>OneStepGPS Devices</h1>
    <p>Connection: {{ devicesStore.connectionState }}</p>
    <DeviceList :devices="sortedDevices" />
  </div>
</template>
