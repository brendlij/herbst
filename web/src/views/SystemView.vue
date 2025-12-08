<script setup lang="ts">
import SystemStats from "../components/SystemStats.vue";
</script>

<template>
  <div class="system-view">
    <h1>System</h1>
    <SystemStats />
  </div>
</template>

<style scoped>
.system-view {
  padding: 0 1rem;
}

h1 {
  margin-bottom: 1.5rem;
  color: var(--color-text);
}
</style>
