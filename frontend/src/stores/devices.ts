import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getDevices, type Device } from '../api/devices'

export const useDevicesStore = defineStore('devices', () => {
  const devices = ref<Device[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchDevices() {
    loading.value = true
    error.value = null
    try {
      devices.value = await getDevices()
    } catch (err) {
      error.value = err instanceof Error ? err.message : String(err)
    } finally {
      loading.value = false
    }
  }

  return { devices, loading, error, fetchDevices }
})
