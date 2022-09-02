import { createApp } from 'vue'
import Antd, { message  } from 'ant-design-vue';
import App from './App.vue'
import store from './store'

import 'ant-design-vue/dist/antd.css';
import './assets/iconfont/iconfont.css';
import './assets/css/index.less';
import { LoadingOutlined } from '@ant-design/icons-vue';

message.config({
  top: `70px`,
  maxCount: 3,
  duration: 2,
});

const app = createApp(App)

app.use(store)
app.use(Antd)

app.component('LoadingOutlined', LoadingOutlined)

app.mount('#app')
