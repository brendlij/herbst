<script setup lang="ts">
import { inject } from "vue";
import type { Ref } from "vue";
import type { HerbstConfig } from "../types/config";
import DockerGrid from "../components/DockerGrid.vue";

const config = inject<Ref<HerbstConfig | null>>("config");
</script>

<template>
  <DockerGrid v-if="config" :docker="config.docker" />
</template>
