<template>
    <v-layout>
        <v-flex xs12 sm12 md6>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-switch
                                label="Use white icon on tray"
                                v-model="o.UseWhiteTrayIcon"
                                ></v-switch>
                            <p>Need to restart</p>
                        </v-flex>
                    </v-layout>
                </v-container>
            </v-card>
            <v-card>
                <v-container fluid>
                    <v-layout column>
                        <v-flex>
                            <v-btn class="info" @click="save">Save</v-btn>
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
            UseWhiteTrayIcon: false,
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
            var s = localStorage.getItem('brook/builtin');
            if (s){
                this.o = JSON.parse(s);
            }
        },
        save () {
            localStorage.setItem('brook/builtin', JSON.stringify(this.o));
            this.girl = "OK";
            this.hey = true;
        },
    },
}
</script>
