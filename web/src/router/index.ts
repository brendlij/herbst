import { createRouter, createWebHistory } from "vue-router";
import ServicesView from "../views/ServicesView.vue";
import DockerLocalView from "../views/DockerLocalView.vue";
import DockerNodes from "../views/DockerNodes.vue";
import SystemView from "../views/SystemView.vue";
import ConfigView from "../views/ConfigView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "services",
      component: ServicesView,
    },
    {
      path: "/docker",
      name: "docker",
      component: DockerLocalView,
    },
    {
      path: "/docker-nodes",
      name: "docker-nodes",
      component: DockerNodes,
    },
    {
      path: "/system",
      name: "system",
      component: SystemView,
    },
    {
      path: "/config",
      name: "config",
      component: ConfigView,
    },
  ],
});

export default router;
