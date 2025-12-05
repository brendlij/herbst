import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import router from "./router";
import "prismjs";
import "prismjs/themes/prism-tomorrow.css"; // oder anderes Theme
import "prismjs/components/prism-toml";

createApp(App).use(router).mount("#app");
