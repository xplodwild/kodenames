<template>
    <v-dialog v-model="startGameDialog" max-width="600px" persistent>
        <v-card>
            <v-card-title>Rejoindre ou créer une partie</v-card-title>

            <v-card-text>
                <v-text-field outlined label="Code de la partie" prepend-icon="mdi-tag" v-model="playerInfo.GameCode"></v-text-field>
                <v-text-field outlined label="Nom du joueur" prepend-icon="mdi-account" v-model="playerInfo.Name"></v-text-field>

                <div class="subtitle-1">Choisissez une équipe :</div>
                <v-item-group mandatory v-model="playerInfo.Team">
                    <div layout="row">
                        <v-item v-slot:default="{ active, toggle }" v-for="t in teams" :key="t" flex :value="t">
                            <v-card :color="active ? t : ''" class="d-flex align-center"
                                    dark height="100" @click="toggle">
                                <div class="display-3 flex-grow-1 text-center text-uppercase" v-if="active">
                                    {{ t }}
                                </div>
                                <div class="flex-grow-1 text-center text-uppercase" v-else>
                                    {{ t }}
                                </div>
                            </v-card>
                        </v-item>
                    </div>
                </v-item-group>
            </v-card-text>

            <v-card-actions layout="row" layout-align="end center">
                <v-btn text :color="playerInfo.Team" @click="join">
                    <v-icon left>mdi-arrow-right</v-icon>
                    Rejoindre la partie
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script>
    import {mapActions} from 'vuex'

    export default {
        name: 'StartGameDialog',
        data() {
            return {
                startGameDialog: true,
                teams: ['red', 'blue'],
                playerInfo: {
                    Name: '',
                    Team: '',
                    GameCode: ''
                }
            }
        },

        sockets: {
            JoinAccepted (data) {
                console.log("INFO: Join Accepted");
                this.setPlayer(data.Player);
                this.setGame(data.Game);
                this.$emit('joined');
                this.startGameDialog = false;
            }
        },

        methods: {
            ...mapActions(['setGame', 'setPlayer']),

            join() {
                this.$socket.client.emit('JoinGame', this.playerInfo);
            }
        }
    }
</script>