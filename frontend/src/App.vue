<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useDevicesStore } from './stores/devices'
import { usePreferencesStore } from './stores/preferences'
import { useSelectionStore } from './stores/selection'
import { useVisibleDevices } from './composables/useVisibleDevices'
import DeviceMap from './components/DeviceMap.vue'
import DeviceAside from './components/DeviceAside.vue'

const devicesStore = useDevicesStore()
const preferencesStore = usePreferencesStore()
const selection = useSelectionStore()
const visibleDevices = useVisibleDevices()

// Only surface "reconnecting" once an error state has persisted a moment —
// EventSource auto-reconnects, and flashing the badge on every transient
// blip would be noisier than useful.
const showReconnecting = ref(false)
let reconnectingTimer: ReturnType<typeof setTimeout> | undefined

watch(
  () => devicesStore.connectionState,
  (state) => {
    clearTimeout(reconnectingTimer)
    if (state === 'error') {
      reconnectingTimer = setTimeout(() => (showReconnecting.value = true), 2000)
    } else {
      showReconnecting.value = false
    }
  },
)

onMounted(async () => {
  devicesStore.connect()
  await preferencesStore.fetchPreferences()
})
</script>

<template>
  <div class="map-layer">
    <DeviceMap :devices="visibleDevices" />
  </div>

  <DeviceAside />

  <button
    type="button"
    class="aside-toggle"
    :aria-label="selection.asideOpen ? 'Close device panel' : 'Open device panel'"
    @click="selection.toggleAside()"
  >
    {{ selection.asideOpen ? '✕' : '☰' }}
  </button>

  <div
    class="connection-badge"
    :class="`is-${devicesStore.connectionState}`"
    :title="`Connection: ${devicesStore.connectionState}`"
  >
    <span class="connection-dot"></span>
    <span v-if="showReconnecting" class="connection-text">reconnecting…</span>
  </div>
</template>

<style scoped>
.map-layer {
  position: fixed;
  inset: 0;
}

.aside-toggle {
  position: fixed;
  top: 12px;
  left: 12px;
  z-index: 30;
  width: 36px;
  height: 36px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--bg);
  color: var(--text-h);
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
  box-shadow: var(--shadow);
}

.connection-badge {
  position: fixed;
  top: 12px;
  right: 12px;
  z-index: 30;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border: 1px solid var(--border);
  border-radius: 999px;
  background: var(--bg);
  box-shadow: var(--shadow);
  font-size: 12px;
}

.connection-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.connection-badge.is-live .connection-dot {
  background: #22c55e;
}

.connection-badge.is-connecting .connection-dot {
  background: #f59e0b;
}

.connection-badge.is-error .connection-dot {
  background: #ef4444;
}
</style>
