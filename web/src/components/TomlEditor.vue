<script setup lang="ts">
import { computed, ref } from "vue";
import Prism from "prismjs";

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const highlightRef = ref<HTMLPreElement | null>(null);

const highlighted = computed(() =>
  Prism.highlight(
    props.modelValue || "",
    Prism.languages.toml as Prism.Grammar,
    "toml"
  )
);

function onInput(e: Event) {
  const value = (e.target as HTMLTextAreaElement).value;
  emit("update:modelValue", value);
}

function onScroll(e: Event) {
  const textarea = e.target as HTMLTextAreaElement;
  if (highlightRef.value) {
    highlightRef.value.scrollTop = textarea.scrollTop;
    highlightRef.value.scrollLeft = textarea.scrollLeft;
  }
}
</script>

<template>
  <div class="toml-editor">
    <!-- Highlight-Layer -->
    <pre
      ref="highlightRef"
      class="highlight"
      aria-hidden="true"
    ><code class="language-toml" v-html="highlighted"></code></pre>

    <!-- Eigentliches Eingabefeld -->
    <textarea
      class="input"
      :value="modelValue"
      @input="onInput"
      @scroll="onScroll"
      spellcheck="false"
    />
  </div>
</template>

<style scoped>
.toml-editor {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 250px;
  font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, Monaco,
    Consolas, "Liberation Mono", "Courier New", monospace;
  overflow: hidden;
}

/* Highlight-Layer */
.highlight {
  margin: 0;
  padding: 0.75rem 1rem;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  white-space: pre-wrap;
  word-wrap: break-word;
  border-radius: 0.75rem;
  background: #111827;
  overflow: auto;
}

/* Textarea oben drüber */
.input {
  position: absolute;
  inset: 0;
  padding: 0.75rem 1rem;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  border-radius: 0.75rem;
  border: none;
  resize: none;
  background: transparent;
  color: transparent; /* Text unsichtbar… */
  caret-color: #f9fafb; /* …Cursor sichtbar */
  font: inherit;
  outline: none;
  overflow: auto;
}

/* sorgt dafür, dass Textposition übereinstimmt */
.input::selection {
  background: rgba(96, 165, 250, 0.35);
}
</style>
