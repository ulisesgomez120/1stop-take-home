import type { Device } from '../api/devices'

// 'visible' is a pseudo-field: it lives in preferences (hidden_device_ids),
// not on Device, so it needs the hidden set to sort by. Ties fall back to
// display_name so the order within each visibility group is deterministic.
export function sortDevices(
  devices: Device[],
  sortBy?: string,
  sortDir?: string,
  hiddenIds?: Set<string>,
): Device[] {
  const dir = sortDir === 'desc' ? -1 : 1

  if (sortBy === 'visible') {
    return [...devices].sort((a, b) => {
      const av = hiddenIds?.has(a.device_id) ? 1 : 0
      const bv = hiddenIds?.has(b.device_id) ? 1 : 0
      if (av !== bv) return (av - bv) * dir
      return a.display_name > b.display_name ? 1 : a.display_name < b.display_name ? -1 : 0
    })
  }

  const key = (sortBy ?? 'display_name') as keyof Device
  return [...devices].sort((a, b) => {
    const av = a[key]
    const bv = b[key]
    if (av === bv) return 0
    return av > bv ? dir : -dir
  })
}
