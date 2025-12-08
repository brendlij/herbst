<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from "vue";

interface SystemStatsData {
  cpu: {
    percent: number;
    model: string;
    cores: number;
    threads: number;
  };
  memory: {
    total: number;
    used: number;
    percent: number;
  };
  disk: {
    total: number;
    used: number;
    percent: number;
  };
  host: {
    hostname: string;
    uptime: number;
    os: string;
    platform: string;
  };
}

const stats = ref<SystemStatsData | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);
let pollInterval: ReturnType<typeof setInterval> | null = null;

async function fetchStats() {
  try {
    const res = await fetch("/api/system/stats");
    if (!res.ok) throw new Error("Failed to fetch stats");
    stats.value = await res.json();
    error.value = null;
  } catch (e) {
    error.value = "Failed to load system stats";
    console.error(e);
  } finally {
    loading.value = false;
  }
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  const parts = [];
  if (days > 0) parts.push(`${days}d`);
  if (hours > 0) parts.push(`${hours}h`);
  if (minutes > 0) parts.push(`${minutes}m`);

  return parts.length > 0 ? parts.join(" ") : "< 1m";
}

const cpuPercent = computed(() => stats.value?.cpu.percent.toFixed(1) ?? "0");
const memPercent = computed(
  () => stats.value?.memory.percent.toFixed(1) ?? "0"
);
const diskPercent = computed(() => stats.value?.disk.percent.toFixed(1) ?? "0");

function getBarColor(percent: number): string {
  if (percent < 60) return "var(--color-success, #4ade80)";
  if (percent < 85) return "var(--color-warning, #fbbf24)";
  return "var(--color-danger, #f87171)";
}

onMounted(() => {
  fetchStats();
  pollInterval = setInterval(fetchStats, 3000);
});

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval);
});
</script>

<template>
  <div class="system-stats">
    <div v-if="loading" class="loading">Loading system stats...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <template v-else-if="stats">
      <!-- Host Info Card -->
      <div class="stats-card host-card">
        <div class="card-header">
          <span class="mdi mdi-server"></span>
          <h3>{{ stats.host.hostname }}</h3>
        </div>
        <div class="host-info">
          <span class="platform">{{ stats.host.platform }}</span>
          <span class="uptime">
            <span class="mdi mdi-clock-outline"></span>
            Uptime: {{ formatUptime(stats.host.uptime) }}
          </span>
        </div>
      </div>

      <!-- Stats Grid -->
      <div class="stats-grid">
        <!-- CPU Card -->
        <div class="stats-card">
          <div class="card-header">
            <span class="mdi mdi-chip"></span>
            <h3>CPU</h3>
            <span class="percent">{{ cpuPercent }}%</span>
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: cpuPercent + '%',
                backgroundColor: getBarColor(stats.cpu.percent),
              }"
            ></div>
          </div>
          <div class="card-details">
            <span class="model">{{ stats.cpu.model.trim() }}</span>
            <span class="cores">{{ stats.cpu.threads }} Threads</span>
          </div>
        </div>

        <!-- Memory Card -->
        <div class="stats-card">
          <div class="card-header">
            <span class="mdi mdi-memory"></span>
            <h3>Memory</h3>
            <span class="percent">{{ memPercent }}%</span>
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: memPercent + '%',
                backgroundColor: getBarColor(stats.memory.percent),
              }"
            ></div>
          </div>
          <div class="card-details">
            <span
              >{{ formatBytes(stats.memory.used) }} /
              {{ formatBytes(stats.memory.total) }}</span
            >
          </div>
        </div>

        <!-- Disk Card -->
        <div class="stats-card">
          <div class="card-header">
            <span class="mdi mdi-harddisk"></span>
            <h3>Disk</h3>
            <span class="percent">{{ diskPercent }}%</span>
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: diskPercent + '%',
                backgroundColor: getBarColor(stats.disk.percent),
              }"
            ></div>
          </div>
          <div class="card-details">
            <span
              >{{ formatBytes(stats.disk.used) }} /
              {{ formatBytes(stats.disk.total) }}</span
            >
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.system-stats {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.loading,
.error {
  text-align: center;
  padding: 2rem;
  color: var(--color-text-muted);
}

.error {
  color: var(--color-danger, #f87171);
}

.stats-card {
  background: var(--color-surface);
  border-radius: 3em;
  padding: 1.5rem;
}

.host-card {
  background: linear-gradient(
    135deg,
    var(--color-surface) 0%,
    var(--color-accent-subtle, var(--color-surface)) 100%
  );
}

.card-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.card-header .mdi {
  font-size: 1.5rem;
  color: var(--color-accent);
}

.card-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-text);
  flex: 1;
}

.card-header .percent {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text);
}

.host-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  color: var(--color-text-muted);
  font-size: 0.9rem;
}

.host-info .platform {
  color: var(--color-text);
  font-weight: 500;
}

.host-info .uptime {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.host-info .uptime .mdi {
  font-size: 1rem;
}

.progress-bar {
  height: 8px;
  background: var(--color-background);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.75rem;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s ease, background-color 0.3s ease;
}

.card-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.card-details .model {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-details .cores {
  color: var(--color-text-muted);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}

@media (max-width: 600px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
