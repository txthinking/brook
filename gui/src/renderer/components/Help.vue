<template>
    <v-layout>
        <v-flex xs12 sm12 md6>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <h4>Step 0</h4>
                            <p>Read the wiki first<v-btn flat small color="primary" href="https://github.com/txthinking/brook/wiki" target="_blank">Wiki</v-btn></p>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <h4>Step 1</h4>
                            <p>Search for anwsers<v-btn flat small color="primary" href="https://github.com/txthinking/brook/issues?utf8=✓&q=" target="_blank">Issues</v-btn><v-btn flat small color="primary" @click="open('https://www.google.com')">Google</v-btn></p>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <h4>Step 2</h4>
                            <p>Submit a new issue<v-btn flat small color="primary" href="https://github.com/txthinking/brook/issues/new" target="_blank">Issue</v-btn></p>
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
