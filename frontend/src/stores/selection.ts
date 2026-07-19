import { defineStore } from 'pinia'
import { ref } from 'vue'

// Single source of truth for the aside: which device's detail view is open
// (null = list view) and whether the aside panel itself is visible. Both the
// map (marker click) and the list (row click) call openDetail, so the two
// entry points can never disagree about what's active.
export const useSelectionStore = defineStore('selection', () => {
  const activeDeviceId = ref<string | null>(null)
  const asideOpen = ref(true)

  function openDetail(id: string) {
    activeDeviceId.value = id
    asideOpen.value = true
  }

  function closeDetail() {
    activeDeviceId.value = null
  }

  function toggleAside() {
    asideOpen.value = !asideOpen.value
  }

  return { activeDeviceId, asideOpen, openDetail, closeDetail, toggleAside }
})
