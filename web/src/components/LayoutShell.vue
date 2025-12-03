<script setup lang="ts">
import HeaderTabs from "./HeaderTabs.vue";
import FooterBar from "./FooterBar.vue";
import HeaderBar from "./HeaderBar.vue";
import Logo from "./Logo.vue";

defineProps<{
  title: string;
}>();
</script>

<template>
  <div class="layout-shell">
    <div class="layout-content">
      <!-- <header class="header">
        <Logo size="120" />
        <h1>{{ title }}</h1>
      </header> -->
      <HeaderBar :title="title" />

      <!-- Tabs unter dem Header -->
      <HeaderTabs active="services" />

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
  padding-top: 2rem;
}
</style>
