import { computed, type ComputedRef } from 'vue'
import { useDevicesStore } from '../stores/devices'
import { usePreferencesStore } from '../stores/preferences'
import type { Device } from '../api/devices'
import { sortDevices } from './sortDevices'

// Feeds the map (Stage 7), not the list — the list shows all devices via
// useSortedDevices so hidden ones stay visible (dimmed) and un-hideable.
export function useVisibleDevices(): ComputedRef<Device[]> {
  const devicesStore = useDevicesStore()
  const preferencesStore = usePreferencesStore()

  return computed(() => {
    const prefs = preferencesStore.preferences
    const hiddenIds = new Set(prefs?.hidden_device_ids ?? [])
    const visible = devicesStore.devices.filter((device) => !hiddenIds.has(device.device_id))
    return sortDevices(visible, prefs?.sort_by, prefs?.sort_dir)
  })
}
