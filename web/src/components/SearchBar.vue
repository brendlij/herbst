<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";

const emit = defineEmits<{
  (e: "search", query: string): void;
  (e: "command", command: string, args: string[]): void;
}>();

// Help modal state
const showHelpModal = ref(false);

// Available commands
const commands: Record<
  string,
  { description: string; icon: string; action: () => void | Promise<void> }
> = {
  reload: {
    description: "Reload configuration",
    icon: "mdi-refresh",
    action: async () => {
      try {
        const res = await fetch("/api/reload", { method: "POST" });
        if (!res.ok) throw new Error("Failed to reload");
        showFeedback("✓ Configuration reloaded");
      } catch (e) {
        showFeedback("✗ Failed to reload config");
      }
    },
  },
  clear: {
    description: "Clear search and filters",
    icon: "mdi-eraser",
    action: () => clearSearch(),
  },
  home: {
    description: "Go to home",
    icon: "mdi-home",
    action: () => {
      window.location.href = "/";
    },
  },
  fullscreen: {
    description: "Toggle fullscreen",
    icon: "mdi-fullscreen",
    action: () => {
      if (document.fullscreenElement) {
        document.exitFullscreen();
      } else {
        document.documentElement.requestFullscreen();
      }
    },
  },
  help: {
    description: "Show available commands",
    icon: "mdi-help-circle",
    action: () => openHelpModal(),
  },
};

const searchQuery = ref("");
const commandFeedback = ref<string | null>(null);
let feedbackTimeout: ReturnType<typeof setTimeout> | null = null;

const placeholder = computed(() => {
  if (searchQuery.value.startsWith(":")) {
    return "Enter command... (try :help)";
  }
  if (searchQuery.value.startsWith("@")) {
    return "Search services…";
  }
  return "Search Google, @service, or :command";
});

const isCommandMode = computed(() => searchQuery.value.startsWith(":"));

// Get matching commands for autocomplete hint
const matchingCommands = computed(() => {
  if (!isCommandMode.value) return [];
  const input = searchQuery.value.slice(1).toLowerCase();
  if (!input) return Object.keys(commands);
  return Object.keys(commands).filter((cmd) => cmd.startsWith(input));
});

// Get unique commands (filter out aliases like reload)
const uniqueCommands = computed(() => {
  const seen = new Set<string>();
  return Object.entries(commands).filter(([_cmd, { description }]) => {
    if (seen.has(description)) return false;
    seen.add(description);
    return true;
  });
});

function showFeedback(message: string, duration = 2000) {
  if (feedbackTimeout) clearTimeout(feedbackTimeout);
  commandFeedback.value = message;
  feedbackTimeout = setTimeout(() => {
    commandFeedback.value = null;
  }, duration);
}

function openHelpModal() {
  showHelpModal.value = true;
}

function closeHelpModal() {
  showHelpModal.value = false;
}

function handleEscapeKey(e: KeyboardEvent) {
  if (e.key === "Escape" && showHelpModal.value) {
    closeHelpModal();
  }
}

function executeCommandFromModal(commandName: string) {
  const command = commands[commandName];
  if (command && commandName !== "help") {
    closeHelpModal();
    command.action();
    emit("command", commandName, []);
    showFeedback(`✓ Executed :${commandName}`);
  }
}

function executeCommand(input: string) {
  const parts = input.slice(1).trim().split(/\s+/);
  const commandName = parts[0]?.toLowerCase();
  const args = parts.slice(1);

  if (!commandName) {
    showFeedback("Enter a command name");
    return;
  }

  const command = commands[commandName];
  if (command) {
    command.action();
    emit("command", commandName, args);
    if (commandName !== "help") {
      showFeedback(`✓ Executed :${commandName}`);
    }
    searchQuery.value = "";
  } else {
    showFeedback(`✗ Unknown command: ${commandName}`);
  }
}

function onSearchInput() {
  const query = searchQuery.value.trim();

  // Commands don't trigger search
  if (query.startsWith(":")) {
    return;
  }

  // Only search services when starting with "@"
  if (query.startsWith("@")) {
    const serviceQuery = query.slice(1); // without "@"
    emit("search", serviceQuery);
  } else {
    // No @ → reset service filter
    emit("search", "");
  }
}

function onSearchKey(e: KeyboardEvent) {
  if (e.key !== "Enter") return;

  const query = searchQuery.value.trim();

  // Enter + : → execute command
  if (query.startsWith(":")) {
    executeCommand(query);
    return;
  }

  // Enter + @ → force service search
  if (query.startsWith("@")) {
    const serviceQuery = query.slice(1);
    emit("search", serviceQuery);
    return;
  }

  // Enter without prefix → Google search
  if (query.length > 0) {
    const encoded = encodeURIComponent(query);
    window.open(`https://www.google.com/search?q=${encoded}`, "_blank");
  }
}

function clearSearch() {
  searchQuery.value = "";
  emit("search", ""); // Reset filter
}

onMounted(() => {
  document.addEventListener("keydown", handleEscapeKey);
});

onUnmounted(() => {
  document.removeEventListener("keydown", handleEscapeKey);
});
</script>

<template>
  <div class="search-wrapper">
    <div class="search-container" :class="{ 'command-mode': isCommandMode }">
      <input
        v-model="searchQuery"
        type="text"
        class="search-input"
        :placeholder="placeholder"
        @input="onSearchInput"
        @keydown="onSearchKey"
      />
      <span
        class="search-icon mdi"
        :class="isCommandMode ? 'mdi-console' : 'mdi-magnify'"
      ></span>

      <!-- Command hint -->
      <div
        v-if="isCommandMode && matchingCommands.length > 0"
        class="command-hint"
      >
        {{ matchingCommands.slice(0, 3).join(", ") }}
        <span v-if="matchingCommands.length > 3">...</span>
      </div>

      <!-- Feedback message -->
      <div v-if="commandFeedback" class="command-feedback">
        {{ commandFeedback }}
      </div>

      <!-- Clear Button -->
      <button
        v-if="searchQuery"
        type="button"
        class="clear-button"
        @click="clearSearch"
      >
        <span class="mdi mdi-close"></span>
      </button>
    </div>

    <!-- Help Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showHelpModal"
          class="modal-overlay"
          @click.self="closeHelpModal"
        >
          <div class="modal-content">
            <div class="modal-header">
              <h2>
                <span class="mdi mdi-console"></span>
                Command Reference
              </h2>
              <button class="modal-close" @click="closeHelpModal">
                <span class="mdi mdi-close"></span>
              </button>
            </div>

            <div class="modal-body">
              <p class="modal-intro">
                Use commands by typing <code>:</code> followed by the command
                name in the search bar.
              </p>

              <div class="command-list">
                <div
                  v-for="[cmd, { description, icon }] in uniqueCommands"
                  :key="cmd"
                  class="command-item"
                  @click="executeCommandFromModal(cmd)"
                >
                  <span class="command-icon mdi" :class="icon"></span>
                  <div class="command-info">
                    <code class="command-name">:{{ cmd }}</code>
                    <span class="command-desc">{{ description }}</span>
                  </div>
                  <span class="command-run mdi mdi-play-circle"></span>
                </div>
              </div>

              <div class="modal-section">
                <h3>Search Prefixes</h3>
                <div class="prefix-list">
                  <div class="prefix-item">
                    <code>@</code>
                    <span>Filter services by name</span>
                  </div>
                  <div class="prefix-item">
                    <code>:</code>
                    <span>Execute a command</span>
                  </div>
                  <div class="prefix-item">
                    <code>text</code>
                    <span>Search Google (press Enter)</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="modal-footer">
              <span class="modal-hint">Press <kbd>Esc</kbd> to close</span>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.search-wrapper {
  display: flex;
  justify-content: center;
  padding: 1.5rem 0;
}

.search-container {
  position: relative;
  width: 100%;
  max-width: 500px;
}

.search-container.command-mode .search-input {
  border-color: var(--color-accent);
  background: rgba(var(--color-accent-rgb, 100, 100, 255), 0.05);
}

.search-input {
  width: 100%;
  padding: 12px 40px;
  padding-left: 44px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 2em;
  color: var(--color-text);
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
}

.search-input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(var(--color-accent-rgb, 100, 100, 255), 0.15);
}

.search-input::placeholder {
  color: var(--color-text);
  opacity: 0.5;
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1.1rem;
  opacity: 0.5;
  pointer-events: none;
  color: var(--color-text);
  transition: color 0.2s;
}

.command-mode .search-icon {
  color: var(--color-accent);
  opacity: 1;
}

.command-hint {
  position: absolute;
  top: 100%;
  left: 16px;
  right: 16px;
  margin-top: 4px;
  padding: 6px 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 0.8rem;
  color: var(--color-text);
  opacity: 0.7;
}

.command-feedback {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-top: 8px;
  padding: 6px 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-accent);
  border-radius: 8px;
  font-size: 0.85rem;
  color: var(--color-text);
  white-space: nowrap;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.clear-button {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  border: none;
  background: transparent;
  padding: 0;
  margin: 0;
  cursor: pointer;
  display: flex;
  opacity: 0.7;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  color: var(--color-text);
}

.clear-button:hover {
  opacity: 0.9;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  width: 100%;
  max-width: 480px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--color-border);
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text);
  display: flex;
  align-items: center;
  gap: 10px;
}

.modal-header h2 .mdi {
  color: var(--color-accent);
}

.modal-close {
  background: transparent;
  border: none;
  color: var(--color-text);
  opacity: 0.6;
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s, background 0.2s;
}

.modal-close:hover {
  opacity: 1;
  background: var(--color-bg);
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
}

.modal-intro {
  margin: 0 0 1.25rem;
  color: var(--color-text);
  opacity: 0.8;
  font-size: 0.9rem;
  line-height: 1.5;
}

.modal-intro code {
  background: var(--color-bg);
  padding: 2px 6px;
  border-radius: 4px;
  color: var(--color-accent);
}

.command-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 1.5rem;
}

.command-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--color-bg);
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s, transform 0.1s;
}

.command-item:hover {
  background: rgba(var(--color-accent-rgb, 100, 100, 255), 0.1);
  transform: translateX(4px);
}

.command-item:active {
  transform: translateX(2px);
}

.command-icon {
  font-size: 1.25rem;
  color: var(--color-accent);
  opacity: 0.8;
}

.command-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.command-name {
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--color-text);
  background: transparent;
  padding: 0;
}

.command-desc {
  font-size: 0.8rem;
  color: var(--color-text);
  opacity: 0.6;
}

.command-run {
  font-size: 1.1rem;
  color: var(--color-text);
  opacity: 0;
  transition: opacity 0.2s;
}

.command-item:hover .command-run {
  opacity: 0.6;
}

.modal-section {
  margin-top: 0.5rem;
}

.modal-section h3 {
  margin: 0 0 0.75rem;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text);
  opacity: 0.7;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.prefix-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.prefix-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 0.9rem;
  color: var(--color-text);
}

.prefix-item code {
  background: var(--color-bg);
  padding: 4px 10px;
  border-radius: 6px;
  color: var(--color-accent);
  min-width: 50px;
  text-align: center;
}

.prefix-item span {
  opacity: 0.7;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: center;
}

.modal-hint {
  font-size: 0.8rem;
  color: var(--color-text);
  opacity: 0.5;
}

.modal-hint kbd {
  background: var(--color-bg);
  padding: 2px 8px;
  border-radius: 4px;
  font-family: inherit;
  font-size: 0.75rem;
}

/* Modal Transition */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .modal-content,
.modal-leave-active .modal-content {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  transform: scale(0.95);
  opacity: 0;
}
</style>
