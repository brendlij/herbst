<script setup lang="ts">
import { computed, ref, onUnmounted } from "vue";
import Prism from "prismjs";

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const textareaRef = ref<HTMLTextAreaElement | null>(null);
const highlightRef = ref<HTMLPreElement | null>(null);
const lineNumbersRef = ref<HTMLDivElement | null>(null);

const highlighted = computed(() =>
  Prism.highlight(
    props.modelValue || "",
    Prism.languages.toml as Prism.Grammar,
    "toml"
  )
);

const lineCount = computed(() => {
  const value = props.modelValue || "";
  if (!value) return 1;

  const lines = value.split("\n");

  // Wenn die letzte Zeile leer ist (Text endet mit "\n") -> nicht mitzählen
  if (lines[lines.length - 1] === "") {
    return Math.max(lines.length - 1, 1);
  }

  return lines.length;
});

function onInput(e: Event) {
  const value = (e.target as HTMLTextAreaElement).value;
  emit("update:modelValue", value);
}

function syncScroll() {
  if (!textareaRef.value) return;
  const { scrollTop, scrollLeft } = textareaRef.value;

  if (highlightRef.value) {
    highlightRef.value.scrollTop = scrollTop;
    highlightRef.value.scrollLeft = scrollLeft;
  }
  if (lineNumbersRef.value) {
    lineNumbersRef.value.scrollTop = scrollTop;
  }
}

// Use requestAnimationFrame for smoother sync
let rafId: number | null = null;
function onScroll() {
  if (rafId) cancelAnimationFrame(rafId);
  rafId = requestAnimationFrame(syncScroll);
}

onUnmounted(() => {
  if (rafId) cancelAnimationFrame(rafId);
});
</script>

<template>
  <div class="toml-editor">
    <!-- Line Numbers -->
    <div ref="lineNumbersRef" class="line-numbers" aria-hidden="true">
      <span v-for="n in lineCount" :key="n" class="line-number">{{ n }}</span>
    </div>

    <div class="editor-content">
      <!-- Highlight-Layer -->
      <pre
        ref="highlightRef"
        class="highlight"
        aria-hidden="true"
      ><code class="language-toml" v-html="highlighted"></code></pre>

      <!-- Eigentliches Eingabefeld -->
      <textarea
        ref="textareaRef"
        class="input"
        wrap="off"
        :value="modelValue"
        @input="onInput"
        @scroll="onScroll"
        spellcheck="false"
      />
    </div>
  </div>
</template>

<style scoped>
.toml-editor {
  display: flex;
  width: 100%;
  height: 100%;
  min-height: 250px;
  font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, Monaco,
    Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.875rem;
  line-height: 1.5rem;
  background: #111827;
  border-radius: 0.75rem;
  overflow: hidden;
}

/* Line Numbers */
.line-numbers {
  display: flex;
  flex-direction: column;
  padding: 0.75rem 0;
  min-width: 3rem;
  background: #0d1117;
  border-right: 1px solid #1f2937;
  overflow: hidden;
  user-select: none;
}

.line-number {
  padding: 0 0.75rem;
  text-align: right;
  color: #4b5563;
  font-size: 0.875rem;
  line-height: 1.5rem;
  height: 1.5rem;
}

/* Editor Content Area */
.editor-content {
  position: relative;
  flex: 1;
  overflow: hidden;
}

/* Shared text styles for perfect alignment */
.highlight,
.input {
  margin: 0;
  padding: 0.75rem 1rem;
  border: none;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  white-space: pre;
  word-wrap: normal;
  overflow-wrap: normal;
  overflow: auto;
  tab-size: 2;
}

/* Highlight-Layer */
.highlight {
  position: absolute;
  inset: 0;
  background: transparent;
  pointer-events: none;
  border-radius: 0;
}

.highlight code {
  display: block;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  background: transparent;
}

/* Textarea */
.input {
  position: relative;
  resize: none;
  background: transparent;
  color: transparent;
  caret-color: #f9fafb;
  outline: none;
  border-radius: 0;
}

.input::selection {
  background: rgba(96, 165, 250, 0.35);
}

/* Hide scrollbar on highlight layer */
.highlight::-webkit-scrollbar {
  display: none;
}
.highlight {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

/* Mehr Platz unten in Text und Highlight */
.highlight,
.input {
  padding-bottom: 2.5rem; /* oder 1.5rem, einfach etwas mehr als eine Zeile */
}

/* Gleiches für die Zeilennummern */
.line-numbers {
  padding-bottom: 2.5rem;
  box-sizing: border-box;
}
</style>
