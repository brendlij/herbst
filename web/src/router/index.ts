import { createRouter, createWebHistory } from "vue-router";
import * as ServicesView from "../views/ServicesView.vue";
import * as DockerLocalView from "../views/DockerLocalView.vue";
import * as DockerNodes from "../views/DockerNodes.vue";
import * as ConfigView from "../views/ConfigView.vue";

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
      path: "/config",
      name: "config",
      component: ConfigView,
    },
  ],
});

export default router;
