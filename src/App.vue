<template>
    <v-app>
        <v-navigation-drawer :value="true" app permanent width="200">
            <sidebar />
        </v-navigation-drawer>

        <v-content>
            <word-grid />
        </v-content>

        <v-footer color="indigo" app>
            <span class="white--text">&copy; 2020 Guigui</span>
        </v-footer>
        <vue-snotify />

        <start-game-dialog />
        <winner-dialog v-model="showWinnerDialog" />
    </v-app>
</template>

<script>

    import Sidebar from "@/components/Sidebar";
    import WordGrid from "@/components/WordGrid";
    import StartGameDialog from "@/components/StartGameDialog";
    import {mapActions, mapGetters} from 'vuex'
    import WinnerDialog from "@/components/WinnerDialog";

    export default {
        name: 'App',
        components: {WinnerDialog, StartGameDialog, WordGrid, Sidebar},

        sockets: {
            GameUpdate (game) {
                this.setGame(game);
            },
            Notice (msg) {
                this.$snotify[msg.Type](msg.Message)
            }
        },

        computed: {
            ...mapGetters(['game', 'player']),

            showWinnerDialog() {
                if (this.game && this.game.Winner) {
                    return true;
                } else {
                    return false;
                }
            }
        },

        watch: {
            game: {
                deep: true,
                handler (val) {
                    // Check if our spy status changed
                    for (let teamName in val.Players) {
                        for (let i = 0; i < val.Players[teamName].length; ++i) {
                            const player = val.Players[teamName][i];
                            if (player.Name === this.player.Name) {
                                this.setPlayer(player);
                                break;
                            }
                        }
                    }
                }
            }
        },

        methods: {
            ...mapActions(['setGame', 'setPlayer'])
        }
    }
</script>
