<script setup lang="ts">
import { computed } from "vue";
import type { Service } from "../types/config";
import ServiceCard from "./ServiceCard.vue";

const props = defineProps<{
  services: Service[];
  searchQuery?: string;
}>();

const filteredServices = computed(() => {
  if (!props.searchQuery || props.searchQuery.trim() === "") {
    return props.services;
  }
  const query = props.searchQuery.toLowerCase();
  return props.services.filter(
    (service) =>
      service.name.toLowerCase().includes(query) ||
      service.url.toLowerCase().includes(query)
  );
});
</script>

<template>
  <div class="service-grid">
    <ServiceCard
      v-for="service in filteredServices"
      :key="service.url"
      :service="service"
    />
    <div v-if="filteredServices.length === 0 && searchQuery" class="no-results">
      No services found for "{{ searchQuery }}"
    </div>
  </div>
</template>

<style scoped>
.service-grid {
  display: grid;
  gap: 2em;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
}

.no-results {
  grid-column: 1 / -1;
  text-align: center;
  padding: 2rem;
  color: var(--color-text);
  opacity: 0.6;
}
</style>
