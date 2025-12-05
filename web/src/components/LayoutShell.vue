<script setup lang="ts">
import FooterBar from "./FooterBar.vue";
import HeaderBar from "./HeaderBar.vue";
import SearchBar from "./SearchBar.vue";
import type { WeatherConfig, DockerConfig } from "../types/config";

defineProps<{
  title: string;
  weather: WeatherConfig;
  docker: DockerConfig;
  activeTab: string;
}>();

const emit = defineEmits<{
  (e: "tabChange", tab: string): void;
  (e: "search", query: string): void;
}>();
</script>

<template>
  <div class="layout-shell">
    <div class="layout-content">
      <!-- United Header with Tabs -->
      <HeaderBar
        :title="title"
        :weather="weather"
        :active-tab="activeTab"
        :docker-enabled="docker.enabled"
        :docker-agents-configured="docker.agentsConfigured"
        @tab-change="emit('tabChange', $event)"
      />

      <!-- Search Bar outside the header -->
      <SearchBar @search="emit('search', $event)" />

      <main class="main">
        <slot />
      </main>
    </div>
    <footer class="footer">
      <FooterBar />
    </footer>
  </div>
</template>

<style scoped>
.layout-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.layout-content {
  flex: 1;
  position: relative;
  z-index: 1;
  padding: var(--spacing, 20px);
}

.footer {
  margin-top: auto;
  padding: var(--spacing, 20px);
  padding-top: 0;
}

.main {
  padding-bottom: 2rem;
}
</style>
