<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
// @ts-ignore - no types available
import CodeEditor from "simple-code-editor";

type ConfigFile = "config" | "themes";

const content = ref("");
const activeFile = ref<ConfigFile>("config");
const saving = ref(false);
const saveStatus = ref<"idle" | "success" | "error">("idle");
const errorMessage = ref("");
const hasUnsavedChanges = ref(false);

const fileOptions: { value: ConfigFile; label: string }[] = [
  { value: "config", label: "config.toml" },
  { value: "themes", label: "themes.toml" },
];

function getApiPath(file: ConfigFile): string {
  return file === "config" ? "/api/config/raw" : "/api/themes/raw";
}

// Handle Ctrl+S / Cmd+S keyboard shortcut
function handleKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === "s") {
    e.preventDefault();
    saveFile();
  }
}

onMounted(async () => {
  await loadFile(activeFile.value);
  window.addEventListener("keydown", handleKeydown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeydown);
});

async function loadFile(file: ConfigFile) {
  const r = await fetch(getApiPath(file));
  content.value = await r.text();
  hasUnsavedChanges.value = false;
}

async function saveFile(): Promise<boolean> {
  saving.value = true;
  saveStatus.value = "idle";
  errorMessage.value = "";

  try {
    const res = await fetch(getApiPath(activeFile.value), {
      method: "PUT",
      headers: { "Content-Type": "text/plain" },
      body: content.value,
    });

    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || "Failed to save");
    }

    hasUnsavedChanges.value = false;
    saveStatus.value = "success";
    setTimeout(() => (saveStatus.value = "idle"), 3000);
    return true;
  } catch (e) {
    saveStatus.value = "error";
    errorMessage.value = (e as Error).message;
    return false;
  } finally {
    saving.value = false;
  }
}

async function switchFile(newFile: ConfigFile) {
  if (newFile === activeFile.value) return;

  // Save current file first if there are unsaved changes
  if (hasUnsavedChanges.value) {
    const saved = await saveFile();
    if (!saved) return; // Don't switch if save failed
  }

  activeFile.value = newFile;
  await loadFile(newFile);
}
</script>

<template>
  <div class="config-page">
    <div class="config-header">
      <div class="header-left">
        <h1>Configuration</h1>
        <select
          class="file-select"
          :value="activeFile"
          @change="
            switchFile(($event.target as HTMLSelectElement).value as ConfigFile)
          "
        >
          <option
            v-for="opt in fileOptions"
            :key="opt.value"
            :value="opt.value"
          >
            {{ opt.label }}
          </option>
        </select>
      </div>
      <div class="actions">
        <span v-if="hasUnsavedChanges" class="status unsaved">● Unsaved</span>
        <span v-if="saveStatus === 'success'" class="status success">
          ✓ Saved & Reloaded
        </span>
        <span v-if="saveStatus === 'error'" class="status error">
          ✗ {{ errorMessage }}
        </span>
        <button class="save-btn" :disabled="saving" @click="saveFile">
          {{ saving ? "Saving..." : "Save & Reload" }}
        </button>
      </div>
    </div>
    <div class="editor-container">
      <CodeEditor
        v-model="content"
        :languages="[['ini', 'TOML']]"
        :line-nums="true"
        :wrap="false"
        :header="false"
        :copy-code="false"
        theme="atom-one-dark"
        width="100%"
        height="100%"
        font-size="14px"
        border-radius="8px"
        @input="hasUnsavedChanges = true"
      />
    </div>
  </div>
</template>

<style scoped>
.config-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 200px);
  min-height: 400px;
}

.config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.config-header h1 {
  margin: 0;
  font-size: 1.5rem;
  color: var(--color-text);
}

.file-select {
  padding: 8px 12px;
  background: var(--color-surface);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 0.9rem;
  cursor: pointer;
  outline: none;
}

.file-select:hover {
  border-color: var(--color-accent);
}

.file-select:focus {
  border-color: var(--color-accent);
}

.actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.status {
  font-size: 0.9rem;
  font-weight: 500;
}

.status.success {
  color: var(--color-success);
}

.status.error {
  color: var(--color-error);
}

.status.unsaved {
  color: var(--color-warning);
}

.save-btn {
  padding: 10px 20px;
  background: var(--color-accent);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s;
}

.save-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.editor-container {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
</style>
