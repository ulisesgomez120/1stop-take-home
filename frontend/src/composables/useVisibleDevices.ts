import { computed, type ComputedRef } from 'vue'
import { useDevicesStore } from '../stores/devices'
import { usePreferencesStore } from '../stores/preferences'
import type { Device } from '../api/devices'

export function useVisibleDevices(): ComputedRef<Device[]> {
  const devicesStore = useDevicesStore()
  const preferencesStore = usePreferencesStore()

  return computed(() => {
    const prefs = preferencesStore.preferences
    const hiddenIds = new Set(prefs?.hidden_device_ids ?? [])
    const sortBy = (prefs?.sort_by ?? 'display_name') as keyof Device
    const dir = prefs?.sort_dir === 'desc' ? -1 : 1

    return devicesStore.devices
      .filter((device) => !hiddenIds.has(device.device_id))
      .sort((a, b) => {
        const av = a[sortBy]
        const bv = b[sortBy]
        if (av === bv) return 0
        return av > bv ? dir : -dir
      })
  })
}
