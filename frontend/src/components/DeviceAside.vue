<script setup lang="ts">
import { useSelectionStore } from '../stores/selection'
import { useSortedDevices } from '../composables/useSortedDevices'
import DeviceList from './DeviceList.vue'
import DeviceDetail from './DeviceDetail.vue'

const selection = useSelectionStore()
const sortedDevices = useSortedDevices()
</script>

<template>
  <aside class="device-aside" :class="{ 'is-open': selection.asideOpen }">
    <DeviceDetail v-if="selection.activeDeviceId" />
    <DeviceList v-else :devices="sortedDevices" />
  </aside>
</template>

<style scoped>
.device-aside {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 20;
  width: min(560px, 100vw);
  box-sizing: border-box;
  /* Leave room for the floating toggle button in the top-left corner. */
  padding-top: 56px;
  overflow-y: auto;
  background: var(--bg);
  border-right: 1px solid var(--border);
  box-shadow: var(--shadow);
  transform: translateX(-100%);
  transition: transform 0.2s ease;
}

.device-aside.is-open {
  transform: translateX(0);
}
</style>
