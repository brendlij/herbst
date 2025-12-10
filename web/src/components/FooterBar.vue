<script setup lang="ts">
import { ref, onMounted } from "vue";

const version = ref<string>("...");

onMounted(async () => {
  try {
    const res = await fetch("/api/version");
    if (res.ok) {
      const data = await res.json();
      version.value = data.version || "unknown";
    }
  } catch {
    version.value = "unknown";
  }
});
</script>

<template>
  <footer class="footer-bar">
    <a
      href="https://github.com/brendlij/herbst"
      target="_blank"
      rel="noopener"
      class="footer-link"
    >
      <span>Powered by Herbst 🍂🍁</span>
      <span class="version">{{ version }}</span>
      <i class="mdi mdi-github"></i>
    </a>
  </footer>
</template>

<style scoped>
.footer-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 16px 24px;
}

.footer-link {
  display: flex;
  align-items: center;
  gap: 8px; /* Abstand zwischen Text & Icon */

  color: var(--color-text-muted);
  text-decoration: none;
  font-size: 0.85rem;

  transition: color 0.2s ease;
}

.footer-link i {
  font-size: 18px;
}

.footer-link:hover {
  color: var(--color-accent);
}

.version {
  opacity: 0.7;
  font-size: 0.75rem;
}
</style>
