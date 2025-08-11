import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App)

// Use environment variable for API URL
const apiUrl = process.env.VUE_APP_API_URL || 'http://localhost:1234/'
app.config.globalProperties.$apiUrl = apiUrl

console.log('API URL:', apiUrl)
console.log('NODE_ENV:', process.env.NODE_ENV)

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