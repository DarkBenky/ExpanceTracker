import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App)

// Add global properties
app.config.globalProperties.$apiUrl = 'http://localhost:1234/'

// Disable development tools and console logs in production
if (process.env.NODE_ENV === 'production') {
    app.config.performance = false
    app.config.devtools = false
    
    // Disable console logs in production
    console.log = () => {}
    console.warn = () => {}
    console.info = () => {}
}

app.mount('#app')