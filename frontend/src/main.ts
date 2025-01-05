import { createApp } from "vue";
import PrimeVue from "primevue/config";
import Aura from '@primevue/themes/aura';
import ToastService from 'primevue/toastservice';

import App from './App.vue'
import router from './router'
import './assets/main.css'


const app = createApp(App)

app.use(router)
app.use(PrimeVue, {
    // Default theme configuration
    theme: {
        preset: Aura,
        options: {
            prefix: 'p',
            darkModeSelector: 'system',
            cssLayer: false
        }
    }
 });
app.use(ToastService);

app.mount('#app')
