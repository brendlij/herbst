<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { HerbstConfig } from "./types/config";
import { applyTheme } from "./lib/theme";
import LayoutShell from "./components/LayoutShell.vue";
import ServiceGrid from "./components/ServiceGrid.vue";

const config = ref<HerbstConfig | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

onMounted(async () => {
  try {
    const res = await fetch("/api/config");
    if (!res.ok) throw new Error("Failed to load config");
    const data: HerbstConfig = await res.json();
    config.value = data;
    applyTheme(data.themeVars);
    document.title = data.title || "herbst";
  } catch (e) {
    error.value = (e as Error).message;
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="app">
    <!-- Loading State -->
    <div v-if="loading" class="state-container">
      <div class="loading">
        <span class="loading-dot"></span>
        <span class="loading-dot"></span>
        <span class="loading-dot"></span>
      </div>
      <p>Loading herbst…</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="state-container">
      <div class="error-box">
        <h2>Oops!</h2>
        <p>{{ error }}</p>
      </div>
    </div>

    <!-- Ready State -->
    <LayoutShell
      v-else-if="config"
      :title="config.title"
      :background="config.ui.background"
    >
      <ServiceGrid :services="config.services" />
    </LayoutShell>
  </div>
</template>

<style scoped>
.app {
  min-height: 100vh;
}

.state-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  color: var(--color-text);
}

.loading {
  display: flex;
  gap: 6px;
}

.loading-dot {
  width: 10px;
  height: 10px;
  background: var(--color-accent);
  border-radius: 50%;
  animation: bounce 1.4s infinite ease-in-out both;
}

.loading-dot:nth-child(1) {
  animation-delay: -0.32s;
}

.loading-dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%,
  80%,
  100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

.error-box {
  background: var(--color-surface);
  border: 1px solid rgba(255, 100, 100, 0.3);
  border-radius: var(--radius, 12px);
  padding: 2rem;
  text-align: center;
  max-width: 400px;
}

.error-box h2 {
  margin: 0 0 0.5rem;
  color: #ff6b6b;
}

.error-box p {
  margin: 0;
  opacity: 0.8;
}
</style>
