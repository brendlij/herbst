<script setup lang="ts">
import type { Service } from "../types/config";
import { resolveIcon } from "../lib/theme";

const props = defineProps<{
  service: Service;
}>();

const iconUrl = resolveIcon(props.service.icon);
</script>

<template>
  <a :href="service.url" target="_blank" rel="noreferrer" class="service-card">
    <div class="service-icon" v-if="iconUrl">
      <img :src="iconUrl" :alt="service.name" />
    </div>
    <div class="service-icon placeholder" v-else>
      <span>{{ service.name.charAt(0).toUpperCase() }}</span>
    </div>
    <span class="service-name">{{ service.name }}</span>
  </a>
</template>

<style scoped>
.service-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px 16px;
  background: var(--color-surface);
  border-radius: var(--radius, 12px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  text-decoration: none;
  color: var(--color-text);
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.service-card:hover {
  transform: translateY(-2px);
  border-color: var(--color-accent);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.service-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-icon img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.service-icon.placeholder {
  background: var(--color-accent);
  border-radius: 12px;
  font-size: 1.25rem;
  font-weight: 600;
  color: #fff;
}

.service-name {
  font-size: 0.9rem;
  font-weight: 500;
  text-align: center;
  line-height: 1.3;
}
</style>
