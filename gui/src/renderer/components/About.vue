<template>
    <v-layout>
        <v-flex xs12 sm12 md6>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <h2>Brook</h2>
                            <p><a href="https://github.com/txthinking/brook" target="_blank">https://github.com/txthinking/brook</a></p>
                            <p>Brook is an open source cross-platform(Linux/MacOS/Windows/Android/iOS) proxy software</p>
                            <p>Licensed under The <a href="https://github.com/txthinking/brook/blob/master/LICENSE" target="_blank">GPLv3</a> License. <a href="https://github.com/txthinking/brook/blob/master/OPENSOURCELICENSES" target="_blank">Open source licenses</a></p>
                            <p>Author: <a href="https://github.com/txthinking" target="_blank">txthinking</a> (<a href="mailto:cloud@txthinking.com" target="_blank">cloud@txthinking.com</a>)</p>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
        </v-flex>
        <v-snackbar v-model="hey">
            {{girl}}
        </v-snackbar>
    </v-layout>
</template>
<script>
export default {
    data: () => ({
        o: {
        },
        hey: false,
        girl: "",
    }),

    computed: {
    },

    watch: {
    },

    created () {
        this.initialize()
    },

    methods: {
        initialize () {
        },
    },
}
</script>
