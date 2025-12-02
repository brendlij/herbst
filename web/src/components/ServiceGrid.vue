<script setup lang="ts">
import type { Service } from "../types/config";
import ServiceCard from "./ServiceCard.vue";

defineProps<{
  services: Service[];
}>();
</script>

<template>
  <div class="service-grid">
    <ServiceCard
      v-for="service in services"
      :key="service.url"
      :service="service"
    />
  </div>
</template>

<style scoped>
.service-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
}
</style>
