<script setup lang="ts">
import { inject } from "vue";
import type { Ref } from "vue";
import type { HerbstConfig } from "../types/config";
import ServiceGrid from "../components/ServiceGrid.vue";

const config = inject<Ref<HerbstConfig | null>>("config");
const searchQuery = inject<Ref<string>>("searchQuery");
</script>

<template>
  <ServiceGrid
    v-if="config"
    :services="config.services"
    :search-query="searchQuery || ''"
  />
</template>
