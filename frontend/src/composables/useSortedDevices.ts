import { computed, type ComputedRef } from 'vue'
import { useDevicesStore } from '../stores/devices'
import { usePreferencesStore } from '../stores/preferences'
import type { Device } from '../api/devices'
import { sortDevices } from './sortDevices'

export function useSortedDevices(): ComputedRef<Device[]> {
  const devicesStore = useDevicesStore()
  const preferencesStore = usePreferencesStore()

  return computed(() => {
    const prefs = preferencesStore.preferences
    return sortDevices(devicesStore.devices, prefs?.sort_by, prefs?.sort_dir)
  })
}
