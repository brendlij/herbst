<script setup lang="ts">
import Logo from "./Logo.vue";
defineProps<{
  title: string;
  background?: string;
}>();
</script>

<template>
  <div
    class="layout-shell"
    :style="background ? { backgroundImage: `url(${background})` } : undefined"
  >
    <div class="layout-overlay" v-if="background"></div>
    <div class="layout-content">
      <header class="header">
        <Logo size="120" />
        <h1>{{ title }}</h1>
      </header>
      <main class="main">
        <slot></slot>
      </main>
    </div>
  </div>
</template>

<style scoped>
.layout-shell {
  min-height: 100vh;
  position: relative;
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
}

.layout-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(2px);
  z-index: 0;
}

.layout-content {
  position: relative;
  z-index: 1;
  max-width: 1100px;
  margin: 0 auto;
  padding: var(--spacing, 20px);
}

.header {
  padding: 2rem 0 1.5rem;
  text-align: center;
}

.header h1 {
  margin: 0;
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--color-text);
  letter-spacing: -0.02em;
}

.main {
  padding-bottom: 2rem;
}
</style>
