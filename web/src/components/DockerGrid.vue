<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import type { DockerContainer, DockerConfig } from "../types/config";

const props = defineProps<{
  docker: DockerConfig;
}>();

const containers = ref<DockerContainer[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
let refreshInterval: ReturnType<typeof setInterval> | null = null;

async function fetchContainers() {
  if (!props.docker.enabled) {
    containers.value = [];
    loading.value = false;
    return;
  }

  try {
    const response = await fetch("/api/docker/containers");
    const data = await response.json();

    if (data.enabled && !data.error) {
      containers.value = data.containers;
      error.value = null;
    } else if (data.error) {
      error.value = data.error;
      containers.value = [];
    }
  } catch (err) {
    error.value = "Failed to fetch containers";
    containers.value = [];
  } finally {
    loading.value = false;
  }
}

function getStateClass(state: string): string {
  switch (state.toLowerCase()) {
    case "running":
      return "state-running";
    case "exited":
    case "dead":
      return "state-stopped";
    case "paused":
      return "state-paused";
    case "restarting":
      return "state-restarting";
    default:
      return "state-unknown";
  }
}

function getStateIcon(state: string): string {
  switch (state.toLowerCase()) {
    case "running":
      return "▶";
    case "exited":
    case "dead":
      return "■";
    case "paused":
      return "⏸";
    case "restarting":
      return "↻";
    default:
      return "?";
  }
}

onMounted(() => {
  fetchContainers();
  // Refresh every 30 seconds
  refreshInterval = setInterval(fetchContainers, 30000);
});

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval);
  }
});

// Refetch when docker config changes
watch(
  () => props.docker.enabled,
  () => {
    fetchContainers();
  }
);
</script>

<template>
  <div class="docker-grid">
    <!-- Loading -->
    <div v-if="loading" class="docker-loading">
      <span>Loading containers...</span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="docker-error">
      <span>{{ error }}</span>
    </div>

    <!-- Empty -->
    <div v-else-if="containers.length === 0" class="docker-empty">
      <span>No containers found</span>
    </div>

    <!-- Container List -->
    <div v-else class="container-list">
      <div
        v-for="container in containers"
        :key="container.id"
        class="container-card"
      >
        <div class="container-icon">
          <span class="container-state-icon">{{
            getStateIcon(container.state)
          }}</span>
        </div>
        <div class="container-content">
          <span class="container-name">{{ container.name }}</span>
          <span class="container-image">{{ container.image }}</span>
        </div>
        <!-- Status Line -->
        <div
          class="status-line"
          :class="getStateClass(container.state)"
          :title="container.status"
        ></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.docker-grid {
  width: 100%;
}

.docker-loading,
.docker-error,
.docker-empty {
  text-align: center;
  padding: 2rem;
  color: var(--color-text-muted);
}

.docker-error {
  color: #ef4444;
}

.container-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}

.container-card {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 16px;
  position: relative;

  padding: 20px;
  background: var(--color-surface);
  border-radius: 3em;
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  /* Backdrop blur for transparent themes */
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);

  /* Fix visual artifacts with border-radius + backdrop-filter */
  overflow: hidden;
  isolation: isolate;
  -webkit-transform: translateZ(0);
  transform: translateZ(0);
}

.container-card:hover {
  transform: translateY(-2px);
}

.container-icon {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  font-size: 1rem;
}

.container-content {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
}

.container-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-image {
  font-size: 0.8rem;
  color: var(--color-text);
  opacity: 0.6;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-state-icon {
  font-size: 1rem;
}

/* Status Line */
.status-line {
  position: absolute;
  bottom: 0;
  left: 20%;
  right: 20%;
  height: 3px;
  border-radius: 3px 3px 0 0;
  transition: background-color 0.3s ease;
}

.status-line.state-running {
  background-color: #86efac; /* pastel green */
}

.status-line.state-stopped {
  background-color: #fca5a5; /* pastel red */
}

.status-line.state-paused {
  background-color: #fcd34d; /* pastel yellow */
}

.status-line.state-restarting {
  background-color: #93c5fd; /* pastel blue */
}

.status-line.state-unknown {
  background-color: #94a3b8; /* gray */
}
</style>
