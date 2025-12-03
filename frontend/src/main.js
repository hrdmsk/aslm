import { createApp } from 'vue'
import { createPinia } from 'pinia' // 追加
import App from './App.vue'

const app = createApp(App)

app.use(createPinia()) // Piniaを使用するように設定
app.mount('#app')