import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        player: {Name: '', GameCode: '', Team: ''},
        game: {}
    },
    getters: {
        player: state => state.player,
        game: state => state.game,
    },
    mutations: {
        SET_PLAYER(state, p) {
            state.player = p;
        },
        SET_GAME(state, game) {
            state.game = game;
        }
    },
    actions: {
        setPlayer({commit}, p) {
            commit('SET_PLAYER', p);
        },
        setGame({commit}, g) {
            commit('SET_GAME', g);
        }
    }

})