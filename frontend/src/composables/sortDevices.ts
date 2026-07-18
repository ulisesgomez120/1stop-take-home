import type { Device } from '../api/devices'

export function sortDevices(devices: Device[], sortBy?: string, sortDir?: string): Device[] {
  const key = (sortBy ?? 'display_name') as keyof Device
  const dir = sortDir === 'desc' ? -1 : 1

  return [...devices].sort((a, b) => {
    const av = a[key]
    const bv = b[key]
    if (av === bv) return 0
    return av > bv ? dir : -dir
  })
}
