<template>
    <div flex="20" style="height: 20%" class="pa-3" @click="onClickWord">
        <v-card layout="column" layout-align="center center"
                :class="['fill-height', 'text-uppercase', textColor+'--text', 'font-weight-bold', 'lighten-2', word.HintedBy.length > 0 ? 'hinted hinted-'+hintedColor : '']"
                :hover="!player.Spy" :outlined="!hinted" :color="cardColor">
            <v-img v-if="word.Found" :src="`/${word.Color}.svg`" style="position: absolute; opacity: 0.3" contain max-height="80%" max-width="80%" />
            <span :style="{opacity: word.Found ? 0.5 : 1}">{{ word.Word }}</span>
            <div layout="row" v-if="word.HintedBy.length > 0">
                <v-chip v-for="h in word.HintedBy" :key="'hint_'+h">{{ h }}</v-chip>
            </div>
        </v-card>
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'

    export default {
        name: 'Word',
        props: {
            word: {
                type: Object,
                required: true
            }
        },

        computed: {
            ...mapGetters(['player', 'game']),

            cardColor() {
                if (this.player.Spy || this.word.Found) {
                    return this.word.Color;
                } else {
                    return '';
                }
            },

            textColor() {
                if ((this.player.Spy || this.word.Found) && this.word.Color === 'black') {
                    return 'white'
                } else {
                    return 'black';
                }
            },

            hintedColor() {
                if (this.word.HintedBy.length <= 0) {
                    return ''
                }

                let hintedName = this.word.HintedBy[0];
                for (let team in this.game.Players) {
                    for (let player in team) {
                        if (player.Name === hintedName) {
                            return player.Team;
                        }
                    }
                }

                return '';
            }
        },

        methods: {
            onClickWord() {
                if (!this.player.Spy && this.player.Team === this.game.CurrentTeam) {
                    this.$socket.client.emit('WordSelected', {Word: this.word.Word})
                }
            }
        }
    }
</script>

<style scoped>
    .hinted {
        border: 2px solid;
    }
    .hinted-blue {
        border-color: #33B5E5;
    }
    .hinted-red {
        border-color: #FF4444;
    }
</style>