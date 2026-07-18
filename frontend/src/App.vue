<script setup lang="ts">
import { onMounted } from 'vue'
import { useDevicesStore } from './stores/devices'
import { usePreferencesStore } from './stores/preferences'
import { useSortedDevices } from './composables/useSortedDevices'
import { useVisibleDevices } from './composables/useVisibleDevices'
import DeviceList from './components/DeviceList.vue'
import DeviceMap from './components/DeviceMap.vue'

const devicesStore = useDevicesStore()
const preferencesStore = usePreferencesStore()
const sortedDevices = useSortedDevices()
const visibleDevices = useVisibleDevices()

onMounted(async () => {
  devicesStore.connect()
  await preferencesStore.fetchPreferences()
})
</script>

<template>
  <div>
    <h1>OneStepGPS Devices</h1>
    <p>Connection: {{ devicesStore.connectionState }}</p>
    <DeviceMap :devices="visibleDevices" />
    <DeviceList :devices="sortedDevices" />
  </div>
</template>
