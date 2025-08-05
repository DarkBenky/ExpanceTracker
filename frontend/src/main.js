import { createApp } from 'vue'
import App from './App.vue'

const app =  createApp(App)

// Add global properties
app.config.globalProperties.$apiUrl = 'http://localhost:1234/'

app.mount('#app')
