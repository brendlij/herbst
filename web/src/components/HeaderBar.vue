<script setup lang="ts">
import { ref, onMounted } from "vue";
import Logo from "./Logo.vue";

defineProps<{
  title: string;
}>();

const time = ref("");
const date = ref("");
const cpuLoad = ref("0%");
const uptime = ref("0h");
const isOnline = ref(true);
const use24h = ref(true);

function updateClock() {
  const now = new Date();

  time.value = now.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: !use24h.value,
  });

  date.value = now.toLocaleDateString([], {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

function toggleHourFormat() {
  use24h.value = !use24h.value;
  updateClock();
}

// Fake CPU/Uptime (du kannst später echte Endpoints reinballern)
function updateSystemInfo() {
  cpuLoad.value = Math.round(Math.random() * 20) + "%";
  uptime.value = `${Math.round(Math.random() * 48)}h`;
}

onMounted(() => {
  updateClock();
  updateSystemInfo();
  setInterval(updateClock, 1000);
  setInterval(updateSystemInfo, 5000);
});
</script>

<template>
  <header class="header-bar">
    <div class="header-left">
      <div class="logo-icon">
        <Logo />
      </div>
      <span class="title">{{ title }}</span>
    </div>

    <div class="status-bar">
      <!-- Online Badge -->
      <span class="dot" :class="{ online: isOnline }"></span>

      <!-- Time -->
      <span class="time" @click="toggleHourFormat">{{ time }}</span>

      <!-- Date -->
      <span class="date">{{ date }}</span>

      <!-- CPU + Uptime -->
      <span class="meta">CPU {{ cpuLoad }} · UP {{ uptime }}</span>
    </div>
  </header>
</template>

<style scoped>
.header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  padding-left: 12px; /* Less padding on left for logo integration */
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  border-radius: 3em;

  /* Backdrop blur for transparent themes */
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  padding: 8px;
  box-sizing: border-box;
  margin-left: -2px; /* Slight overlap with header edge */
}

.title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--color-text);
}

.status-bar {
  display: flex;
  align-items: center;
  gap: 12px; /* perfekt */
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-text-muted);
  transition: background 0.2s;
}

.dot.online {
  background: #4ade80; /* green-400 */
}

.time {
  font-weight: 600;
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
}
.time:hover {
  color: var(--color-accent);
}

.date {
  color: var(--color-text-muted);
}

.meta {
  color: var(--color-text-muted);
  font-size: 0.78rem;
  opacity: 0.8;
}
</style>
