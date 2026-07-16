import { apiGet, apiPut } from './client'

export interface Preferences {
  sort_by: string
  sort_dir: string
  hidden_device_ids: string[]
  device_icons: Record<string, string>
  updated_at: string
}

export function getPreferences(): Promise<Preferences> {
  return apiGet<Preferences>('/api/preferences')
}

export function putPreferences(prefs: Preferences): Promise<Preferences> {
  return apiPut<Preferences>('/api/preferences', prefs)
}
