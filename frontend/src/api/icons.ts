import { apiURL } from './client'

export async function uploadDeviceIcon(deviceId: string, file: File): Promise<string> {
  const form = new FormData()
  form.append('icon', file)

  const res = await fetch(apiURL(`/api/devices/${encodeURIComponent(deviceId)}/icon`), {
    method: 'POST',
    body: form,
  })
  if (!res.ok) {
    throw new Error(`upload icon failed: ${res.status}`)
  }
  const data = (await res.json()) as { url: string }
  return data.url
}
