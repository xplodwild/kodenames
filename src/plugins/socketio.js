import VueSocketIOExt from 'vue-socket.io-extended';
import io from 'socket.io-client';
import Vue from 'vue'

const socket = io('http://localhost:8000/');

Vue.use(VueSocketIOExt, socket);