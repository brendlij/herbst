<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";

interface DockerContainer {
  id: string;
  name: string;
  image: string;
  state: string;
  status: string;
  created: number;
}

interface DockerAgent {
  name: string;
  token: string;
  connected: boolean;
  lastSeen: string | null;
  containers: DockerContainer[];
}

const agents = ref<DockerAgent[]>([]);
const loading = ref(true);
const serverHost = ref(window.location.host);
const agentProtocol = ref("ws");
const copiedAgent = ref<string | null>(null);
const singleLineMode = ref<Record<string, boolean>>({});
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
    // Use protocol from backend config (docker.agent-protocol) if available
    if (json.agentProtocol) {
      agentProtocol.value = json.agentProtocol;
    }
  } catch (e) {
    console.error("Failed to load docker agents:", e);
  } finally {
    loading.value = false;
  }
}

function getCommand(agent: DockerAgent, singleLine: boolean): string {
  if (singleLine) {
    return `docker run -d --name herbst-docker-agent -v /var/run/docker.sock:/var/run/docker.sock -e HERBST_URL="${agentProtocol.value}://${serverHost.value}/api/agents/ws" -e HERBST_TOKEN="${agent.token}" -e NODE_NAME="${agent.name}" ghcr.io/brendlij/herbst-docker-agent:latest`;
  }
  return `docker run -d \\
  --name herbst-docker-agent \\
  -v /var/run/docker.sock:/var/run/docker.sock \\
  -e HERBST_URL="${agentProtocol.value}://${serverHost.value}/api/agents/ws" \\
  -e HERBST_TOKEN="${agent.token}" \\
  -e NODE_NAME="${agent.name}" \\
  ghcr.io/brendlij/herbst-docker-agent:latest`;
}

function copyCommand(agent: DockerAgent) {
  const cmd = getCommand(agent, singleLineMode.value[agent.name] || false);
  navigator.clipboard.writeText(cmd);
  copiedAgent.value = agent.name;
  setTimeout(() => {
    copiedAgent.value = null;
  }, 2000);
}

function toggleSingleLine(agentName: string) {
  singleLineMode.value[agentName] = !singleLineMode.value[agentName];
}

function getStateClass(state: string): string {
  switch (state.toLowerCase()) {
    case "running":
      return "state-running";
    case "exited":
    case "dead":
      return "state-stopped";
    case "paused":
      return "state-paused";
    case "restarting":
      return "state-restarting";
    default:
      return "state-unknown";
  }
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
          <div class="setup-header">
            <p class="setup-title">Run this command on your Docker host:</p>
            <label class="single-line-toggle">
              <input
                type="checkbox"
                :checked="singleLineMode[agent.name]"
                @change="toggleSingleLine(agent.name)"
              />
              Single line
            </label>
          </div>

          <pre class="command-block">{{
            getCommand(agent, singleLineMode[agent.name] || false)
          }}</pre>

          <button
            class="copy-btn"
            :class="{ copied: copiedAgent === agent.name }"
            @click="copyCommand(agent)"
          >
            <span
              class="mdi"
              :class="
                copiedAgent === agent.name ? 'mdi-check' : 'mdi-content-copy'
              "
            ></span>
            {{ copiedAgent === agent.name ? "Copied!" : "Copy Command" }}
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
          <div class="container-grid">
            <div
              v-for="c in agent.containers"
              :key="c.id || c.name"
              class="container-card"
            >
              <div class="container-content">
                <span class="container-name">{{ c.name }}</span>
                <span class="container-status">{{ c.status }}</span>
              </div>
              <div
                class="status-line"
                :class="getStateClass(c.state)"
                :title="c.status"
              ></div>
            </div>
          </div>
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
  color: var(--color-success);
}

.status.connected .status-dot {
  background: var(--color-success);
}

.status.disconnected {
  color: var(--color-error);
}

.status.disconnected .status-dot {
  background: var(--color-error);
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

.setup-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.setup-title,
.containers-title {
  margin: 0 0 1rem 0;
  font-size: 0.9rem;
  color: var(--color-text);
  opacity: 0.8;
}

.single-line-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8rem;
  color: var(--color-text);
  opacity: 0.7;
  cursor: pointer;
  user-select: none;
}

.single-line-toggle input {
  cursor: pointer;
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
  transition: background 0.2s, opacity 0.2s;
}

.copy-btn:hover {
  opacity: 0.9;
}

.copy-btn.copied {
  background: var(--color-success);
}

.containers-section {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.container-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}

.container-card {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 14px;
  position: relative;

  padding: 16px;
  background: var(--color-bg);
  border-radius: 2em;
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);

  overflow: hidden;
  isolation: isolate;
  -webkit-transform: translateZ(0);
  transform: translateZ(0);
}

.container-card:hover {
  transform: translateY(-2px);
}

.container-content {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
}

.container-name {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-status {
  font-size: 0.75rem;
  color: var(--color-text);
  opacity: 0.6;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Status Line */
.status-line {
  position: absolute;
  bottom: 0;
  left: 10%;
  right: 10%;
  height: 3px;
  border-radius: 3px 3px 0 0;
  transition: background-color 0.3s ease;
}

.status-line.state-running {
  background-color: var(--color-success);
}

.status-line.state-stopped {
  background-color: var(--color-error);
}

.status-line.state-paused {
  background-color: var(--color-warning);
}

.status-line.state-restarting {
  background-color: var(--color-info);
}

.status-line.state-unknown {
  background-color: var(--color-text-muted);
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
