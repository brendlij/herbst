<script setup lang="ts">
import { inject, computed } from "vue";
import type { Ref } from "vue";
import type { HerbstConfig } from "../types/config";
import ServiceGrid from "../components/ServiceGrid.vue";

const config = inject<Ref<HerbstConfig | null>>("config");
const searchQuery = inject<Ref<string>>("searchQuery");

// Check if we have sections configured
const hasSections = computed(() => {
  return config?.value?.sections && config.value.sections.length > 0;
});
</script>

<template>
  <div v-if="config" class="services-view">
    <!-- Sections mode: grouped services with titles -->
    <template v-if="hasSections">
      <div
        v-for="section in config.sections"
        :key="section.title"
        class="service-section"
      >
        <h2 class="section-title">{{ section.title }}</h2>
        <ServiceGrid
          :services="section.services"
          :search-query="searchQuery || ''"
        />
      </div>
    </template>

    <!-- Legacy mode: flat services list -->
    <ServiceGrid
      v-else
      :services="config.services"
      :search-query="searchQuery || ''"
    />
  </div>
</template>

<style scoped>
.services-view {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.service-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--color-border);
}
</style>
