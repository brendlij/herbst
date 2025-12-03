<template>
  <nav class="tabs">
    <div class="tab active">Services</div>
  </nav>
</template>

<style scoped>
.tabs {
  display: flex;
  justify-content: center;
  margin-top: 12px;

  /* Tabs sollen schweben */
  padding: 4px 0;
}

.tab {
  padding: 10px 18px;
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--color-text-muted);
  border-radius: 10px;
  cursor: default;
  position: relative;
  transition: 0.2s ease;
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
