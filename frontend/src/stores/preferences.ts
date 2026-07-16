import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getPreferences, putPreferences, type Preferences } from '../api/preferences'

export const usePreferencesStore = defineStore('preferences', () => {
  const preferences = ref<Preferences | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchPreferences() {
    loading.value = true
    error.value = null
    try {
      preferences.value = await getPreferences()
    } catch (err) {
      error.value = err instanceof Error ? err.message : String(err)
    } finally {
      loading.value = false
    }
  }

  async function update(partial: Partial<Preferences>) {
    if (!preferences.value) return
    const merged = { ...preferences.value, ...partial }
    error.value = null
    try {
      preferences.value = await putPreferences(merged)
    } catch (err) {
      error.value = err instanceof Error ? err.message : String(err)
    }
  }

  return { preferences, loading, error, fetchPreferences, update }
})
