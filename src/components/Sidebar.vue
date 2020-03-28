<template>
    <v-list>
        <v-subheader>Salle "{{ game.Code }}"</v-subheader>
        <template v-for="(team, teamName) in game.Players">
            <div :key="teamName">
                <v-subheader :class="`${teamName} white--text text-uppercase`" layout="row">
                    <span flex>{{ teamName }}</span>
                    <v-chip v-if="game.CurrentTeam === teamName" small :color="`${teamName} lighten-3`">
                        A VOTRE TOUR !
                    </v-chip>
                </v-subheader>
                <v-list-item v-for="p in team" :key="p.ID">
                    <v-list-item-avatar>
                        <v-icon :color="teamName">mdi-account</v-icon>
                    </v-list-item-avatar>
                    <v-list-item-content>
                        <v-list-item-title :class="{'font-weight-bold': p.Name === player.Name}">{{ p.Name }}</v-list-item-title>
                    </v-list-item-content>
                    <v-list-item-action>
                        <v-icon color="orange" v-if="p.Spy">mdi-incognito</v-icon>
                    </v-list-item-action>
                </v-list-item>
            </div>
        </template>

        <v-divider />

        <v-list-item @click="nextTurn" v-if="game.CurrentTeam === player.Team && player.Spy" class="primary--text">
            <v-list-item-avatar>
                <v-icon>mdi-close</v-icon>
            </v-list-item-avatar>
            Finir le tour
        </v-list-item>
    </v-list>
</template>

<script>
    import {mapGetters} from 'vuex'

    export default {
        name: 'Sidebar',
        computed: {
            ...mapGetters(['game', 'player'])
        },

        data() {
            return {
            }
        },

        methods: {
            nextTurn() {
                this.$socket.client.emit('NextTurn');
            }
        }
    }
</script>