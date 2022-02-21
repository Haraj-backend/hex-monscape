import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./routes";
import "./assets/app.css";

createApp(App).use(router).use(createPinia()).mount("#app");
