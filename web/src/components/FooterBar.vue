<script setup lang="ts"></script>

<template>
  <footer class="footer-bar">
    <a
      href="https://github.com/brendlij/herbst"
      target="_blank"
      rel="noopener"
      class="footer-link"
    >
      <span>Powered by Herbst 🍂🍁</span>
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
</style>
