import Vue from 'vue'
import App from './App.vue'
import vuetify from '@/plugins/vuetify'
import store from '@/plugins/store'
import "@/plugins/socketio";
import '@/plugins/snotify'
import "vue-snotify/styles/material.scss"

import './styles/flexbox.css'

Vue.config.productionTip = false

new Vue({
  vuetify,
  store,
  render: h => h(App),
}).$mount('#app')
