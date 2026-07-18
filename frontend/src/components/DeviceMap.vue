<script setup lang="ts">
import { ref, watch } from 'vue'
import { GoogleMap, Marker, InfoWindow } from 'vue3-google-map'
import type { Device } from '../api/devices'
import { usePreferencesStore } from '../stores/preferences'
import { apiURL } from '../api/client'

const props = defineProps<{
  devices: Device[]
}>()

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY as string
const preferencesStore = usePreferencesStore()
const mapRef = ref<InstanceType<typeof GoogleMap>>()
const activeDeviceId = ref<string | null>(null)

function markerIcon(deviceId: string): google.maps.Icon | undefined {
  const path = preferencesStore.preferences?.device_icons[deviceId]
  if (!path) return undefined
  return {
    url: apiURL(path),
    scaledSize: new google.maps.Size(32, 32),
  }
}

function fitBounds(devices: Device[]) {
  const map = mapRef.value?.map
  if (!map || devices.length === 0) return

  const bounds = new google.maps.LatLngBounds()
  for (const device of devices) {
    bounds.extend({ lat: device.lat, lng: device.lng })
  }
  map.fitBounds(bounds)
}

// Fit bounds once on initial load and whenever the visible device *set*
// changes (hide/show), not on every SSE position update — otherwise the
// map would jump around on every poll tick.
watch(
  () => props.devices.map((d) => d.device_id).sort().join(','),
  () => fitBounds(props.devices),
)
</script>

<template>
  <GoogleMap
    ref="mapRef"
    :api-key="apiKey"
    style="width: 100%; height: 480px"
    :center="{ lat: 39.5, lng: -98.35 }"
    :zoom="4"
  >
    <Marker
      v-for="device in props.devices"
      :key="device.device_id"
      :options="{
        position: { lat: device.lat, lng: device.lng },
        title: device.display_name,
        icon: markerIcon(device.device_id),
      }"
    >
      <InfoWindow
        :model-value="activeDeviceId === device.device_id"
        @update:model-value="(open) => (activeDeviceId = open ? device.device_id : null)"
      >
        <div>
          <strong>{{ device.display_name }}</strong>
          <div>{{ device.drive_status }}</div>
        </div>
      </InfoWindow>
    </Marker>
  </GoogleMap>
</template>
