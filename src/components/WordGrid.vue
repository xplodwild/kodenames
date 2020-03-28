<template>
    <div layout="row" layout-wrap class="fill-height">
        <word :word="w" v-for="w in words" :key="w.Word"></word>
    </div>
</template>

<script>
    import Word from "@/components/Word";
    import {mapGetters} from 'vuex';

    export default {
        components: {Word},
        computed: {
            ...mapGetters(['game']),

            words() {
                return this.game.Words;
            }
        },

        data() {
            return {}
        }
    }
</script>