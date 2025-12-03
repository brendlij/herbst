<script setup lang="ts">
defineProps<{
  activeTab: string;
  dockerEnabled: boolean;
}>();

const emit = defineEmits<{
  (e: "tabChange", tab: string): void;
}>();
</script>

<template>
  <nav class="tabs">
    <div
      class="tab"
      :class="{ active: activeTab === 'services' }"
      @click="emit('tabChange', 'services')"
    >
      Services
    </div>
    <div
      v-if="dockerEnabled"
      class="tab"
      :class="{ active: activeTab === 'docker' }"
      @click="emit('tabChange', 'docker')"
    >
      Docker
    </div>
  </nav>
</template>

<style scoped>
.tabs {
  display: flex;
  justify-content: flex-start;
  gap: 8px;
  background-color: var(--color-surface);
  border-radius: 0 0 2em 2em;

  /* Tabs sollen schweben */
  padding: 4px 0 1em 1em;
}

.tab {
  padding: 10px 18px;
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--color-text-muted);
  border-radius: 10px;
  cursor: pointer;
  position: relative;
  transition: 0.2s ease;
  user-select: none;
}

.tab:hover {
  color: var(--color-text);
  background: var(--color-surface);
}

/* ACTIVE TAB */
.tab.active {
  color: var(--color-text);
}

.tab.active::after {
  content: "";
  position: absolute;
  bottom: -2px;
  left: 20%;
  right: 20%;
  height: 2px;
  background: var(--color-accent);
  border-radius: 4px;
}
</style>
