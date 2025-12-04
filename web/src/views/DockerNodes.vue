<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";

interface DockerAgent {
  name: string;
  token: string;
  connected: boolean;
  lastSeen: string | null;
  containers: Array<{ name: string; image: string; state: string }>;
}

const agents = ref<DockerAgent[]>([]);
const loading = ref(true);
const serverHost = ref(window.location.host);
let pollInterval: ReturnType<typeof setInterval> | null = null;

async function loadAgents() {
  try {
    const res = await fetch("/api/docker/agents");
    const json = await res.json();

    agents.value = json.agents || [];
    // Use serverHost from backend config (docker.host) if available
    if (json.serverHost) {
      serverHost.value = json.serverHost;
    }
  } catch (e) {
    console.error("Failed to load docker agents:", e);
  } finally {
    loading.value = false;
  }
}

function copyCommand(agent: DockerAgent) {
  const cmd = `docker run -d \\
  --name herbst-docker-agent \\
  -v /var/run/docker.sock:/var/run/docker.sock \\
  -e HERBST_URL="ws://${serverHost.value}/api/agents/ws" \\
  -e HERBST_TOKEN="${agent.token}" \\
  -e NODE_NAME="${agent.name}" \\
  ghcr.io/brendlij/herbst-docker-agent:latest`;

  navigator.clipboard.writeText(cmd);
}

onMounted(() => {
  loadAgents();
  pollInterval = setInterval(loadAgents, 5000);
});

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval);
  }
});
</script>

<template>
  <div class="docker-nodes-page">
    <h1>Docker Nodes</h1>

    <div v-if="loading" class="loading">Loading…</div>

    <div v-else-if="agents.length === 0" class="empty">
      <p>No docker agents configured.</p>
      <p class="hint">Add agents in your <code>config.toml</code>:</p>
      <pre class="config-example">
[[docker.agents]]
name  = "my-docker-host"
token = "your-secret-token"</pre
      >
    </div>

    <div v-else class="nodes-grid">
      <div class="node-card" v-for="agent in agents" :key="agent.name">
        <div class="node-header">
          <h2>{{ agent.name }}</h2>
          <span
            class="status"
            :class="agent.connected ? 'connected' : 'disconnected'"
          >
            <span class="status-dot"></span>
            {{ agent.connected ? "Connected" : "Not connected" }}
          </span>
        </div>

        <div class="node-info">
          <p v-if="agent.lastSeen">
            <span class="label">Last seen:</span>
            {{ new Date(agent.lastSeen).toLocaleString() }}
          </p>
        </div>

        <!-- Setup instructions for disconnected agents -->
        <div v-if="!agent.connected" class="setup-section">
          <p class="setup-title">Run this command on your Docker host:</p>

          <pre class="command-block">
docker run -d \
  --name herbst-docker-agent \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e HERBST_URL="ws://{{ serverHost }}/api/agents/ws" \
  -e HERBST_TOKEN="{{ agent.token }}" \
  -e NODE_NAME="{{ agent.name }}" \
  ghcr.io/brendlij/herbst-docker-agent:latest</pre
          >

          <button class="copy-btn" @click="copyCommand(agent)">
            <span class="mdi mdi-content-copy"></span>
            Copy Command
          </button>
        </div>

        <!-- Container list for connected agents -->
        <div
          v-if="agent.connected && agent.containers?.length"
          class="containers-section"
        >
          <p class="containers-title">
            Containers ({{ agent.containers.length }})
          </p>
          <ul class="container-list">
            <li
              v-for="c in agent.containers"
              :key="c.name"
              class="container-item"
            >
              <span class="container-name">{{ c.name }}</span>
              <span class="container-image">{{ c.image }}</span>
              <span class="container-state" :class="c.state">{{
                c.state
              }}</span>
            </li>
          </ul>
        </div>

        <div
          v-if="agent.connected && !agent.containers?.length"
          class="no-containers"
        >
          <p>No containers running</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.docker-nodes-page {
  padding: 1rem 0;
}

.docker-nodes-page h1 {
  margin: 0 0 1.5rem;
  font-size: 1.5rem;
  color: var(--color-text);
}

.loading,
.empty {
  text-align: center;
  padding: 2rem;
  color: var(--color-text);
  opacity: 0.7;
}

.empty .hint {
  margin-top: 1rem;
  font-size: 0.9rem;
}

.empty code {
  background: var(--color-surface);
  padding: 2px 6px;
  border-radius: 4px;
}

.config-example {
  background: var(--color-surface);
  padding: 1rem;
  border-radius: 8px;
  font-size: 0.85rem;
  text-align: left;
  display: inline-block;
  margin-top: 0.5rem;
}

.nodes-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
}

.node-card {
  background: var(--color-surface);
  padding: 1.25rem;
  border-radius: 12px;
  border: 1px solid var(--color-border);
}

.node-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.node-header h2 {
  margin: 0;
  font-size: 1.1rem;
  color: var(--color-text);
}

.status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 500;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status.connected {
  color: #00c896;
}

.status.connected .status-dot {
  background: #00c896;
}

.status.disconnected {
  color: #ff6a6a;
}

.status.disconnected .status-dot {
  background: #ff6a6a;
}

.node-info {
  margin-bottom: 1rem;
}

.node-info p {
  margin: 0;
  font-size: 0.9rem;
  color: var(--color-text);
  opacity: 0.8;
}

.label {
  opacity: 0.6;
}

.setup-section {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.setup-title,
.containers-title {
  margin: 0 0 0.75rem;
  font-size: 0.9rem;
  color: var(--color-text);
  opacity: 0.8;
}

.command-block {
  background: var(--color-bg);
  padding: 1rem;
  border-radius: 8px;
  font-size: 0.8rem;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--color-text);
  margin-bottom: 0.75rem;
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: var(--color-accent);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: opacity 0.2s;
}

.copy-btn:hover {
  opacity: 0.9;
}

.containers-section {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.container-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.container-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: var(--color-bg);
  border-radius: 8px;
  font-size: 0.85rem;
}

.container-name {
  font-weight: 500;
  color: var(--color-text);
}

.container-image {
  flex: 1;
  color: var(--color-text);
  opacity: 0.6;
  font-size: 0.8rem;
}

.container-state {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
}

.container-state.running {
  background: rgba(0, 200, 150, 0.2);
  color: #00c896;
}

.container-state.exited {
  background: rgba(255, 106, 106, 0.2);
  color: #ff6a6a;
}

.no-containers {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
  color: var(--color-text);
  opacity: 0.6;
  font-size: 0.9rem;
}
</style>
