import Snotify, { SnotifyPosition } from 'vue-snotify'
import Vue from 'vue'

Vue.use(Snotify, {
    toast: {
        position: SnotifyPosition.centerBottom,
        timeout: 3000
    }
})