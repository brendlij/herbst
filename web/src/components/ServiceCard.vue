<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import type { Service } from "../types/config";
import { resolveIcon } from "../lib/theme";

const props = defineProps<{
  service: Service;
}>();

const iconUrl = resolveIcon(props.service.icon);

// Online status
const isOnline = ref<boolean | null>(null); // null = checking, true = online, false = offline
let checkInterval: number | null = null;

async function checkHealth() {
  if (!props.service.onlineBadge) return;

  try {
    const res = await fetch(
      `/api/health?url=${encodeURIComponent(props.service.url)}`
    );
    const data = await res.json();
    isOnline.value = data.online;
  } catch {
    isOnline.value = false;
  }
}

onMounted(() => {
  if (props.service.onlineBadge) {
    checkHealth();
    checkInterval = window.setInterval(checkHealth, 60000); // Check every minute
  }
});

onUnmounted(() => {
  if (checkInterval) {
    clearInterval(checkInterval);
  }
});
</script>

<template>
  <a
    :href="service.url"
    target="_blank"
    rel="noreferrer"
    class="service-card"
    :class="{ 'has-status': service.onlineBadge }"
  >
    <div class="service-icon">
      <img v-if="iconUrl" :src="iconUrl" :alt="service.name" />
      <span v-else>{{ service.name.charAt(0).toUpperCase() }}</span>
    </div>
    <div class="service-content">
      <span class="service-name">{{ service.name }}</span>
    </div>
    <!-- Online Status Line -->
    <div
      v-if="service.onlineBadge"
      class="status-line"
      :class="{
        online: isOnline === true,
        offline: isOnline === false,
        checking: isOnline === null,
      }"
      :title="
        isOnline === true
          ? 'Online'
          : isOnline === false
          ? 'Offline'
          : 'Checking...'
      "
    ></div>
  </a>
</template>

<style scoped>
.service-card {
  display: flex;
  flex-direction: row; /* icon left, text right */
  align-items: center;
  gap: 16px;
  position: relative;

  padding: 20px;
  background: var(--color-surface);
  border-radius: 3em;
  border: 1px solid rgba(255, 255, 255, 0.06);
  text-decoration: none;
  color: var(--color-text);
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;

  /* Backdrop blur for transparent themes */
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);

  /* Fix visual artifacts with border-radius + backdrop-filter */
  overflow: hidden;
  isolation: isolate;
  -webkit-transform: translateZ(0);
  transform: translateZ(0);
}

.service-content {
  display: flex;
  flex-direction: column;
}

.service-card:hover {
  transform: translateY(-2px);
}

.service-icon {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center; /* icon selbst auch zentrieren */
}

.service-icon img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.service-icon.placeholder {
  background: var(--color-accent);
  border-radius: 12px;
  font-size: 1.25rem;
  font-weight: 600;
  color: #fff;
}

.service-name {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 4px; /* title oben */
  display: block;
}

/* Online Status Line */
.status-line {
  position: absolute;
  bottom: 0;
  left: 20%;
  right: 20%;
  height: 3px;
  border-radius: 3px 3px 0 0;
  transition: background-color 0.3s ease, opacity 0.3s ease;
}

.status-line.online {
  background-color: var(--color-success);
}

.status-line.offline {
  background-color: var(--color-error);
}

.status-line.checking {
  background-color: var(--color-text-muted);
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 0.4;
  }
  50% {
    opacity: 1;
  }
}
</style>
