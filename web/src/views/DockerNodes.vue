<script setup lang="ts">
import { ref, onMounted } from "vue";

const nodes = ref<any[]>([]);
const loading = ref(true);
const host = window.location.host;
async function loadNodes() {
  loading.value = true;
  const res = await fetch("/api/docker/nodes");
  const json = await res.json();

  nodes.value = Object.values(json).map((n: any) => ({
    ...n,
    connected:
      n.lastSeen && Date.now() - new Date(n.lastSeen).getTime() < 15000,
  }));

  loading.value = false;
}
function copy(n: any) {
  const cmd = `
docker run -d \\
  --name herbst-docker-agent \\
  -v /var/run/docker.sock:/var/run/docker.sock \\
  -e HERBST_URL="ws://${host}/api/agents/ws" \\
  -e HERBST_TOKEN="${n.token}" \\
  -e NODE_NAME="${n.name}" \\
  ghcr.io/brendlij/herbst-docker-agent:latest
`;

  navigator.clipboard.writeText(cmd);
}

onMounted(() => {
  loadNodes();
  setInterval(loadNodes, 5000);
});
</script>

<template>
  <div class="page">
    <h1>Docker Nodes</h1>

    <div v-if="loading">Loading…</div>

    <div class="nodes">
      <div class="node" v-for="n in nodes" :key="n.name">
        <div class="header">
          <h2>{{ n.name }}</h2>
          <span :class="n.connected ? 'ok' : 'bad'">
            {{ n.connected ? "Connected" : "Not connected" }}
          </span>
        </div>

        <div class="info">
          <p><strong>Kind:</strong> {{ n.kind }}</p>
          <p><strong>Last seen:</strong> {{ n.lastSeen || "—" }}</p>
        </div>

        <div v-if="!n.connected" class="setup">
          <p>Run this on your Docker machine:</p>

          <pre class="cmd">
docker run -d \
  --name herbst-docker-agent \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e HERBST_URL="ws://{{ host }}/api/agents/ws" \
  -e HERBST_TOKEN="{{ n.token }}" \
  -e NODE_NAME="{{ n.name }}" \
  ghcr.io/brendlij/herbst-docker-agent:latest
          </pre>

          <button @click="copy(n)">Copy</button>
        </div>

        <div v-if="n.connected" class="containers">
          <p><strong>Containers:</strong></p>
          <ul>
            <li v-for="c in n.containers">{{ c.name }} — {{ c.image }}</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page {
  padding: 2rem;
}

.nodes {
  display: grid;
  gap: 1rem;
}

.node {
  background: var(--color-surface);
  padding: 1rem;
  border-radius: 0.8rem;
}

.header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.6rem;
}

.ok {
  color: #00c896;
  font-weight: bold;
}

.bad {
  color: #ff6a6a;
  font-weight: bold;
}

.cmd {
  background: var(--color-bg);
  padding: 0.5rem;
  font-size: 0.9rem;
  border-radius: 0.4rem;
  white-space: pre-wrap;
}

button {
  margin-top: 0.5rem;
}
</style>
