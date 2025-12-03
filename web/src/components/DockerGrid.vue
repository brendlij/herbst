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
        :class="getStateClass(container.state)"
      >
        <div class="container-header">
          <span class="container-state-icon">{{
            getStateIcon(container.state)
          }}</span>
          <span class="container-name">{{ container.name }}</span>
        </div>
        <div class="container-details">
          <span class="container-image">{{ container.image }}</span>
          <span class="container-status">{{ container.status }}</span>
        </div>
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
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 14px 16px;
  transition: all 0.2s ease;
  border-left: 3px solid var(--color-border);
}

.container-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.container-card.state-running {
  border-left-color: #22c55e;
}

.container-card.state-stopped {
  border-left-color: #ef4444;
}

.container-card.state-paused {
  border-left-color: #f59e0b;
}

.container-card.state-restarting {
  border-left-color: #3b82f6;
}

.container-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.container-state-icon {
  font-size: 0.75rem;
  opacity: 0.8;
}

.state-running .container-state-icon {
  color: #22c55e;
}

.state-stopped .container-state-icon {
  color: #ef4444;
}

.state-paused .container-state-icon {
  color: #f59e0b;
}

.state-restarting .container-state-icon {
  color: #3b82f6;
}

.container-name {
  font-weight: 600;
  color: var(--color-text);
  font-size: 0.95rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.container-image {
  font-size: 0.8rem;
  color: var(--color-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-status {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  opacity: 0.8;
}
</style>
