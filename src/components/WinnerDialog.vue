<template>
    <v-dialog v-model="value" max-width="600px" persistent>
        <v-card>
            <v-card-title>Fin de la partie !</v-card-title>

            <v-card-text>
                <v-img :src="`/${game.Winner}.svg`" max-width="500px" max-height="300px" contain />
                <div :class="`title text-center ${game.Winner}--text`">L'équipe {{ game.Winner }} a gagné !</div>

                <div class="text-center">
                    La partie sera remise à zéro dans 10 secondes, les maître-espions seront changés.
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script>
    import {mapGetters} from 'vuex'

    export default {
        name: 'WinnerDialog',
        props: {
            value: {
                type: Boolean,
                required: true
            }
        },

        computed: {
            ...mapGetters(['game'])
        }
    }
</script>