import { defineStore } from 'pinia'
import { ref } from 'vue'
import { apiURL } from '../api/client'
import type { Device } from '../api/devices'

export type ConnectionState = 'connecting' | 'live' | 'error'

export const useDevicesStore = defineStore('devices', () => {
  const devices = ref<Device[]>([])
  const connectionState = ref<ConnectionState>('connecting')
  let source: EventSource | null = null

  function connect() {
    if (source) return
    connectionState.value = 'connecting'
    source = new EventSource(apiURL('/api/devices/stream'))

    source.onopen = () => {
      connectionState.value = 'live'
    }

    source.onmessage = (event) => {
      devices.value = JSON.parse(event.data) as Device[]
    }

    source.onerror = () => {
      connectionState.value = 'error'
    }
  }

  function disconnect() {
    source?.close()
    source = null
  }

  return { devices, connectionState, connect, disconnect }
})
