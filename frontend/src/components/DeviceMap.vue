<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { GoogleMap, AdvancedMarker } from 'vue3-google-map'
import type { Device } from '../api/devices'
import { usePreferencesStore } from '../stores/preferences'
import { useSelectionStore } from '../stores/selection'
import { apiURL } from '../api/client'

const props = defineProps<{
  devices: Device[]
}>()

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_API_KEY as string
const preferencesStore = usePreferencesStore()
const selection = useSelectionStore()
const mapRef = ref<InstanceType<typeof GoogleMap>>()

// Gate everything that touches google.maps.* (and the markers themselves)
// until the Maps script has actually loaded — rendering before that point
// throws in the render function and leaves the map blank until an unrelated
// reactive change forces a re-render.
const mapReady = computed(() => mapRef.value?.ready ?? false)

function iconUrl(deviceId: string): string | null {
  const path = preferencesStore.preferences?.device_icons[deviceId]
  return path ? apiURL(path) : null
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

// Fit bounds once the map is ready and whenever the visible device *set*
// changes (hide/show), not on every SSE position update — otherwise the
// map would jump around on every poll tick.
watch(
  [mapReady, () => props.devices.map((d) => d.device_id).sort().join(',')],
  () => {
    if (mapReady.value) fitBounds(props.devices)
  },
)
</script>

<template>
  <GoogleMap
    ref="mapRef"
    :api-key="apiKey"
    map-id="DEMO_MAP_ID"
    class="device-map"
    :map-type-control="false"
    :fullscreen-control="false"
    :center="{ lat: 39.5, lng: -98.35 }"
    :zoom="4"
  >
    <template v-if="mapReady">
      <AdvancedMarker
        v-for="device in props.devices"
        :key="device.device_id"
        :options="{
          position: { lat: device.lat, lng: device.lng },
          title: device.display_name,
          gmpClickable: true,
        }"
        @click="selection.openDetail(device.device_id)"
      >
        <template #content>
          <!-- Outer wrapper is never rotated so the label stays upright;
               only the inner icon rotates to the device's heading. -->
          <div
            class="marker"
            :class="{ 'is-selected': selection.activeDeviceId === device.device_id }"
          >
            <span class="marker-label">{{ device.display_name }}</span>
            <img
              v-if="iconUrl(device.device_id)"
              :src="iconUrl(device.device_id)!"
              alt=""
              class="marker-icon"
              :style="{ transform: `rotate(${device.heading}deg)` }"
            />
            <svg
              v-else
              class="marker-icon marker-arrow"
              :style="{ transform: `rotate(${device.heading}deg)` }"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path d="M12 2 L19 21 L12 17 L5 21 Z" />
            </svg>
          </div>
        </template>
      </AdvancedMarker>
    </template>
  </GoogleMap>
</template>

<style scoped>
.device-map {
  display: block;
  width: 100%;
  height: 100%;
}

.marker {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  cursor: pointer;
}

.marker-label {
  padding: 2px 8px;
  border-radius: 999px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text-h);
  font: 500 12px/1.4 var(--sans);
  letter-spacing: 0;
  white-space: nowrap;
  box-shadow: var(--shadow);
}

.marker-icon {
  width: 32px;
  height: 32px;
  object-fit: cover;
  border-radius: 4px;
}

.marker-arrow {
  fill: var(--accent);
  stroke: var(--bg);
  stroke-width: 1;
}

.marker.is-selected .marker-label {
  background: var(--accent-bg);
  border-color: var(--accent-border);
  color: var(--accent);
}

.marker.is-selected .marker-icon {
  filter: drop-shadow(0 0 4px var(--accent));
}
</style>
