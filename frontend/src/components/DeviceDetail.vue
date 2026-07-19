<script setup lang="ts">
import { computed, ref } from "vue";
import { useDevicesStore } from "../stores/devices";
import { usePreferencesStore } from "../stores/preferences";
import { useSelectionStore } from "../stores/selection";
import { apiURL } from "../api/client";
import { uploadDeviceIcon } from "../api/icons";

const devicesStore = useDevicesStore();
const preferencesStore = usePreferencesStore();
const selection = useSelectionStore();
const uploadError = ref("");

const device = computed(() => devicesStore.devices.find((d) => d.device_id === selection.activeDeviceId) ?? null);

const iconUrl = computed(() => {
  if (!device.value) return null;
  const path = preferencesStore.preferences?.device_icons[device.value.device_id];
  return path ? apiURL(path) : null;
});

function formatLastUpdate(iso: string): string {
  return new Date(iso).toLocaleString(undefined, {
    dateStyle: "medium",
    timeStyle: "medium",
  });
}

async function onIconChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file || !device.value) return;

  const deviceId = device.value.device_id;
  uploadError.value = "";
  try {
    const url = await uploadDeviceIcon(deviceId, file);
    const icons = { ...(preferencesStore.preferences?.device_icons ?? {}), [deviceId]: url };
    await preferencesStore.update({ device_icons: icons });
  } catch (err) {
    uploadError.value = err instanceof Error ? err.message : String(err);
  } finally {
    input.value = "";
  }
}
</script>

<template>
  <div class="device-detail">
    <button type="button" class="back-button" @click="selection.closeDetail()">← Back to list</button>

    <template v-if="device">
      <h2 class="detail-name">{{ device.display_name }}</h2>
      <dl class="detail-fields">
        <dt>Active state</dt>
        <dd>{{ device.active_state }}</dd>
        <dt>Drive status</dt>
        <dd>{{ device.drive_status }}</dd>
        <dt>Heading</dt>
        <dd>{{ Math.round(device.heading) }}°</dd>
        <dt>Position</dt>
        <dd>{{ device.lat.toFixed(5) }}, {{ device.lng.toFixed(5) }}</dd>
        <dt>Last update</dt>
        <dd>{{ formatLastUpdate(device.dt_tracker) }}</dd>
      </dl>

      <div class="icon-section">
        <h3 class="icon-heading">Map icon</h3>
        <div class="icon-row">
          <img v-if="iconUrl" :src="iconUrl" alt="" class="device-icon" />
          <span v-else class="icon-placeholder">Default arrow</span>
          <input type="file" accept="image/png,image/jpeg,image/webp" @change="onIconChange" />
        </div>
        <p v-if="uploadError" class="upload-error">{{ uploadError }}</p>
      </div>

      <h4>
        *notes from me: You could add more device details here. A section for 'Cameras View' could go here if any,
        instead of a camera page -> camera view -> select device -> select camera. Maybe a History section showing the
        last N days as a quick view and leave the detailed history with date range in the history page.
      </h4>
    </template>
    <p v-else class="detail-missing">Device not found in the current feed.</p>
  </div>
</template>

<style scoped>
.device-detail {
  padding: 16px;
}

.back-button {
  border: 1px solid var(--border);
  border-radius: 6px;
  background: transparent;
  color: var(--text);
  font: 14px/1.4 var(--sans);
  padding: 6px 12px;
  cursor: pointer;
}

.back-button:hover {
  color: var(--text-h);
  border-color: var(--accent-border);
}

.detail-name {
  margin: 16px 0 8px;
}

.detail-fields {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 8px 16px;
  margin: 0;
  font-size: 15px;
}

.detail-fields dt {
  color: var(--text-h);
  font-weight: 500;
}

.detail-fields dd {
  margin: 0;
}

.icon-section {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--border);
}

.icon-heading {
  margin: 0 0 8px;
  font-size: 15px;
  font-weight: 500;
  color: var(--text-h);
}

.icon-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.device-icon {
  width: 32px;
  height: 32px;
  object-fit: cover;
  border-radius: 4px;
}

.icon-placeholder {
  font-size: 14px;
}

.upload-error {
  margin: 8px 0 0;
  color: var(--danger, #c0392b);
  font-size: 0.8em;
}

.detail-missing {
  margin-top: 16px;
  font-size: 15px;
}
</style>
